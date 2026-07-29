package finder

import (
	"fmt"
	"strconv"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/log"
	"git.intra.weibo.com/search_fe/wbutil-go/internal/config"
	_ "git.intra.weibo.com/search_fe/wbutil-go/tools/debug"
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
)

// searchChannelScenes 配置缓存
var searchChannelScenes []string

func init() {
	// 尝试从 feed.yaml 加载配置，如果失败则使用空切片
	v, err := config.FeedViper()
	if err == nil && v != nil {
		searchChannelScenes = v.GetStringSlice("search_channel_scenes.data")
	}
}

// contains 检查切片中是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type Result map[string]any

// TimelineResult 下拉刷新（热门流）场景的输出。
type TimelineResult map[string]any

// PartialDiscoverResult 异步局部刷新场景下的返回结构
type PartialDiscoverResult map[string]any

// Handle 首次进入发现页的容器层处理入口
func Handle(ctx *context.Context, resp map[string]any) Result {
	// 如果resp为nil，触发panic，由recovery中间件返回502错误响应
	if resp == nil {
		panic("[finderHandle] resp is nil, upstream service may have failed")
	}

	channelInfo := utils.ArrayValue(resp, "channelInfo", map[string]any{})
	channelInfoVal, ok := channelInfo.(map[string]any)
	if !ok {
		panic("[finderHandle] channelInfo type assertion failed")
	}
	channels := utils.ArrayValue(channelInfo, "channels", []any{})
	channelsVal, ok := channels.([]any)
	if !ok {
		panic("[finderHandle] channels type assertion failed")
	}

	// 如果channels为空，触发panic，由recovery中间件返回502错误响应
	if len(channelsVal) == 0 {
		panic("[finderHandle] channels is empty")
	}

	// 1) 选择当前频道索引
	curKey := getCurKey(ctx, channelInfoVal)
	if curKey < 0 || curKey >= len(channelsVal) {
		curKey = 0
	}

	ch, _ := channelsVal[curKey].(map[string]any)

	// 2) 计算 curPageId（含 square_remake 特殊逻辑）
	curPageId := getCurPageID(ctx, ch)
	pageType, _ := ch["pageDataType"].(string)

	// 3) 拆出当前频道的 payload / search / loadedInfo / params
	var curPayload map[string]any
	if payload := utils.ArrayValue(ch, "payload", nil); payload != nil {
		if p, ok := payload.(map[string]any); ok {
			curPayload = p
		}
	}

	loadedInfo := utils.ArrayValue(curPayload, "loadedInfo", map[string]any{})

	// 4) 设置 is_bigday_info 参数（兼容多种类型）
	params := utils.ArrayValue(ch, "params", map[string]any{})
	if isBigDayInfo := utils.ArrayValue(params, "is_bigday_info", nil); isBigDayInfo != nil {
		ctx.SetRequestParam("is_bigday_info", fmt.Sprint(isBigDayInfo))
	}

	// 删除当前频道上的原始 payload，后续按新结构回填
	delete(ch, "payload")
	channelsVal[curKey] = ch
	channelInfoVal["channels"] = channels
	resp["channelInfo"] = channelInfo

	// 5) 青少年模式：若是青少年且 pageType 为 flow，则不走热门流逻辑，直接透传搜索结果
	if ctx.IsTeenager() {
		log.Append(ctx, "wbutil-go teenager")
		if pageType == "flow" {
			ctx.SetRequestParam("folwId", curPageId)
			ch["payload"] = map[string]any{
				"items": curPayload["data"],
			}
		} else {
			ch["payload"] = curPayload["data"]
		}
		channelsVal[curKey] = ch
		channelInfoVal["channels"] = channels
		resp["channelInfo"] = channelInfo
		searchLog(ctx, curPageId, 0, nil)
		return resp
	}

	// 6) scenes 设置：当前频道 pageId 命中配置或当前为游客时设置为 1
	isScenes := contains(searchChannelScenes, curPageId)
	if isScenes || ctx.IsGuest() {
		ctx.SetRequestParam("scenes", "1")
	}

	// 7) payload_param 透传到请求参数
	if pp := utils.ArrayValue(params, "payload_param", nil); pp != nil {
		if ppMap, ok := pp.(map[string]any); ok {
			for k, v := range ppMap {
				if sv, ok2 := v.(string); ok2 {
					ctx.SetRequestParam(k, sv)
				} else {
					// 保底转字符串
					ctx.SetRequestParam(k, fmt.Sprint(v))
				}
			}
		}
	}

	// 8) pageType 为 flow 时设置 flowId
	if pageType == "flow" {
		ctx.SetRequestParam("flowId", curPageId)
	}

	// 8.1) 设置 fid 到请求参数中，确保3156日志能获取到fid
	ctx.SetRequestParam("fid", curPageId)

	// 9) 对当前频道 payload 调用 getHotFeedRenderFlow
	payload := map[string]any{}
	if len(curPayload) > 0 {
		payload = getHotFeedRenderFlow(ctx, curPayload)
	}

	// 10) 合并 loadedInfo
	payload["loadedInfo"] = loadedInfo

	// 11) 记录 searchLog(curPageId, scenes)
	scenesStr := ctx.DefaultRequestParam("scenes", "0")
	scenes := 0
	if v, err := strconv.Atoi(scenesStr); err == nil {
		scenes = v
	}
	searchLog(ctx, curPageId, scenes, nil)

	// 12) 将处理后的 payload 回填到当前频道
	if len(payload) > 0 {
		ch["payload"] = payload
		channelsVal[curKey] = ch
		channelInfoVal["channels"] = channels
		resp["channelInfo"] = channelInfo
	}

	return resp
}

