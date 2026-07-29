package cardlist

import (
	"strconv"
	"strings"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/log"
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
	"git.intra.weibo.com/search_fe/wbutil-go/view/container"
	"git.intra.weibo.com/search_fe/wbutil-go/view/template"

	"github.com/bytedance/sonic"
)

// Handle cardlist 业务返回结果处理
func Handle(ctx *context.Context, response map[string]any) (map[string]any, error) {
	isContainer := IsContainer(ctx)
	if isContainer {
		ctx.SetRequestParam("isFlowRequest", "1")
	}

	response = ProcessLog(ctx, response)

	containerId := ctx.DefaultRequestParam("containerid", "")
	cardlistInfo := getMapValue(response, "cardlistInfo")
	if cardlistInfo == nil {
		cardlistInfo = make(map[string]any)
		response["cardlistInfo"] = cardlistInfo
	}

	if _, hasMenus := cardlistInfo["cardlist_menus"]; !hasMenus {
		cardlistInfo["cardlist_menus"] = AddMenus()
		if cid := getStringValue(cardlistInfo, "containerid"); cid == "" {
			cardlistInfo["containerid"] = containerId
		}
	}

	response = FilterHandle(response)

	if isContainer {
		response = container.Handle(ctx, response)
		response = template.CardlistContainerRender(ctx, response)
	} else {
		response = template.CardlistRender(ctx, response)
	}

	return response, nil
}

// RequestInit 请求阶段初始化
func RequestInit(ctx *context.Context) {
	containerId := ctx.DefaultRequestParam("containerid", "")
	ctx.SetRequestParam("fid", containerId)

	partnerId := ""
	if len(containerId) >= 6 {
		partnerId = containerId[:6]
	}

	// 老版发现页兼容逻辑
	if partnerId == "102803" {
		extparam := ctx.DefaultRequestParam("extparam", "")
		ext := ctx.DefaultRequestParam("ext", "")
		if extparam == "discover" || ext == "discover" {
			ctx.SetRequestParam("need_head_cards", "0")
		}
	}
}

// IsContainer 判断本次请求客户端是否开启容器化开关
func IsContainer(ctx *context.Context) bool {
	isContainer := ctx.DefaultRequestParam("is_container", "")
	return isContainer != "" && isContainer != "0"
}

// AddMenus 处理菜单
func AddMenus() []map[string]any {
	return []map[string]any{
		{"type": "button_menus_refresh", "name": "刷新"},
		{"type": "gohome", "name": "返回首页", "params": map[string]any{"scheme": "sinaweibo://gotohome"}},
	}
}

// logData 日志数据收集结构
type logData struct {
	mids             []string
	weibos           []map[string]any
	avatar           []string
	scheduleLiveCard []int
	hotspot          []int
	redpacket        []string
	searchCate       []string
	jingxi           []int
	jingxiAdid       []string
	keywords         []int
	keywordsAdid     []string
	highLights       []string
}

func newLogData() *logData {
	return &logData{
		mids:             make([]string, 0, 32),
		weibos:           make([]map[string]any, 0, 32),
		avatar:           make([]string, 0, 32),
		scheduleLiveCard: make([]int, 0, 32),
		hotspot:          make([]int, 0, 32),
		redpacket:        make([]string, 0, 32),
		searchCate:       make([]string, 0, 32),
		jingxi:           make([]int, 0, 32),
		jingxiAdid:       make([]string, 0, 32),
		keywords:         make([]int, 0, 32),
		keywordsAdid:     make([]string, 0, 32),
		highLights:       make([]string, 0, 32),
	}
}

