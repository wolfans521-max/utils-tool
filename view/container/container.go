package container

import (
	"strconv"
	"strings"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
)

// Flow 表示容器层统一的数据结构，用于承载 Discover 等页面的渲染结果
type Flow struct {
	Items       []any          `json:"items"`
	PageData    map[string]any `json:"pageData"`
	RefreshInfo map[string]any `json:"refreshInfo"`
	MoreInfo    map[string]any `json:"moreInfo"`
	LoadedInfo  map[string]any `json:"loadedInfo"`
	Config      map[string]any `json:"config"`
}

// NewFlow 创建一个带有空字段的 Flow，避免上层访问时需要频繁判空。
func NewFlow() *Flow {
	return &Flow{
		Items:       []any{},
		PageData:    map[string]any{},
		RefreshInfo: map[string]any{},
		MoreInfo:    map[string]any{},
		LoadedInfo:  map[string]any{},
		Config:      map[string]any{},
	}
}

// WithItems 设置 items 字段并返回自身，便于链式调用。
func (f *Flow) WithItems(items []any) *Flow {
	if f == nil {
		return f
	}
	if items == nil {
		f.Items = []any{}
		return f
	}
	f.Items = items
	return f
}

// WithPageData 设置 pageData 字段。
func (f *Flow) WithPageData(pageData map[string]any) *Flow {
	if f == nil {
		return f
	}
	if pageData == nil {
		f.PageData = map[string]any{}
		return f
	}
	f.PageData = pageData
	return f
}

// WithRefreshInfo 设置 refreshInfo 字段。
func (f *Flow) WithRefreshInfo(refreshInfo map[string]any) *Flow {
	if f == nil {
		return f
	}
	if refreshInfo == nil {
		f.RefreshInfo = map[string]any{}
		return f
	}
	f.RefreshInfo = refreshInfo
	return f
}

// WithMoreInfo 设置 moreInfo 字段。
func (f *Flow) WithMoreInfo(moreInfo map[string]any) *Flow {
	if f == nil {
		return f
	}
	if moreInfo == nil {
		f.MoreInfo = map[string]any{}
		return f
	}
	f.MoreInfo = moreInfo
	return f
}

// WithLoadedInfo 设置 loadedInfo 字段。
func (f *Flow) WithLoadedInfo(loadedInfo map[string]any) *Flow {
	if f == nil {
		return f
	}
	if loadedInfo == nil {
		f.LoadedInfo = map[string]any{}
		return f
	}
	f.LoadedInfo = loadedInfo
	return f
}

// WithConfig 设置 config 字段。
func (f *Flow) WithConfig(config map[string]any) *Flow {
	if f == nil {
		return f
	}
	if config == nil {
		f.Config = map[string]any{}
		return f
	}
	f.Config = config
	return f
}

// ToMap 将 Flow 转换为 map[string]any 形式，便于与现有 map 结构互操作。
func (f *Flow) ToMap() map[string]any {
	if f == nil {
		return map[string]any{}
	}
	return map[string]any{
		"items":       f.Items,
		"pageData":    f.PageData,
		"refreshInfo": f.RefreshInfo,
		"moreInfo":    f.MoreInfo,
		"loadedInfo":  f.LoadedInfo,
		"config":      f.Config,
	}
}

type Container struct{}

func NewContainer() *Container { return &Container{} }

// Handle 容器化处理入口
func Handle(ctx *context.Context, timeline map[string]any) map[string]any {
	c := NewContainer()

	if _, ok := timeline["statuses"]; ok {
		// Feed 流格式化
		flow := c.FormatFeed(ctx, timeline)
		return flow.ToMap()
	} else if _, ok := timeline["cards"]; ok {
		// Cardlist 流格式化
		return c.FormatCardList(ctx, timeline)
	}

	return timeline
}

