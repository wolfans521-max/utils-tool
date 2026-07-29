// Package api 提供API请求管理的核心实现
package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"math"
	"net/http"
	"net/url"
	urlpkg "net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"

	wbcontext "git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/discovery"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/log"
	"git.intra.weibo.com/search_fe/wbutil-go/internal/api/client"
	"git.intra.weibo.com/search_fe/wbutil-go/internal/config"
	"github.com/goccy/go-yaml"
	"github.com/spf13/viper"
)

var (
	motanInitOnce sync.Once
	motanInitErr  error
	// daNodeIndex 用于DA节点的轮询负载均衡
	daNodeIndex uint32
)

// Protocol 协议类型
type Protocol string

const (
	ProtocolHTTP   Protocol = "http"
	ProtocolHTTPS  Protocol = "https"
	ProtocolMotan  Protocol = "motan"
	ProtocolMotan2 Protocol = "motan2"
	ProtocolDA     Protocol = "da"
)

// PlatformInterfaceServices 平台接口服务名列表（Motan协议使用）
var PlatformInterfaceServices = []string{
	"com.weibo.api.FeedCoreStatusHttpService",
}

// PlatformAPIDomains 平台接口域名列表（HTTP/HTTPS协议使用）
var PlatformAPIDomains = []string{"i.api.weibo.com", "i2.api.weibo.com"}

// APIConfig 单个API配置
type APIConfig struct {
	URL         string            `yaml:"url"`          // 请求URL
	GrayURL     string            `yaml:"gray_url"`     // 灰度URL
	Auth        bool              `yaml:"auth"`         // 是否需要认证
	Type        string            `yaml:"type"`         // 协议类型
	Method      string            `yaml:"method"`       // 请求方法
	Headers     map[string]string `yaml:"headers"`      // 请求头
	Timeout     int               `yaml:"timeout"`      // 超时时间（毫秒）
	Service     string            `yaml:"service"`      // 服务名
	Group       string            `yaml:"group"`        // 服务组
	Version     int               `yaml:"version"`      // 版本号
	TryTimes    int               `yaml:"try_times"`    // 重试次数
	Keepalive   bool              `yaml:"keepalive"`    // 是否保持连接
	ContentJSON bool              `yaml:"content_json"` // 是否JSON格式
	HTTPMethod  string            `yaml:"http_method"`  // HTTP 方法，仅 Motan 协议使用
}

// Node 服务节点
type Node struct {
	Host    string                 `json:"host"`
	Port    int                    `json:"port"`
	ExtInfo map[string]interface{} `json:"extInfo"`
}

// GrayRule 灰度规则
type GrayRule struct {
	UIDs  []int64 `json:"uids"`  // 指定用户ID列表
	Pos   int     `json:"pos"`   // 位置
	Num   int     `json:"num"`   // 数字位数
	Value []int   `json:"value"` // 匹配值
}

// RequestOption 请求选项
type RequestOption struct {
	URLKey      string                 // API配置键
	Protocol    Protocol               // 协议
	URL         string                 // 请求URL
	Method      string                 // 请求方法
	Headers     map[string]string      // 请求头
	Args        map[string]interface{} // 请求参数
	Timeout     time.Duration          // 超时时间
	Service     string                 // 服务名
	Group       string                 // 服务组
	Version     int                    // 版本号
	TryTimes    int                    // 重试次数
	Keepalive   bool                   // 是否保持连接
	ContentJSON bool                   // 是否JSON格式
	Nodes       []Node                 // 服务节点列表
	OriginNodes []Node                 // 原始节点列表（灰度前）
	WebDegrade  int                    // 降级等级
	SeqID       string                 // seqid
	NodeLog     string                 // 节点日志（host:port 或 "-"）
	HTTPMethod  string                 // HTTP 方法，仅 Motan 协议使用
}

// RequestTiming 请求时间记录结构
// 注意：不同协议的时间记录精度不同
// - HTTP/DA 协议：可以精确追踪 ConnectTime、SendTime、WaitTime、ReadTime
// - Motan 协议：由于通过Mesh代理且BaseCall是同步阻塞调用，只能记录 TotalTime
type RequestTiming struct {
	TotalTime   time.Duration // 总耗时
	ConnectTime time.Duration // 连接时间（包含DNS解析、TCP连接、TLS握手），仅HTTP/DA有效
	SendTime    time.Duration // 请求发送时间，仅HTTP/DA有效
	WaitTime    time.Duration // 等待响应时间（TTFB - Time To First Byte），仅HTTP/DA有效
	ReadTime    time.Duration // 响应读取时间，仅HTTP/DA有效
}

// RequestResult 请求结果信息（用于日志记录）
type RequestResult struct {
	StatusCode  int            // HTTP状态码
	ResponseLen int            // 响应长度
	Cost        time.Duration  // 请求耗时
	TryTimes    int            // 实际重试次数
	Err         string         // 错误信息
	NodeLog     string         // 节点日志（host:port 或 "-"）
	Timing      *RequestTiming // 详细时间记录
}

// BatchRequest 批量请求项
type BatchRequest struct {
	URLKey string                 // API配置键
	Params map[string]interface{} // 请求参数
	Other  *OtherOption           // 其他选项
}

// OtherOption 其他请求选项
type OtherOption struct {
	Header      map[string]string // 额外请求头
	Timeout     time.Duration     // 超时时间
	ContentJSON bool              // 是否使用JSON格式发送数据
	TryTimes    int               // 重试次数（优先级高于配置文件）
	HTTPMethod  string            // HTTP 方法，仅 Motan 协议使用
}

// 回调函数类型定义
type (
	// DegradeFunc 降级函数
	DegradeFunc func(key string) int
)

// Manager API管理器
type Manager struct {
	apisConfig          map[string]*APIConfig // API配置
	source              string                // 来源标识
	tokenFile           string                // token文件路径
	defaultMethod       string                // 默认请求方法
	envGroup            string                // 环境组
	env                 string                // 环境
	degradeKeyPrefix    string                // 降级键前缀
	degradeFunc         DegradeFunc           // 降级函数
	httpClient          *client.HTTPClient    // HTTP客户端（复用连接）
	motanConfig         *client.MotanConfig   // Motan配置（延迟初始化）
	batchMaxConcurrency int                   // 批量请求最大并发数
	breakerGroup        *client.BreakerGroup  // 下游熔断器组（每个 urlKey 一个熔断器）
	mu                  sync.RWMutex          // 读写锁
}