func getCurPageID(ctx *context.Context, ch map[string]any) string {
	curPageId := ""
	if v, ok := ch["pageId"].(string); ok {
		curPageId = v
	}
	// square_remake == 1 时强制改写
	if ctx.DefaultRequestParam("square_remake", "0") == "1" {
		curPageId = "102803_ctg1_1780_-_ctg1_1780"
	}
	return curPageId
}

func getCurKey(ctx *context.Context, channelInfo map[string]any) int {
	// 默认从请求 tab 获取
	curPageId := ctx.DefaultRequestParam("tab", "")

	if cc, ok := channelInfo["channelConfig"].(map[string]any); ok {
		if sel, ok2 := cc["selectInfo"].(map[string]any); ok2 {
			if pid, ok3 := sel["pageId"].(string); ok3 && pid != "" {
				curPageId = pid
			}
		}
		if defSel, ok2 := cc["defaultSelectInfo"].(map[string]any); ok2 {
			if pid, ok3 := defSel["pageId"].(string); ok3 && pid != "" {
				// 记下默认 pageId
				// 若未找到匹配，将 fallback 到该默认索引
				defaultPageId := pid
				if chs, ok4 := channelInfo["channels"].([]any); ok4 {
					defaultKey := 0
					curKey := -1
					for idx, it := range chs {
						if m, ok5 := it.(map[string]any); ok5 {
							if p, ok6 := m["pageId"].(string); ok6 {
								if p == defaultPageId {
									defaultKey = idx
								}
								if curPageId != "" && p == curPageId {
									curKey = idx
								}
							}
						}
					}
					if curKey >= 0 {
						return curKey
					}
					return defaultKey
				}
			}
		}
	}

	// 如果未能通过配置找到，尝试直接匹配 channelInfo.channels
	if chs, ok := channelInfo["channels"].([]any); ok {
		for idx, it := range chs {
			if m, ok2 := it.(map[string]any); ok2 {
				if p, ok3 := m["pageId"].(string); ok3 && p == curPageId && p != "" {
					return idx
				}
			}
		}
	}

	return 0
}

// TimelineHandle 发现页下拉刷新（热门流）的容器层处理入口
func TimelineHandle(ctx *context.Context, resp map[string]any) TimelineResult {
	// 1) 使用 Common.getHotFeedRenderFlow 统一处理热门流 + 主搜索流
	flow := getHotFeedRenderFlow(ctx, resp)

	// 2) 记录行为日志 searchLog(flowId, scenes)
	flowId := ctx.DefaultRequestParam("flowId", "")
	scenesStr := ctx.DefaultRequestParam("scenes", "0")
	scenes := 0
	if v, err := strconv.Atoi(scenesStr); err == nil {
		scenes = v
	}
	searchLog(ctx, flowId, scenes, nil)

	return flow
}

// PartialDiscoverHandle 发现页异步局部刷新（partial_refresh）的容器层处理入口
func PartialDiscoverHandle(ctx *context.Context, resp map[string]any) PartialDiscoverResult {
	// 1) 取出 data 与 loadedInfo
	var rawData any
	if v, ok := resp["data"]; ok {
		rawData = v
	} else {
		// 兼容：如果没有 data 字段，则直接把整个 resp 当作 items
		rawData = resp
	}

	var loadedInfo map[string]any
	if v, ok := resp["loadedInfo"]; ok {
		if m, ok2 := v.(map[string]any); ok2 {
			loadedInfo = m
		}
	}

	// 2) payload 兼容处理
	flowEnable := ctx.DefaultRequestParam("discover_flow_enable", "") != ""
	processed := processPayload(ctx, rawData, flowEnable)

	// 3) iOS 特定版本的 fid 特判：仅当 iOS 且版本等于 F10 时设置
	if ctx.From != nil && ctx.IsIos() && ctx.Compare("=", "F10") {
		ctx.SetRequestParam("fid", "102803_ctg1_1780_-_ctg1_1780")
	}
	fid := ctx.DefaultRequestParam("fid", "")
	// 临时兼容 2026-06-24
	if fid == "5t8eh2pvtr" {
		fid = "102803_ctg1_1780_-_ctg1_1780"
	}

	// 4) 记录一次 searchLog（行为日志 26），scenes 固定为 1
	searchLog(ctx, fid, 1, nil)

	// 5) 返回局部刷新结构
	return map[string]any{
		"type":       "replace",
		"items":      processed,
		"loadedInfo": loadedInfo,
	}
}