// ProcessLog 处理日志
func ProcessLog(ctx *context.Context, result map[string]any) map[string]any {
	if !log.IsRecord(ctx) {
		return result
	}

	cardlistInfo := getMapValue(result, "cardlistInfo")
	if cardlistInfo != nil {
		if mapi := getMapValue(cardlistInfo, "mapi"); mapi != nil {
			if notRecordLog, ok := mapi["not_record_log"].(bool); ok && notRecordLog {
				return result
			}
		}
	}

	containerId := ctx.DefaultRequestParam("containerid", "")
	partnerId := ""
	if len(containerId) >= 6 {
		partnerId = containerId[:6]
	}

	result = log.TopnClick(ctx, result)
	cards, _ := result["cards"].([]any)

	if len(cards) > 0 {
		ld := collectLogData(cards)
		logSuffixes := buildLogSuffixes(ctx, containerId, partnerId, cardlistInfo, ld)

		reqType := ctx.DefaultRequestParam("type", "")
		aiTabNativeEnable := ctx.DefaultRequestParam("ai_tab_native_enable", "")
		if reqType != "200" || (reqType == "200" && aiTabNativeEnable == "2") {
			log.FeedExpo(ctx, 900, ld.mids, logSuffixes, "")
		}

		log.CardExp(ctx, cards, 900, false)

		partnerIdInt, _ := strconv.Atoi(partnerId)
		log.RepostAdExpo(ctx, convertWeibosToAny(ld.weibos), partnerIdInt)

		logActionCode(ctx, result, ld.mids)
	} else {
		logActionCode(ctx, result, nil)
	}

	pageAttr := getStringValue(cardlistInfo, "page_attr")
	objectType := getStringValue(cardlistInfo, "object_type")
	log.PageVisit(ctx, "get_index", pageAttr, objectType, nil)

	return result
}

// collectLogData 收集日志数据
func collectLogData(cards []any) *logData {
	ld := newLogData()

	for _, cardItem := range cards {
		card := utils.StructToMap(cardItem)
		if card == nil || len(card) == 0 {
			continue
		}

		// 处理 statuses
		if statuses, ok := card["statuses"].([]any); ok {
			for _, statusItem := range statuses {
				if mblog, ok := statusItem.(map[string]any); ok {
					if mid := getStringValue(mblog, "mid"); mid != "" {
						ld.mids = append(ld.mids, mid)
						ld.searchCate = append(ld.searchCate, getStringValue(mblog, "cate_id"))
						ld.weibos = append(ld.weibos, mblog)
					}
				}
			}
		}

		// 处理 card_group 或 pics
		groups := card["card_group"]
		if groups == nil {
			groups = card["pics"]
		}

		if groupsSlice, ok := groups.([]any); ok {
			for _, groupItem := range groupsSlice {
				if group, ok := groupItem.(map[string]any); ok {
					collectMblogData(ld, group, "mblog")
					collectElementMblog(ld, group, "left_element.mblog")
					collectElementMblog(ld, group, "right_element.mblog")
				}
			}
		}

		// 处理 card 级别的 left/right element
		collectElementMblog(ld, card, "left_element.mblog")
		collectElementMblog(ld, card, "right_element.mblog")

		// 处理 playlist.statuses
		if playlist := getMapValue(card, "playlist"); playlist != nil {
			if statuses, ok := playlist["statuses"].([]any); ok {
				for _, statusItem := range statuses {
					if mblog, ok := statusItem.(map[string]any); ok {
						if mid := getStringValue(mblog, "mid"); mid != "" {
							ld.mids = append(ld.mids, mid)
							ld.searchCate = append(ld.searchCate, getStringValue(card, "cate_id"))
							ld.weibos = append(ld.weibos, mblog)
						}
					}
				}
			}
		}

		// 处理 card.mblog
		collectMblogData(ld, card, "mblog")
	}

	return ld
}

