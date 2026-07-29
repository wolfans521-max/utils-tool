package render

import (
	"fmt"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/log"
	"git.intra.weibo.com/search_fe/wbutil-go/view/page/cardlist"
	"git.intra.weibo.com/search_fe/wbutil-go/view/page/discover"
	"git.intra.weibo.com/search_fe/wbutil-go/view/page/pageinfo"
)

const CardListType = "cardlistInfo"
const pageInfoType = "pageInfo"
const DiscoverType = "discover"

// discover_source 有效值
var validDiscoverSources = map[string]bool{
	"discover":        true,
	"pull_refresh":    true,
	"partial_refresh": true,
}

// Handle 页面返回结果转换及处理
func Handle(ctx *context.Context, response map[string]any) (map[string]any, error) {
	pageType := GetPageType(ctx, response)
	startTime := time.Now()

	var result map[string]any
	var err error

	switch pageType {
	case CardListType:
		result, err = cardlist.Handle(ctx, response)
	case pageInfoType:
		result, err = pageinfo.Handle(ctx, response)
	case DiscoverType:
		result, err = discover.Handle(ctx, response)
	default:
		// 未知类型，返回原始数据
		result = response
	}

	// 记录处理耗时日志
	elapsed := time.Since(startTime).Seconds()
	log.AddCustomMark(ctx, "mapiSoraHandle", fmt.Sprintf("%.4f", elapsed))

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetPageType 根据响应数据和请求参数判断页面类型
func GetPageType(ctx *context.Context, data map[string]any) string {
	// 1. 检查是否存在 cardlistInfo
	if _, ok := data[CardListType]; ok {
		return CardListType
	}

	// 2. 检查是否存在 pageInfo
	if _, ok := data[pageInfoType]; ok {
		return pageInfoType
	}

	// 3. 发现页通过传参 discover_source 来确定返回数据
	// taskType == 'loadMore' 时强制 discover_source = 'pull_refresh'
	taskType := ctx.DefaultRequestParam("taskType", "pull_refresh")
	if taskType == "loadMore" {
		ctx.SetRequestParam("discover_source", "pull_refresh")
	}

	discoverSource := ctx.DefaultRequestParam("discover_source", "")
	if validDiscoverSources[discoverSource] {
		return DiscoverType
	}

	return ""
}
