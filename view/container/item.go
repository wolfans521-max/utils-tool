package container

import (
	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
)

/**
card 转换为 item 结构
card转换文档：https://gitlab.weibo.cn/pagecard/task/-/issues/206
*/

// CardlistToItem 将 cards 转换为 items
func CardlistToItem(ctx *context.Context, cards []any, showStyle int, needSpan bool) []any {
	var items []any

	for _, cardItem := range cards {
		// 尝试将 cardItem 转换为 map[string]any
		// 支持 map[string]any 类型和结构体类型（如 *C11）
		var card map[string]any
		if m, ok := cardItem.(map[string]any); ok {
			card = m
		} else {
			// 使用 utils.StructToMap 将结构体转换为 map[string]any
			card = utils.StructToMap(cardItem)
			if len(card) == 0 {
				items = append(items, cardItem)
				continue
			}
		}

		cardType := getCardTypeValue(card)

		// flow 类型的直接透传
		if cardType == "flow" {
			delete(card, "card_type")
			items = append(items, card)
			continue
		}
		// dynamic 类型的直接透传
		if cardType == "dynamic" {
			delete(card, "card_type")
			items = append(items, card)
			continue
		}

		// card_type == 11 的处理
		if getCardTypeInt(card) == 11 {
			// is_asyn 异步卡片
			if getBool(card, "is_asyn") {
				items = append(items, getAsynCardItem(card, showStyle))
			}

			// to_grid_item 转换为 grid
			if getBool(card, "to_grid_item") {
				items = append(items, cardGroupToGrid(card))
			} else if hasCardGroup(card) {
				// 有 title 时先添加 title
				if title := getString(card, "title"); title != "" {
					items = append(items, getGroupTitle(title, showStyle))
				}
				// 转换为 vertical group
				items = append(items, cardGroupToVer(card, showStyle))
			}
			continue
		}

		// card_type 为 63 或 9 的处理（feed 类型）
		cardTypeInt := getCardTypeInt(card)
		if cardTypeInt == 63 || cardTypeInt == 9 {
			feedItem := getFeed(ctx, card)
			if showStyle == 1 {
				if data, ok := feedItem["data"].(map[string]any); ok {
					delete(data, "style")
				}
			}
			items = append(items, feedItem)
			continue
		}

		// 其他类型使用 getCardItem
		items = append(items, getCardItem(card, showStyle))
	}

	// 第1位31搜索框的card margin为0 @张宁
	if len(items) > 0 {
		if first, ok := items[0].(map[string]any); ok {
			if data, ok := first["data"].(map[string]any); ok {
				if getCardTypeInt(data) == 31 && showStyle == 1 {
					if first["style"] == nil {
						first["style"] = make(map[string]any)
					}
					first["style"].(map[string]any)["margin"] = []any{0, 0, 0, 0}
				}
			}
		}
	}

	// 流里第1位加span占位@文旭, 第一个是card31（搜索框）的不下发span @张宁
	if needSpan && len(items) > 0 {
		taskType := ""
		if ctx != nil {
			taskType = ctx.DefaultRequestParam("taskType", "")
		}

		if taskType != "loadMore" {
			firstCardType := 0
			if first, ok := items[0].(map[string]any); ok {
				if data, ok := first["data"].(map[string]any); ok {
					firstCardType = getCardTypeInt(data)
				}
			}

			if firstCardType == 156 {
				// card_type 156 特殊处理
				span := map[string]any{
					"category": "cell",
					"type":     "span",
					"span": map[string]any{
						"style": map[string]any{
							"height": 8,
						},
					},
				}
				if first, ok := items[0].(map[string]any); ok {
					if first["style"] == nil {
						first["style"] = make(map[string]any)
					}
					first["style"].(map[string]any)["margin"] = []any{0, 0, 0, 8}
				}
				items = append([]any{span}, items...)
			} else if firstCardType != 31 {
				// 非 card31 且非 filterGroupStyle
				filterGroupStyle := ""
				if ctx != nil {
					filterGroupStyle = ctx.DefaultRequestParam("filterGroupStyle", "")
				}
				if filterGroupStyle == "" {
					span := map[string]any{
						"category": "cell",
						"type":     "span",
						"span": map[string]any{
							"style": map[string]any{
								"height": 10,
							},
						},
					}
					items = append([]any{span}, items...)
				}
			}
		}
	}

	return items
}

