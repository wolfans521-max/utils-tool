package finder

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/log"
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
	"git.intra.weibo.com/search_fe/wbutil-go/view/container"
)

// processPayload 对发现页payload做兼容处理
func processPayload(ctx *context.Context, payload any, isFlow bool) any {
	_ = ctx
	_ = isFlow

	// 目前仅处理切片形态的 payload：[]any
	items, ok := payload.([]any)
	if !ok {
		return payload
	}

	for i, v := range items {
		m, ok := v.(map[string]any)
		if !ok {
			// 非 map[string]any 类型（如 *card.ItemGroup），直接保留原样
			continue
		}

		// 1) group 场景：遍历内部 items，处理 card_type=89 的广告卡片
		if cat, _ := m["category"].(string); cat == "group" {
			if rawItems, ok := m["items"].([]any); ok {
				for j, it := range rawItems {
					im, ok := it.(map[string]any)
					if !ok {
						continue
					}
					dataMap, ok := im["data"].(map[string]any)
					if !ok {
						continue
					}

					if ct, ok := toInt(dataMap["card_type"]); ok && ct == 89 {
						if mblog, ok := dataMap["mblog"].(map[string]any); ok {
							mblog = compatibleAdSuperAdWindowSquare(mblog)
							mergeMediaInfo(mblog)
							dataMap["mblog"] = mblog
							im["data"] = dataMap
							rawItems[j] = im
						}
					}
				}
				m["items"] = rawItems
			}
		}

		// 2) 非 group：如果有 data.mblog 且 mblog.mark 不为空，则做同样处理
		if dataMap, ok := m["data"].(map[string]any); ok {
			if mblog, ok := dataMap["mblog"].(map[string]any); ok {
				if mark, _ := mblog["mark"].(string); mark != "" {
					mergeMediaInfo(mblog)
					mblog = compatibleAdSuperAdWindowSquare(mblog)
					dataMap["mblog"] = mblog
					m["data"] = dataMap
				}
			}
		}

		items[i] = m
	}

	return items
}

// compatibleAdSuperAdWindowSquare 兼容超级大视窗广告可能出现丢失 readtimetype 字段的问题
func compatibleAdSuperAdWindowSquare(weibo map[string]any) map[string]any {
	if weibo == nil {
		return weibo
	}

	mark, _ := weibo["mark"].(string)
	if mark == "" {
		return weibo
	}

	const suffix = "pos5f4f39c30c2c3"
	if strings.HasSuffix(mark, suffix) {
		if rt, _ := weibo["readtimetype"].(string); rt == "" {
			weibo["readtimetype"] = "super_ad_window_square"
		}
	}

	return weibo
}

// mergeMediaInfo 将 mblog.customized.media_info 合并进 mblog.page_info.media_info
func mergeMediaInfo(mblog map[string]any) {
	if mblog == nil {
		return
	}

	customized, ok := mblog["customized"].(map[string]any)
	if !ok || customized == nil {
		return
	}
	pageInfo, ok := mblog["page_info"].(map[string]any)
	if !ok || pageInfo == nil {
		return
	}

	cMI, cOK := customized["media_info"]
	pMI, pOK := pageInfo["media_info"]
	if !cOK || !pOK || cMI == nil || pMI == nil {
		return
	}

	// 数组场景：按顺序拼接
	if cSlice, ok := cMI.([]any); ok {
		if pSlice, ok2 := pMI.([]any); ok2 {
			pageInfo["media_info"] = append(pSlice, cSlice...)
			mblog["page_info"] = pageInfo
			return
		}
	}

	// map 场景：后者覆盖前者
	if cMap, ok := cMI.(map[string]any); ok {
		if pMap, ok2 := pMI.(map[string]any); ok2 {
			for k, v := range cMap {
				pMap[k] = v
			}
			pageInfo["media_info"] = pMap
			mblog["page_info"] = pageInfo
		}
	}
}

