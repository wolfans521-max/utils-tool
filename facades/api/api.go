package api

import (
	"fmt"
	"sync"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"

	internalAPI "git.intra.weibo.com/search_fe/wbutil-go/internal/api"
)

var (
	manager     *internalAPI.Manager
	managerOnce sync.Once
	initErr     error
)

// OtherOption 其他请求选项
type OtherOption struct {
	Header      map[string]string // 额外请求头
	Timeout     time.Duration     // 超时时间（毫秒）
	ContentJSON bool              // 是否使用JSON格式发送数据
	TryTimes    int               // 重试次数（优先级高于配置文件）
	HTTPMethod  string            // HTTP 方法，仅 Motan 协议使用
}

// BatchRequestItem 批量请求项
type BatchRequestItem struct {
	URLKey string                 // API配置键
	Params map[string]interface{} // 请求参数
	Other  *OtherOption           // 其他选项
}

// Request 单个请求
//
//	other := &api.OtherOption{
//	    Timeout: 100 * time.Millisecond,
//	}
//
// result := api.Request(ctx, "api_name", params, other)
func Request(ctx *context.Context, urlKey string, params map[string]interface{}, other ...*OtherOption) map[string]interface{} {
	mgr, err := GetManager()
	if err != nil {
		return nil
	}
	var opt *internalAPI.OtherOption
	if len(other) > 0 && other[0] != nil {
		// 转换为 internal 包的 OtherOption
		opt = &internalAPI.OtherOption{
			Header:      other[0].Header,
			Timeout:     other[0].Timeout,
			ContentJSON: other[0].ContentJSON,
			HTTPMethod:  other[0].HTTPMethod,
			TryTimes:    other[0].TryTimes,
		}
	}

	result, _ := mgr.Request(ctx, urlKey, params, opt)
	return result
}

// MRequest 批量并发API请求
//
//	results := api.MRequest(ctx, []*api.BatchRequestItem{
//	    {URLKey: "mblog", Params: map[string]interface{}{"uid": 123}},
//	    {URLKey: "user_info", Params: map[string]interface{}{"uid": 456}},
//	})
func MRequest(ctx *context.Context, requests []*BatchRequestItem) []map[string]interface{} {
	mgr, err := GetManager()
	if err != nil {
		return nil
	}

	// 转换为 internal 包的 BatchRequest
	internalRequests := make([]*internalAPI.BatchRequest, len(requests))
	for i, req := range requests {
		var internalOther *internalAPI.OtherOption
		if req.Other != nil {
			internalOther = &internalAPI.OtherOption{
				Header:      req.Other.Header,
				Timeout:     req.Other.Timeout,
				ContentJSON: req.Other.ContentJSON,
				TryTimes:    req.Other.TryTimes,
			}
		}
		internalRequests[i] = &internalAPI.BatchRequest{
			URLKey: req.URLKey,
			Params: req.Params,
			Other:  internalOther,
		}
	}

	results, _ := mgr.BatchRequest(ctx, internalRequests)
	return results
}

func GetManager() (*internalAPI.Manager, error) {
	managerOnce.Do(func() {
		config, err := internalAPI.LoadConfig()
		if err != nil {
			initErr = fmt.Errorf("加载配置失败: %w", err)
			return
		}

		manager, initErr = internalAPI.NewManager(config)
	})

	return manager, initErr
}