// collectMblogData 收集 mblog 数据
func collectMblogData(ld *logData, parent map[string]any, key string) {
	mblog := getMapValue(parent, key)
	if mblog == nil {
		return
	}

	mid := getStringValue(mblog, "id")
	if mid == "" {
		mid = getStringValue(mblog, "mid")
	}
	if mid == "" {
		return
	}

	ld.mids = append(ld.mids, mid)
	ld.searchCate = append(ld.searchCate, getStringValue(parent, "cate_id"))
	ld.weibos = append(ld.weibos, mblog)

	// avatar type
	avatarType := ""
	if user := getMapValue(mblog, "user"); user != nil {
		if avatar := getMapValue(user, "avatar"); avatar != nil {
			avatarType = getStringValue(avatar, "type")
		}
	}
	ld.avatar = append(ld.avatar, avatarType)

	// schedule_live_card
	scheduleLive := 0
	if pageInfo := getMapValue(mblog, "page_info"); pageInfo != nil {
		if pageInfo["story_live_subscription_info"] != nil {
			scheduleLive = 1
		}
	}
	ld.scheduleLiveCard = append(ld.scheduleLiveCard, scheduleLive)

	// hotspot 和 redpacket
	ld.hotspot = append(ld.hotspot, getIntValue(mblog, "hotspot"))
	ld.redpacket = append(ld.redpacket, getStringValue(mblog, "redpacket"))

	// jingxi
	jingxiFlag := 0
	jingxiAdid := ""
	if extendInfo := getMapValue(mblog, "extend_info"); extendInfo != nil {
		if attitudeDynamicAd := getMapValue(extendInfo, "attitude_dynamic_ad"); attitudeDynamicAd != nil {
			if getStringValue(attitudeDynamicAd, "adType") == "jingxi" {
				jingxiFlag = 1
				jingxiAdid = getStringValue(attitudeDynamicAd, "adid")
			}
		}
	}
	ld.jingxi = append(ld.jingxi, jingxiFlag)
	ld.jingxiAdid = append(ld.jingxiAdid, jingxiAdid)

	// keywords
	keywordsFlag := 0
	keywordsAdid := ""
	if extendInfo := getMapValue(mblog, "extend_info"); extendInfo != nil {
		if bgCard := getMapValue(extendInfo, "bg_card"); bgCard != nil {
			if bgCard["keywords"] != nil {
				keywordsFlag = 1
			}
			keywordsAdid = getStringValue(bgCard, "adid")
		}
	}
	ld.keywords = append(ld.keywords, keywordsFlag)
	ld.keywordsAdid = append(ld.keywordsAdid, keywordsAdid)

	// highLights
	ld.highLights = append(ld.highLights, log.HighlightsKeywordFeedExpo(mblog))
}

// collectElementMblog 收集 element 中的 mblog 数据
func collectElementMblog(ld *logData, parent map[string]any, path string) {
	if mblogVal := utils.ArrayValue(parent, path, nil); mblogVal != nil {
		if mblog, ok := mblogVal.(map[string]any); ok {
			if mid := getStringValue(mblog, "id"); mid != "" {
				ld.mids = append(ld.mids, mid)
				ld.searchCate = append(ld.searchCate, getStringValue(parent, "cate_id"))
				ld.weibos = append(ld.weibos, mblog)
			}
		}
	}
}

// buildLogSuffixes 构建日志后缀
func buildLogSuffixes(ctx *context.Context, containerId, partnerId string, cardlistInfo map[string]any, ld *logData) []string {
	logSuffixes := []string{
		"containerid=>" + partnerId,
		"full_containerid=>" + containerId,
		"avatar_type=>" + strings.Join(ld.avatar, ";"),
		"schedule_live_card=>" + joinInts(ld.scheduleLiveCard, ";"),
		"uicode=>" + ctx.DefaultRequestParam("uicode", ""),
		"luicode=>" + ctx.DefaultRequestParam("luicode", ""),
		"from=>" + ctx.GetFrom(),
		"aid=>" + ctx.DefaultRequestParam("aid", ""),
		"wm=>" + ctx.Wm(),
		"temp_mids=>" + getStringValue(cardlistInfo, "temp_mids"),
		"lcardid=>" + ctx.DefaultRequestParam("lcardid", ""),
		"lfid=>" + ctx.DefaultRequestParam("lfid", ""),
		"scenes=>" + ctx.DefaultRequestParam("scenes", ""),
		"page=>" + ctx.DefaultRequestParam("page", ""),
		"extparam=>" + log.FeedExpoExtparam(ctx),
		"hotspot=>" + joinInts(ld.hotspot, ";"),
		"redpacket=>" + strings.Join(ld.redpacket, ";"),
		"cate_list=>" + strings.Join(ld.searchCate, ";"),
		"jingxi=>" + joinInts(ld.jingxi, ";"),
		"jingxi_adid=>" + strings.Join(ld.jingxiAdid, ";"),
		"keywords=>" + joinInts(ld.keywords, ";"),
		"keywords_adid=>" + strings.Join(ld.keywordsAdid, ";"),
		"feed_expo_ext=>" + getStringValue(cardlistInfo, "feed_expo_ext"),
		"high_lights=>" + strings.Join(ld.highLights, ";"),
		"mblog_types=>" + strings.Join(mblogTypes(ld.weibos), ";"),
		"search_ext=>" + getSearchExt(ctx),
	}

	// 处理 extension 参数
	extensionStr := ctx.DefaultRequestParam("extension", "{}")
	if extensionStr != "" && extensionStr != "{}" {
		var extension map[string]any
		if err := sonic.Unmarshal([]byte(extensionStr), &extension); err == nil {
			for k, v := range extension {
				logSuffixes = append(logSuffixes, k+"=>"+toString(v))
			}
		}
	}

	// 添加可选字段
	addOptionalSuffix := func(key, value string) {
		if value != "" {
			logSuffixes = append(logSuffixes, key+"=>"+value)
		}
	}

	addOptionalSuffix("search_request_id", getStringValue(cardlistInfo, "search_request_id"))
	addOptionalSuffix("search_ssid", getStringValue(cardlistInfo, "search_ssid"))
	addOptionalSuffix("search_vsid", getStringValue(cardlistInfo, "search_vsid"))
	addOptionalSuffix("searchbar_source", ctx.DefaultRequestParam("searchbar_source", ""))
	addOptionalSuffix("feed_cate", ctx.DefaultRequestParam("feed_cate", ""))
	addOptionalSuffix("source_mid", log.GetSourceMid(ctx))
	addOptionalSuffix("search_layout", ctx.DefaultRequestParam("search_layout", ""))
	addOptionalSuffix("search_mode", ctx.DefaultRequestParam("search_mode_info", ""))

	return logSuffixes
}

