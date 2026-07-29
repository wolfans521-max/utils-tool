package api_test

import (
	"git.intra.weibo.com/search_fe/wbutil-go/facades/api"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/log"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// newTestContext 创建测试用的 Context，包含真实的 gin.Context
func newTestContext() *context.Context {
	// 创建一个 gin 测试引擎
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建一个模拟的 HTTP Request
	req, _ := http.NewRequest("GET", "/test", nil)
	c.Request = req

	// 使用真实的 gin.Context 创建 Context
	return context.New(c)
}

// TestRequest 测试单个API请求
func TestRequest(t *testing.T) {
	// 创建测试上下文
	ctx := newTestContext()

	// 准备请求参数
	params := map[string]any{
		"uid":        "2771122141",
		"access_key": "searchdiscover",
		"deep":       3,
		"type":       1,
	}

	other := &api.OtherOption{
		TryTimes: 1,
	}

	// 执行请求
	result := api.Request(ctx, "lbs", params, other)
	_ = log.Flush(ctx)

	// 验证结果
	if result == nil {
		t.Error("返回结果为空")
		return
	}

	// 验证结果 - 由于test_api配置不存在，返回空map是预期行为
	t.Logf("请求完成，返回结果: %+v", result)
}

// TestRequestWithDefaultOptions 测试使用默认选项的请求
func TestRequestWithDefaultOptions(t *testing.T) {
	ctx := newTestContext()

	params := map[string]any{
		"uid": "123456",
	}

	abtestOther := &api.OtherOption{
		TryTimes: 0,
	}

	abtestResult := api.Request(ctx, "abtest", params, abtestOther)
	_ = log.Flush(ctx)

	// 验证 abtest 结果
	if abtestResult == nil {
		t.Error("abtest 返回结果为空")
		return
	}

	t.Logf("abtest 请求成功，返回结果: %+v", abtestResult)
}

// TestRequestWithCustomHeaders 测试带自定义请求头的请求
func TestRequestWithCustomHeaders(t *testing.T) {
	ctx := newTestContext()

	params := map[string]any{
		"uid": "123456",
	}

	other := &api.OtherOption{
		Header: map[string]string{
			"X-Custom-Header": "custom-value",
			"User-Agent":      "test-client/1.0",
		},
		Timeout: 100 * time.Millisecond,
	}

	result := api.Request(ctx, "test_api", params, other)

	t.Logf("请求成功，返回结果: %+v", result)
}

// TestRequestWithJSONContent 测试使用JSON格式发送数据
func TestRequestWithJSONContent(t *testing.T) {
	ctx := newTestContext()

	params := map[string]any{
		"uid": "123456",
		"data": map[string]any{
			"key1": "value1",
			"key2": 123,
		},
	}

	other := &api.OtherOption{
		ContentJSON: true, // 使用JSON格式
		Timeout:     100 * time.Millisecond,
	}

	result := api.Request(ctx, "test_api", params, other)

	t.Logf("请求成功，返回结果: %+v", result)
}

// TestMRequest 测试批量并发请求
func TestMRequest(t *testing.T) {
	ctx := newTestContext()

	// 准备批量请求
	requests := []*api.BatchRequestItem{
		{
			URLKey: "lbs",
			Params: map[string]any{
				"uid":        "2771122141",
				"access_key": "searchdiscover",
				"deep":       3,
				"type":       1,
			},
			Other: &api.OtherOption{
				Timeout: 100 * time.Millisecond,
			},
		},
		{
			URLKey: "abtest",
			Params: map[string]any{
				"businessId": "22",
				"userId":     "2771122141",
				"from":       "10G2293010",
			},
			Other: &api.OtherOption{
				Timeout: 50 * time.Millisecond,
			},
		},
	}

	// 执行批量请求
	results := api.MRequest(ctx, requests)
	_ = log.Flush(ctx)

	// 验证结果
	if len(results) != len(requests) {
		t.Errorf("返回结果数量不匹配，期望: %d, 实际: %d", len(requests), len(results))
		return
	}

	// 检查每个请求的结果
	for i, result := range results {
		t.Logf("请求 %d 结果: %+v", i, result)

		// 检查是否有错误
		if errMsg, ok := result["error"]; ok {
			t.Logf("请求 %d 失败: %v", i, errMsg)
		}
	}
}

// TestMRequestConcurrent 测试批量并发请求
func TestMRequestConcurrent(t *testing.T) {
	// 总共100次请求，每次5个并发
	totalRequests := 200
	concurrency := 15

	t.Logf("开始并发测试：总共 %d 次请求，每次 %d 个并发", totalRequests, concurrency)

	successCount := 0
	failCount := 0
	var mu sync.Mutex

	// 创建信号量控制并发数
	sem := make(chan struct{}, concurrency)

	// 使用 WaitGroup 等待所有请求完成
	var wg sync.WaitGroup

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{} // 获取信号量

		go func(requestIndex int) {
			defer wg.Done()
			defer func() { <-sem }() // 释放信号量

			ctx := newTestContext()

			// 准备批量请求
			requests := []*api.BatchRequestItem{
				{
					URLKey: "lbs",
					Params: map[string]any{
						"uid":        "2771122141",
						"access_key": "searchdiscover",
						"deep":       3,
						"type":       1,
					},
					Other: &api.OtherOption{
						Timeout: 100 * time.Millisecond,
					},
				},
				{
					URLKey: "abtest",
					Params: map[string]any{
						"businessId": "22",
						"userId":     "2771122141",
						"from":       "10G2293010",
					},
					Other: &api.OtherOption{
						Timeout: 200 * time.Millisecond,
					},
				},
			}

			// 执行批量请求
			results := api.MRequest(ctx, requests)
			_ = log.Flush(ctx)

			// 验证结果
			mu.Lock()
			if len(results) != len(requests) {
				t.Errorf("第 %d 次请求返回结果数量不匹配，期望: %d, 实际: %d", requestIndex+1, len(requests), len(results))
				failCount++
			} else {
				// 检查是否有错误
				hasError := false
				for _, result := range results {
					if _, ok := result["error"]; ok {
						// t.Logf("第 %d 次请求失败: %v", requestIndex+1, errMsg)
						hasError = true
						break
					}
				}
				if hasError {
					failCount++
				} else {
					successCount++
				}
			}
			mu.Unlock()
		}(i)
	}

	// 等待所有请求完成
	wg.Wait()

	t.Logf("并发测试完成：成功 %d 次，失败 %d 次，总计 %d 次", successCount, failCount, totalRequests)
}