// FormatCardList cardlist 流内容格式化转换
func (c *Container) FormatCardList(ctx *context.Context, timeline map[string]any) map[string]any {
	showStyle := 2
	if cl, ok := timeline["cardlistInfo"].(map[string]any); ok {
		if ss, ok := getInt(cl, "show_style"); ok {
			showStyle = ss
		}
	} else if pi, ok := timeline["pageInfo"].(map[string]any); ok {
		if ss, ok := getInt(pi, "show_style"); ok {
			showStyle = ss
		}
	}

	needSpan := true
	if mapi, ok := timeline["mapi"].(map[string]any); ok {
		if noSpan, ok := mapi["no_span"].(bool); ok && noSpan {
			needSpan = false
		}
	}

	var items []any
	if cards, ok := timeline["cards"].([]any); ok && len(cards) > 0 {
		// 使用新的 CardlistToItem 函数
		items = CardlistToItem(ctx, cards, showStyle, needSpan)
	}

	data := map[string]any{
		"items": items,
	}

	// 处理 cards_xxx 变成 items_xxx
	for key := range timeline {
		if strings.HasPrefix(key, "cards_") {
			itemKey := "items_" + strings.TrimPrefix(key, "cards_")
			if cardsArr, ok := timeline[key].([]any); ok {
				itemsArr := CardlistToItem(ctx, cardsArr, showStyle, needSpan)
				// 去掉第一个 span
				if len(itemsArr) > 0 {
					if first, ok := itemsArr[0].(map[string]any); ok {
						if first["type"] == "span" {
							itemsArr = itemsArr[1:]
						}
					}
				}
				data[itemKey] = itemsArr
			}
			delete(timeline, key)
		}
	}

	data["pageData"] = c.dealPageData(ctx, timeline, "cardlistInfo")
	data["refreshInfo"] = c.dealRefreshInfo(timeline, "cardlistInfo")
	data["loadedInfo"] = c.dealLoadedInfo(timeline, "cardlistInfo")
	data["config"] = c.dealConfig()

	// 处理 moreInfo
	count := 20
	if cl, ok := timeline["cardlistInfo"].(map[string]any); ok {
		if ps, ok := getInt(cl, "pagesize"); ok {
			count = ps
		}
	}

	sinceId := ""
	total := 0
	if cl, ok := timeline["cardlistInfo"].(map[string]any); ok {
		if sid := getString(cl, "since_id"); sid != "" {
			sinceId = sid
		}
		if t, ok := getInt(cl, "total"); ok {
			total = t
		}
	}

	page := 1
	if ctx != nil {
		if p := ctx.RequestInt("page"); p > 0 {
			page = p
		}
	}

	if sinceId != "" || (ctx != nil && ctx.DefaultRequestParam("since_id", "") == "" && total > page*count) {
		data["moreInfo"] = c.dealMoreInfo(timeline, "cardlistInfo")
	}

	// 处理 ext 数据
	extData := c.dealExt(timeline, "cardlistInfo")
	if len(extData) > 0 {
		for k, v := range extData {
			data[k] = v
		}
	}

	return data
}

// FormatFeed 格式化feed博文流
func (c *Container) FormatFeed(ctx *context.Context, timeline map[string]any) *Flow {
	flow := NewFlow()
	if timeline == nil {
		return flow
	}

	flow.Items = c.dealStatuses(timeline)

	if !ctx.RequestHas("load_mode") {
		flow.PageData = c.dealPageData(ctx, timeline, "cardlistInfo")
		flow.RefreshInfo = c.dealRefreshInfo(timeline, "cardlistInfo")
		flow.LoadedInfo = c.dealLoadedInfo(timeline, "cardlistInfo")
		flow.Config = c.dealConfig()
		if maxID, ok := getInt(timeline, "max_id"); ok && maxID != 0 {
			flow.MoreInfo = c.dealMoreInfo(timeline, "cardlistInfo")
		}
		// gsid/sut 冗余下发
		if auth, ok := flow.LoadedInfo["authInfo"].(map[string]any); ok {
			if g := getString(timeline, "gsid"); g != "" {
				auth["gsid"] = g
			}
			if s := getString(timeline, "sut"); s != "" {
				auth["sut"] = s
			}
			flow.LoadedInfo["authInfo"] = auth
			if g, _ := auth["gsid"].(string); g != "" {
				flow.PageData["gsid"] = g
			}
			if s, _ := auth["sut"].(string); s != "" {
				flow.PageData["sut"] = s
			}
		}
		// 广告间隔 preAdInterval/lastAdInterval
		if hasKey(timeline, "preMarkInterval") || hasKey(timeline, "lastMarkInterval") {
			num := len(flow.Items)
			var idx []int
			for i, itemAny := range flow.Items {
				item, _ := itemAny.(map[string]any)
				if item["category"] == "group" {
					idx = append(idx, i)
					continue
				}
				if data, ok := item["data"].(map[string]any); ok {
					if mt, ok2 := getInt(data, "mblogtype"); ok2 && mt == 1 {
						idx = append(idx, i)
					}
				}
			}
			preAd := -1
			if len(idx) > 0 {
				preAd = idx[0]
			}
			if flow.RefreshInfo == nil {
				flow.RefreshInfo = map[string]any{}
			}
			params, _ := flow.RefreshInfo["params"].(map[string]any)
			if params == nil {
				params = map[string]any{}
			}
			params["preAdInterval"] = preAd
			flow.RefreshInfo["params"] = params
			if len(flow.MoreInfo) > 0 {
				lastAd := -1
				if len(idx) > 0 {
					lastAd = num - 1 - idx[len(idx)-1]
				}
				moreParams, _ := flow.MoreInfo["params"].(map[string]any)
				if moreParams == nil {
					moreParams = map[string]any{}
				}
				moreParams["lastAdInterval"] = lastAd
				flow.MoreInfo["params"] = moreParams
			}
		}
	} else {
		flow.Config = map[string]any{
			"action": "replace",
			"offset": 0,
			"target": "cell",
		}
	}

	return flow
}