// cardGroupToGrid card11的group转换为category=>grid类型，header和style都由业务方控制
func cardGroupToGrid(group map[string]any) map[string]any {
	data := map[string]any{
		"category": "group",
		"type":     "grid",
		"style":    group["group_style"],
	}

	if groupHeader := group["group_header"]; groupHeader != nil {
		data["header"] = groupHeader
	}
	if groupFooter := group["group_footer"]; groupFooter != nil {
		data["footer"] = groupFooter
	}
	if itemid := getString(group, "itemid"); itemid != "" {
		data["itemId"] = itemid
	}
	var items []any
	cardGroup := getCardGroup(group)
	if len(cardGroup) > 0 {
		showType := getIntOrDefault(group, "show_type", 0)

		for _, cardItem := range cardGroup {
			card, ok := cardItem.(map[string]any)
			if !ok {
				continue
			}

			// dynamic 类型的直接透传
			if getCardTypeValue(card) == "dynamic" {
				delete(card, "card_type")
				items = append(items, card)
				continue
			}

			cardTypeInt := getCardTypeInt(card)
			if cardTypeInt == 63 || cardTypeInt == 9 {
				feedItem := getFeed(nil, card)
				if itemStyle, ok := card["item_style"].(map[string]any); ok {
					if data, ok := feedItem["data"].(map[string]any); ok {
						data["style"] = itemStyle
					}
				}
				items = append(items, feedItem)
			} else {
				cardItem := getCardItem(card, showType)
				items = append(items, cardItem)
			}
		}
	}

	data["items"] = items
	return data
}

// cardGroupToVer card11的group转换为vertical类型
func cardGroupToVer(group map[string]any, showStyle int) map[string]any {
	var items []any
	showType := getIntOrDefault(group, "show_type", 0)

	cardGroup := getCardGroup(group)
	if cardGroup == nil {
		cardGroup = []any{}
	}
	num := len(cardGroup)
	for k, cardItem := range cardGroup {
		card, ok := cardItem.(map[string]any)
		if !ok {
			continue
		}

		cardTypeValue := getCardTypeValue(card)

		// flow 类型的直接透传
		if cardTypeValue == "flow" {
			delete(card, "card_type")
			items = append(items, card)
			continue
		}

		// dynamic 类型的直接透传
		if cardTypeValue == "dynamic" {
			delete(card, "card_type")
			items = append(items, card)
			continue
		}

		cardTypeInt := getCardTypeInt(card)
		if cardTypeInt == 63 || cardTypeInt == 9 {
			feedItem := getFeed(nil, card)
			if showType == 1 && k != num-1 {
				if data, ok := feedItem["data"].(map[string]any); ok {
					delete(data, "style")
				}
			} else {
				if data, ok := feedItem["data"].(map[string]any); ok {
					if data["style"] == nil {
						data["style"] = make(map[string]any)
					}
					style := data["style"].(map[string]any)
					style["hide_top_divide_line"] = 1
					style["divider_size"] = 0

					// 只影响最后一条
					if itemStyle, ok := card["item_style"].(map[string]any); ok {
						if mblogStyle, ok := itemStyle["mblog_style"].(map[string]any); ok {
							data["style"] = mblogStyle
						}
					}
				}
			}
			items = append(items, feedItem)
		} else {
			cardItem := getCardItem(card, showType)
			items = append(items, cardItem)

			// showType == 2 时添加分割线
			if showType == 2 && k != num-1 {
				items = append(items, map[string]any{
					"category": "cell",
					"type":     "span",
					"span": map[string]any{
						"style": map[string]any{
							"margin": []any{10, 0, 10, 0},
							"background": map[string]any{
								"type":     "color",
								"color":    "#eeeeee",
								"colorKey": "MainFeedBackgroundColor",
							},
						},
					},
					"style": map[string]any{
						"height": 0.5,
						"background": map[string]any{
							"type":     "color",
							"color":    "#ffffff",
							"colorKey": "CommonCardBackground",
						},
					},
				})
			}
		}
	}

	groupStruct := map[string]any{
		"category": "group",
		"type":     "vertical",
		"items":    []any{},
	}

	if len(items) == 0 {
		return groupStruct
	}

	// group内最后一个card 设置成（无）无分割线, 如果是feed设置"hide_divide_line": true
	k := len(items) - 1
	if lastItem, ok := items[k].(map[string]any); ok {
		category := getString(lastItem, "category")
		if category != "feed" && category != "dynamic" {
			showDivideLine := 0
			if data, ok := lastItem["data"].(map[string]any); ok {
				showDivideLine = getIntOrDefault(data, "show_divide_line", 0)
				itemStyle := data["item_style"]
				if showDivideLine == 0 && itemStyle == nil {
					if lastItem["style"] == nil {
						lastItem["style"] = make(map[string]any)
					}
					lastItem["style"].(map[string]any)["margin"] = []any{0, 0, 0, 0}
				}
			}
		}
	}

	// 设置 group style
	var groupStyle map[string]any
	if gs, ok := group["group_style"].(map[string]any); ok {
		groupStyle = gs
	} else {
		groupStyle = getShowStyle(showStyle)
	}

	groupStruct["itemId"] = getString(group, "itemid")
	groupStruct["items"] = items
	groupStruct["style"] = groupStyle

	// flow_struct 合并
	if flowStruct, ok := group["flow_struct"].(map[string]any); ok {
		for k, v := range flowStruct {
			groupStruct[k] = v
		}
	}

	return groupStruct
}

