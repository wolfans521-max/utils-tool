package api

import (
	"strings"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/log"
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/spf13/cast"
)

// 业务号常量定义
const (
	BizIdDiscoverPage = "102803" // 发现页
	BizIdSearchPage1  = "100103" // 搜索落地页
	BizIdSearchPage2  = "100303" // 搜索落地页
	BizIdSearchPage3  = "231522" // 搜索落地页
)

func defaultMblogParamsByBizId(ctx *context.Context) map[string]any {
	params := make(map[string]any)
	containerId := ctx.DefaultRequestParam("containerid", "")
	bizId := strutil.Substring(containerId, 0, 6)

	switch bizId {
	case BizIdDiscoverPage:
		params["uid"] = ctx.GetUserId()
		params["filter_content"] = 0
		params["show_more_pic"] = 1
		params["is_globalfilter"] = 1
		params["ad_process"] = 1
		params["ui_type"] = 1
		params["with_positive_flag"] = 1
		params["shallow_consume"] = "1"
		params["add_readcount"] = 1
		params["flowId"] = ctx.DefaultRequestParam("containerid", "")
		params["sceneCardId"] = strutil.Substring(containerId, 0, 6) + "000"
		params["card_opts"] = "{\"mobile\":true}"

	case BizIdSearchPage1, BizIdSearchPage2, BizIdSearchPage3:
		params["isGetLongText"] = 1
		params["need_new_fast"] = 1
		params["show_card_fields"] = 7
	}

	return params
}

// 组装微博接口参数
// isRender: true表示用于render渲染模式，false表示用于普通请求模式
func buildMblogParams(ctx *context.Context, ids string, isRender bool, extParams map[string]any) map[string]any {
	params := defaultMblogParamsByBizId(ctx)
	params["ids"] = ids
	params["rip"] = ctx.ClientIP()
	params["isvideoad"] = 1
	params["source"] = 2936099636

	// 新网关优先使用spr透传，新showabtch会根据spr内的fid等信息，做文章、视频等点击布码使用
	params["spr"] = ctx.GetSpr()

	if val, ok := extParams["uid"]; ok {
		params["uid"] = val
	}
	if val, ok := extParams["filter_content"]; ok {
		params["filter_content"] = val
	}
	if val, ok := extParams["show_more_pic"]; ok {
		params["show_more_pic"] = val
	}
	if val, ok := extParams["is_globalfilter"]; ok {
		params["is_globalfilter"] = val
	}
	if val, ok := extParams["ad_process"]; ok {
		params["ad_process"] = val
	}
	if val, ok := extParams["ui_type"]; ok {
		params["ui_type"] = val
	}
	if val, ok := extParams["with_positive_flag"]; ok {
		params["with_positive_flag"] = val
	}
	if val, ok := extParams["shallow_consume"]; ok {
		params["shallow_consume"] = val
	}
	if val, ok := extParams["card_opts"]; ok {
		params["card_opts"] = val
	}
	if val, ok := extParams["sceneCardId"]; ok {
		params["sceneCardId"] = val
	}
	if val, ok := extParams["isGetLongText"]; ok {
		params["isGetLongText"] = val
	}
	if val, ok := extParams["need_new_fast"]; ok {
		params["need_new_fast"] = val
	}
	if val, ok := extParams["is_agg_card"]; ok {
		params["is_agg_card"] = val
	}
	if val, ok := extParams["use_new_timeline_status_style"]; ok {
		params["use_new_timeline_status_style"] = val
	}

	// 增加阅读数
	sid := ctx.DefaultRequestParam("sid", "")
	if sid == "t_wap_ios" || sid == "t_wap_android" || sid == "t_wap" || sid == "t_search" || sid == "t_search_nologin" {
		params["add_readcount"] = 1
	}

	// render模式：无论是否为空都添加到params
	if isRender {
		params["c"] = ctx.C()
		params["skin"] = ctx.DefaultRequestParam("skin", "default")
		params["v_p"] = ctx.DefaultRequestParam("vp", "90")
		if imageType := ctx.DefaultRequestParam("image_type", ""); imageType != "" {
			params["image_type"] = imageType
		}
		if vp := ctx.DefaultRequestParam("v_p", ""); vp != "" {
			params["v_p"] = vp
		}
		if fs := ctx.DefaultRequestParam("f_s", ""); fs != "" {
			params["f_s"] = fs
		}
		if ispush := ctx.DefaultRequestParam("ispush", ""); ispush != "" {
			params["ispush"] = ispush
		}
		if refresh := ctx.DefaultRequestParam("refresh", ""); refresh != "" {
			params["refresh"] = refresh
		}
		if ptl := ctx.DefaultRequestParam("ptl", ""); ptl != "" {
			params["ptl"] = ptl
		}
		if containerid := ctx.DefaultRequestParam("containerid", ""); containerid != "" {
			params["flowId"] = containerid
		}

		// show_card_fields 参数，是一个按位的开关，默认只出url_struct、page_info等渲染后的字段。
		if val, ok := extParams["show_card_fields"]; ok {
			params["show_card_fields"] = val
		}
	}

	return params
}

// RenderMblogOption 微博render showbatch options，用于接口并发请求
func RenderMblogOption(ctx *context.Context, ids string, extParams map[string]any) *BatchRequestItem {
	params := buildMblogParams(ctx, ids, true, extParams)
	return &BatchRequestItem{
		URLKey: "render_show_batch",
		Params: params,
	}
}

// ProcessApiResult 处理博文接口返回结果，提取 statuses 字段
func ProcessApiResult(ctx *context.Context, ids string, apiResult map[string]any) []any {
	if apiResult == nil {
		return []any{}
	}

	statuses, exists := apiResult["statuses"]
	if !exists || statuses == nil {
		return []any{}
	}

	statusesList, ok := statuses.([]any)
	if !ok {
		return []any{}
	}

	recordMblogResultLog(ctx, ids, statusesList)

	return statusesList
}

// recordMblogResultLog 比较ids和statusesList中元素的mid，记录diff到req日志
func recordMblogResultLog(ctx *context.Context, ids string, statusesList []any) {
	if ctx == nil || ids == "" {
		return
	}

	responseMids := make(map[string]bool)
	for _, item := range statusesList {
		if itemMap, ok := item.(map[string]any); ok {
			if mid := utils.MapValue[string](itemMap, "mid"); mid != "" {
				responseMids[mid] = true
			}
		}
	}

	idList := strings.Split(ids, ",")
	var diffMids []string
	for _, id := range idList {
		if id != "" && !responseMids[id] {
			diffMids = append(diffMids, id)
		}
	}

	log.Append(ctx, "mids:"+cast.ToString(len(idList)))
	log.Append(ctx, "fetch:"+cast.ToString(len(responseMids)))
	if len(diffMids) > 0 {
		log.Append(ctx, "showBatchDiff:"+strings.Join(diffMids, ","))
	}
}

// RenderMblogShowBatch 微博render showbatch请求
func RenderMblogShowBatch(ctx *context.Context, ids string, extParams map[string]any) []any {
	params := buildMblogParams(ctx, ids, false, extParams)
	apiResult := Request(ctx, "render_show_batch", params)
	result := ProcessApiResult(ctx, ids, apiResult)

	return result
}