func (c *Container) dealStatuses(timeline map[string]any) []any {
	items := make([]any, 0)

	// groupInfo（taskType != loadMore 时才输出）
	if gi, ok := timeline["groupInfo"].(map[string]any); ok && len(gi) > 0 {
		if getString(timeline, "taskType") != "loadMore" {
			items = append(items, map[string]any{
				"category": "feedBiz",
				"type":     "groupInfo",
				"data":     gi,
			})
		}
	}

	// special_follow_push
	if sf, ok := timeline["special_follow_push"].(map[string]any); ok && len(sf) > 0 {
		items = append(items, map[string]any{
			"category": "feedBiz",
			"type":     "specialFollowPush",
			"data":     sf,
		})
	}

	raw, ok := timeline["statuses"].([]any)
	if !ok || len(raw) == 0 {
		return items
	}

	lastTime := int64(0)
	if ts, ok := getInt(timeline, "lastItemTime"); ok {
		lastTime = int64(ts)
	}
	isBigday := getIntOrDefault(timeline, "is_bigday_info", 0)
	isMiniFlow := getIntOrDefault(timeline, "isMiniFlow", 0) == 1

	for _, v := range raw {
		weibo, ok := v.(map[string]any)
		if !ok {
			// 非 map[string]any 类型（如 *card.ItemGroup），直接保留原样
			items = append(items, v)
			continue
		}

		// new_come 标记
		createdAt := getString(weibo, "created_at")
		if createdAt != "" {
			if t, err := time.Parse(time.RubyDate, createdAt); err == nil {
				if getString(weibo, "mark") == "" && t.Unix() > lastTime && !getBool(weibo, "is_disable_highlight") && getString(timeline, "taskType") != "loadMore" {
					weibo["new_come"] = 1
				}
			}
		}

		// is_bigday 样式
		if isBigday == 1 && getString(weibo, "item_category") == "" {
			weibo["style"] = map[string]any{
				"divider_size":       10,
				"divider_color":      "FFEEEEEE",
				"divider_color_dark": "FF151515",
			}
		}

		var weiboData map[string]any

		switch {
		case getBool(weibo, "hot_feed_positive_feedback"):
			weiboData = c.createWeiboFeedbackItem(weibo)
		case getString(weibo, "item_category") == "group":
			weiboData = c.wrapGroupStatuses(weibo)
		case getString(weibo, "search_debug_info") != "":
			weiboData = c.createWeiboDebugItem(weibo)
		case getString(weibo, "item_category") == "insert_item":
			weiboData = weibo
		default:
			weiboData = map[string]any{
				"category": "feed",
				"data":     weibo,
			}
		}

		// markInterval：mblogtype==1 且版本 >=C42（此处直接按 mblogtype 判断）
		if mt, ok := getInt(weibo, "mblogtype"); ok && mt == 1 {
			itemExt, _ := weiboData["itemExt"].(map[string]any)
			if itemExt == nil {
				itemExt = map[string]any{}
			}
			itemExt["markInterval"] = 1
			weiboData["itemExt"] = itemExt
		}

		if isMiniFlow {
			weiboData["type"] = "mini"
		}

		items = append(items, weiboData)
	}

	// insert_struct 插入
	if ins, ok := timeline["insert_struct"].(map[string]any); ok && len(ins) > 0 {
		insertData := c.getInsert(ins, timeline)
		// 尝试找到插入位置
		pos := -1
		if sts, ok := timeline["statuses"].([]any); ok {
			targetID := getString(ins, "insert_before_id")
			for i, v := range sts {
				if w, ok := v.(map[string]any); ok && getString(w, "id") == targetID {
					pos = i
					break
				}
			}
		}
		if pos >= 0 && pos <= len(items) {
			// 插入到 pos 位置
			newItems := append([]any{}, items[:pos]...)
			newItems = append(newItems, insertData)
			newItems = append(newItems, items[pos:]...)
			items = newItems
		} else {
			items = append(items, insertData)
		}
	}

	return items
}