// getShowStyle 获取展示样式
func getShowStyle(showStyle int) map[string]any {
	var margin []any
	switch showStyle {
	case 0:
		margin = []any{0, 0, 0, 0.5}
	case 1:
		margin = []any{0, 0, 0, 10}
	default:
		margin = []any{0, 0, 0, 0}
	}

	return map[string]any{
		"margin":  margin,
		"padding": []any{0, 0, 0, 0},
	}
}

// getCardItem 获取卡片 item 结构
func getCardItem(card map[string]any, showStyle int) map[string]any {
	style := getShowStyle(showStyle)
	style["background"] = map[string]any{
		"type":     "color",
		"color":    "#ffffff",
		"colorKey": "CommonCardBackground",
	}

	// 支持业务方透传style样式
	if itemStyle, ok := card["item_style"].(map[string]any); ok {
		style = itemStyle
	} else {
		cardType := getCardTypeInt(card)
		if cardType == 112 && showStyle != 1 {
			// 112特殊样式，023时都是0，1根据样式显示
			style["margin"] = []any{0, 0, 0, 0}
			style["padding"] = []any{0, 0, 0, 0}
		} else if cardType == 180 && getIntOrDefault(card, "layout_style", 0) == 1 {
			style["background"].(map[string]any)["color"] = "#F2F2F2"
		} else {
			style = getSpecialCardStyle(cardType, style)
		}
	}

	item := map[string]any{
		"category": "card",
		"data":     card,
		"style":    style,
		"itemExt": map[string]any{
			"anchorId": getStringDefault(card, "anchorId", getString(card, "itemid")),
		},
	}

	// 合并 itemExt
	if itemExt, ok := card["itemExt"].(map[string]any); ok {
		existingItemExt := item["itemExt"].(map[string]any)
		for k, v := range itemExt {
			existingItemExt[k] = v
		}
	}

	return item
}

// getGroupTitle 获取 group 标题
func getGroupTitle(title string, styleId int) map[string]any {
	left := 10
	if styleId == 1 {
		left = 0
	}

	return map[string]any{
		"category": "cell",
		"type":     "text",
		"text": map[string]any{
			"style": map[string]any{
				"textSize":     12,
				"textColor":    "#939393",
				"textColorKey": "CommonGray93",
				"ellipsize":    "end",
				"maxLines":     1,
				"margin":       []any{12, left, 12, 10},
			},
			"content": title,
		},
	}
}

// getAsynCardItem 获取异步卡片 item
func getAsynCardItem(card map[string]any, showStyle int) map[string]any {
	card["card_type"] = 197
	card["height"] = 80
	card["asyn_api_path"] = "container/asyn"
	card["asyn_api_params"] = map[string]any{
		"itemid":     getString(card, "itemid"),
		"show_style": showStyle,
	}

	if asynApiExtparams := card["asyn_api_extparams"]; asynApiExtparams != nil {
		card["asyn_api_params"].(map[string]any)["extparam"] = asynApiExtparams
	}

	showCardStyle := getShowStyle(showStyle)
	style := map[string]any{
		"background": map[string]any{
			"type":     "color",
			"color":    "#ffffff",
			"colorKey": "CommonCardBackground",
		},
	}
	for k, v := range showCardStyle {
		style[k] = v
	}

	return map[string]any{
		"category": "card",
		"data":     card,
		"style":    style,
	}
}