// logActionCode 行为码日志 26
func logActionCode(ctx *context.Context, result map[string]any, mids []string) {
	containerId := ctx.DefaultRequestParam("containerid", "")
	cardlistInfo := getMapValue(result, "cardlistInfo")

	actionLogParams := map[string]any{
		"page":              ctx.DefaultRequestParam("page", ""),
		"since_id":          ctx.DefaultRequestParam("since_id", ""),
		"mid":               ctx.DefaultRequestParam("mid", ""),
		"page_pro":          getStringValue(cardlistInfo, "page_attr"),
		"extparam":          log.ActionCodeExtparam(ctx),
		"last":              cardlistInfo["last"],
		"empty":             cardlistInfo["empty"],
		"orifid":            ctx.DefaultRequestParam("orifid", ""),
		"oriuicode":         ctx.DefaultRequestParam("oriuicode", ""),
		"request_referer":   ctx.DefaultRequestParam("request_referer", ""),
		"filter_label_word": ctx.DefaultRequestParam("filter_label_word", ""),
		"search_ext":        getSearchExt(ctx),
	}
	// 添加可选字段
	addOptionalParam := func(key, value string) {
		if value != "" {
			actionLogParams[key] = value
		}
	}

	// 添加可选字段
	addOptionalParamNotEmpty := func(key, value string) {
		if value != "" && value != "0" {
			actionLogParams[key] = value
		}
	}

	addOptionalParam("volume", getStringValue(cardlistInfo, "volume"))
	addOptionalParamNotEmpty("scenes", ctx.DefaultRequestParam("scenes", ""))
	addOptionalParam("srid", getStringValue(cardlistInfo, "search_request_id"))
	addOptionalParam("search_ssid", getStringValue(cardlistInfo, "search_ssid"))
	addOptionalParam("search_vsid", getStringValue(cardlistInfo, "search_vsid"))
	addOptionalParam("feed_cate", ctx.DefaultRequestParam("feed_cate", ""))
	addOptionalParam("root_mid", ctx.DefaultRequestParam("root_mid", ""))
	addOptionalParam("is_positivity", ctx.DefaultRequestParam("is_positivity", ""))
	addOptionalParam("reply_type", ctx.DefaultRequestParam("reply_type", ""))
	addOptionalParam("is_half_screen", ctx.DefaultRequestParam("is_half_screen", ""))
	addOptionalParam("q_extra", ctx.DefaultRequestParam("q_extra", ""))
	addOptionalParam("source_mid", log.GetSourceMid(ctx))
	addOptionalParam("searchbar_source", ctx.DefaultRequestParam("searchbar_source", ""))
	addOptionalParam("search_mode", ctx.DefaultRequestParam("search_mode_info", ""))
	addOptionalParam("search_layout", ctx.DefaultRequestParam("search_layout", ""))
	addOptionalParam("is_blink_flag", ctx.DefaultRequestParam("is_blink_flag", ""))

	if len(mids) > 0 {
		actionLogParams["have_mblog"] = 1
	}

	reqType := ctx.DefaultRequestParam("type", "")
	aiTabNativeEnable := ctx.DefaultRequestParam("ai_tab_native_enable", "")
	if reqType == "200" && aiTabNativeEnable != "2" {
		actionLogParams["stats_exclude"] = 1
	}

	log.Action(26, containerId, actionLogParams, ctx)
}