// Config 管理器配置
type Config struct {
	APIsConfig          map[string]*APIConfig      // API配置
	Source              string                     // 来源标识
	TokenFile           string                     // token文件路径
	DefaultMethod       string                     // 默认请求方法
	EnvGroup            string                     // 环境组
	Env                 string                     // 环境
	DegradeKeyPrefix    string                     // 降级键前缀
	DegradeFunc         DegradeFunc                // 降级函数
	HTTPTimeout         time.Duration              // HTTP超时时间
	MotanConfig         *client.MotanConfig        // Motan MeshClient配置
	BatchMaxConcurrency int                        // 批量请求最大并发数，默认10
	HTTPTransport       client.HTTPTransportConfig // HTTP Transport 连接池配置
	BreakerConfig       client.BreakerConfig       // 熔断器配置
}

// NewManager 创建API管理器
func NewManager(config *Config) (*Manager, error) {
	if config == nil {
		config = &Config{}
	}

	// 设置默认值
	if config.DefaultMethod == "" {
		config.DefaultMethod = "GET"
	}
	if config.Source == "" {
		config.Source = "2936099636"
	}
	if config.EnvGroup == "" {
		config.EnvGroup = "yf"
	}
	if config.DegradeKeyPrefix == "" {
		config.DegradeKeyPrefix = "api_degrade_"
	}
	if config.HTTPTimeout == 0 {
		config.HTTPTimeout = 30 * time.Second
	}
	if config.BatchMaxConcurrency <= 0 {
		config.BatchMaxConcurrency = 10
	}
	// HTTP Transport 连接池配置默认值
	if config.HTTPTransport.MaxIdleConns <= 0 {
		config.HTTPTransport.MaxIdleConns = 500
	}
	if config.HTTPTransport.MaxIdleConnsPerHost <= 0 {
		config.HTTPTransport.MaxIdleConnsPerHost = 100
	}
	if config.HTTPTransport.ResponseHeaderTimeout <= 0 {
		config.HTTPTransport.ResponseHeaderTimeout = 1000
	}

	// 初始化熔断器配置（未配置时使用默认值）
	breakerConfig := config.BreakerConfig
	if breakerConfig.FailureThreshold <= 0 {
		breakerConfig = client.DefaultBreakerConfig()
	}

	manager := &Manager{
		apisConfig:          config.APIsConfig,
		source:              config.Source,
		tokenFile:           config.TokenFile,
		defaultMethod:       config.DefaultMethod,
		envGroup:            config.EnvGroup,
		env:                 config.Env,
		degradeKeyPrefix:    config.DegradeKeyPrefix,
		degradeFunc:         config.DegradeFunc,
		motanConfig:         config.MotanConfig,
		httpClient:          client.NewHTTPClient(config.HTTPTimeout, config.HTTPTransport),
		batchMaxConcurrency: config.BatchMaxConcurrency,
		breakerGroup:        client.NewBreakerGroup(breakerConfig),
	}

	return manager, nil
}

// Request 单个请求
func (m *Manager) Request(ctx *wbcontext.Context, urlKey string, params map[string]interface{}, other *OtherOption) (map[string]interface{}, error) {
	if other == nil {
		other = &OtherOption{}
	}

	// 🔥 请求级超时检查：如果 context 已取消（request_timeout 到期），跳过下游调用
	// 避免在 deadline 已过后还发起新的 HTTP/Motan/DA 调用
	if ctx != nil && ctx.GinContext() != nil && ctx.GinContext().Request != nil {
		if err := ctx.GinContext().Request.Context().Err(); err != nil {
			return nil, fmt.Errorf("request context cancelled: %w", err)
		}
	}

	// 构建请求选项
	opt, err := m.setRequest(ctx, urlKey, params, other)

	if err != nil {
		return nil, err
	}

	// 🔥 熔断器检查：如果该下游服务已熔断，快速失败
	breaker := m.breakerGroup.Get(urlKey)
	if err := breaker.Allow(); err != nil {
		// 熔断器已打开，快速失败
		reqResult := &RequestResult{
			Cost:       0,
			Err:        "circuit_breaker_open",
			TryTimes:   0,
			StatusCode: 503,
		}
		m.logRequest(ctx, opt, reqResult)
		log.Append(ctx, fmt.Sprintf(" custom_counter:%s_total", urlKey))
		return nil, fmt.Errorf("circuit breaker [%s] open: %w", urlKey, err)
	}

	// 记录请求开始时间
	startTime := time.Now()

	// 执行请求
	result, reqResult, err := m.doRequestWithResult(ctx, opt)

	// 计算请求耗时
	cost := time.Since(startTime)
	reqResult.Cost = cost

	// 🔥 熔断器结果记录：只有基础设施级故障才触发熔断
	// 故障定义：
	//   - 网络错误（timeout, connection refused, DNS failure）→ 记录失败
	//   - HTTP 5xx → 记录失败（下游服务不可用）
	//   - HTTP 4xx → 不记录失败（客户端错误，非下游故障）
	//   - 业务级错误（auth failed, rate limit）→ 不记录失败
	//   - 响应解析错误 → 不记录失败（可能是自身 bug）
	if isInfrastructureFailure(err, reqResult) {
		breaker.RecordFailure()
	} else {
		breaker.RecordSuccess()
	}

	// 记录网络请求日志（与 Lua 版本格式一致）
	// 格式: network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
	m.logRequest(ctx, opt, reqResult)

	// 记录自定义计数器日志
	logMsg := fmt.Sprintf(" custom_counter:%s_total", urlKey)
	log.Append(ctx, logMsg)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// isInfrastructureFailure 判断是否为基础设施级故障（应触发熔断）
// 只有这类故障才说明下游服务不可用，需要熔断保护
func isInfrastructureFailure(err error, reqResult *RequestResult) bool {
	if err == nil {
		// 请求成功，但需要检查状态码
		// HTTP 5xx 表示下游服务端错误，属于基础设施故障
		return reqResult.StatusCode >= 500
	}

	// 请求返回了错误
	// 排除业务级错误：如果状态码为 4xx，说明是客户端问题，不是下游故障
	if reqResult.StatusCode >= 400 && reqResult.StatusCode < 500 {
		return false
	}

	// 其他错误（timeout, connection refused, DNS failure 等）都是基础设施故障
	return true
}

// BatchRequest 批量并发请求
func (m *Manager) BatchRequest(ctx *wbcontext.Context, requests []*BatchRequest) ([]map[string]interface{}, error) {
	if len(requests) == 0 {
		return []map[string]interface{}{}, nil
	}

	results := make([]map[string]interface{}, len(requests))

	// 🔥 关键优化：使用带缓冲的channel来控制并发数，避免同时创建过多连接
	// 从配置读取最大并发数，默认10
	maxConcurrency := m.getBatchMaxConcurrency()
	if len(requests) < maxConcurrency {
		maxConcurrency = len(requests)
	}
	semaphore := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup
	var mu sync.Mutex

	// 🔥 关键优化：为批量请求设置整体超时context
	// 计算最大超时时间（取所有请求中最大的超时时间，如果没有则默认500ms）
	maxTimeout := 500 * time.Millisecond
	for _, req := range requests {
		if req.Other != nil && req.Other.Timeout > maxTimeout {
			maxTimeout = req.Other.Timeout
		}
	}

	// 创建带超时的context - 移除+5s缓冲，使用精确超时
	batchCtx := context.Background()
	if ctx != nil && ctx.GinContext() != nil && ctx.GinContext().Request != nil {
		batchCtx = ctx.GinContext().Request.Context()
	}
	batchCtx, cancel := context.WithTimeout(batchCtx, maxTimeout)
	defer cancel()

	for i, req := range requests {
		wg.Add(1)
		go func(index int, request *BatchRequest) {
			defer wg.Done()

			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-batchCtx.Done():
				// 整体超时，直接返回错误
				mu.Lock()
				results[index] = map[string]any{}
				mu.Unlock()
				return
			}

			// 执行单个请求
			result, err := m.Request(ctx, request.URLKey, request.Params, request.Other)
			mu.Lock()
			if err != nil {
				// 当请求失败时，返回空map
				results[index] = map[string]any{}
			} else {
				results[index] = result
			}
			mu.Unlock()
		}(i, req)
	}

	wg.Wait()
	return results, nil
}

