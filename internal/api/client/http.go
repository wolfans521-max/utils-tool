// Package client 提供API客户端实现
// HTTP客户端用于执行HTTP/HTTPS请求
package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

// HTTPClient HTTP客户端
type HTTPClient struct {
	client    *http.Client
	transport *http.Transport
}

// HTTPTransportConfig HTTP Transport 连接池配置
type HTTPTransportConfig struct {
	MaxIdleConns          int // 最大空闲连接数，默认500
	MaxIdleConnsPerHost   int // 每个主机最大空闲连接数，默认100
	ResponseHeaderTimeout int // 响应头超时时间（毫秒），默认1000
}

// HTTPRequestOption HTTP请求选项
type HTTPRequestOption struct {
	URL         string                 // 请求URL
	Method      string                 // 请求方法
	Headers     map[string]string      // 请求头
	Args        map[string]interface{} // 请求参数
	Timeout     time.Duration          // 超时时间
	ContentJSON bool                   // 是否JSON格式
}

// HTTPRequestTiming HTTP请求时间记录结构
type HTTPRequestTiming struct {
	TotalTime   time.Duration // 总耗时
	ConnectTime time.Duration // 连接时间（包含DNS解析、TCP连接、TLS握手）
	SendTime    time.Duration // 请求发送时间
	WaitTime    time.Duration // 等待响应时间（TTFB）
	ReadTime    time.Duration // 数据读取时间
}

// HTTPRequestResult HTTP请求结果
type HTTPRequestResult struct {
	StatusCode  int                // HTTP状态码
	ResponseLen int                // 响应长度
	NodeLog     string             // 节点日志
	Err         string             // 错误信息
	TryTimes    int                // 实际重试次数
	ConnectTime time.Duration      // 连接时间
	ReadTime    time.Duration      // 数据读取时间
	Timing      *HTTPRequestTiming // 详细时间记录
}

// NewHTTPClient 创建HTTP客户端
// transportConfig: HTTP Transport 连接池配置，如果为 nil 则使用默认值
func NewHTTPClient(timeout time.Duration, transportConfig HTTPTransportConfig) *HTTPClient {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	// 设置默认值
	if transportConfig.MaxIdleConns <= 0 {
		transportConfig.MaxIdleConns = 500
	}
	if transportConfig.MaxIdleConnsPerHost <= 0 {
		transportConfig.MaxIdleConnsPerHost = 100
	}
	if transportConfig.ResponseHeaderTimeout <= 0 {
		transportConfig.ResponseHeaderTimeout = 1000
	}

	// 创建共享的 Transport，支持连接复用
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,  // 默认连接超时
			KeepAlive: 30 * time.Second, // TCP keepalive
		}).DialContext,
		MaxIdleConns:          transportConfig.MaxIdleConns,                                            // 最大空闲连接数
		MaxIdleConnsPerHost:   transportConfig.MaxIdleConnsPerHost,                                     // 每个主机最大空闲连接数
		IdleConnTimeout:       90 * time.Second,                                                        // 空闲连接超时
		TLSHandshakeTimeout:   10 * time.Second,                                                        // TLS握手超时
		ExpectContinueTimeout: 1 * time.Second,                                                         // 100-continue超时
		ResponseHeaderTimeout: time.Duration(transportConfig.ResponseHeaderTimeout) * time.Millisecond, // 响应头超时
		DisableKeepAlives:     false,                                                                   // 启用连接复用
	}

	return &HTTPClient{
		transport: transport,
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

// GetClient 获取底层HTTP客户端
func (c *HTTPClient) GetClient() *http.Client {
	return c.client
}

// GetTransport 获取底层HTTP Transport
func (c *HTTPClient) GetTransport() *http.Transport {
	return c.transport
}

// DoRequest 执行HTTP请求（支持重试）
func (c *HTTPClient) DoRequest(ctx context.Context, opt *HTTPRequestOption, maxRetries int) (map[string]interface{}, *HTTPRequestResult, error) {
	var lastErr error
	if maxRetries < 0 {
		maxRetries = 0
	}

	reqResult := &HTTPRequestResult{
		TryTimes: 1,
		Timing:   &HTTPRequestTiming{}, // 初始化详细时间记录
	}

	// 尝试执行请求（初次请求 + 重试次数）
	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, respLen, statusCode, nodeLog, timing, err := c.DoRequestOnce(ctx, opt)
		reqResult.ResponseLen = respLen
		reqResult.StatusCode = statusCode
		reqResult.NodeLog = nodeLog
		reqResult.TryTimes = attempt + 1

		// 设置详细时间记录
		if timing != nil {
			reqResult.Timing = timing
		}

		if err == nil {
			// 请求成功，返回结果
			reqResult.Err = "success"
			return result, reqResult, nil
		}

		lastErr = err

		// 如果还有重试机会，短暂延迟后重试
		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt+1) * 10 * time.Millisecond)
		}
	}

	reqResult.Err = lastErr.Error()
	return nil, reqResult, lastErr
}