// FilterHandle 过滤冗余信息，避免影响客户端页面解析
func FilterHandle(response map[string]any) map[string]any {
	cards, ok := response["cards"].([]any)
	if !ok || len(cards) == 0 {
		return response
	}

	for i, cardItem := range cards {
		card, ok := cardItem.(map[string]any)
		if !ok {
			continue
		}

		if cardGroup, ok := card["card_group"].([]any); ok && len(cardGroup) > 0 {
			for j, groupItem := range cardGroup {
				if group, ok := groupItem.(map[string]any); ok {
					cardGroup[j] = filterPics(group)
				}
			}
			card["card_group"] = cardGroup
		}

		card = filterPics(card)
		cards[i] = card
	}

	response["cards"] = cards
	return response
}

// filterPics card3 等图片类卡片处理
func filterPics(card map[string]any) map[string]any {
	pics, ok := card["pics"].([]any)
	if !ok || len(pics) == 0 {
		return card
	}

	cardType := getIntValue(card, "card_type")
	readtimetype := getStringValue(card, "readtimetype")

	if cardType == 3 && readtimetype == "mblog" {
		for i, picItem := range pics {
			if pic, ok := picItem.(map[string]any); ok {
				delete(pic, "mblog")
				delete(pic, "cate_id")
				pics[i] = pic
			}
		}
		card["pics"] = pics
	}

	return card
}

// mblogTypes 微博类型合集，图片和视频混排算作视频
// 类型：0，文本，1，图片，2，视频，3，直播
func mblogTypes(weibos []map[string]any) []string {
	types := make([]string, 0, len(weibos))

	for _, weibo := range weibos {
		mblogType := 0

		if pageInfo := getMapValue(weibo, "page_info"); pageInfo != nil {
			objectType := getStringValue(pageInfo, "object_type")
			if objectType == "live" {
				mblogType = 3
			} else if objectType == "video" || objectType == "adFeedVideo" {
				mblogType = 2
			}
		}

		if mblogType < 2 {
			if mixMediaInfo := getMapValue(weibo, "mix_media_info"); mixMediaInfo != nil {
				if items, ok := mixMediaInfo["items"].([]any); ok && len(items) > 0 {
					hasVideo, hasPic := false, false
					for _, item := range items {
						if mediaInfo, ok := item.(map[string]any); ok {
							mediaType := getStringValue(mediaInfo, "type")
							if mediaType == "video" || mediaType == "adFeedVideo" {
								hasVideo = true
							}
							if mediaType == "pic" {
								hasPic = true
							}
						}
					}
					if hasVideo {
						mblogType = 2
					} else if hasPic {
						mblogType = 1
					}
				}
			}
		}

		if mblogType < 1 {
			if picIds, ok := weibo["pic_ids"].([]any); ok && len(picIds) > 0 {
				mblogType = 1
			}
		}

		types = append(types, strconv.Itoa(mblogType))
	}

	return types
}

// getSearchExt 获取搜索扩展参数
func getSearchExt(ctx *context.Context) string {
	var searchExt []string

	if searchLayout := ctx.DefaultRequestParam("search_layout", ""); searchLayout != "" {
		searchExt = append(searchExt, "search_layout="+searchLayout)
	}

	if searchPid := ctx.DefaultRequestParam("search_pid", ""); searchPid != "" {
		picSearchMid := ctx.DefaultRequestParam("pic_search_mid", "")
		searchExt = append(searchExt, "pic_search=图搜_"+picSearchMid+"_"+searchPid)
	}

	return strings.Join(searchExt, "&")
}

// 辅助函数

func convertWeibosToAny(weibos []map[string]any) []any {
	result := make([]any, len(weibos))
	for i, w := range weibos {
		result[i] = w
	}
	return result
}

func getMapValue(m map[string]any, key string) map[string]any {
	if m == nil {
		return nil
	}
	if v, ok := m[key].(map[string]any); ok {
		return v
	}
	return nil
}

func getStringValue(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getIntValue(m map[string]any, key string) int {
	if m == nil {
		return 0
	}
	switch v := m[key].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	return 0
}

func joinInts(s []int, sep string) string {
	strs := make([]string, len(s))
	for i, v := range s {
		strs[i] = strconv.Itoa(v)
	}
	return strings.Join(strs, sep)
}

func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}