// setRequest 设置请求参数
func (m *Manager) setRequest(ctx *wbcontext.Context, urlKey string, params map[string]interface{}, other *OtherOption) (*RequestOption, error) {
	m.mu.RLock()
	apiConfig, exists := m.apisConfig[urlKey]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("API配置不存在: %s", urlKey)
	}

	// 直接使用用户传入的参数
	args := make(map[string]interface{})
	if params != nil {
		for k, v := range params {
			args[k] = v
		}
	}

	// 解析URL和协议
	opt := m.parseURL(apiConfig, other)
	opt.URLKey = urlKey
	opt.Args = args
	opt.Method = apiConfig.Method
	if opt.Method == "" {
		opt.Method = m.defaultMethod
	}
	// 设置超时 - 优先使用用户传入的值，否则使用配置文件的值
	if other.Timeout > 0 {
		opt.Timeout = other.Timeout
	} else if apiConfig.Timeout > 0 {
		opt.Timeout = time.Duration(apiConfig.Timeout) * time.Millisecond
	} else {
		opt.Timeout = 100 * time.Millisecond // 默认100ms
	}

	// 设置重试次数 - 优先使用用户传入的值，否则使用配置文件的值
	if other.TryTimes > 0 {
		opt.TryTimes = other.TryTimes
	} else {
		opt.TryTimes = apiConfig.TryTimes
	}
	opt.Keepalive = apiConfig.Keepalive
	// ContentJSON 优先使用用户传入的值，否则使用配置文件的值
	if other.ContentJSON {
		opt.ContentJSON = true
	} else {
		opt.ContentJSON = apiConfig.ContentJSON
	}
	if other.HTTPMethod != "" {
		opt.HTTPMethod = other.HTTPMethod
	} else {
		opt.HTTPMethod = apiConfig.HTTPMethod
	}
	opt.Version = apiConfig.Version
	if opt.Version == 0 {
		opt.Version = 1
	}

	// 只使用用户传入的请求头
	opt.Headers = make(map[string]string)
	if other.Header != nil {
		for k, v := range other.Header {
			opt.Headers[k] = v
		}
	}

	// 处理认证
	if apiConfig.Auth {
		args["source"] = m.source
		token, err := m.getToken(ctx)
		if err == nil && token != "" {
			opt.Headers["Authorization"] = token
			if opt.Protocol == ProtocolDA {
				args["token"] = token
			}
		}
	}

	// 根据UID选择节点（灰度）
	if uid := m.getUIDInt64(ctx, args); uid > 0 && len(opt.Nodes) > 0 {
		opt.OriginNodes = opt.Nodes
		opt.Nodes = m.parseNodes(opt.Nodes, uid)
	}

	return opt, nil
}

// parseURL 解析URL和协议
func (m *Manager) parseURL(apiConfig *APIConfig, other *OtherOption) *RequestOption {
	opt := &RequestOption{
		Service: apiConfig.Service,
		Group:   apiConfig.Group,
	}

	// 1. 优先使用配置文件中的 type 字段确定协议
	if apiConfig.Type != "" {
		protocol := strings.ToLower(apiConfig.Type)
		opt.Protocol = Protocol(protocol)

		// 对于 motan2 协议，统一转换为 motan
		if protocol == "motan2" {
			opt.Protocol = ProtocolMotan
		}

		// 对于 HTTP/HTTPS 协议，处理 URL
		if protocol == "http" || protocol == "https" {
			// 使用配置的URL
			apiURL := apiConfig.URL
			// 替换环境组
			apiURL = strings.ReplaceAll(apiURL, "%s", m.envGroup)
			opt.URL = apiURL

			// 如果配置了服务发现，获取节点列表
			if opt.Service != "" && opt.Group != "" {
				opt.Group = strings.ReplaceAll(opt.Group, "%s", m.envGroup)
				nodes, _ := discovery.NamingGet(opt.Service, opt.Group)
				opt.Nodes = convertNodes(nodes)
			}
		} else if protocol == "motan" || protocol == "motan2" {
			// Motan 协议不需要 URL，只需要 service 和 method
			// Service 已经在上面从 apiConfig 中获取
		} else if protocol == "da" {
			// DA 协议处理：从URL解析service和group
			// URL格式：da://service/group/version
			// 例如：da://ks-search/search.ks/0
			if opt.Service == "" || opt.Group == "" {
				// 从URL解析service和group
				urlParts := strings.SplitN(apiConfig.URL, "://", 2)
				if len(urlParts) == 2 {
					pathParts := strings.Split(urlParts[1], "/")
					if len(pathParts) >= 2 {
						opt.Service = pathParts[0]
						opt.Group = pathParts[1]
					}
				}
			}
			// 通过服务发现获取节点列表
			if opt.Service != "" && opt.Group != "" {
				nodes, _ := discovery.NamingGet(opt.Service, opt.Group)
				opt.Nodes = convertNodes(nodes)
			}
		}

		return opt
	}

	// 2. 如果没有 type 字段，则从 URL 解析协议
	apiURL := apiConfig.URL
	// 替换环境组
	apiURL = strings.ReplaceAll(apiURL, "%s", m.envGroup)
	opt.URL = apiURL

	// 从 URL 解析协议
	if strings.Contains(apiURL, "://") {
		parts := strings.SplitN(apiURL, "://", 2)
		protocol := strings.ToLower(parts[0])
		opt.Protocol = Protocol(protocol)

		// HTTP/HTTPS协议
		if protocol == "http" || protocol == "https" {
			if opt.Service != "" && opt.Group != "" {
				opt.Group = strings.ReplaceAll(opt.Group, "%s", m.envGroup)
				nodes, _ := discovery.NamingGet(opt.Service, opt.Group)
				opt.Nodes = convertNodes(nodes)
			}
			return opt
		}

		// 其他协议解析
		remaining := parts[1]
		pathParts := strings.Split(remaining, "/")

		switch protocol {
		case "motan2":
			opt.Protocol = ProtocolMotan
			// 只有在配置文件中没有明确指定 service 和 group 时，才从 URL 解析
			if opt.Service == "" && len(pathParts) >= 3 {
				opt.Service = pathParts[1]
			}
			if opt.Group == "" && len(pathParts) >= 3 {
				opt.Group = pathParts[2]
			}
		case "da":
			opt.Protocol = ProtocolDA
			// 只有在配置文件中没有明确指定时，才从 URL 解析
			if opt.Service == "" && len(pathParts) >= 2 {
				opt.Service = pathParts[0]
			}
			if opt.Group == "" && len(pathParts) >= 2 {
				opt.Group = pathParts[1]
			}
			if opt.Service != "" && opt.Group != "" {
				nodes, _ := discovery.NamingGet(opt.Service, opt.Group)
				opt.Nodes = convertNodes(nodes)
			}
		}
	} else {
		// 没有协议前缀，默认为 HTTP
		opt.Protocol = ProtocolHTTP
	}

	return opt
}