// getSpecialCardStyle 获取特殊卡片样式
func getSpecialCardStyle(cardType int, style map[string]any) map[string]any {
	switch cardType {
	case 8, 10, 7, 3, 30:
		style["padding"] = []any{12, 0, 0, 0}
	case 42:
		// 版本判断 >C61，这里简化处理
		style["padding"] = []any{0, 0, 8, 0}
	case 4:
		// 版本判断 >C61，这里简化处理
		style["padding"] = []any{13, 0, 11, 0}
	case 127, 159:
		style["grid"] = true
	case 174:
		// @客户端蔡丰 要求
		style["padding"] = []any{13, 13, 13, 13}
	case 216:
		style["margin"] = []any{0, 0, 0, 0}
		style["padding"] = []any{0, 0, 0, 0}
		delete(style, "background")
	case 58, 156, 164:
		// 只保留margin padding
		delete(style, "background")
	}

	return style
}

// getFeed 获取 feed 结构
func getFeed(ctx *context.Context, card map[string]any) map[string]any {
	mblog, ok := card["mblog"].(map[string]any)
	if !ok {
		mblog = make(map[string]any)
	}

	// 处理 flag_img
	if titleSource, ok := mblog["title_source"].(map[string]any); ok {
		if flagImg := titleSource["flag_img"]; flagImg != nil {
			if headerInfo, ok := mblog["header_info"].(map[string]any); ok {
				if avatar, ok := headerInfo["avatar"].(map[string]any); ok {
					avatar["flag_img"] = flagImg
				}
			}
		}
	}

	// 过滤 mblog
	mblog = filterMblog(mblog)

	// 设置 style
	if itemStyle, ok := card["item_style"].(map[string]any); ok {
		if mblogStyle, ok := itemStyle["mblog_style"].(map[string]any); ok {
			mblog["style"] = mblogStyle
		} else {
			mblog["style"] = map[string]any{
				"hide_divide_line": true,
			}
		}
	} else {
		mblog["style"] = map[string]any{
			"hide_divide_line": true,
		}
	}

	// 加打点结构
	if actionlog := card["actionlog"]; actionlog != nil {
		mblog["click_actionlog"] = actionlog
	}

	weiboData := map[string]any{
		"category": "feed",
		"data":     mblog,
	}

	// isMiniFlow 处理
	isMiniFlow := 0
	if ctx != nil {
		isMiniFlow = ctx.RequestInt("isMiniFlow")
	}
	if isMiniFlow == 1 {
		weiboData["type"] = "mini"
	}

	// 版本判断 <C42，这里简化处理，假设都是新版本

	// itemExt 处理
	if itemid := getString(card, "itemid"); itemid != "" {
		if weiboData["itemExt"] == nil {
			weiboData["itemExt"] = make(map[string]any)
		}
		itemExt := weiboData["itemExt"].(map[string]any)
		itemExt["anchorId"] = getStringDefault(card, "anchorId", itemid)

		if mblogId := getString(mblog, "id"); itemid != mblogId {
			itemExt["extraItemId"] = itemid
		}
	}

	// is_container_scheme 处理
	if getBool(card, "is_container_scheme") {
		if scheme := getString(card, "scheme"); scheme != "" {
			if weiboData["itemExt"] == nil {
				weiboData["itemExt"] = make(map[string]any)
			}
			weiboData["itemExt"].(map[string]any)["redirectScheme"] = scheme
		}
	}

	return weiboData
}