// getHotFeedRenderFlow 处理热门流 + 主搜索流，生成统一容器结构
func getHotFeedRenderFlow(ctx *context.Context, resp map[string]any) map[string]any {
	if resp == nil {
		return map[string]any{}
	}

	// 1) 先处理主搜索流 data
	var searchItems []any
	if rawData, ok := resp["data"]; ok {
		if processed := processPayload(ctx, rawData, true); processed != nil {
			if arr, ok := processed.([]any); ok {
				searchItems = arr
			}
		}
	}
	if searchItems == nil {
		searchItems = []any{}
	}

	// 初始化结果，使用容器层生成完整的Flow结构
	var result map[string]any
	// 2) 热门流 hot_mix_rank_feed
	if hotRaw, ok := resp["hot_mix_rank_feed"].(map[string]any); ok && len(hotRaw) > 0 {
		// 2.1 屏效优化
		optimized := improveEfficiency(hotRaw)

		// 2.2 使用容器层 FormatFeed 做容器化，生成完整的Flow结构
		var flowData *container.Flow
		if flow := container.NewContainer().FormatFeed(ctx, optimized); flow != nil {
			flowData = flow
		}

		if flowData == nil {
			flowData = container.NewFlow()
		}

		// 2.3 获取热门流items
		var hotItems []any
		if flowData != nil && flowData.Items != nil {
			hotItems = flowData.Items
		}
		if hotItems == nil {
			hotItems = []any{}
		}

		// 2.4 合并搜索流与热门流 items
		merged := make([]any, 0, len(searchItems)+len(hotItems))
		merged = append(merged, searchItems...)
		merged = append(merged, hotItems...)

		// 2.5 使用容器层生成的完整Flow结构作为结果
		result = flowData.ToMap()
		result["items"] = merged

		// 2.6 曝光日志
		typeCode := 900
		if v, ok := hotRaw["finder_hot_feed_back"]; ok {
			if n, ok2 := toInt(v); ok2 && n == 1 {
				typeCode = 16
			}
		}

		containerID := ctx.DefaultRequestParam("containerid", "")
		createFeedExposureLog(ctx, hotRaw, containerID, typeCode)
		sts, _ := hotRaw["statuses"].([]any)
		log.CardExp(ctx, sts, typeCode, true)
		log.RepostAdExpoContainer(ctx, sts, 202)

		// 2.7 添加打码日志3156 @施源 需求 2018.09.25
		actionLogExt := map[string]any{
			"extparam":        ctx.DefaultRequestParam("extparam", ""),
			"request_referer": ctx.DefaultRequestParam("request_referer", ""),
		}
		if pushMid := ctx.DefaultRequestParam("push_mid", ""); pushMid != "" {
			actionLogExt["push_mid"] = pushMid
		}
		log.Action(3156, "", actionLogExt, ctx)

		// 2.8 合并热门流返回的 loadedInfo（若存在）到 result.loadedInfo
		if flowData != nil && flowData.LoadedInfo != nil {
			loaded, _ := result["loadedInfo"].(map[string]any)
			if loaded == nil {
				loaded = make(map[string]any, len(flowData.LoadedInfo))
			}
			for k, val := range flowData.LoadedInfo {
				loaded[k] = val
			}
			result["loadedInfo"] = loaded
		}
	} else if subBandRaw, ok := resp["sub_band_feed"].(map[string]any); ok && len(subBandRaw) > 0 {
		result = subBandRaw
	} else {
		// 如果没有热门流，创建一个空的Flow结构
		flowData := container.NewFlow()
		result = flowData.ToMap()
		result["items"] = searchItems
	}

	// 3) 合并 loadedInfo（主搜索流返回的 loadedInfo 覆盖/补充热门流的部分）
	if v, ok := resp["loadedInfo"].(map[string]any); ok && v != nil {
		loaded, _ := result["loadedInfo"].(map[string]any)
		if loaded == nil {
			loaded = make(map[string]any, len(v))
		}
		for k, val := range v {
			loaded[k] = val
		}
		result["loadedInfo"] = loaded
	}

	return result
}