// parseNodes 根据UID解析节点（灰度逻辑）
func (m *Manager) parseNodes(nodes []Node, uid int64) []Node {
	var normalNodes, grayNodes []Node

	for _, node := range nodes {
		rule := m.extractGrayRule(node.ExtInfo)
		if rule == nil {
			normalNodes = append(normalNodes, node)
			continue
		}

		// 检查是否在指定UID列表中
		if len(rule.UIDs) > 0 {
			for _, ruleUID := range rule.UIDs {
				if ruleUID == uid {
					return []Node{node} // 命中指定UID，直接返回
				}
			}
		}

		// 按位置和数值匹配
		if rule.Pos > 0 && rule.Num > 0 && len(rule.Value) > 0 {
			offset := math.Pow(10, float64(rule.Pos-rule.Num))
			posNum := int(math.Floor(float64(uid)/offset)) % int(math.Pow(10, float64(rule.Num)))
			for _, v := range rule.Value {
				if v == posNum {
					grayNodes = append(grayNodes, node)
					break
				}
			}
		}
	}

	// 优先返回灰度节点
	if len(grayNodes) > 0 {
		return grayNodes
	}
	if len(normalNodes) > 0 {
		return normalNodes
	}
	return nodes
}

// extractGrayRule 提取灰度规则
func (m *Manager) extractGrayRule(extInfo map[string]interface{}) *GrayRule {
	if extInfo == nil {
		return nil
	}

	ruleData, ok := extInfo["rule"]
	if !ok {
		return nil
	}

	ruleMap, ok := ruleData.(map[string]interface{})
	if !ok {
		return nil
	}

	rule := &GrayRule{}
	if uids, ok := ruleMap["uids"].([]interface{}); ok {
		for _, uid := range uids {
			if uidInt, ok := uid.(int64); ok {
				rule.UIDs = append(rule.UIDs, uidInt)
			} else if uidFloat, ok := uid.(float64); ok {
				rule.UIDs = append(rule.UIDs, int64(uidFloat))
			}
		}
	}

	if pos, ok := ruleMap["pos"].(float64); ok {
		rule.Pos = int(pos)
	}
	if num, ok := ruleMap["num"].(float64); ok {
		rule.Num = int(num)
	}
	if values, ok := ruleMap["value"].([]interface{}); ok {
		for _, v := range values {
			if vInt, ok := v.(float64); ok {
				rule.Value = append(rule.Value, int(vInt))
			}
		}
	}

	return rule
}

// doRequestWithResult 执行请求并返回结果信息（用于日志记录）
func (m *Manager) doRequestWithResult(ctx *wbcontext.Context, opt *RequestOption) (map[string]interface{}, *RequestResult, error) {
	switch opt.Protocol {
	case ProtocolHTTP, ProtocolHTTPS:
		return m.doHTTPRequestWithResult(ctx, opt)
	case ProtocolMotan, ProtocolMotan2:
		return m.doMotanRequestWithResult(ctx, opt)
	case ProtocolDA:
		return m.doDARequestWithResult(ctx, opt)
	default:
		return m.doHTTPRequestWithResult(ctx, opt)
	}
}

// logRequest 记录请求日志
// 格式: network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
// 如果有详细时间记录，则追加: |msg:{err}(conn:{time},send:{time},wait:{time},read:{time})
func (m *Manager) logRequest(ctx *wbcontext.Context, opt *RequestOption, result *RequestResult) {
	// 获取 web_degrade 参数
	webDegrade := opt.WebDegrade
	if webDegrade == 0 {
		if wd, ok := opt.Args["web_degrade"]; ok {
			if wdInt, ok := wd.(int); ok {
				webDegrade = wdInt
			} else if wdStr, ok := wd.(string); ok {
				webDegrade, _ = strconv.Atoi(wdStr)
			}
		}
	}

	// 节点日志，如果没有则使用 "-"
	nodeLog := result.NodeLog
	if nodeLog == "" {
		nodeLog = opt.NodeLog
	}
	if nodeLog == "" {
		nodeLog = "-"
	}

	// 状态码，成功时为 200，失败时为 502
	statusCode := result.StatusCode
	if statusCode == 0 {
		if result.Err == "" {
			statusCode = 200
		} else {
			statusCode = 502
		}
	}

	// 错误信息，成功时为 "success"
	errMsg := result.Err
	if errMsg == "" {
		errMsg = "success"
	} else if statusCode != 200 {
		// 只有请求失败时才简化错误信息展示
		errMsg = simplifyErrorMessage(errMsg)
	}

	var logMsg string
	var timingInfo string
	if result.Timing != nil {
		if opt.Protocol == ProtocolHTTP || opt.Protocol == ProtocolHTTPS || opt.Protocol == ProtocolDA {
			// HTTP/DA协议包含完整的时间信息
			timingInfo = fmt.Sprintf("(conn:%.3f,send:%.3f,wait:%.3f,read:%.3f)",
				result.Timing.ConnectTime.Seconds(),
				result.Timing.SendTime.Seconds(),
				result.Timing.WaitTime.Seconds(),
				result.Timing.ReadTime.Seconds(),
			)
		}
	} else {
		// 没有详细时间记录时，不添加时间信息
		timingInfo = ""
	}

	// 格式化日志消息
	// network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
	logMsg = fmt.Sprintf(" network:%s|%s|%.3f|%s|%d|%d|%d|%d|msg:%s%s",
		opt.Protocol,
		opt.URLKey,
		result.Cost.Seconds(),
		nodeLog,
		result.ResponseLen,
		result.TryTimes,
		statusCode,
		webDegrade,
		errMsg,
		timingInfo,
	)

	// 追加到请求日志
	log.Append(ctx, logMsg)
}