// DoRequestOnce 执行单次HTTP请求
// 返回值: result, responseLen, statusCode, nodeLog, timing, error
func (c *HTTPClient) DoRequestOnce(ctx context.Context, opt *HTTPRequestOption) (map[string]interface{}, int, int, string, *HTTPRequestTiming, error) {
	// 构建请求URL
	reqURL := opt.URL
	var reqBody io.Reader

	// 从URL中提取节点日志信息
	nodeLog := "-"
	if parsedURL, err := url.Parse(reqURL); err == nil {
		nodeLog = parsedURL.Host
	}
	// http请求nodeLog默认为"-"
	nodeLog = "-"

	// 复制headers以避免修改原始map
	headers := make(map[string]string)
	for k, v := range opt.Headers {
		headers[k] = v
	}

	// 处理参数
	if opt.ContentJSON {
		// JSON格式
		jsonData, err := sonic.Marshal(opt.Args)
		if err != nil {
			return nil, 0, 0, nodeLog, &HTTPRequestTiming{}, fmt.Errorf("序列化参数失败: %w", err)
		}
		reqBody = strings.NewReader(string(jsonData))
		headers["Content-Type"] = "application/json"
	} else {
		// URL编码格式
		values := url.Values{}
		for k, v := range opt.Args {
			values.Add(k, fmt.Sprintf("%v", v))
		}
		if opt.Method == "GET" {
			if strings.Contains(reqURL, "?") {
				reqURL += "&" + values.Encode()
			} else {
				reqURL += "?" + values.Encode()
			}
		} else {
			reqBody = strings.NewReader(values.Encode())
			headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
	}

	// 创建带超时的context
	reqCtx, cancel := context.WithTimeout(ctx, opt.Timeout)
	defer cancel()

	// 记录请求开始时间
	requestStart := time.Now()

	// 用于记录详细时间的变量
	timing := &HTTPRequestTiming{}

	// 记录各个阶段的时间点
	var dnsStart, dnsDone, connectStart, connectDone, tlsStart, tlsDone, wroteRequest, gotFirstByte time.Time

	// 创建 httptrace 来追踪连接时间
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			dnsDone = time.Now()
		},
		ConnectStart: func(network, addr string) {
			connectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			connectDone = time.Now()
		},
		TLSHandshakeStart: func() {
			tlsStart = time.Now()
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			tlsDone = time.Now()
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			wroteRequest = time.Now()
			// 计算请求发送时间
			if !connectDone.IsZero() {
				timing.SendTime = wroteRequest.Sub(connectDone)
			} else {
				timing.SendTime = wroteRequest.Sub(requestStart)
			}
		},
		GotFirstResponseByte: func() {
			gotFirstByte = time.Now()
			// 计算等待响应时间（TTFB）
			if !wroteRequest.IsZero() {
				timing.WaitTime = gotFirstByte.Sub(wroteRequest)
			} else {
				timing.WaitTime = gotFirstByte.Sub(requestStart)
			}
		},
	}

	req, err := http.NewRequestWithContext(httptrace.WithClientTrace(reqCtx, trace), opt.Method, reqURL, reqBody)
	if err != nil {
		return nil, 0, 0, nodeLog, &HTTPRequestTiming{}, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 执行请求
	resp, err := c.client.Do(req)
	if err != nil {
		// 计算总耗时
		totalTime := time.Since(requestStart)
		timing.TotalTime = totalTime

		// 计算连接时间（包括DNS解析、TCP连接和TLS握手）
		var totalConnectTime time.Duration
		if !dnsStart.IsZero() && !dnsDone.IsZero() {
			totalConnectTime += dnsDone.Sub(dnsStart)
		}
		if !connectStart.IsZero() && !connectDone.IsZero() {
			totalConnectTime += connectDone.Sub(connectStart)
		}
		if !tlsStart.IsZero() && !tlsDone.IsZero() {
			totalConnectTime += tlsDone.Sub(tlsStart)
		}
		// 如果没有记录到任何连接阶段时间，使用总时间作为连接时间
		if totalConnectTime == 0 {
			totalConnectTime = totalTime
		}
		timing.ConnectTime = totalConnectTime

		// 计算等待响应时间（如果已有记录则使用，否则用总时间减去连接时间）
		if !gotFirstByte.IsZero() {
			if !wroteRequest.IsZero() {
				timing.WaitTime = gotFirstByte.Sub(wroteRequest)
			} else {
				timing.WaitTime = gotFirstByte.Sub(requestStart)
			}
		} else {
			// 没有收到响应，总时间减去连接时间即为等待时间
			if totalTime > timing.ConnectTime {
				timing.WaitTime = totalTime - timing.ConnectTime
			}
		}

		// readTime = 总时间 - connectTime - 等待时间
		if totalTime > timing.ConnectTime+timing.WaitTime {
			timing.ReadTime = totalTime - timing.ConnectTime - timing.WaitTime
		}

		return nil, 0, 502, nodeLog, timing, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	// 计算连接时间（包括DNS解析、TCP连接和TLS握手）
	var totalConnectTime time.Duration
	if !dnsStart.IsZero() && !dnsDone.IsZero() {
		totalConnectTime += dnsDone.Sub(dnsStart)
	}
	if !connectStart.IsZero() && !connectDone.IsZero() {
		totalConnectTime += connectDone.Sub(connectStart)
	}
	if !tlsStart.IsZero() && !tlsDone.IsZero() {
		totalConnectTime += tlsDone.Sub(tlsStart)
	}
	timing.ConnectTime = totalConnectTime

	// 记录读取开始时间
	readStart := time.Now()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// 读取失败时，计算已经花费的读取时间
		timing.ReadTime = time.Since(readStart)
		timing.TotalTime = time.Since(requestStart)
		return nil, 0, statusCode, nodeLog, timing, fmt.Errorf("读取响应失败: %w", err)
	}

	// 计算数据读取时间（从收到第一个字节到读取完成）
	timing.ReadTime = time.Since(readStart)
	timing.TotalTime = time.Since(requestStart)

	responseLen := len(body)

	// 解析JSON - 兼容返回数据可能是 map 或 slice 的情况
	result, err := ParseJSONResponse(body)

	if err != nil {
		return nil, responseLen, statusCode, nodeLog, timing, fmt.Errorf("解析响应失败: %w", err)
	}

	return result, responseLen, statusCode, nodeLog, timing, nil
}

// ParseJSONResponse 解析JSON响应，兼容返回数据可能是 map 或 slice 的情况
// 如果是 slice，则包装为 {"data": slice} 的形式返回
func ParseJSONResponse(body []byte) (map[string]interface{}, error) {
	// 先尝试解析为 map
	var mapResult map[string]interface{}
	if err := sonic.Unmarshal(body, &mapResult); err == nil {
		return mapResult, nil
	}

	// 如果解析为 map 失败，尝试解析为 slice
	var sliceResult []interface{}
	if err := sonic.Unmarshal(body, &sliceResult); err == nil {
		// 将 slice 包装为 map 返回
		return map[string]interface{}{
			"result": sliceResult,
		}, nil
	} else {
		// 两种解析都失败，返回真实错误（非 timeout）
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
}