// improveEfficiency 提升热门流屏效
func improveEfficiency(feedData map[string]any) map[string]any {
	if feedData == nil {
		return map[string]any{}
	}

	// filter_improve_screen_shallow_publisher: 删除 headers
	if v, ok := feedData["filter_improve_screen_shallow_publisher"]; ok {
		if b, ok2 := v.(bool); (ok2 && b) || (!ok2 && v != nil) {
			delete(feedData, "headers")
		}
	}

	rawStatuses, ok := feedData["statuses"].([]any)
	if !ok || len(rawStatuses) == 0 {
		return feedData
	}

	// filter_improve_screen_tag
	var tag map[string]any
	if v, ok := feedData["filter_improve_screen_tag"].(map[string]any); ok {
		tag = v
	} else {
		tag = map[string]any{}
	}

	for i, st := range rawStatuses {
		status, ok := st.(map[string]any)
		if !ok {
			// 非 map[string]any 类型（如 *card.ItemGroup），直接保留原样
			continue
		}

		// insert_item 场景需要深入 items
		if cat, _ := status["item_category"].(string); cat == "insert_item" {
			if items, ok := status["items"].([]any); ok {
				for idx, it := range items {
					im, ok := it.(map[string]any)
					if !ok {
						continue
					}
					if imCat, _ := im["category"].(string); imCat == "feed" {
						if dataMap, ok := im["data"].(map[string]any); ok {
							im["data"] = improveEfficiencyDetail(tag, dataMap)
							items[idx] = im
						}
					}
				}
				status["items"] = items
			}
			rawStatuses[i] = status
			continue
		}

		// 普通 status
		rawStatuses[i] = improveEfficiencyDetail(tag, status)
	}

	feedData["statuses"] = rawStatuses
	return feedData
}

// improveEfficiencyDetail 对单条 status 做精简
func improveEfficiencyDetail(tag map[string]any, status map[string]any) map[string]any {
	if status == nil {
		return status
	}

	// comment: 去掉评论引导
	if v, ok := tag["comment"]; ok {
		if b, ok2 := v.(bool); (ok2 && b) || (!ok2 && v != nil) {
			delete(status, "enable_comment_guide")
		}
	}

	// attitude: 去掉 attitude bar
	if v, ok := tag["attitude"]; ok {
		if b, ok2 := v.(bool); (ok2 && b) || (!ok2 && v != nil) {
			delete(status, "show_attitude_bar")
		}
	}

	adRenderVersion := 0
	if v, ok := toInt(status["render_version"]); ok {
		adRenderVersion = v
	}

	// darwinTag: 精简 tag_struct
	if v, ok := tag["darwinTag"]; ok {
		if b, ok2 := v.(bool); (ok2 && b) || (!ok2 && v != nil) {
			if ts, ok := status["tag_struct"].([]any); ok && len(ts) > 0 {
				filtered := make([]any, 0, len(ts))
				for _, t := range ts {
					tagItem, ok := t.(map[string]any)
					if !ok {
						continue
					}

					typeName, _ := tagItem["otype"].(string)
					bdType, _ := tagItem["bd_object_type"].(string)
					if typeName == "place" || typeName == "adFeedDarwinTag" || typeName == "long_ad_message_tag" ||
						typeName == "adFeedDarwinTag_nature" || typeName == "hotWeiboTag" || bdType == "jbptag" || adRenderVersion == 2 {
						filtered = append(filtered, tagItem)
					}
				}
				status["tag_struct"] = filtered
			}
		}
	}

	return status
}

// searchLog 记录发现页相关的曝光/行为日志
func searchLog(ctx *context.Context, fid string, scenes int, cardListInfo map[string]any) {
	if fid != "" {
		ctx.SetRequestParam("fid", fid)
	}

	params := make(map[string]any)
	params["page"] = ctx.DefaultRequestParam("page", "")
	params["since_id"] = ctx.DefaultRequestParam("since_id", "")
	params["mid"] = ctx.DefaultRequestParam("mid", "")

	// page_pro、last、empty 字段始终输出（即使为空）
	if cardListInfo != nil {
		if v, ok := cardListInfo["page_attr"]; ok {
			params["page_pro"] = v
		} else {
			params["page_pro"] = ""
		}
		if v, ok := cardListInfo["last"]; ok {
			params["last"] = v
		} else {
			params["last"] = ""
		}
		if v, ok := cardListInfo["empty"]; ok {
			params["empty"] = v
		} else {
			params["empty"] = ""
		}
	} else {
		params["page_pro"] = ""
		params["last"] = ""
		params["empty"] = ""
	}

	// ext 字段：始终包含两个元素 [ext, extparam]
	params["ext"] = []string{
		ctx.DefaultRequestParam("ext", ""),
		ctx.DefaultRequestParam("extparam", ""),
	}

	params["orifid"] = ctx.DefaultRequestParam("oriuicode", "")
	params["request_referer"] = ctx.DefaultRequestParam("request_referer", "")

	if scenes != 0 {
		params["scenes"] = scenes
	}

	if ctx.DefaultRequestParam("discover_flow_enable", "") != "" {
		params["discover_flow_enable"] = ctx.DefaultRequestParam("discover_flow_enable", "")
	}

	// push搜索词
	if pushQuery := ctx.DefaultRequestParam("push_query", ""); pushQuery != "" {
		params["push_query"] = url.QueryEscape(pushQuery)
	}

	log.Action(26, ctx.DefaultRequestParam("fid", ""), params, ctx)
}