func (c *Container) dealPageData(ctx *context.Context, timeline map[string]any, infoContentType string) map[string]any {
	if infoContentType == "" {
		infoContentType = "cardlistInfo"
	}
	if pd, ok := tryMap(timeline, infoContentType, "pageData"); ok {
		return pd
	}
	flowId := ctx.DefaultRequestParam("flowId", "")
	if flowId == "" {
		if cl, ok := timeline["cardlistInfo"].(map[string]any); ok {
			flowId = getString(cl, "containerid")
		}
	}
	return map[string]any{
		"pageDataType": getStringDefault(timeline, "pageDataType", "feedStream"),
		"flowId":       flowId,
		"title":        "feedflow",
		"style": map[string]any{
			"flowType": "flow",
			"padding":  []any{0, 0, 0, 0},
		},
		"apiPath":        strings.TrimLeft(ctx.Path(), "/"),
		"is_first_level": 1,
	}
}

func (c *Container) dealLoadedInfo(timeline map[string]any, infoContentType string) map[string]any {
	if infoContentType == "" {
		infoContentType = "cardlistInfo"
	}
	// feedMiniFlowInfo 配置
	data := map[string]any{
		"feedMiniFlowInfo": map[string]any{
			"switch_title_mini":          "        ",
			"switch_btn_title_mini":      "精简模式",
			"switch_title_classical":     "        ",
			"switch_btn_title_classical": "标准模式",
		},
	}
	// bottomTab.showPageUp - 始终设置，即使为 false
	showPageUp := getBool(timeline, "show_pageup_bubble")
	data["bottomTab"] = map[string]any{"showPageUp": showPageUp}
	// topToast
	if getString(timeline, "remind_text_click_jump_scheme") == "" && getString(timeline, "taskType") != "loadMore" {
		remind := getString(timeline, "remind_text")
		if remind == "" {
			if cl, ok := timeline[infoContentType].(map[string]any); ok {
				remind = getString(cl, "remind_text")
			}
		}
		if remind != "" {
			tt := map[string]any{
				"content":         remind,
				"playAudio":       true,
				"height":          32,
				"duration":        2000,
				"showType":        0,
				"backgroundColor": "#F4AD54",
				"textColor":       "#FFFFFF",
			}
			if ul, ok := timeline["user_list"].([]any); ok {
				tt["user_list"] = ul
			}
			data["topToast"] = tt
		}
	}
	// isMiniFlow 透传
	if imf := getIntOrDefault(timeline, "isMiniFlow", 0); imf == 1 {
		fm, _ := data["feedMiniFlowInfo"].(map[string]any)
		fm["isMiniFlow"] = 1
		data["feedMiniFlowInfo"] = fm
	}
	// empty 清空
	if sts, ok := timeline["statuses"].([]any); ok && len(sts) == 0 && getIntOrDefault(timeline, "no_cache", 0) == 1 {
		data["empty"] = map[string]any{"clearOnEmpty": 1}
	}
	if mapi, ok := timeline["mapi"].(map[string]any); ok {
		if getIntOrDefault(mapi, "no_cache", 0) == 1 {
			data["empty"] = map[string]any{"clearOnEmpty": 1}
		}
	}
	// mapKey 透传 gsid/sut/headers/top_bubble/follow_guide_info/presentAnimation/serviceMap/popup_scheme/push_callback/autorefresh_interval/disable_polling/polling_url/can_show_refresh_guide/searchBarStyleInfo
	mapKeys := []string{"gsid", "sut", "is_push_invisible_tip", "headers", "top_bubble", "follow_guide_info", "presentAnimation", "serviceMap", "popup_scheme", "push_callback", "autorefresh_interval", "disable_polling", "polling_url", "can_show_refresh_guide", "searchBarStyleInfo"}
	for _, k := range mapKeys {
		if val, ok := timeline[k]; ok {
			if k == "gsid" || k == "sut" {
				auth, _ := data["authInfo"].(map[string]any)
				if auth == nil {
					auth = map[string]any{}
				}
				auth[k] = val
				data["authInfo"] = auth
			} else if k == "push_callback" {
				sm, _ := data["serviceMap"].(map[string]any)
				if sm == nil {
					sm = map[string]any{}
				}
				sm[k] = val
				data["serviceMap"] = sm
			} else {
				data[k] = val
			}
			continue
		}
		if cl, ok := timeline[infoContentType].(map[string]any); ok {
			if val, ok2 := cl[k]; ok2 {
				data[k] = val
			}
		}
	}
	// shareData
	if cl, ok := timeline[infoContentType].(map[string]any); ok {
		if getIntOrDefault(cl, "can_shared", 1) == 1 {
			sdKeys := []string{"share_containerid", "object_id", "cardlist_title", "containerid", "more_btns"}
			sd := map[string]any{}
			if sc, ok := cl["share_content"].(map[string]any); ok {
				sd["desc"] = getString(sc, "description")
				sd["url"] = getString(sc, "custom_share_path")
				sd["icon"] = getString(sc, "icon")
			}
			for _, k := range sdKeys {
				if v, ok := cl[k]; ok {
					if k == "cardlist_title" {
						sd["title"] = v
					} else {
						sd[k] = v
					}
				}
			}
			if len(sd) > 0 {
				sm, _ := data["serviceMap"].(map[string]any)
				if sm == nil {
					sm = map[string]any{}
				}
				sm["shareData"] = sd
				data["serviceMap"] = sm
			}
		}
	}
	return data
}