// filterMblog 端上新流用工具解析json,类型/结构不对会导致解析错误，这里一个个字段处理
func filterMblog(mblog map[string]any) map[string]any {
	// 删除空字符串字段
	for k, v := range mblog {
		if str, ok := v.(string); ok && str == "" {
			delete(mblog, k)
		}
	}

	// 删除空 title
	if title, ok := mblog["title"]; ok {
		if str, ok := title.(string); ok && str == "" {
			delete(mblog, "title")
		}
	}

	// 处理 retweeted_status.geo
	if retweetedStatus, ok := mblog["retweeted_status"].(map[string]any); ok {
		if geo, ok := retweetedStatus["geo"]; ok {
			if str, ok := geo.(string); ok && str == "" {
				retweetedStatus["geo"] = nil
			}
		}
	}

	// 过滤 user
	if user, ok := mblog["user"].(map[string]any); ok {
		mblog["user"] = filterUser(user)
	}

	// 评论聚合结构过滤user
	if commentSummary, ok := mblog["comment_summary"].(map[string]any); ok {
		if commentList, ok := commentSummary["comment_list"].([]any); ok {
			for k, v := range commentList {
				if comment, ok := v.(map[string]any); ok {
					if user, ok := comment["user"].(map[string]any); ok {
						comment["user"] = filterUser(user)
					}
					if replyComment, ok := comment["reply_comment"].(map[string]any); ok {
						if user, ok := replyComment["user"].(map[string]any); ok {
							replyComment["user"] = filterUser(user)
						}
					}
					commentList[k] = comment
				}
			}
			commentSummary["comment_list"] = commentList
		}
	}

	// 类型转换
	btnFields := []string{"buttons", "mblog_menus", "mblog_menus_new", "mblog_buttons"}
	for _, b := range btnFields {
		if buttons, ok := mblog[b].([]any); ok {
			for m, button := range buttons {
				if btn, ok := button.(map[string]any); ok {
					// show_loading 转换
					if showLoading, ok := btn["show_loading"].(string); ok {
						if showLoading == "1" || showLoading == "true" {
							btn["show_loading"] = 1
						} else {
							btn["show_loading"] = 0
						}
					}
					// sub_type 转换
					if subType, ok := btn["sub_type"].(string); ok {
						if n, ok := getInt(map[string]any{"v": subType}, "v"); ok {
							btn["sub_type"] = n
						}
					}
					buttons[m] = btn
				}
			}
			mblog[b] = buttons
		}
	}

	return mblog
}

// filterUser 过滤用户信息
func filterUser(user map[string]any) map[string]any {
	if user == nil {
		return user
	}

	// allow_all_comment 转换为 bool
	if allowAllComment, ok := user["allow_all_comment"].(int); ok {
		user["allow_all_comment"] = allowAllComment != 0
	}

	// verified_type_ext 空字符串转为 0
	if verifiedTypeExt, ok := user["verified_type_ext"].(string); ok && verifiedTypeExt == "" {
		user["verified_type_ext"] = 0
	}

	// 删除空 badge
	if badge, ok := user["badge"]; ok {
		if badgeMap, ok := badge.(map[string]any); ok && len(badgeMap) == 0 {
			delete(user, "badge")
		}
		if badgeSlice, ok := badge.([]any); ok && len(badgeSlice) == 0 {
			delete(user, "badge")
		}
	}

	return user
}

// getCardTypeValue 获取卡片类型值（字符串形式）
func getCardTypeValue(card map[string]any) string {
	if cardType, ok := card["card_type"].(string); ok {
		return cardType
	}
	return ""
}

// getCardTypeInt 获取卡片类型值（整数形式）
func getCardTypeInt(card map[string]any) int {
	if cardType, ok := card["card_type"].(int); ok {
		return cardType
	}
	if cardType, ok := card["card_type"].(int64); ok {
		return int(cardType)
	}
	if cardType, ok := card["card_type"].(float64); ok {
		return int(cardType)
	}
	return 0
}

// hasCardGroup 检查 card 是否有非空的 card_group
// 支持 []any 和 []map[string]any 两种类型
func hasCardGroup(card map[string]any) bool {
	if card == nil {
		return false
	}
	cg := card["card_group"]
	if cg == nil {
		return false
	}
	switch v := cg.(type) {
	case []any:
		return len(v) > 0
	case []map[string]any:
		return len(v) > 0
	}
	return false
}

// getCardGroup 获取 card_group 并转换为 []any
// 支持 []any 和 []map[string]any 两种类型
func getCardGroup(card map[string]any) []any {
	if card == nil {
		return nil
	}
	cg := card["card_group"]
	if cg == nil {
		return nil
	}
	switch v := cg.(type) {
	case []any:
		return v
	case []map[string]any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = item
		}
		return result
	}
	return nil
}