// TestBatchRequestEmpty 测试空批量请求
func TestBatchRequestEmpty(t *testing.T) {
	ctx := newTestContext()

	// 空请求列表
	results := api.MRequest(ctx, []*api.BatchRequestItem{})
	if len(results) != 0 {
		t.Errorf("空批量请求应返回空结果，实际返回: %d 个结果", len(results))
	}
}

// TestBatchRequestWithRetry 测试带重试的批量请求
func TestBatchRequestWithRetry(t *testing.T) {
	ctx := newTestContext()

	requests := []*api.BatchRequestItem{
		{
			URLKey: "test_api",
			Params: map[string]any{
				"uid": "123456",
			},
			Other: &api.OtherOption{
				TryTimes: 3, // 重试3次
				Timeout:  100 * time.Millisecond,
			},
		},
	}

	results := api.MRequest(ctx, requests)

	t.Logf("批量请求结果: %+v", results)
}

// TestGetManager 测试获取Manager实例
func TestGetManager(t *testing.T) {
	manager, err := api.GetManager()
	if err != nil {
		t.Errorf("获取Manager失败: %v", err)
		return
	}

	if manager == nil {
		t.Error("Manager实例为空")
		return
	}

	t.Log("成功获取Manager实例")
}

// TestGetManagerSingleton 测试Manager单例模式
func TestGetManagerSingleton(t *testing.T) {
	manager1, err1 := api.GetManager()
	if err1 != nil {
		t.Errorf("第一次获取Manager失败: %v", err1)
		return
	}

	manager2, err2 := api.GetManager()
	if err2 != nil {
		t.Errorf("第二次获取Manager失败: %v", err2)
		return
	}

	// 验证是否为同一个实例
	if manager1 != manager2 {
		t.Error("Manager不是单例模式")
		return
	}

	t.Log("Manager单例模式验证成功")
}