func (c *Container) dealRefreshInfo(timeline map[string]any, infoContentType string) map[string]any {
	if infoContentType == "" {
		infoContentType = "cardlistInfo"
	}
	newTime := int64(getIntOrDefault(timeline, "lastItemTime", 0))
	sinceID := ""
	count := 25
	if sts, ok := timeline["statuses"].([]any); ok && len(sts) > 0 {
		for _, v := range sts {
			if weibo, ok := v.(map[string]any); ok {
				if ct := getString(weibo, "created_at"); ct != "" {
					if t, err := time.Parse(time.RubyDate, ct); err == nil {
						if getString(weibo, "mark") == "" && t.Unix() > newTime {
							newTime = t.Unix()
						}
					}
				}
			}
		}
		if sinceIDVal := utils.ArrayValue(timeline, "since_id", ""); sinceIDVal != "" {
			sinceID = sinceIDVal.(string)
		}
	} else {
		count = 20
	}
	data := map[string]any{
		"pagingType":  "page",
		"refreshtype": "clear",
		"params": map[string]any{
			"count":        count,
			"lastItemTime": newTime,
		},
	}
	if sinceID != "" {
		params, _ := data["params"].(map[string]any)
		params["since_id"] = sinceID
		data["params"] = params
	}
	if rp, ok := timeline["refreshParams"].(map[string]any); ok {
		params, _ := data["params"].(map[string]any)
		for k, v := range rp {
			params[k] = v
		}
		data["params"] = params
	}
	return data
}

func (c *Container) dealMoreInfo(timeline map[string]any, infoContentType string) map[string]any {
	if infoContentType == "" {
		infoContentType = "cardlistInfo"
	}
	data := map[string]any{
		"pagingType": "cursor",
		"params": map[string]any{
			"max_id":        getIntOrDefault(timeline, "max_id", 0),
			"page":          getIntOrDefault(timeline, "page", 1) + 1,
			"count":         "25",
			"hot_feed_push": getIntOrDefault(timeline, "hot_feed_push", 0),
		},
		"error":    "网络加载失败，请稍后尝试",
		"loading":  "正在加载更多",
		"moreType": "default",
	}
	if _, ok := timeline["statuses"]; ok {
		delete(data["params"].(map[string]any), "page")
	}
	if cl, ok := timeline[infoContentType].(map[string]any); ok {
		if ps, ok2 := getInt(cl, "pagesize"); ok2 {
			data["params"].(map[string]any)["count"] = ps
		}
		if sid, ok2 := getInt(cl, "since_id"); ok2 && sid != 0 {
			data["params"].(map[string]any)["since_id"] = sid
			delete(data["params"].(map[string]any), "page")
		}
		if getIntOrDefault(cl, "is_interrupt", 0) == 1 {
			data["noMore"] = true
			delete(data, "error")
		}
		if et, ok2 := cl["ext_trans"].(map[string]any); ok2 {
			for k, v := range et {
				data["params"].(map[string]any)[k] = v
			}
		}
		if al, ok2 := getInt(cl, "autoLoadMoreIndex"); ok2 {
			data["autoLoadMoreIndex"] = al
		}
	}
	if mp, ok := timeline["moreParams"].(map[string]any); ok {
		for k, v := range mp {
			data["params"].(map[string]any)[k] = v
		}
	}
	if phri := utils.ArrayValue(timeline, "pre_hot_request_id", nil); phri != nil {
		data["params"].(map[string]any)["pre_hot_request_id"] = phri
	}
	return data
}

func (c *Container) dealConfig() map[string]any {
	return map[string]any{
		"paging": map[string]any{
			"threshold": 10,
		},
	}
}