// doHTTPRequestWithResult 执行HTTP/HTTPS请求并返回结果信息（支持重试）
func (m *Manager) doHTTPRequestWithResult(ctx *wbcontext.Context, opt *RequestOption) (map[string]interface{}, *RequestResult, error) {
	reqResult := &RequestResult{
		TryTimes: 1,
		Timing:   &RequestTiming{}, // 初始化Timing结构
	}

	// 构建HTTP请求选项
	httpOpt := &client.HTTPRequestOption{
		URL:         opt.URL,
		Method:      opt.Method,
		Headers:     opt.Headers,
		Args:        opt.Args,
		Timeout:     opt.Timeout,
		ContentJSON: opt.ContentJSON,
	}

	// 获取标准context
	stdCtx := context.Background()
	if ctx != nil && ctx.GinContext() != nil && ctx.GinContext().Request != nil {
		stdCtx = ctx.GinContext().Request.Context()
	}

	// 执行请求
	result, httpResult, err := m.httpClient.DoRequest(stdCtx, httpOpt, opt.TryTimes)

	// 转换结果
	if httpResult != nil {
		reqResult.ResponseLen = httpResult.ResponseLen
		reqResult.StatusCode = httpResult.StatusCode
		reqResult.NodeLog = httpResult.NodeLog
		reqResult.Err = httpResult.Err
		reqResult.TryTimes = httpResult.TryTimes

		// 如果HTTP客户端提供了详细时间记录，则使用它
		if httpResult.Timing != nil && reqResult.Timing != nil {
			reqResult.Timing.TotalTime = httpResult.Timing.TotalTime
			reqResult.Timing.ConnectTime = httpResult.Timing.ConnectTime
			reqResult.Timing.SendTime = httpResult.Timing.SendTime
			reqResult.Timing.WaitTime = httpResult.Timing.WaitTime
			reqResult.Timing.ReadTime = httpResult.Timing.ReadTime
		}
	}

	if err != nil {
		return nil, reqResult, err
	}

	// 结果干预
	return m.interveneResult(result, nil, reqResult), reqResult, nil
}

// doMotanRequestWithResult 执行Motan请求并返回结果信息（通过MeshClient，支持重试）
func (m *Manager) doMotanRequestWithResult(ctx *wbcontext.Context, opt *RequestOption) (map[string]interface{}, *RequestResult, error) {
	reqResult := &RequestResult{
		TryTimes: 1,
		NodeLog:  opt.Service,      // Motan 使用服务名作为节点日志
		Timing:   &RequestTiming{}, // 初始化Timing结构
	}

	// 延迟初始化Motan客户端：只有在第一次使用motan协议时才初始化
	if !client.IsMotanInitialized() {
		if m.motanConfig == nil {
			reqResult.Err = "Motan MeshClient未配置"
			return nil, reqResult, fmt.Errorf("Motan MeshClient未配置，请在app.yaml中配置motan相关参数")
		}
		// 使用sync.Once确保只初始化一次
		motanInitOnce.Do(func() {
			motanInitErr = client.InitMotanMeshClient(m.motanConfig)
		})
		if motanInitErr != nil {
			reqResult.Err = motanInitErr.Error()
			return nil, reqResult, fmt.Errorf("初始化Motan MeshClient失败: %w", motanInitErr)
		}
	}

	// 确保Service已配置
	if opt.Service == "" {
		reqResult.Err = "Motan服务名(Service)未配置"
		return nil, reqResult, fmt.Errorf("Motan服务名(Service)未配置")
	}

	var lastErr error
	maxRetries := opt.TryTimes
	if maxRetries < 0 {
		maxRetries = 0
	}

	// 尝试执行请求（初次请求 + 重试次数）
	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, respLen, statusCode, nodeLog, errMsg, timing, err := m.doMotanRequestOnceWithDetails(ctx, opt)

		reqResult.ResponseLen = respLen
		reqResult.StatusCode = statusCode
		reqResult.NodeLog = nodeLog
		reqResult.TryTimes = attempt + 1

		// 更新时间记录
		if timing != nil && reqResult.Timing != nil {
			reqResult.Timing.TotalTime = timing.TotalTime
			// Motan协议只记录TotalTime，其他字段设为0
			reqResult.Timing.SendTime = 0
			reqResult.Timing.WaitTime = 0
		}

		if err == nil {
			// 请求成功，检查是否有接口级别的错误（如 Auth failed）
			if errMsg != "" {
				// 接口返回了错误信息，记录错误但不重试
				reqResult.Err = errMsg
				return result, reqResult, nil
			}
			reqResult.Err = "success"
			return result, reqResult, nil
		}

		lastErr = err

		// 如果还有重试机会，短暂延迟后重试
		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt+1) * 10 * time.Millisecond)
		}
	}

	// 所有重试都失败
	reqResult.Err = lastErr.Error()
	return nil, reqResult, lastErr
}