// BenchmarkRequest 性能测试：单个请求
func BenchmarkRequest(b *testing.B) {
	ctx := newTestContext()
	params := map[string]any{
		"uid": "123456",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = api.Request(ctx, "test_api", params)
	}
}

// BenchmarkBatchRequest 性能测试：批量请求
func BenchmarkBatchRequest(b *testing.B) {
	ctx := newTestContext()
	requests := []*api.BatchRequestItem{
		{
			URLKey: "test_api_1",
			Params: map[string]any{"uid": "123456"},
		},
		{
			URLKey: "test_api_2",
			Params: map[string]any{"uid": "789012"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = api.MRequest(ctx, requests)
	}
}

// ExampleRequest 示例：基本的API请求
func ExampleRequest() {
	// 创建上下文
	ctx := newTestContext()

	// 准备请求参数
	params := map[string]any{
		"uid":  "123456",
		"page": 1,
	}

	// 配置请求选项
	other := &api.OtherOption{
		Timeout: 100 * time.Millisecond,
	}

	// 执行请求
	result := api.Request(ctx, "api_name", params, other)

	// 使用返回结果
	_ = result
}

// ExampleRequest_withDefaultOptions 示例：使用默认选项的请求
func ExampleRequest_withDefaultOptions() {
	ctx := newTestContext()
	params := map[string]any{
		"uid": "123456",
	}

	// 不传递 other 参数，使用配置文件中的默认值
	result := api.Request(ctx, "api_name", params)

	_ = result
}

// ExampleRequest_withJSONContent 示例：使用JSON格式发送数据
func ExampleRequest_withJSONContent() {
	ctx := newTestContext()
	params := map[string]any{
		"uid": "123456",
		"data": map[string]any{
			"key1": "value1",
			"key2": 123,
		},
	}

	other := &api.OtherOption{
		ContentJSON: true, // 使用JSON格式
		Timeout:     100 * time.Millisecond,
	}

	result := api.Request(ctx, "api_name", params, other)

	_ = result
}

// ExampleRequest_withCustomHeaders 示例：带自定义请求头
func ExampleRequest_withCustomHeaders() {
	ctx := newTestContext()
	params := map[string]any{
		"uid": "123456",
	}

	other := &api.OtherOption{
		Header: map[string]string{
			"X-Custom-Header": "custom-value",
			"Authorization":   "Bearer token",
		},
		Timeout: 100 * time.Millisecond,
	}

	result := api.Request(ctx, "api_name", params, other)

	_ = result
}

// ExampleRequest_withRetry 示例：配置重试次数
func ExampleRequest_withRetry() {
	ctx := newTestContext()
	params := map[string]any{
		"uid": "123456",
	}

	other := &api.OtherOption{
		TryTimes: 1, // 失败后重试3次
		Timeout:  100 * time.Millisecond,
	}

	result := api.Request(ctx, "api_name", params, other)

	_ = result
}

// ExampleMRequest 示例：批量并发请求
func ExampleMRequest() {
	ctx := newTestContext()

	// 准备多个请求
	requests := []*api.BatchRequestItem{
		{
			URLKey: "api_user_info",
			Params: map[string]any{
				"uid": "123456",
			},
			Other: &api.OtherOption{
				Timeout: 100 * time.Millisecond,
			},
		},
		{
			URLKey: "api_user_posts",
			Params: map[string]any{
				"uid":  "123456",
				"page": 1,
				"size": 20,
			},
			Other: &api.OtherOption{
				Timeout: 100 * time.Millisecond,
			},
		},
		{
			URLKey: "api_user_followers",
			Params: map[string]any{
				"uid": "123456",
			},
			Other: nil, // 使用默认配置
		},
	}

	// 执行批量请求（并发执行）
	results := api.MRequest(ctx, requests)

	// 处理每个请求的结果
	for i, result := range results {
		// 检查是否有错误
		if errMsg, ok := result["error"]; ok {
			// 处理错误
			_ = errMsg
			continue
		}

		// 使用正常结果
		_ = result
		_ = i
	}
}

// ExampleMRequest_withDifferentAPIs 示例：批量请求不同的API
func ExampleMRequest_withDifferentAPIs() {
	ctx := newTestContext()

	requests := []*api.BatchRequestItem{
		{
			URLKey: "mblog",
			Params: map[string]any{"uid": 123},
		},
		{
			URLKey: "user_info",
			Params: map[string]any{"uid": 456},
		},
		{
			URLKey: "search",
			Params: map[string]any{
				"q":    "关键词",
				"page": 1,
			},
		},
	}

	results := api.MRequest(ctx, requests)

	_ = results
}

// ExampleGetManager 示例：获取Manager实例
func ExampleGetManager() {
	// 获取全局Manager实例（单例模式）
	manager, err := api.GetManager()
	if err != nil {
		// 处理初始化错误
		return
	}

	// 使用manager实例
	_ = manager
}

// ==================== 高并发测试 ====================

// TestHighConcurrencyRequest 测试高并发单个请求 - 使用真实接口 finder_ad_banner
func TestHighConcurrencyRequest(t *testing.T) {
	const (
		concurrency = 100  // 并发数
		totalReqs   = 1000 // 总请求数
	)

	var (
		wg           sync.WaitGroup
		successCount int64
		failCount    int64
	)

	// 创建信号量控制并发数
	sem := make(chan struct{}, concurrency)

	startTime := time.Now()

	for i := 0; i < totalReqs; i++ {
		wg.Add(1)
		go func(reqID int) {
			defer wg.Done()

			// 获取信号量
			sem <- struct{}{}
			defer func() { <-sem }()

			// 每个goroutine创建独立的上下文
			ctx := newTestContext()

			// 使用真实接口 finder_ad_banner 的参数
			params := map[string]any{
				"page":   1,
				"count":  10,
				"source": "search",
			}

			other := &api.OtherOption{
				Timeout: 500 * time.Millisecond,
			}

			result := api.Request(ctx, "finder_ad_banner", params, other)

			// 统计结果
			if _, ok := result["error"]; ok {
				atomic.AddInt64(&failCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	t.Logf("高并发测试完成 (finder_ad_banner):")
	t.Logf("  总请求数: %d", totalReqs)
	t.Logf("  并发数: %d", concurrency)
	t.Logf("  成功数: %d", successCount)
	t.Logf("  失败数: %d", failCount)
	t.Logf("  总耗时: %v", duration)
	t.Logf("  QPS: %.2f", float64(totalReqs)/duration.Seconds())
}

// TestHighConcurrencyBatchRequest 测试高并发批量请求 - 使用真实接口
func TestHighConcurrencyBatchRequest(t *testing.T) {
	const (
		concurrency  = 50  // 并发数
		totalBatches = 200 // 总批次数
	)

	var (
		wg           sync.WaitGroup
		successCount int64
		failCount    int64
	)

	sem := make(chan struct{}, concurrency)

	startTime := time.Now()

	for i := 0; i < totalBatches; i++ {
		wg.Add(1)
		go func(batchID int) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			ctx := newTestContext()

			// 构建批量请求 - 使用真实接口配置
			requests := []*api.BatchRequestItem{
				{
					// mblog 接口 - 微博showbatch
					URLKey: "mblog",
					Params: map[string]any{
						"ids":    "4500000000000001,4500000000000002,4500000000000003",
						"source": "search",
					},
					Other: &api.OtherOption{
						Timeout: 500 * time.Millisecond,
					},
				},
				{
					// finder_ad_banner 接口
					URLKey: "finder_ad_banner",
					Params: map[string]any{
						"page":   1,
						"count":  10,
						"source": "search",
					},
					Other: &api.OtherOption{
						Timeout: 500 * time.Millisecond,
					},
				},
				{
					// wenda_box 接口 - 问答推荐
					URLKey: "wenda_box",
					Params: map[string]any{
						"q":      "测试关键词",
						"uid":    "1234567890",
						"source": "search",
					},
					Other: &api.OtherOption{
						Timeout:     500 * time.Millisecond,
						ContentJSON: true,
					},
				},
			}

			results := api.MRequest(ctx, requests)

			// 统计结果
			for _, result := range results {
				if _, ok := result["error"]; ok {
					atomic.AddInt64(&failCount, 1)
				} else {
					atomic.AddInt64(&successCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	totalReqs := totalBatches * 3 // 每批3个请求
	t.Logf("高并发批量请求测试完成 (mblog, finder_ad_banner, wenda_box):")
	t.Logf("  总批次数: %d", totalBatches)
	t.Logf("  每批请求数: 3")
	t.Logf("  总请求数: %d", totalReqs)
	t.Logf("  并发数: %d", concurrency)
	t.Logf("  成功数: %d", successCount)
	t.Logf("  失败数: %d", failCount)
	t.Logf("  总耗时: %v", duration)
	t.Logf("  QPS: %.2f", float64(totalReqs)/duration.Seconds())
}

// TestHighConcurrencyGetManager 测试高并发获取Manager实例（验证单例模式线程安全）
func TestHighConcurrencyGetManager(t *testing.T) {
	const concurrency = 100

	var (
		wg       sync.WaitGroup
		managers = make([]any, concurrency)
		errors   = make([]error, concurrency)
	)

	// 同时启动所有goroutine
	startCh := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			<-startCh // 等待信号同时开始

			manager, err := api.GetManager()
			managers[idx] = manager
			errors[idx] = err
		}(i)
	}

	// 发送开始信号
	close(startCh)
	wg.Wait()

	// 验证所有获取的Manager是否为同一实例
	var firstManager any
	var errorCount int

	for i := 0; i < concurrency; i++ {
		if errors[i] != nil {
			errorCount++
			t.Logf("goroutine %d 获取Manager失败: %v", i, errors[i])
			continue
		}

		if firstManager == nil {
			firstManager = managers[i]
		} else if managers[i] != firstManager {
			t.Errorf("goroutine %d 获取的Manager实例不一致", i)
		}
	}

	if errorCount == 0 {
		t.Logf("高并发获取Manager测试通过: %d 个goroutine都获取到相同的实例", concurrency)
	} else {
		t.Logf("高并发获取Manager测试完成: %d 个成功, %d 个失败", concurrency-errorCount, errorCount)
	}
}

// TestHighConcurrencyMixedRequests 测试高并发混合请求（单个请求和批量请求混合）- 使用真实接口
func TestHighConcurrencyMixedRequests(t *testing.T) {
	const (
		concurrency = 100
		totalReqs   = 500
	)

	var (
		wg                 sync.WaitGroup
		singleSuccessCount int64
		singleFailCount    int64
		batchSuccessCount  int64
		batchFailCount     int64
	)

	sem := make(chan struct{}, concurrency)

	startTime := time.Now()

	for i := 0; i < totalReqs; i++ {
		wg.Add(1)
		go func(reqID int) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			ctx := newTestContext()

			// 奇数执行单个请求（finder_ad_banner），偶数执行批量请求
			if reqID%2 == 1 {
				// 使用真实接口 finder_ad_banner
				params := map[string]any{
					"page":   1,
					"count":  10,
					"source": "search",
				}

				other := &api.OtherOption{
					Timeout: 500 * time.Millisecond,
				}

				result := api.Request(ctx, "finder_ad_banner", params, other)

				if _, ok := result["error"]; ok {
					atomic.AddInt64(&singleFailCount, 1)
				} else {
					atomic.AddInt64(&singleSuccessCount, 1)
				}
			} else {
				// 使用真实接口批量请求
				requests := []*api.BatchRequestItem{
					{
						// mblog 接口
						URLKey: "mblog",
						Params: map[string]any{
							"ids":    "4500000000000001,4500000000000002",
							"source": "search",
						},
						Other: &api.OtherOption{Timeout: 500 * time.Millisecond},
					},
					{
						// wenda_box 接口
						URLKey: "wenda_box",
						Params: map[string]any{
							"q":      "测试关键词",
							"uid":    "1234567890",
							"source": "search",
						},
						Other: &api.OtherOption{
							Timeout:     500 * time.Millisecond,
							ContentJSON: true,
						},
					},
				}

				results := api.MRequest(ctx, requests)

				for _, result := range results {
					if _, ok := result["error"]; ok {
						atomic.AddInt64(&batchFailCount, 1)
					} else {
						atomic.AddInt64(&batchSuccessCount, 1)
					}
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	t.Logf("高并发混合请求测试完成 (finder_ad_banner, mblog, wenda_box):")
	t.Logf("  总请求数: %d", totalReqs)
	t.Logf("  并发数: %d", concurrency)
	t.Logf("  单个请求(finder_ad_banner) - 成功: %d, 失败: %d", singleSuccessCount, singleFailCount)
	t.Logf("  批量请求(mblog+wenda_box) - 成功: %d, 失败: %d", batchSuccessCount, batchFailCount)
	t.Logf("  总耗时: %v", duration)
}

// TestHighConcurrencyWithDifferentTimeouts 测试高并发下不同超时配置 - 使用真实接口 finder_ad_banner
func TestHighConcurrencyWithDifferentTimeouts(t *testing.T) {
	const (
		concurrency = 50
		totalReqs   = 200
	)

	timeouts := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		500 * time.Millisecond,
		1000 * time.Millisecond,
	}

	var (
		wg      sync.WaitGroup
		results = make(map[time.Duration]*struct {
			success int64
			fail    int64
		})
		mu sync.Mutex
	)

	// 初始化结果map
	for _, timeout := range timeouts {
		results[timeout] = &struct {
			success int64
			fail    int64
		}{}
	}

	sem := make(chan struct{}, concurrency)

	startTime := time.Now()

	for i := 0; i < totalReqs; i++ {
		wg.Add(1)
		go func(reqID int) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			ctx := newTestContext()

			// 轮询使用不同的超时配置
			timeout := timeouts[reqID%len(timeouts)]

			// 使用真实接口 finder_ad_banner
			params := map[string]any{
				"page":   1,
				"count":  10,
				"source": "search",
			}

			other := &api.OtherOption{
				Timeout: timeout,
			}

			result := api.Request(ctx, "finder_ad_banner", params, other)

			mu.Lock()
			if _, ok := result["error"]; ok {
				results[timeout].fail++
			} else {
				results[timeout].success++
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	t.Logf("高并发不同超时配置测试完成 (finder_ad_banner):")
	t.Logf("  总请求数: %d", totalReqs)
	t.Logf("  并发数: %d", concurrency)
	t.Logf("  总耗时: %v", duration)
	for _, timeout := range timeouts {
		r := results[timeout]
		t.Logf("  超时 %v - 成功: %d, 失败: %d", timeout, r.success, r.fail)
	}
}

// TestHighConcurrencyWithRetry 测试高并发下的重试机制 - 使用真实接口 mblog
func TestHighConcurrencyWithRetry(t *testing.T) {
	const (
		concurrency = 50
		totalReqs   = 200
	)

	var (
		wg           sync.WaitGroup
		successCount int64
		failCount    int64
	)

	sem := make(chan struct{}, concurrency)

	startTime := time.Now()

	for i := 0; i < totalReqs; i++ {
		wg.Add(1)
		go func(reqID int) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			ctx := newTestContext()

			// 使用真实接口 mblog
			params := map[string]any{
				"ids":    "4500000000000001,4500000000000002,4500000000000003",
				"source": "search",
			}

			other := &api.OtherOption{
				Timeout:  500 * time.Millisecond,
				TryTimes: 3, // 重试3次
			}

			result := api.Request(ctx, "mblog", params, other)

			if _, ok := result["error"]; ok {
				atomic.AddInt64(&failCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	t.Logf("高并发重试机制测试完成 (mblog):")
	t.Logf("  总请求数: %d", totalReqs)
	t.Logf("  并发数: %d", concurrency)
	t.Logf("  重试次数: 3")
	t.Logf("  成功数: %d", successCount)
	t.Logf("  失败数: %d", failCount)
	t.Logf("  总耗时: %v", duration)
}

// TestStressRequest 压力测试：持续高负载 - 使用真实接口 finder_ad_banner
func TestStressRequest(t *testing.T) {
	const (
		concurrency = 200
		duration    = 5 * time.Second // 持续5秒
	)

	var (
		wg           sync.WaitGroup
		successCount int64
		failCount    int64
		running      int64 = 1
	)

	sem := make(chan struct{}, concurrency)

	startTime := time.Now()

	// 启动请求生成器
	go func() {
		for atomic.LoadInt64(&running) == 1 {
			wg.Add(1)
			go func() {
				defer wg.Done()

				sem <- struct{}{}
				defer func() { <-sem }()

				if atomic.LoadInt64(&running) == 0 {
					return
				}

				ctx := newTestContext()

				// 使用真实接口 finder_ad_banner
				params := map[string]any{
					"page":   1,
					"count":  10,
					"source": "search",
				}

				other := &api.OtherOption{
					Timeout: 500 * time.Millisecond,
				}

				result := api.Request(ctx, "finder_ad_banner", params, other)

				if _, ok := result["error"]; ok {
					atomic.AddInt64(&failCount, 1)
				} else {
					atomic.AddInt64(&successCount, 1)
				}
			}()

			// 控制请求生成速率
			time.Sleep(time.Millisecond)
		}
	}()

	// 等待指定时间
	time.Sleep(duration)
	atomic.StoreInt64(&running, 0)

	// 等待所有请求完成
	wg.Wait()

	totalReqs := successCount + failCount
	actualDuration := time.Since(startTime)

	t.Logf("压力测试完成 (finder_ad_banner):")
	t.Logf("  测试时长: %v", actualDuration)
	t.Logf("  并发数: %d", concurrency)
	t.Logf("  总请求数: %d", totalReqs)
	t.Logf("  成功数: %d", successCount)
	t.Logf("  失败数: %d", failCount)
	t.Logf("  QPS: %.2f", float64(totalReqs)/actualDuration.Seconds())
}

// BenchmarkHighConcurrencyRequest 基准测试：高并发请求 - 使用真实接口 finder_ad_banner
func BenchmarkHighConcurrencyRequest(b *testing.B) {
	const concurrency = 50

	sem := make(chan struct{}, concurrency)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sem <- struct{}{}

			ctx := newTestContext()
			params := map[string]any{
				"page":   1,
				"count":  10,
				"source": "search",
			}

			other := &api.OtherOption{
				Timeout: 500 * time.Millisecond,
			}

			_ = api.Request(ctx, "finder_ad_banner", params, other)

			<-sem
		}
	})
}

// BenchmarkHighConcurrencyBatchRequest 基准测试：高并发批量请求 - 使用真实接口
func BenchmarkHighConcurrencyBatchRequest(b *testing.B) {
	const concurrency = 50

	sem := make(chan struct{}, concurrency)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sem <- struct{}{}

			ctx := newTestContext()
			requests := []*api.BatchRequestItem{
				{
					// mblog 接口
					URLKey: "mblog",
					Params: map[string]any{
						"ids":    "4500000000000001,4500000000000002,4500000000000003",
						"source": "search",
					},
					Other: &api.OtherOption{Timeout: 500 * time.Millisecond},
				},
				{
					// finder_ad_banner 接口
					URLKey: "finder_ad_banner",
					Params: map[string]any{
						"page":   1,
						"count":  10,
						"source": "search",
					},
					Other: &api.OtherOption{Timeout: 500 * time.Millisecond},
				},
				{
					// wenda_box 接口
					URLKey: "wenda_box",
					Params: map[string]any{
						"q":      "测试关键词",
						"uid":    "1234567890",
						"source": "search",
					},
					Other: &api.OtherOption{
						Timeout:     500 * time.Millisecond,
						ContentJSON: true,
					},
				},
			}

			_ = api.MRequest(ctx, requests)

			<-sem
		}
	})
}