func (c *Container) dealExt(timeline map[string]any, infoContentType string) map[string]any {
	if infoContentType == "" {
		infoContentType = "cardlistInfo"
	}
	data := map[string]any{}
	if getString(timeline, "commonPage") == "" {
		keys := []string{"header", "footer", "channelInfo", "title_top", "cardlist_menus", "settings", "navigationBar", "headerBack", "globalChannelConfig"}
		if cl, ok := timeline[infoContentType].(map[string]any); ok {
			for _, k := range keys {
				if v, ok2 := cl[k]; ok2 {
					if k == "title_top" || k == "cardlist_menus" || k == "settings" {
						nav, _ := data["navInfo"].(map[string]any)
						if nav == nil {
							nav = map[string]any{}
						}
						nav[k] = v
						data["navInfo"] = nav
					} else {
						data[k] = v
					}
				}
			}
		}
		return data
	}
	// commonPage 结构
	if cl, ok := timeline[infoContentType].(map[string]any); ok {
		data["pageTitle"] = map[string]any{
			"type": "default",
			"data": map[string]any{
				"title": getStringDefault(cl, "title_top", getString(cl, "cardlist_title")),
				"right_button": map[string]any{
					"text":          "",
					"icon":          "https://h5.sinaimg.cn/upload/1059/799/2022/03/30/share_pic.png",
					"iconHighlight": "https://h5.sinaimg.cn/upload/1059/799/2022/03/30/share_pic_hightlight.png",
					"scheme":        "",
				},
			},
		}
		if menus, ok2 := cl["cardlist_menus"].([]any); ok2 {
			data["pageTitle"].(map[string]any)["data"].(map[string]any)["right_button"].(map[string]any)["cardlist_menus"] = menus
		}
		data["pageFooter"] = map[string]any{
			"type":  "default",
			"style": map[string]any{"type": "flow"},
		}
		if tb, ok2 := cl["toolbar_menus"].([]any); ok2 {
			data["pageFooter"].(map[string]any)["data"] = map[string]any{"menus": tb}
		}
		if ch, ok2 := cl["cardlist_head_cards"].([]any); ok2 && len(ch) > 0 {
			if head, ok3 := ch[0].(map[string]any); ok3 {
				if clist, ok4 := head["channel_list"].([]any); ok4 && len(clist) > 0 {
					data["channelInfo"] = map[string]any{
						"channelConfig": map[string]any{"selectInfo": map[string]any{"pageDataType": "flow", "tabKey": "commonPage", "flowId": getString(timeline, "flowId")}},
						"channels":      c.toChannels(clist),
					}
				}
			}
		}
		data["pageTabBar"] = map[string]any{
			"type": "default",
			"data": map[string]any{"show_menu": 0, "menu_scheme": ""},
			"style": map[string]any{
				"height":              34,
				"backgroundColor":     "",
				"backgroundDarkColor": "",
				"isDivideEqually":     false,
				"gravity":             "center",
			},
		}
		data["pageRefresh"] = c.pageRefresh()
	}
	return data
}

func (c *Container) pageRefresh() map[string]any {
	return map[string]any{
		"type":        "default",
		"refreshType": "flow",
		"style":       map[string]any{"type": "default"},
	}
}

// toChannels 将 channelList 转换为客户端需要的 channels 结构。
func (c *Container) toChannels(channelList []any) []any {
	res := make([]any, 0, len(channelList))
	for _, v := range channelList {
		cMap, ok := v.(map[string]any)
		if !ok {
			continue
		}
		channel := map[string]any{
			"pageDataType": "flow",
			"containerid":  getString(cMap, "containerid"),
			"flowId":       getString(cMap, "containerid"),
			"apiPath":      "/flowlist",
			"title":        getString(cMap, "name"),
			"titleInfo": map[string]any{
				"icon":     "",
				"iconDark": "",
				"subTitle": "",
				"style": map[string]any{
					"font":                     "",
					"selectFont":               "Bold",
					"fontSize":                 16,
					"selectFontSize":           16,
					"textColor":                "#939393",
					"textDarkColor":            "#D3D3D3",
					"selectTextColor":          "#333333",
					"selectTextDarkColor":      "#AAAAAA",
					"backgroundIcon":           "",
					"backgroundDarkIcon":       "",
					"selectBackgroundIcon":     "",
					"selectBackgroundDarkIcon": "",
					"iconWidth":                24,
					"iconHeight":               24,
					"padding":                  []any{0, 10, 10, 0},
					"sliderColor":              []any{"#FFA300", "#FF6A00"},
					"sliderDarkColor":          []any{"#FF6A00", "#FF6A00"},
				},
			},
		}
		res = append(res, channel)
	}
	return res
}