// doMotanRequestOnceWithDetails 执行单次Motan请求并返回详细信息
// 返回值: result, responseLen, statusCode, nodeLog, errMsg, timing, error
func (m *Manager) doMotanRequestOnceWithDetails(ctx *wbcontext.Context, opt *RequestOption) (map[string]interface{}, int, int, string, string, *RequestTiming, error) {
	// 转换参数为 map[string]string 格式（MeshClient要求）
	args := make(map[string]string)
	for k, v := range opt.Args {
		args[k] = fmt.Sprintf("%v", v)
	}

	// 调用Motan服务（带时间记录）
	// service: 要调用的RPC服务名称
	// method: 要调用rpc的方法（URL路径）
	// args: 参数（map[string]string类型）
	// headers: 请求头（如Authorization）
	// timeout: 超时时间
	var motanArgs interface{}
	if opt.ContentJSON {
		rawAgrs, _ := sonic.Marshal(args)
		motanArgs = string(rawAgrs)
		if opt.Headers == nil {
			opt.Headers = make(map[string]string)
		}
		if _, ok := opt.Headers["Content-Type"]; !ok {
			opt.Headers["Content-Type"] = "application/json"
		}
	} else {
		motanArgs = args
	}

	if opt.HTTPMethod != "" {
		if opt.Headers == nil {
			opt.Headers = make(map[string]string)
		}
		if _, ok := opt.Headers["HTTP_Method"]; !ok {
			opt.Headers["HTTP_Method"] = opt.HTTPMethod
		}
	}

	resp, _, timing, err := client.MotanCallWithTiming(
		opt.Service,
		opt.Method,
		motanArgs,
		opt.Headers,
		opt.Timeout,
	)
	if err != nil {
		return nil, 0, 502, opt.Service, "", nil, err
	}

	responseLen := len(resp)

	// 解析响应 - 兼容返回数据可能是 map 或 slice 的情况
	result, parseErr := client.ParseJSONResponse(resp)
	if parseErr != nil {
		return nil, responseLen, 502, opt.Service, "", nil, fmt.Errorf("解析Motan响应失败: %w", parseErr)
	}

	// 将Motan的时间记录转换为主API的时间记录格式
	reqTiming := &RequestTiming{}
	if timing != nil {
		reqTiming.TotalTime = timing.TotalTime
		// Motan协议只记录TotalTime，其他字段设为0
		reqTiming.ConnectTime = 0
		reqTiming.SendTime = 0
		reqTiming.WaitTime = 0
		reqTiming.ReadTime = 0
	}

	tempReqResult := &RequestResult{StatusCode: 200, NodeLog: opt.Service}
	intervenedResult := m.interveneResult(result, opt, tempReqResult)
	return intervenedResult, responseLen, tempReqResult.StatusCode, opt.Service, tempReqResult.Err, reqTiming, nil
}

// doDARequestWithResult 执行DA请求并返回结果信息（支持重试）
// DA协议是基于TCP的二进制协议
func (m *Manager) doDARequestWithResult(ctx *wbcontext.Context, opt *RequestOption) (map[string]interface{}, *RequestResult, error) {
	reqResult := &RequestResult{
		TryTimes: 1,
		Timing:   &RequestTiming{}, // 初始化Timing结构
	}

	// 确保有可用节点
	if len(opt.Nodes) == 0 {
		reqResult.Err = "没有可用节点"
		return nil, reqResult, fmt.Errorf("DA协议请求失败: 没有可用节点 (service=%s, group=%s)", opt.Service, opt.Group)
	}

	var lastErr error
	maxRetries := opt.TryTimes
	if maxRetries < 0 {
		maxRetries = 0
	}

	// 尝试执行请求（初次请求 + 重试次数）
	for attempt := 0; attempt <= maxRetries; attempt++ {
		startTime := time.Now()
		result, respLen, statusCode, nodeLog, errMsg, err := m.doDARequestOnceWithDetails(ctx, opt)
		duration := time.Since(startTime)

		reqResult.ResponseLen = respLen
		reqResult.StatusCode = statusCode
		reqResult.NodeLog = nodeLog
		reqResult.TryTimes = attempt + 1

		// 记录DA请求的时间信息
		if reqResult.Timing != nil {
			reqResult.Timing.TotalTime = duration
			// DA协议的时间分配大致如下：
			// - ConnectTime: 连接建立时间
			// - SendTime: 请求发送时间
			// - WaitTime: 等待响应时间
			// - ReadTime: 响应读取时间
		}

		if err == nil {
			// 请求成功，检查是否有接口级别的错误（如 Auth failed）
			if errMsg != "" {
				// 接口返回了错误信息，记录错误但不重试
				reqResult.Err = errMsg
				return result, reqResult, nil
			}
			reqResult.Err = "success"
			return result, reqResult, nil
		}

		lastErr = err

		// 如果还有重试机会，短暂延迟后重试
		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt+1) * 10 * time.Millisecond)
		}
	}

	// 所有重试都失败
	reqResult.Err = lastErr.Error()
	return nil, reqResult, fmt.Errorf("DA请求失败（已重试%d次）: %w", maxRetries, lastErr)
}