// createFeedExposureLog 发现页热门流曝光日志。
func createFeedExposureLog(ctx *context.Context, data map[string]any, cid string, typeCode int) {
	var (
		mids             []string
		fromCateIds      []string
		recommendSources []string
		contributors     []string
		cityId           string
		avatar           []string
		scheduleLiveCard []string
		jinxiFlag        bool
		jingxiFlags      []string
		jingxiAdids      []string
		keywords         []string
		keywordsAdids    []string
		redpacketFlags   []string
	)

	if sts, ok := data["statuses"].([]any); ok {
		for _, s := range sts {
			st, ok := s.(map[string]any)
			if !ok {
				continue
			}
			mids = append(mids, toStringWithDefault(st["mid"], ""))
			fromCateIds = append(fromCateIds, toStringWithDefault(st["from_cateid"], ""))
			recommendSources = append(recommendSources, toStringWithDefault(st["recommend_source"], ""))
			contributors = append(contributors, toStringWithDefault(st["contributor"], ""))
			cityId = toStringWithDefault(st["cityid"], "")
			avatarVal := 0
			if u, ok := st["user"].(map[string]any); ok {
				if av, ok := u["avatar"].(map[string]any); ok {
					if t, ok := toInt(av["type"]); ok {
						avatarVal = t
					}
				}
			}
			avatar = append(avatar, strconv.Itoa(avatarVal))
			scheduleLiveCardVal := 0
			if pi, ok := st["page_info"].(map[string]any); ok {
				if _, ok := pi["story_live_subscription_info"]; ok {
					scheduleLiveCardVal = 1
				}
			}
			scheduleLiveCard = append(scheduleLiveCard, strconv.Itoa(scheduleLiveCardVal))
			jinxiFlag = false
			var adaAdid string
			if ext, ok := st["extend_info"].(map[string]any); ok {
				if ada, ok := ext["attitude_dynamic_ad"].(map[string]any); ok {
					adType := toString(ada["adType"])
					jinxiFlag = (adType == "jingxi")
					adaAdid = toStringWithDefault(ada["adid"], "")
				}
			}

			jingxiFlags = append(jingxiFlags, strconv.Itoa(boolToInt(jinxiFlag)))
			if jinxiFlag {
				jingxiAdids = append(jingxiAdids, adaAdid)
			} else {
				jingxiAdids = append(jingxiAdids, "")
			}

			keywordsVal := 0
			var bgAdid string
			if ext, ok := st["extend_info"].(map[string]any); ok {
				if bg, ok := ext["bg_card"].(map[string]any); ok {
					if _, ok := bg["keywords"]; ok {
						keywordsVal = 1
					}
					bgAdid = toStringWithDefault(bg["adid"], "")
				}
			}
			keywords = append(keywords, strconv.Itoa(keywordsVal))
			keywordsAdids = append(keywordsAdids, bgAdid)

			hotspot, _ := st["hotspot"].(bool)
			rp, _ := st["redpacket"].(bool)
			var redpacketVal int
			if hotspot {
				if rp {
					redpacketVal = 2
				} else {
					redpacketVal = 1
				}
			} else {
				redpacketVal = 0
			}
			redpacketFlags = append(redpacketFlags, strconv.Itoa(redpacketVal))
		}
	}

	fromCateIdStr := utils.JoinSemicolonStrings(fromCateIds)
	recommendSourcesStr := utils.JoinSemicolonStrings(recommendSources)
	contributorStr := utils.JoinSemicolonStrings(contributors)
	avatarStr := utils.JoinSemicolonStrings(avatar)
	scheduleLiveCardStr := utils.JoinSemicolonStrings(scheduleLiveCard)
	jingxiStr := utils.JoinSemicolonStrings(jingxiFlags)
	jingxiAdIdStr := utils.JoinSemicolonStrings(jingxiAdids)
	keywordsAdIdStr := utils.JoinSemicolonStrings(keywordsAdids)

	maxIdVal := fmt.Sprint(data["max_id"])
	hotRequestIdVal := toStringWithDefault(data["hot_request_id"], "")

	logSuffixes := []string{
		"containerid=>" + substring(cid, 0, 6),
		"full_containerid=>" + cid,
		"source=>" + fromCateIdStr,
		"recommend_source=>" + recommendSourcesStr,
		"isPageUp=>" + ctx.DefaultRequestParam("scenes", ""),
		"from=>" + ctx.GetFrom(),
		"scenes=>" + ctx.DefaultRequestParam("scenes", ""),
		"wm=>" + ctx.Wm(),
		"luicode=>" + ctx.DefaultRequestParam("luicode", ""),
		"refresh=>" + ctx.DefaultRequestParam("refresh", ""),
		"networktype=>" + ctx.DefaultRequestParam("networktype", ""),
		"uicode=>" + ctx.DefaultRequestParam("uicode", ""),
		"ext=>" + ctx.DefaultRequestParam("ext", ""),
		"featurecode=>" + ctx.DefaultRequestParam("featurecode", ""),
		"contributor=>" + contributorStr,
		"refresh_sourceid=>" + ctx.DefaultRequestParam("refresh_sourceid", ""),
		"cityid=>" + cityId,
		"max_id=>" + maxIdVal,
		"hot_request_id=>" + hotRequestIdVal,
		"avatar_type=>" + avatarStr,
		"schedule_live_card=>" + scheduleLiveCardStr,
		"ugshare_type=>" + utils.JoinSemicolonStrings(redpacketFlags),
		"jingxi=>" + jingxiStr,
		"jingxi_adid=>" + jingxiAdIdStr,
		"keywords_adid=>" + keywordsAdIdStr,
		"stream_visible=>" + ctx.DefaultRequestParam("stream_visible", ""),
		"search_hot_mix_feed=>1",
	}

	lFidParam := ctx.DefaultRequestParam("lfid", "")
	logSuffixes = append(logSuffixes, "lfid=>"+lFidParam)

	if strings.HasPrefix(cid, "102803") {
		pushMidParam := ctx.DefaultRequestParam("push_mid", "")
		logSuffixes = append(logSuffixes, "push_mid=>"+pushMidParam)
	}

	log.FeedExpo(ctx, typeCode, mids, logSuffixes, "")
}

func toStringWithDefault(v any, def string) string {
	if v == nil {
		return def
	}
	s := toString(v)
	if s == "" {
		return def
	}
	return s
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func substring(s string, start, length int) string {
	if start >= len(s) {
		return ""
	}
	end := start + length
	if end > len(s) {
		end = len(s)
	}
	return s[start:end]
}

func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case fmt.Stringer:
		return x.String()
	}
	return fmt.Sprint(v)
}

// toInt 将 card_type 等动态类型安全地转换为 int。
func toInt(v any) (int, bool) {
	switch x := v.(type) {
	case int:
		return x, true
	case int8:
		return int(x), true
	case int16:
		return int(x), true
	case int32:
		return int(x), true
	case int64:
		return int(x), true
	case uint:
		return int(x), true
	case uint8:
		return int(x), true
	case uint16:
		return int(x), true
	case uint32:
		return int(x), true
	case uint64:
		return int(x), true
	case float32:
		return int(x), true
	case float64:
		return int(x), true
	case string:
		if x == "" {
			return 0, false
		}
		if n, err := strconv.Atoi(x); err == nil {
			return n, true
		}
	}
	return 0, false
}