// wrapHotTagCard 热门流顶部插入模块。
func (c *Container) wrapHotTagCard(weibo map[string]any) []any {
	// 处理 insert_items，支持 []any、[]map[string]any 或 []*card.ItemCard
	var result []any
	if items, ok := weibo["insert_items"]; ok {
		switch v := items.(type) {
		case []any:
			result = v
		case []map[string]any:
			result = make([]any, len(v))
			for i, item := range v {
				result[i] = item
			}
		default:
			// 处理其他类型（如 []*card.ItemCard），使用工具函数转换
			result = utils.ConvertToAnySlice(items)
		}
		delete(weibo, "insert_items")
	}
	return result
}

// wrapGroupStatuses group 类型的组合样式。
func (c *Container) wrapGroupStatuses(weibo map[string]any) map[string]any {
	var items []any
	hotTag := c.wrapHotTagCard(weibo)
	debugInfo := c.createWeiboDebugCardInfo(weibo, []any{})
	if len(hotTag) > 0 {
		items = append(items, hotTag...)
	}
	if debugInfo != nil {
		weibo["style"] = map[string]any{"line_color": "00000000", "divider_color": "FFFFFF", "hide_divide_line": true, "divider_size": 0}
	}
	items = append(items, map[string]any{
		"category": "feed",
		"data":     weibo,
	})
	if debugInfo != nil {
		items = append(items, debugInfo)
	}
	return map[string]any{
		"category": "group",
		"items":    items,
		"type":     "vertical",
	}
}

// createWeiboDebugCardInfo 附带 debug 信息的卡片。
func (c *Container) createWeiboDebugCardInfo(weibo map[string]any, margin []any) map[string]any {
	if getString(weibo, "search_debug_info") == "" {
		return nil
	}
	debugInfoArr := strings.Split(getString(weibo, "search_debug_info"), "|")
	tags := make([]any, 0, len(debugInfoArr))
	for _, info := range debugInfoArr {
		tags = append(tags, map[string]any{
			"type":  0,
			"title": info,
			"style": map[string]any{
				"font_size":          14,
				"height":             30,
				"left_right_padding": 12,
				"font_color":         "#3E3E3E",
				"background_color":   "#F7F7F7",
			},
		})
	}
	if len(margin) == 0 {
		margin = []any{0, 0, 0, 10}
	}
	return map[string]any{
		"category": "card",
		"style": map[string]any{
			"background": map[string]any{"type": "color", "color": "#ffffff", "colorKey": "CommonCardBackground"},
			"margin":     margin,
			"padding":    []any{0, 0, 0, 0},
		},
		"data": map[string]any{
			"card_type": 170,
			"tag_info": map[string]any{
				"max_lines": 20,
				"itemSpace": 10,
				"lineSpace": 8,
				"tags":      tags,
			},
		},
	}
}

// createWeiboFeedbackItem 热点流正反馈数据拼装。
func (c *Container) createWeiboFeedbackItem(weibo map[string]any) map[string]any {
	var items []any
	debugData := c.createWeiboDebugCardInfo(weibo, []any{0, 0, 0, 0})
	fbClick, _ := weibo["hot_feed_positive_feedback_click"].([]any)
	delete(weibo, "hot_feed_positive_feedback_click")
	delete(weibo, "hot_feed_positive_feedback")
	items = append(items, map[string]any{
		"category": "feed",
		"data":     weibo,
		"itemExt":  map[string]any{"anchorId": getString(weibo, "id")},
	})
	if debugData != nil {
		items = append(items, debugData)
	}
	if v, ok := weibo["hot_feed_trend_entrance"]; ok {
		items = append(items, v)
		delete(weibo, "hot_feed_trend_entrance")
	}
	items = append(items, map[string]any{
		"category": "cell",
		"type":     "span",
		"span": map[string]any{
			"style": map[string]any{
				"background": map[string]any{"type": "color", "color": "#eeeeee", "colorKey": "MainFeedBackgroundColor"},
			},
		},
		"style": map[string]any{
			"height": 10,
			"background": map[string]any{
				"type":     "color",
				"colorKey": "CommonCardBackground",
				"color":    "#ffffff",
			},
		},
	})
	header := map[string]any{
		"category": "cell",
		"type":     "title",
		"style": map[string]any{
			"height":     35,
			"padding":    []any{12, 0, 10, 0},
			"background": map[string]any{"type": "color", "color": "#FFFFFF", "colorKey": "CommonCardBackground"},
		},
		"title": map[string]any{
			"style": map[string]any{
				"textSize":     14,
				"textColor":    "#676767",
				"textColorKey": "CommonGray93",
				"ellipsize":    "end",
				"maxLines":     1,
			},
			"content": "热点相关",
			"type":    "text",
		},
		"content": map[string]any{
			"style": map[string]any{
				"textSize":     12,
				"textColor":    "#507DAF",
				"textColorKey": "CommonGray93",
				"ellipsize":    "end",
				"maxLines":     1,
			},
			"type":    "text",
			"content": "",
		},
		"dot": map[string]any{
			"type":         "dot",
			"unreadKey":    "",
			"dotSource":    "unread",
			"active":       false,
			"activeType":   "icon",
			"inactiveType": "icon",
			"dotText":      "",
			"iconUrl":      "http://h5.sinaimg.cn/upload/2016/09/27/555/timeline_icon_delete.png",
			"style": map[string]any{
				"dotColor":     "#F43530",
				"dotColorKey":  "UserProfileCardManagementColor",
				"textSize":     10,
				"textColor":    "#FFFFFF",
				"textColorKey": "CommonButtonText",
				"width":        15,
				"height":       15,
			},
			"click": map[string]any{
				"type":   "dialog",
				"title":  "",
				"cancel": "取消",
				"items":  fbClick,
			},
		},
	}
	return map[string]any{
		"category": "group",
		"type":     "vertical",
		"items":    items,
		"style": map[string]any{
			"margin":  []any{0, 0, 0, 0},
			"padding": []any{0, 0, 0, 0},
		},
		"header": header,
	}
}