// doDARequestOnceWithDetails 执行单次DA请求并返回详细信息
// 返回值: result, responseLen, statusCode, nodeLog, errMsg, error
func (m *Manager) doDARequestOnceWithDetails(ctx *wbcontext.Context, opt *RequestOption) (map[string]interface{}, int, int, string, string, error) {
	// 使用轮询负载均衡选择节点
	node := m.selectDANode(opt.Nodes)

	// 构建节点日志
	nodeLog := fmt.Sprintf("%s-%d", node.Host, node.Port)

	// 创建DA客户端，使用统一的超时时间
	daClient := client.NewDAClient(
		node.Host,
		node.Port,
		opt.Timeout,
		opt.Version,
	)

	// 执行DA请求（带时间记录）
	respData, timing, err := daClient.RequestWithTiming(opt.Args)
	if err != nil {
		return nil, 0, 502, nodeLog, "", fmt.Errorf("DA协议请求失败 (service=%s, group=%s, node=%s:%d): %w",
			opt.Service, opt.Group, node.Host, node.Port, err)
	}

	responseLen := len(respData)

	// 解析JSON响应 - 兼容返回数据可能是 map 或 slice 的情况
	result, err := client.ParseJSONResponse(respData)
	if err != nil {
		return nil, responseLen, 502, nodeLog, "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 结果干预（需要创建临时的 RequestResult 用于检测错误）
	tempReqResult := &RequestResult{StatusCode: 200}
	if tempReqResult.Timing == nil {
		tempReqResult.Timing = &RequestTiming{}
	}
	// 将DA的时间记录转换为主API的时间记录格式
	if timing != nil {
		tempReqResult.Timing.TotalTime = timing.TotalTime
		tempReqResult.Timing.ConnectTime = timing.ConnectTime
		tempReqResult.Timing.SendTime = timing.SendTime
		tempReqResult.Timing.WaitTime = timing.WaitTime
		tempReqResult.Timing.ReadTime = timing.ReadTime
	}

	intervenedResult := m.interveneResult(result, nil, tempReqResult)
	return intervenedResult, responseLen, tempReqResult.StatusCode, nodeLog, tempReqResult.Err, nil
}

// selectDANode 使用轮询算法选择DA节点
func (m *Manager) selectDANode(nodes []Node) Node {
	if len(nodes) == 1 {
		return nodes[0]
	}

	// 使用原子操作实现轮询负载均衡
	index := atomic.AddUint32(&daNodeIndex, 1)
	return nodes[index%uint32(len(nodes))]
}

// hasPlatformError 检查结果是否包含平台接口错误
// error_code 只判断存在且不为空，不判断具体值
func hasPlatformError(result map[string]any) bool {
	// 检查是否存在 error 字段
	_, hasError := result["error"]
	if !hasError {
		return false
	}

	// 检查 error_code 是否存在
	errorCode, hasErrorCode := result["error_code"]
	if !hasErrorCode {
		return false
	}

	// 判断 error_code 是否为空值
	switch v := errorCode.(type) {
	case string:
		return v != ""
	case int, int64, float64:
		return v != 0
	default:
		return false
	}
}

// isPlatformAPI 判断是否为平台接口
// 对于 Motan 协议，检查 service 是否在平台接口服务列表中
// 对于 HTTP 协议，检查 URL 域名是否为平台域名
func isPlatformAPI(protocol Protocol, service string, url string) bool {
	// Motan 协议：检查 service 是否在平台接口服务列表中
	if protocol == ProtocolMotan {
		for _, platformService := range PlatformInterfaceServices {
			if service == platformService {
				return true
			}
		}
		return false
	}

	// HTTP/HTTPS 协议：检查域名
	if protocol == ProtocolHTTP || protocol == ProtocolHTTPS {
		if url == "" {
			return false
		}
		// 解析 URL 获取主机名
		u, err := urlpkg.Parse(url)
		if err != nil {
			return false
		}
		host := u.Hostname()
		// 检查是否在平台域名列表中
		for _, domain := range PlatformAPIDomains {
			if host == domain {
				return true
			}
		}
	}

	return false
}

// 对于平台接口（HTTP和Motan），检查返回结果中的 error 和 error_code
// 如果存在且 error_code 不为空，则将状态码更改为 502
func (m *Manager) interveneResult(result map[string]any, opt *RequestOption, reqResult *RequestResult) map[string]any {
	if result == nil {
		return nil
	}

	// 平台接口特殊错误处理（HTTP 和 Motan）
	if opt != nil && isPlatformAPI(opt.Protocol, opt.Service, opt.URL) {
		if hasPlatformError(result) {
			reqResult.StatusCode = 502

			// 提取错误信息
			if errMsg, ok := result["error"].(string); ok {
				reqResult.Err = errMsg
			}

			return result
		}
	}

	return result
}

// getUIDInt64 获取UID（int64）
func (m *Manager) getUIDInt64(ctx *wbcontext.Context, args map[string]any) int64 {
	uidStr := ctx.GetUserId()

	if uidStr == "" {
		return 0
	}
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	return uid
}

// getToken 获取认证token
func (m *Manager) getToken(ctx *wbcontext.Context) (string, error) {
	if m.tokenFile == "" {
		return "", fmt.Errorf("tokenFile未配置")
	}

	// 读取token文件
	token, err := GetTAuthToken(m.tokenFile, ctx.GetUserId())
	if err != nil {
		return "", err
	}

	return token, nil
}

// convertNodes 转换discovery节点为内部节点
func convertNodes(discoveryNodes []discovery.NamingNode) []Node {
	nodes := make([]Node, len(discoveryNodes))
	for i, dn := range discoveryNodes {
		nodes[i] = Node{
			Host:    dn.Host,
			Port:    dn.Port,
			ExtInfo: dn.Ext, // discovery.NamingNode 使用 Ext 字段
		}
	}
	return nodes
}

// GenerateSign 生成签名（工具方法）
func GenerateSign(secret, param string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(param))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// LoadConfigFromViper 使用 viper 实例加载配置
func LoadConfigFromViper(v *viper.Viper) (*Config, error) {
	// 从 viper 加载应用配置
	var appConfig AppConfig
	if err := config.LoadSectionFromViper(v, "api", &appConfig.API); err != nil {
		// api 配置不存在时使用默认值
		appConfig.API = AppAPIConfig{}
	}
	if err := config.LoadSectionFromViper(v, "motan", &appConfig.Motan); err != nil {
		// motan 配置不存在时使用默认值
		appConfig.Motan = AppMotanConfig{}
	}
	// 加载顶层 env 配置
	appConfig.Env = v.GetString("env")

	// 确定接口配置文件路径
	interfacePath := appConfig.API.InterfacePath
	if interfacePath == "" {
		interfacePath = "config/interface.yaml"
	}
	// 将相对路径转为相对配置文件目录的绝对路径
	interfacePath = config.AbsPathFromConfigFile(v, interfacePath)
	// 读取interface.yaml配置文件
	data, err := os.ReadFile(interfacePath)
	if err != nil {
		return nil, fmt.Errorf("读取接口配置文件失败 (%s): %w", interfacePath, err)
	}

	// 解析配置
	var apisConfig map[string]*APIConfig
	if err := yaml.Unmarshal(data, &apisConfig); err != nil {
		return nil, fmt.Errorf("解析接口配置文件失败: %w", err)
	}
	// 构建管理器配置
	// 构建熔断器配置
	breakerConfig := client.BreakerConfig{
		FailureThreshold:  appConfig.API.Breaker.FailureThreshold,
		OpenTimeout:       time.Duration(appConfig.API.Breaker.OpenTimeout) * time.Second,
		HalfOpenMaxProbes: appConfig.API.Breaker.HalfOpenMaxProbes,
		SuccessThreshold:  appConfig.API.Breaker.SuccessThreshold,
	}

	// 构建管理器配置
	cfg := &Config{
		APIsConfig:          apisConfig,
		Source:              getConfigValue(appConfig.API.Source, "2936099636"),
		TokenFile:           config.AbsPathFromConfigFile(v, appConfig.API.TokenFile),
		DefaultMethod:       getConfigValue(appConfig.API.DefaultMethod, "GET"),
		EnvGroup:            getConfigValue(appConfig.API.EnvGroup, "yf"),
		Env:                 getConfigValue(appConfig.Env, "prod"),
		DegradeKeyPrefix:    getConfigValue(appConfig.API.DegradeKeyPrefix, "api_degrade_"),
		BatchMaxConcurrency: appConfig.API.BatchMaxConcurrency,
		HTTPTransport: client.HTTPTransportConfig{
			MaxIdleConns:          appConfig.API.HTTPTransport.MaxIdleConns,
			MaxIdleConnsPerHost:   appConfig.API.HTTPTransport.MaxIdleConnsPerHost,
			ResponseHeaderTimeout: appConfig.API.HTTPTransport.ResponseHeaderTimeout,
		},
		BreakerConfig: breakerConfig,
	}
	if appConfig.Motan.MeshAddress != "" || appConfig.Motan.Application != "" {
		cfg.MotanConfig = &client.MotanConfig{
			MeshAddress:       appConfig.Motan.MeshAddress,
			Application:       appConfig.Motan.Application,
			Serialization:     appConfig.Motan.Serialization,
			RequestTimeout:    time.Duration(appConfig.Motan.RequestTimeout) * time.Millisecond,
			MinRequestTimeout: time.Duration(appConfig.Motan.MinRequestTimeout) * time.Millisecond,
			MaxRequestTimeout: time.Duration(appConfig.Motan.MaxRequestTimeout) * time.Millisecond,
		}
	}

	// 设置 Token 缓存 TTL（从配置读取，默认 300 秒 = 5 分钟）
	tokenCacheTTL := appConfig.API.TokenCacheTTL
	if tokenCacheTTL <= 0 {
		tokenCacheTTL = 300
	}
	SetTokenCacheTTL(time.Duration(tokenCacheTTL) * time.Second)

	return cfg, nil
}

// LoadConfig 加载配置（使用 ComponentViper 懒加载方式）
func LoadConfig() (*Config, error) {
	v, err := config.ComponentViper()
	if err != nil {
		return nil, fmt.Errorf("加载配置文件失败: %w", err)
	}
	return LoadConfigFromViper(v)
}

// AppConfig 应用配置结构
type AppConfig struct {
	Env   string         `yaml:"env"`
	API   AppAPIConfig   `yaml:"api"`
	Motan AppMotanConfig `yaml:"motan"`
}

// AppAPIConfig API相关配置
type AppAPIConfig struct {
	Source              string                     `yaml:"source" mapstructure:"source"`
	TokenFile           string                     `yaml:"token_file" mapstructure:"token_file"`
	TokenCacheTTL       int                        `yaml:"token_cache_ttl" mapstructure:"token_cache_ttl"` // Token 缓存 TTL（秒），默认 300（5分钟）
	DefaultMethod       string                     `yaml:"default_method" mapstructure:"default_method"`
	EnvGroup            string                     `yaml:"env_group" mapstructure:"env_group"`
	DegradeKeyPrefix    string                     `yaml:"degrade_key_prefix" mapstructure:"degrade_key_prefix"`
	InterfacePath       string                     `yaml:"interface_path" mapstructure:"interface_path"`               // 接口配置文件路径
	BatchMaxConcurrency int                        `yaml:"batch_max_concurrency" mapstructure:"batch_max_concurrency"` // 批量请求最大并发数，默认10
	HTTPTransport       client.HTTPTransportConfig `yaml:"http_transport" mapstructure:"http_transport"`               // HTTP Transport 连接池配置
	Breaker             AppBreakerConfig           `yaml:"breaker" mapstructure:"breaker"`                             // 熔断器配置
}

// AppBreakerConfig 熔断器配置（YAML 映射）
type AppBreakerConfig struct {
	FailureThreshold  int `yaml:"failure_threshold" mapstructure:"failure_threshold"`       // 连续失败次数阈值，默认 5
	OpenTimeout       int `yaml:"open_timeout" mapstructure:"open_timeout"`                 // 熔断打开持续时间（秒），默认 10
	HalfOpenMaxProbes int `yaml:"half_open_max_probes" mapstructure:"half_open_max_probes"` // 半开状态最大探测数，默认 1
	SuccessThreshold  int `yaml:"success_threshold" mapstructure:"success_threshold"`       // 半开状态连续成功次数，默认 2
}

// AppMotanConfig Motan相关配置
type AppMotanConfig struct {
	MeshAddress       string `yaml:"mesh_address" mapstructure:"mesh_address"`               // Mesh代理地址，默认 "127.0.0.1:9981"
	Application       string `yaml:"application" mapstructure:"application"`                 // 调用服务的业务方
	Serialization     string `yaml:"serialization" mapstructure:"serialization"`             // 序列化方式，默认 "simple"
	RequestTimeout    int    `yaml:"request_timeout" mapstructure:"request_timeout"`         // 请求超时时间（毫秒），默认 1000
	MinRequestTimeout int    `yaml:"min_request_timeout" mapstructure:"min_request_timeout"` // 最小请求超时时间（毫秒），默认为 RequestTimeout/2
	MaxRequestTimeout int    `yaml:"max_request_timeout" mapstructure:"max_request_timeout"` // 最大请求超时时间（毫秒），默认为 RequestTimeout*2
}

// getConfigValue 获取配置值，如果为空则返回默认值
func getConfigValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// getBatchMaxConcurrency 获取批量请求最大并发数
func (m *Manager) getBatchMaxConcurrency() int {
	if m.batchMaxConcurrency <= 0 {
		return 10
	}
	return m.batchMaxConcurrency
}

// GetHTTPClient 获取底层HTTP客户端（用于外部直接使用）
func (m *Manager) GetHTTPClient() *http.Client {
	return m.httpClient.GetClient()
}

// GetBreakerStats 获取所有熔断器状态（用于监控接口）
func (m *Manager) GetBreakerStats() map[string]client.BreakerStats {
	if m.breakerGroup == nil {
		return map[string]client.BreakerStats{}
	}
	return m.breakerGroup.AllStats()
}

// ResetBreaker 重置指定服务的熔断器（用于运维手动恢复）
func (m *Manager) ResetBreaker(urlKey string) {
	if m.breakerGroup == nil {
		return
	}
	cb := m.breakerGroup.Get(urlKey)
	cb.Reset()
}

// GetTokenCacheStats 获取 Token 缓存监控指标
func (m *Manager) GetTokenCacheStats() TokenCacheStats {
	return GetTokenCacheStats()
}

// simplifyErrorMessage 简化错误信息展示
// 根据错误信息内容返回简化的错误类型：
// - 连接超时类错误返回命中的关键词（空格替换为下划线）
// - 读取超时类错误返回命中的关键词（空格替换为下划线）
// - 其他类型返回 errLower（URL编码）
func simplifyErrorMessage(errMsg string) string {
	errLower := strings.ToLower(errMsg)

	// 连接超时类错误关键词
	connectTimeoutKeywords := []string{
		"connection timeout",
		"connect timeout",
		"dial timeout",
		"connection refused",
		"no route to host",
		"network is unreachable",
		"i/o timeout",
		"connect: connection timed out",
		"dial tcp",
		"connection reset",
		"connection request timeout",
	}

	// 读取超时类错误关键词
	readTimeoutKeywords := []string{
		"read timeout",
		"read: connection reset",
		"response timeout",
		"context deadline exceeded",
		"read: connection timed out",
		"eof",
		"unexpected eof",
		"receive request timeout",
	}

	// 检查是否为连接超时类错误
	for _, keyword := range connectTimeoutKeywords {
		if strings.Contains(errLower, keyword) {
			// 将空格替换为下划线
			return strings.ReplaceAll(keyword, " ", "_")
		}
	}

	// 检查是否为读取超时类错误
	for _, keyword := range readTimeoutKeywords {
		if strings.Contains(errLower, keyword) {
			// 将空格替换为下划线
			return strings.ReplaceAll(keyword, " ", "_")
		}
	}

	// 其他类型错误，返回 errLower 并进行 URL 编码
	return url.QueryEscape(errLower)
}
