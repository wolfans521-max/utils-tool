package discover

import (
	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/view/container/finder"
	"git.intra.weibo.com/search_fe/wbutil-go/view/template"
)

const (
	// sourceDiscover 首次进入发现页
	sourceDiscover = "discover"
	// sourcePullRefresh 下拉刷新
	sourcePullRefresh = "pull_refresh"
	// sourcePartialRefresh 异步局部刷新
	sourcePartialRefresh = "partial_refresh"
)

// Handle 根据 discover_source 分发发现页三种形态的输出
// 入参 resp 为下游接口返回的原始结构（通常是解码后的 JSON map），由调用方保证结构正确。
func Handle(ctx *context.Context, resp map[string]any) (map[string]any, error) {
	// taskType 默认 pull_refresh，当 taskType=loadMore 时强制 discover_source=pull_refresh
	taskType := ctx.DefaultRequestParam("taskType", "")
	if taskType == "loadMore" {
		ctx.SetRequestParam("discover_source", sourcePullRefresh)
	}

	discoverSource := ctx.DefaultRequestParam("discover_source", "")

	switch discoverSource {
	case sourcePullRefresh:
		// 下拉刷新：容器层使用 Timeline 逻辑处理，再走 containerTimelineRender
		timeline := finder.TimelineHandle(ctx, resp)
		return template.ContainerTimelineRender(ctx, timeline), nil
	case sourcePartialRefresh:
		// 异步局部刷新：容器层使用 Discover 逻辑处理，再走 containerDiscoverRender
		partial := finder.PartialDiscoverHandle(ctx, resp)
		return template.ContainerDiscoverRender(ctx, partial), nil
	case sourceDiscover:
		fallthrough
	default:
		// 首次进入发现页（或 discover_source 缺失时的默认行为）：
		// 使用 Finder 逻辑处理整个频道结构，再走 finderRender
		first := finder.Handle(ctx, resp)
		return template.Render(ctx, first), nil
	}
}