// createWeiboDebugItem 附带 debug 信息的博文样式。
func (c *Container) createWeiboDebugItem(weibo map[string]any) map[string]any {
	delete(weibo, "style")
	items := []any{
		map[string]any{
			"category": "feed",
			"data":     weibo,
		},
	}
	debugData := c.createWeiboDebugCardInfo(weibo, []any{0, 0, 0, 0})
	if debugData != nil {
		items = append(items, debugData)
	}
	items = append(items, map[string]any{
		"category": "cell",
		"type":     "span",
		"style": map[string]any{
			"height": 10,
			"background": map[string]any{
				"type":     "color",
				"colorKey": "CommonBackground",
				"color":    "#f0f0f0",
			},
		},
	})
	return map[string]any{
		"category": "group",
		"type":     "vertical",
		"items":    items,
		"style": map[string]any{
			"margin":  []any{0, 0, 0, 0},
			"padding": []any{0, 0, 0, 0},
			"background": map[string]any{
				"type":     "color",
				"color":    "#FFFFFF",
				"colorKey": "CommonCardBackground",
			},
		},
	}
}

// getInsert 根据 insert_struct 构造插入数据。
func (c *Container) getInsert(insert map[string]any, timeline map[string]any) map[string]any {
	// insert_struct 的具体格式依赖于下游返回，此处保持原样，只补充必要字段。
	return insert
}

// dealConfig 已实现；其他函数已完成。

// 辅助函数
func getString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok2 := v.(string); ok2 {
			return s
		}
	}
	return ""
}

func getStringDefault(m map[string]any, key, def string) string {
	if s := getString(m, key); s != "" {
		return s
	}
	return def
}

func getInt(m map[string]any, key string) (int, bool) {
	if v, ok := m[key]; ok {
		switch t := v.(type) {
		case int:
			return t, true
		case int64:
			return int(t), true
		case float64:
			return int(t), true
		case string:
			if t == "" {
				return 0, false
			}
			if n, err := strconv.Atoi(t); err == nil {
				return n, true
			}
		}
	}
	return 0, false
}

func getIntOrDefault(m map[string]any, key string, def int) int {
	if v, ok := getInt(m, key); ok {
		return v
	}
	return def
}

func getBool(m map[string]any, key string) bool {
	if v, ok := m[key]; ok {
		switch t := v.(type) {
		case bool:
			return t
		case int:
			return t != 0
		case int64:
			return t != 0
		case float64:
			return t != 0
		case string:
			return t == "1" || strings.ToLower(t) == "true"
		}
	}
	return false
}

func hasKey(m map[string]any, key string) bool {
	_, ok := m[key]
	return ok
}

// tryMap 从 map[string]any 中按 key 提取子 map，若存在则返回。
func tryMap(m map[string]any, key string, subKey string) (map[string]any, bool) {
	if mm, ok := m[key].(map[string]any); ok {
		if subKey == "" {
			return mm, true
		}
		if vv, ok2 := mm[subKey].(map[string]any); ok2 {
			return vv, true
		}
	}
	return nil, false
}
