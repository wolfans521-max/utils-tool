package log

// kpi业务日志相关

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
)

const CustomMarkKey = "custom_mark"

// FeedExpo Feed曝光日志
func FeedExpo(ctx *context.Context, typ int, mids []string, suffixes []string, logCid string) {
	if ctx == nil || len(mids) == 0 {
		return
	}

	// 复制一份，避免调用方切片被修改
	logSuffixes := append([]string{}, suffixes...)

	// 聂钰需求：曝光日志加 launchid/ul_hid/ul_sid 等
	if v := ctx.DefaultRequestParam("launchid", ""); v != "" {
		logSuffixes = append(logSuffixes, "launchid=>"+v)
	}
	if v := ctx.DefaultRequestParam("ul_hid", ""); v != "" {
		logSuffixes = append(logSuffixes, "ul_hid=>"+v)
	}
	if v := ctx.DefaultRequestParam("ul_sid", ""); v != "" {
		logSuffixes = append(logSuffixes, "ul_sid=>"+v)
	}

	if _, ok := ctx.GetRequestParam("f_s"); ok {
		logSuffixes = append(logSuffixes, "f_s=>"+ctx.RequestParam("f_s"))
	}
	if _, ok := ctx.GetRequestParam("a_e"); ok {
		logSuffixes = append(logSuffixes, "a_e=>"+ctx.RequestParam("a_e"))
	}
	if ctx.RequestHas("is_positivity") {
		logSuffixes = append(logSuffixes, "is_positivity=>"+ctx.RequestParam("is_positivity"))
	}
	if ctx.RequestHas("reply_type") {
		logSuffixes = append(logSuffixes, "reply_type=>"+ctx.RequestParam("reply_type"))
	}
	if ctx.RequestHas("is_half_screen") {
		logSuffixes = append(logSuffixes, "is_half_screen=>"+ctx.RequestParam("is_half_screen"))
	}
	if ctx.RequestHas("q_extra") {
		logSuffixes = append(logSuffixes, "q_extra=>"+ctx.RequestParam("q_extra"))
	}
	if ctx.RequestHas("ispush") {
		logSuffixes = append(logSuffixes, "ispush=>"+ctx.RequestParam("ispush"))
	}
	if v := ctx.DefaultRequestParam("root_mid", ""); v != "" {
		logSuffixes = append(logSuffixes, "root_mid=>"+v)
	}

	logSuffixes = append(logSuffixes, "go_server=>1")

	appid := "6"
	if ctx.IsWeiboPro() {
		if v := ctx.DefaultRequestParam("appid", ""); v != "" {
			appid = v
		}
	}

	fields := []string{
		strconv.FormatInt(time.Now().Unix(), 10), // time()
		ctx.GetUserId(),
		strings.Join(mids, ","),
		appid,
		logCid,
		strconv.Itoa(typ),
		strconv.Itoa(len(mids)),
		strings.Join(logSuffixes, ","),
	}

	Write("feedexpo_mob", fields, ctx)
}

// CardExp Card曝光日志
func CardExp(ctx *context.Context, cards []any, interfaceID int, isOriStatuses bool) {
	if len(cards) == 0 {
		return
	}

	cardLog := []string{
		ctx.Request().Host,
		ctx.GetUserId(),
		"3439264077",
		strconv.Itoa(interfaceID),
		ctx.ClientIP(),
	}

	var cardLogExt map[string]any
	if isOriStatuses {
		cardLogExt = getCardLogExtByOriStatuses(cards)
	} else {
		cardLogExt = getCardLogExtByOri(cards)
	}

	cardObjectId := cardLogExt["cardObjectId"].([]string)
	cardMid := cardLogExt["cardMid"].([]string)
	cardObjectTotal := cardLogExt["cardObjectTotal"].(int)
	oasisCard := cardLogExt["oasisCard"].([]int)
	uids := cardLogExt["uids"].([]string)
	weibotype := cardLogExt["weibotype"].([]int)
	sinaCard := cardLogExt["sinaCard"].([]int)
	mediaAuthor := cardLogExt["mediaAuthor"].([]string)
	tencentCard := cardLogExt["tencentCard"].([]int)
	funiuProps := cardLogExt["funiuProps"].([]string)
	productType := cardLogExt["productType"].([]string)
	shopWindow := cardLogExt["shopWindow"].(map[string]string)
	kolCard := cardLogExt["kolCard"].([]string)
	movieButton := cardLogExt["movieButton"].([]string)

	oasisType := utils.JoinSemicolonInts(oasisCard)
	isOasisUser := 2

	sinaCards := utils.JoinSemicolonInts(sinaCard)
	funiuPropsJoined := utils.JoinSemicolonStrings(funiuProps)
	tencentCards := utils.JoinSemicolonInts(tencentCard)
	uidsJoined := utils.JoinSemicolonStrings(uids)
	weibotypeJoined := utils.JoinSemicolonInts(weibotype)
	mediaAuthorJoined := utils.JoinSemicolonStrings(mediaAuthor)
	productTypeJoined := utils.JoinSemicolonStrings(productType)

	var shopWindowsParts []string
	for k, v := range shopWindow {
		shopWindowsParts = append(shopWindowsParts, k+":"+v)
	}
	shopWindows := utils.JoinSemicolonStrings(shopWindowsParts)

	kolCardJoined := utils.JoinSemicolonStrings(kolCard)
	movieButtons := utils.JoinSemicolonStrings(movieButton)

	// 将 cardObjectId, cardMid, 和 suffixes 数组合并
	cardObjectIdsJoined := utils.JoinCommaStrings(cardObjectId)
	cardMidsJoined := utils.JoinCommaStrings(cardMid)

	cardObjectTotalStr := strconv.Itoa(cardObjectTotal)

	// 使用 strings.Builder 构建 suffixes，避免多次字符串分配
	var sb strings.Builder
	sb.Grow(512) // 预分配足够的容量
	sb.WriteString("spr=>")
	// spr参数需要URL解码
	spr := ctx.DefaultRequestParam("spr", "")
	if decoded, err := url.QueryUnescape(spr); err == nil {
		spr = decoded
	}
	sb.WriteString(spr)
	sb.WriteString(",mid_num=>")
	sb.WriteString(cardObjectTotalStr)
	sb.WriteString(",oid_num=>")
	sb.WriteString(cardObjectTotalStr)
	sb.WriteString(",oasis_type=>")
	sb.WriteString(oasisType)
	sb.WriteString(",r_o_u=>")
	sb.WriteString(strconv.Itoa(isOasisUser))
	sb.WriteString(",s_c=>")
	sb.WriteString(sinaCards)
	sb.WriteString(",funiu_props=>")
	sb.WriteString(funiuPropsJoined)
	sb.WriteString(",tencent_cards=>")
	sb.WriteString(tencentCards)
	sb.WriteString(",uids=>")
	sb.WriteString(uidsJoined)
	sb.WriteString(",weibo_type=>")
	sb.WriteString(weibotypeJoined)
	sb.WriteString(",media_author=>")
	sb.WriteString(mediaAuthorJoined)
	sb.WriteString(",product_type=>")
	sb.WriteString(productTypeJoined)
	sb.WriteString(",shop_window_scene=>")
	sb.WriteString(shopWindows)
	sb.WriteString(",kol_card=>")
	sb.WriteString(kolCardJoined)
	sb.WriteString(",movie_button=>")
	sb.WriteString(movieButtons)
	suffixesJoined := sb.String()

	allFields := append(cardLog, cardObjectIdsJoined, cardMidsJoined, suffixesJoined)

	if cardObjectTotal > 0 {
		Write("cardexp_mob", allFields, ctx)
	}
}

// getCardLogExtByOriStatuses 处理原始show_batch结果
func getCardLogExtByOriStatuses(cards []any) map[string]any {
	var (
		cardObjectId    []string
		cardMid         []string
		cardObjectTotal int
		oasisCard       []int
		uids            []string
		weibotype       []int
		sinaCard        []int
		mediaAuthor     []string
		tencentCard     []int
		movieButton     []string
		funiuProps      []string
		shopWindow      = make(map[string]string)
		kolCard         []string
		productType     []string
	)

	for _, cardAny := range cards {
		card, ok := cardAny.(map[string]any)
		if !ok {
			continue
		}

		if pageInfo, ok := card["page_info"].(map[string]any); ok {
			oid, _ := utils.ArrayValue(card, "page_info.actionlog.oid", "").(string)
			cardObjectId = append(cardObjectId, oid)
			mid, _ := card["mid"].(string)
			cardMid = append(cardMid, mid)
			uid, _ := utils.ArrayValue(card, "user.id", "").(string)
			uids = append(uids, uid)

			if _, ok := card["retweeted_status"]; ok {
				weibotype = append(weibotype, 2)
			} else {
				weibotype = append(weibotype, 1)
			}

			pageId, _ := pageInfo["page_id"].(string)
			if strings.HasPrefix(pageId, "231854") || strings.HasPrefix(pageId, "231878") {
				oasisCard = append(oasisCard, 1)
			} else {
				oasisCard = append(oasisCard, 0)
			}

			if _, ok := pageInfo["sina_multicard"]; ok {
				sinaCard = append(sinaCard, 1)
			} else {
				sinaCard = append(sinaCard, 0)
			}

			cardInfo, _ := utils.ArrayValue(card, "page_info.card_info").(map[string]any)
			lcInfos, _ := cardInfo["lc_infos"].([]any)
			if len(lcInfos) > 0 {
				text, _ := utils.ArrayValue(lcInfos[0], "text", "").(string)
				funiuProps = append(funiuProps, text)
			} else {
				funiuProps = append(funiuProps, "")
			}

			cardsArr, _ := pageInfo["cards"].([]any)
			if len(cardsArr) > 1 {
				bizType, _ := utils.ArrayValue(cardsArr[1], "biz_type", "").(string)
				if bizType == "tencent_video" {
					tencentCard = append(tencentCard, 1)
				} else {
					tencentCard = append(tencentCard, 0)
				}
			} else {
				tencentCard = append(tencentCard, 0)
			}

			if authorId, ok := pageInfo["authorid"].(string); ok {
				mediaAuthor = append(mediaAuthor, authorId)
			}

			movieTitle, _ := utils.ArrayValue(card, "page_info.movie_button_title", "").(string)
			movieButton = append(movieButton, movieTitle)

			cardObjectTotal++
		}

		if commonStruct, ok := card["common_struct"].([]any); ok && len(commonStruct) > 0 {
			firstItem := commonStruct[0].(map[string]any)
			oid, _ := utils.ArrayValue(firstItem, "actionlog.oid", "").(string)
			cardObjectId = append(cardObjectId, oid)
			mid, _ := card["mid"].(string)
			cardMid = append(cardMid, mid)
			uid, _ := utils.ArrayValue(card, "user.id", "").(string)
			uids = append(uids, uid)

			if _, ok := card["retweeted_status"]; ok {
				weibotype = append(weibotype, 2)
			} else {
				weibotype = append(weibotype, 1)
			}

			pageId, _ := utils.ArrayValue(firstItem, "page_id", "").(string)
			if strings.HasPrefix(pageId, "231854") || strings.HasPrefix(pageId, "231878") {
				oasisCard = append(oasisCard, 2) // Note: Different value than page_info branch
			} else {
				oasisCard = append(oasisCard, 0)
			}

			prodType, _ := utils.ArrayValue(firstItem, "actionlog.product_type", "").(string)
			if prodType != "" {
				productType = append(productType, prodType)
			}

			scene, _ := utils.ArrayValue(firstItem, "actionlog.shop_window_scene", "").(string)
			shopWindow[oid] = scene

			sinaCard = append(sinaCard, 0)

			cardObjectTotal++
		}

		if tagStruct, ok := card["tag_struct"].([]any); ok {
			for _, vAny := range tagStruct {
				v, ok := vAny.(map[string]any)
				if !ok {
					continue
				}
				mid, _ := card["mid"].(string)
				cardMid = append(cardMid, mid)
				uid, _ := utils.ArrayValue(card, "user.id", "").(string)
				uids = append(uids, uid)
				oid, _ := v["oid"].(string)
				cardObjectId = append(cardObjectId, oid)
			}
		}
	}

	logInfo := map[string]any{
		"cardObjectId":    cardObjectId,
		"cardMid":         cardMid,
		"cardObjectTotal": cardObjectTotal,
		"oasisCard":       oasisCard,
		"uids":            uids,
		"weibotype":       weibotype,
		"sinaCard":        sinaCard,
		"mediaAuthor":     mediaAuthor,
		"tencentCard":     tencentCard,
		"movieButton":     movieButton,
		"funiuProps":      funiuProps,
		"shopWindow":      shopWindow,
		"kolCard":         kolCard, // Not populated in ori_statuses path
		"productType":     productType,
	}

	return logInfo
}

// getCardLogExtByOri 处理旧流博文数据
func getCardLogExtByOri(cards []any) map[string]any {
	var (
		cardObjectId    []string
		cardMid         []string
		cardObjectTotal int
		oasisCard       []int
		uids            []string
		weibotype       []int
		sinaCard        []int
		mediaAuthor     []string
		tencentCard     []int
		movieButton     []string
		funiuProps      []string
		shopWindow      = make(map[string]string)
		kolCard         []string
		productType     []string
	)

	for _, cardAny := range cards {
		card, ok := cardAny.(map[string]any)
		if !ok {
			continue
		}

		// Check for 'card_group' first
		if group, ok := utils.ArrayValue(card, "card_group", []any{}).([]any); ok && len(group) > 0 {
			for _, valueAny := range group {
				value, ok := valueAny.(map[string]any)
				if !ok {
					continue
				}

				// Process mblog.page_info inside card_group
				if utils.ArrayValue(value, "mblog.page_info") != nil {
					oid, _ := utils.ArrayValue(value, "mblog.page_info.actionlog.oid", "").(string)
					cardObjectId = append(cardObjectId, oid)
					mid, _ := utils.ArrayValue(value, "mblog.mid", "").(string)
					cardMid = append(cardMid, mid)
					uid, _ := utils.ArrayValue(value, "mblog.user.id", "").(string)
					uids = append(uids, uid)

					if utils.ArrayValue(value, "mblog.retweeted_status") != nil {
						weibotype = append(weibotype, 2)
					} else {
						weibotype = append(weibotype, 1)
					}

					pageId, _ := utils.ArrayValue(value, "mblog.page_info.page_id", "").(string)
					if strings.HasPrefix(pageId, "231854") || strings.HasPrefix(pageId, "231878") {
						oasisCard = append(oasisCard, 1)
					} else {
						oasisCard = append(oasisCard, 0)
					}

					if utils.ArrayValue(value, "mblog.page_info.sina_multicard") != nil {
						sinaCard = append(sinaCard, 1)
					} else {
						sinaCard = append(sinaCard, 0)
					}

					if authorId, ok := utils.ArrayValue(value, "mblog.page_info.authorid", "").(string); ok {
						mediaAuthor = append(mediaAuthor, authorId)
					}

					cardsArr, _ := utils.ArrayValue(value, "mblog.page_info.cards", []any{}).([]any)
					if len(cardsArr) > 1 {
						bizType, _ := utils.ArrayValue(cardsArr[1], "biz_type", "").(string)
						if bizType == "tencent_video" {
							tencentCard = append(tencentCard, 1)
						} else {
							tencentCard = append(tencentCard, 0)
						}
					} else {
						tencentCard = append(tencentCard, 0)
					}

					movieTitle, _ := utils.ArrayValue(value, "mblog.page_info.movie_button_title", "").(string)
					movieButton = append(movieButton, movieTitle)

					cardObjectTotal++
				}

				// Process mblog.common_struct inside card_group
				if utils.ArrayValue(value, "mblog.common_struct") != nil {
					commonStructArr, ok := utils.ArrayValue(value, "mblog.common_struct", []any{}).([]any)
					if !ok || len(commonStructArr) == 0 {
						continue
					}
					firstItem, ok := commonStructArr[0].(map[string]any)
					if !ok {
						continue
					}

					oid, _ := utils.ArrayValue(firstItem, "actionlog.oid", "").(string)
					cardObjectId = append(cardObjectId, oid)

					mid, _ := utils.ArrayValue(value, "mblog.mid", "").(string)
					cardMid = append(cardMid, mid)
					uid, _ := utils.ArrayValue(value, "mblog.user.id", "").(string)
					uids = append(uids, uid)

					if utils.ArrayValue(value, "mblog.retweeted_status") != nil {
						weibotype = append(weibotype, 2)
					} else {
						weibotype = append(weibotype, 1)
					}

					pageId, _ := utils.ArrayValue(firstItem, "page_id", "").(string)
					if strings.HasPrefix(pageId, "231854") || strings.HasPrefix(pageId, "231878") {
						oasisCard = append(oasisCard, 2)
					} else {
						oasisCard = append(oasisCard, 0)
					}

					prodType, _ := utils.ArrayValue(firstItem, "actionlog.product_type", "").(string)
					if prodType != "" {
						productType = append(productType, prodType)
					}

					scene, _ := utils.ArrayValue(firstItem, "actionlog.shop_window_scene", "").(string)
					shopWindow[oid] = scene

					sinaCard = append(sinaCard, 0)

					cardObjectTotal++
				}

				// Process mblog.tag_struct inside card_group
				tagStruct, ok := utils.ArrayValue(value, "mblog.tag_struct", []any{}).([]any)
				if ok {
					mid, _ := utils.ArrayValue(value, "mblog.mid", "").(string)
					uid, _ := utils.ArrayValue(value, "mblog.user.id", "").(string)

					for _, vAny := range tagStruct {
						v, ok := vAny.(map[string]any)
						if !ok {
							continue
						}
						cardMid = append(cardMid, mid)
						uids = append(uids, uid)
						oid, _ := utils.ArrayValue(v, "oid", "").(string)
						cardObjectId = append(cardObjectId, oid)
					}
				}
			}
			continue // Skip other checks if card_group was processed
		}

		// Process 'page_info' (top-level card)
		if utils.ArrayValue(card, "page_info") != nil {
			oid, _ := utils.ArrayValue(card, "mblog.page_info.actionlog.oid", "").(string)
			cardObjectId = append(cardObjectId, oid)
			mid, _ := utils.ArrayValue(card, "mid", "").(string)
			cardMid = append(cardMid, mid)
			uid, _ := utils.ArrayValue(card, "user.id", "").(string)
			uids = append(uids, uid)

			if utils.ArrayValue(card, "retweeted_status") != nil {
				weibotype = append(weibotype, 2)
			} else {
				weibotype = append(weibotype, 1)
			}

			pageId, _ := utils.ArrayValue(card, "page_info.page_id", "").(string)
			if strings.HasPrefix(pageId, "231854") || strings.HasPrefix(pageId, "231878") {
				oasisCard = append(oasisCard, 1)
			} else {
				oasisCard = append(oasisCard, 0)
			}

			if utils.ArrayValue(card, "page_info.sina_multicard") != nil {
				sinaCard = append(sinaCard, 1)
			} else {
				sinaCard = append(sinaCard, 0)
			}

			cardInfo, _ := utils.ArrayValue(card, "page_info.card_info").(map[string]any)
			if cardInfo != nil {
				lcInfos, _ := cardInfo["lc_infos"].([]any)
				if len(lcInfos) > 0 {
					text, _ := utils.ArrayValue(lcInfos[0], "text", "").(string)
					funiuProps = append(funiuProps, text)
				} else {
					funiuProps = append(funiuProps, "")
				}
			} else {
				funiuProps = append(funiuProps, "")
			}

			cardsArr, _ := utils.ArrayValue(card, "page_info.cards", []any{}).([]any)
			if len(cardsArr) > 1 {
				bizType, _ := utils.ArrayValue(cardsArr[1], "biz_type", "").(string)
				if bizType == "tencent_video" {
					tencentCard = append(tencentCard, 1)
				} else {
					tencentCard = append(tencentCard, 0)
				}
			} else {
				tencentCard = append(tencentCard, 0)
			}

			if authorId, ok := utils.ArrayValue(card, "page_info.authorid", "").(string); ok {
				mediaAuthor = append(mediaAuthor, authorId)
			}

			movieTitle, _ := utils.ArrayValue(card, "page_info.movie_button_title", "").(string)
			movieButton = append(movieButton, movieTitle)

			cardObjectTotal++
		}

		// Process 'mblog.common_struct' (top-level card)
		if utils.ArrayValue(card, "mblog.common_struct") != nil {
			commonStructArr, ok := utils.ArrayValue(card, "mblog.common_struct", []any{}).([]any)
			if !ok || len(commonStructArr) == 0 {
				// Skip to next iteration if commonStruct is empty
			} else {
				firstItem, ok := commonStructArr[0].(map[string]any)
				if !ok {
					// Skip to next iteration if first item is not a map
				} else {

					oid, _ := utils.ArrayValue(firstItem, "actionlog.oid", "").(string)
					cardObjectId = append(cardObjectId, oid)

					mid, _ := utils.ArrayValue(card, "mblog.mid", "").(string)
					cardMid = append(cardMid, mid)
					uid, _ := utils.ArrayValue(card, "user.id", "").(string) // Note: Uses card['user'], not card['mblog']['user']
					uids = append(uids, uid)

					if utils.ArrayValue(card, "retweeted_status") != nil { // Note: Checks card['retweeted_status'], not card['mblog']['retweeted_status']
						weibotype = append(weibotype, 2)
					} else {
						weibotype = append(weibotype, 1)
					}

					pageId, _ := utils.ArrayValue(firstItem, "page_id", "").(string)
					if strings.HasPrefix(pageId, "231854") || strings.HasPrefix(pageId, "231878") {
						oasisCard = append(oasisCard, 2)
					} else {
						oasisCard = append(oasisCard, 0)
					}

					prodType, _ := utils.ArrayValue(firstItem, "actionlog.product_type", "").(string)
					if prodType != "" {
						productType = append(productType, prodType)
					}

					scene, _ := utils.ArrayValue(firstItem, "actionlog.shop_window_scene", "").(string)
					shopWindow[oid] = scene

					sinaCard = append(sinaCard, 0)

					cardObjectTotal++
				}
			}
		}

		// Process 'common_struct' (top-level card)
		if utils.ArrayValue(card, "common_struct") != nil {
			commonStructArr, ok := utils.ArrayValue(card, "common_struct", []any{}).([]any)
			if !ok || len(commonStructArr) == 0 {
				// Skip to next iteration if commonStruct is empty
			} else {
				firstItem, ok := commonStructArr[0].(map[string]any)
				if !ok {
					// Skip to next iteration if first item is not a map
				} else {

					oid, _ := utils.ArrayValue(firstItem, "actionlog.oid", "").(string)
					cardObjectId = append(cardObjectId, oid)

					mid, _ := utils.ArrayValue(card, "mid", "").(string)
					cardMid = append(cardMid, mid)
					uid, _ := utils.ArrayValue(card, "user.id", "").(string)
					uids = append(uids, uid)

					if utils.ArrayValue(card, "retweeted_status") != nil {
						weibotype = append(weibotype, 2)
					} else {
						weibotype = append(weibotype, 1)
					}

					pageId, _ := utils.ArrayValue(firstItem, "page_id", "").(string)
					if strings.HasPrefix(pageId, "231854") || strings.HasPrefix(pageId, "231878") {
						oasisCard = append(oasisCard, 2)
					} else {
						oasisCard = append(oasisCard, 0)
					}

					prodType, _ := utils.ArrayValue(firstItem, "actionlog.product_type", "").(string)
					if prodType != "" {
						productType = append(productType, prodType)
					}

					kolCardVal, _ := utils.ArrayValue(firstItem, "kol_card", "").(string)
					if kolCardVal != "" {
						kolCard = append(kolCard, kolCardVal)
					}

					scene, _ := utils.ArrayValue(firstItem, "actionlog.shop_window_scene", "").(string)
					shopWindow[oid] = scene

					sinaCard = append(sinaCard, 0)
					cardObjectTotal++
				}
			}
		}

		// Process 'mblog.tag_struct' (top-level card)
		tagStruct, ok := utils.ArrayValue(card, "mblog.tag_struct", []any{}).([]any)
		if ok {
			mid, _ := utils.ArrayValue(card, "mblog.mid", "").(string)
			uid, _ := utils.ArrayValue(card, "mblog.user.id", "").(string)

			for _, vAny := range tagStruct {
				v, ok := vAny.(map[string]any)
				if !ok {
					continue
				}
				cardMid = append(cardMid, mid)
				uids = append(uids, uid)
				oid, _ := utils.ArrayValue(v, "oid", "").(string)
				cardObjectId = append(cardObjectId, oid)
			}
		}

		// Process 'tag_struct' (top-level card)
		tagStruct, ok = utils.ArrayValue(card, "tag_struct", []any{}).([]any)
		if ok {
			mid, _ := utils.ArrayValue(card, "mid", "").(string)
			uid, _ := utils.ArrayValue(card, "user.id", "").(string)

			for _, vAny := range tagStruct {
				v, ok := vAny.(map[string]any)
				if !ok {
					continue
				}
				cardMid = append(cardMid, mid)
				uids = append(uids, uid)
				oid, _ := utils.ArrayValue(v, "oid", "").(string)
				cardObjectId = append(cardObjectId, oid)
			}
		}
	}

	logInfo := map[string]any{
		"cardObjectId":    cardObjectId,
		"cardMid":         cardMid,
		"cardObjectTotal": cardObjectTotal,
		"oasisCard":       oasisCard,
		"uids":            uids,
		"weibotype":       weibotype,
		"sinaCard":        sinaCard,
		"mediaAuthor":     mediaAuthor,
		"tencentCard":     tencentCard,
		"movieButton":     movieButton,
		"funiuProps":      funiuProps,
		"shopWindow":      shopWindow,
		"kolCard":         kolCard,
		"productType":     productType,
	}

	return logInfo
}

// RepostAdExpo 广告被二次转发后的曝光日志
func RepostAdExpo(ctx *context.Context, weibos []any, typ int) {
	if ctx == nil || len(weibos) == 0 {
		return
	}

	serverIP := os.Getenv("SERVER_ADDR")
	if serverIP == "" {
		serverIP = ctx.RemoteIP()
		if serverIP == "" {
			serverIP = ctx.ClientIP()
		}
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	from := ctx.GetFrom()
	network := ctx.DefaultRequestParam("networktype", "")
	viewerID := safeUserID(ctx)

	for _, w := range weibos {
		st, ok := w.(map[string]any)
		if !ok {
			continue
		}
		mark := toString(st["mark"])
		if mark == "" || strings.Count(mark, "|") != 2 {
			continue
		}
		userMap, _ := st["user"].(map[string]any)
		author := toString(userMap["id"])
		mid := toString(st["id"])

		fields := []string{
			serverIP,
			now,
			author,
			mid,
			mark,
			viewerID,
			from,
			network,
			strconv.Itoa(typ),
		}
		Write("repostadexpo_mob", fields, ctx)
	}
}

// RepostAdExpoContainer 广告被二次转发后的曝光日志,统一流结构
func RepostAdExpoContainer(ctx *context.Context, weibos []any, typ int) {
	if ctx == nil || len(weibos) == 0 {
		return
	}

	serverIP := os.Getenv("SERVER_ADDR")
	if serverIP == "" {
		serverIP = ctx.RemoteIP()
		if serverIP == "" {
			serverIP = ctx.ClientIP()
		}
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	from := ctx.GetFrom()
	network := ctx.DefaultRequestParam("networktype", "")
	viewerID := safeUserID(ctx)

	for _, w := range weibos {
		st, ok := w.(map[string]any)
		if !ok {
			continue
		}
		mark := toString(st["mark"])
		if mark == "" || strings.Count(mark, "|") != 2 {
			continue
		}
		userMap, _ := st["user"].(map[string]any)
		author := toString(userMap["id"])
		mid := toString(st["id"])

		fields := []string{
			serverIP,
			now,
			author,
			mid,
			mark,
			viewerID,
			from,
			network,
			strconv.Itoa(typ),
		}
		Write("repostadexpo_mob", fields, ctx)
	}
}

// PageVisit page的通用日志
func PageVisit(ctx *context.Context, function, pageAttr, objectType string, vp *int) {
	if ctx == nil {
		return
	}

	// page_id/containerid 的选择依赖 v_p
	vpVal := ""
	if vp != nil {
		vpVal = strconv.Itoa(*vp)
	} else {
		vpVal = ctx.DefaultRequestParam("v_p", "")
	}

	pageKey := "page_id"
	if n, err := strconv.Atoi(vpVal); err == nil && n >= 5 {
		pageKey = "containerid"
	}
	pageID := ctx.DefaultRequestParam(pageKey, "")

	vpFinal := vpVal
	if vpFinal == "" {
		vpFinal = ctx.DefaultRequestParam("v_p", "")
	}

	fields := []any{
		safeUserID(ctx),
		pageID,
		ctx.GetFrom(),
		ctx.Wm(),
		ctx.ClientIP(),
		function,
		vpFinal,
		ctx.DefaultRequestParam("card_id", ""),
		ctx.DefaultRequestParam("status", ""),
		ctx.Ua(),
		ctx.C(),
		ctx.DefaultRequestParam("sourcetype", ""),
		"page_attr" + "=>" + pageAttr,
		ctx.DefaultRequestParam("luicode", ""),
		"object_type" + "=>" + objectType,
	}

	Write("page_visit", fields, ctx)
}

// HighlightsKeywordFeedExpo 星光词信息流曝光埋点
func HighlightsKeywordFeedExpo(mblog map[string]any) string {
	if mblog == nil {
		return ""
	}
	ext, ok := mblog["extend_info"].(map[string]any)
	if !ok {
		return ""
	}
	bgCards, ok := ext["bg_cards"].([]any)
	if !ok || len(bgCards) == 0 {
		return ""
	}
	displayInfo, ok := mblog["display_info"].(map[string]any)
	if !ok || len(displayInfo) == 0 {
		return ""
	}

	highlightsAdIDs := make([]string, 0, len(bgCards))
	for _, b := range bgCards {
		m, ok := b.(map[string]any)
		if !ok {
			continue
		}
		if adid, _ := m["adid"].(string); adid != "" {
			highlightsAdIDs = append(highlightsAdIDs, adid)
		}
	}
	if len(highlightsAdIDs) == 0 {
		return ""
	}
	mid, _ := mblog["mid"].(string)
	if mid == "" {
		return ""
	}
	return mid + "_" + strings.Join(highlightsAdIDs, "$$")
}

// ActionCodeExtparam action_code行为码日志内extparam获取，这里面会检查并剔除影响日志本身的违法value值
func ActionCodeExtparam(ctx *context.Context) string {
	if ctx == nil {
		return ""
	}
	return filterActionCodeExtparam(ctx.DefaultRequestParam("extparam", ""))
}

// FeedExpoExtparam feed博文曝光日志内extparam获取，这里面会检查并剔除影响日志本身的违法value值
func FeedExpoExtparam(ctx *context.Context) string {
	if ctx == nil {
		return ""
	}

	extparam := ctx.DefaultRequestParam("extparam", "")
	if extparam == "" {
		return ""
	}
	pairs := strings.Split(extparam, "&")
	kept := make([]string, 0, len(pairs))
	for _, p := range pairs {
		if p == "" {
			continue
		}
		if strings.Contains(p, ",") || strings.Contains(p, "=>") {
			continue
		}
		kept = append(kept, p)
	}
	return strings.Join(kept, "&")
}

func GetSourceMid(ctx *context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.DefaultRequestParam("menu_source_mid", ""); v != "" {
		return v
	}
	if v := ctx.DefaultRequestParam("prove_mid", ""); v != "" {
		return v
	}
	if v := ctx.DefaultRequestParam("search_source_mid", ""); v != "" {
		return v
	}
	return ""
}

func filterActionCodeExtparam(extparam string) string {
	if extparam == "" {
		return ""
	}
	pairs := strings.Split(extparam, "&")
	var kept []string
	for _, p := range pairs {
		if p == "" {
			continue
		}
		if strings.Contains(p, "`") || strings.Contains(p, ":") || strings.Contains(p, "|") {
			continue
		}
		kept = append(kept, p)
	}
	return strings.Join(kept, "&")
}

// safeUserID 与 view/container/finder/common.go 中保持一致。
func safeUserID(ctx *context.Context) string {
	if ctx != nil && ctx.User != nil {
		return ctx.User.Id
	}
	return ""
}

// toString 是本文件内部使用的辅助函数。
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

// FinderActionLog 发现页行为码日志
// curFid: 当前 fid，如果传入则会设置到 request 中
// scenes: 场景值，可选
// cardListInfo: 卡片列表信息，可选，包含 page_attr, last, empty 等字段
func FinderActionLog(ctx *context.Context, curFid string, scenes int, cardListInfo map[string]any) {
	if ctx == nil {
		return
	}

	// 如果传入了 curFid，设置到 request 中
	if curFid != "" {
		ctx.SetRequestParam("fid", curFid)
	}

	// 初始化 cardListInfo
	if cardListInfo == nil {
		cardListInfo = make(map[string]any)
	}

	// 构建 actionLogParams
	actionLogParams := map[string]any{
		"page":            ctx.DefaultRequestParam("page", ""),
		"since_id":        ctx.DefaultRequestParam("since_id", ""),
		"mid":             ctx.DefaultRequestParam("mid", ""),
		"page_pro":        cardListInfo["page_attr"],
		"last":            cardListInfo["last"],
		"empty":           cardListInfo["empty"],
		"orifid":          ctx.DefaultRequestParam("oriuicode", ""),
		"request_referer": ctx.DefaultRequestParam("request_referer", ""),
	}

	// 处理 page_pro 默认值
	if actionLogParams["page_pro"] == nil {
		actionLogParams["page_pro"] = ""
	}

	// 构建 ext_info 数组
	var extInfo []string
	if ext := ctx.DefaultRequestParam("ext", ""); ext != "" {
		extInfo = append(extInfo, ext)
	}
	if extparam := ctx.DefaultRequestParam("extparam", ""); extparam != "" {
		extInfo = append(extInfo, extparam)
	}
	if len(extInfo) > 0 {
		actionLogParams["ext"] = extInfo
	}

	// 如果有 scenes 参数
	if scenes != 0 {
		actionLogParams["scenes"] = scenes
	}

	// 如果有 discover_flow_enable 参数
	if discoverFlowEnable := ctx.DefaultRequestParam("discover_flow_enable", ""); discoverFlowEnable != "" {
		actionLogParams["discover_flow_enable"] = discoverFlowEnable
	}

	actionLogParams["go_server"] = 1

	// 调用 Action 发送 26 行为码日志
	// 第四个参数是额外的 options，包含 fid
	Action(26, ctx.DefaultRequestParam("fid", ""), actionLogParams, ctx)
}

func AddCustomMark(ctx *context.Context, markName string, markValue ...any) {
	// 构建标记字符串
	var markStr string
	if len(markValue) > 0 {
		switch v := markValue[0].(type) {
		case string:
			markStr = fmt.Sprintf("%s[%s]", markName, v)
		case int:
			markStr = fmt.Sprintf("%s[%d]", markName, v)
		case float64:
			markStr = fmt.Sprintf("%s[%v]", markName, v) // 使用 %v 保留原始格式
		default:
			markStr = markName
		}
	} else {
		markStr = markName
	}

	// 从 Context 获取现有标记
	existingMark, exists := ctx.Get(CustomMarkKey)
	var markString string
	if exists && existingMark != "" {
		markString = existingMark.(string) + "-" + markStr
	} else {
		markString = markStr
	}

	ctx.Set(CustomMarkKey, markString)
}

// IsRecord 是否记录日志。如离线缓存刷接口类的，不需要写入日志避免影响正常用户统计
func IsRecord(ctx *context.Context) bool {
	if ctx == nil {
		return false
	}

	// 离线缓存刷接口不记录日志
	if ctx.DefaultRequestParam("gateway_offline_cache", "") != "" {
		return false
	}

	// 刷站环境不记录日志（除非是异步搜索请求）
	envGroup := os.Getenv("SEARCH_ENV_GROUP")
	if envGroup == "shuazhan" && !isAsyncSearch(ctx) {
		return false
	}

	// 命中抓站规则不记录日志
	if (ctx.DefaultRequestParam("hit_dproxy_zhuazhan_rule", "") != "" ||
		ctx.DefaultRequestParam("hit_local_zhuazhan", "") != "") && !isAsyncSearch(ctx) {
		return false
	}

	return true
}

// isAsyncSearch 是否是异步请求
func isAsyncSearch(ctx *context.Context) bool {
	if ctx == nil {
		return false
	}
	uri := ctx.Request().URL.Path
	return uri == "/mi/api/async.php"
}

// TopnClick topn 点击计算口径日志
func TopnClick(ctx *context.Context, result map[string]any) map[string]any {
	cards, ok := result["cards"].([]any)
	if !ok || len(cards) == 0 {
		return result
	}

	topn := newTopn(ctx)
	result = topn.handle(result)

	return result
}

// topn Topn 点击计算口径日志处理器
type topn struct {
	ctx *context.Context
	pos int

	// 标题类卡片类型
	titleCardTypes []int
	// 分割线类卡片类型
	spanCardTypes []int
	// 博文类卡片类型
	mblogCardTypes []int
	// 需要更新的 itemid
	updateItemid map[string]string
}

// newTopn 创建 topn 处理器
func newTopn(ctx *context.Context) *topn {
	return &topn{
		ctx:            ctx,
		pos:            0,
		titleCardTypes: []int{6, 42, 58, 101, 112, 120},
		spanCardTypes:  []int{143, 197},
		mblogCardTypes: []int{9, 59, 81, 88, 89, 164, 165, 243},
		updateItemid:   make(map[string]string),
	}
}

// handle 处理 topn 点击计算
func (t *topn) handle(result map[string]any) map[string]any {
	cards, ok := result["cards"].([]any)
	if !ok || len(cards) == 0 {
		return result
	}

	// topn_pos 优先继承上一页的使用
	if t.ctx != nil {
		if topnPos := t.ctx.RequestInt("topn_pos"); topnPos > 0 {
			t.pos = topnPos
		}
	}

	for key, cardItem := range cards {
		card, ok := cardItem.(map[string]any)
		if !ok {
			continue
		}

		if cardGroup, ok := card["card_group"].([]any); ok && len(cardGroup) > 0 {
			pos := t.processGroupPos(card)
			for i, item := range cardGroup {
				if itemMap, ok := item.(map[string]any); ok {
					cardGroup[i] = t.cardItem(itemMap, pos, true)
				}
			}
			card["card_group"] = cardGroup
		} else {
			card = t.cardItem(card, 0, false)
		}

		cards[key] = card
	}

	// 翻页透传 topn_pos，下一页使用
	if t.pos > 0 {
		cardlistInfo, ok := result["cardlistInfo"].(map[string]any)
		if !ok {
			cardlistInfo = make(map[string]any)
			result["cardlistInfo"] = cardlistInfo
		}
		extTrans, ok := cardlistInfo["ext_trans"].(map[string]any)
		if !ok {
			extTrans = make(map[string]any)
			cardlistInfo["ext_trans"] = extTrans
		}
		extTrans["topn_pos"] = t.pos
	}

	result["cards"] = cards

	// 将 itemid 更新为最新的
	result = t.updateNestedValues(result)

	return result
}

// processGroupPos 组合类 card 位置字段处理
func (t *topn) processGroupPos(card map[string]any) int {
	pos := 0
	groups, ok := card["card_group"].([]any)
	if !ok || len(groups) == 0 {
		return pos
	}

	var cardTypes []int
	for i, item := range groups {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		cardType := getCardType(itemMap)

		// 该卡片类型指定不做位置计算
		if wboxParam, ok := itemMap["wboxParam"].(map[string]any); ok {
			if noTopnPos, ok := wboxParam["no_topn_pos"].(bool); ok && noTopnPos {
				continue
			}
		}
		if noTopnPos, ok := itemMap["no_topn_pos"].(bool); ok && noTopnPos {
			continue
		}
		if noTopnPos, ok := itemMap["no_topn_pos"].(int); ok && noTopnPos != 0 {
			continue
		}

		// 去掉 title、更多入口、分割线、间距类型卡片
		if !t.containsInt(t.titleCardTypes, cardType) && !t.containsInt(t.spanCardTypes, cardType) {
			cardTypes = append(cardTypes, cardType)
		}
		_ = i
	}

	// 如果只剩一个卡片，则整体按一个位置计算
	// show_type，0:标示灰色分隔线连一起，3:分割线不展示
	count := len(cardTypes)
	showType := getShowType(card)
	if count == 1 || (count > 1 && (showType == 0 || showType == 3) && !t.hasMultiComputePosCard(cardTypes)) {
		t.pos++
		pos = t.pos
	}

	return pos
}

// hasMultiComputePosCard group 内是否包含多个计算位置的卡片
func (t *topn) hasMultiComputePosCard(cardTypes []int) bool {
	counts := make(map[int]int)
	for _, ct := range cardTypes {
		counts[ct]++
	}

	total := 0
	for _, mblogType := range t.mblogCardTypes {
		if counts[mblogType] > 1 {
			return true
		}
		total += counts[mblogType]
	}

	return total > 1
}

// cardItem 单个卡片元素处理
func (t *topn) cardItem(card map[string]any, pos int, isFromGroup bool) map[string]any {
	if t.isCalculatePos(card, pos, isFromGroup) {
		t.pos++
	}

	if !t.isAppend(card) || t.pos <= 0 {
		return card
	}

	// 异步轮询 card 格式：轮询次数_$this->pos
	loopNum := 0
	if t.ctx != nil {
		loopNum = t.ctx.RequestInt("loop_num")
	}

	var topnStr string
	if loopNum > 0 {
		topnStr = fmt.Sprintf("|topn_pos:%d_%d", loopNum, t.pos)
	} else {
		topnStr = fmt.Sprintf("|topn_pos:%d", t.pos)
	}

	card = t.append(card, topnStr)

	// 定义需要处理的字段映射表
	listFields := []string{"groups", "sub_cards", "group", "pics", "pic_items", "items", "elements"}
	for _, field := range listFields {
		if fieldVal, ok := card[field].([]any); ok && len(fieldVal) > 0 {
			for key, value := range fieldVal {
				if valueMap, ok := value.(map[string]any); ok {
					fieldVal[key] = t.append(valueMap, topnStr)
				}
			}
			card[field] = fieldVal
		}
	}

	// 处理 wboxParam
	if wboxParam, ok := card["wboxParam"].(map[string]any); ok {
		if data, ok := wboxParam["data"].([]any); ok && len(data) > 0 {
			for key, value := range data {
				if valueMap, ok := value.(map[string]any); ok {
					data[key] = t.append(valueMap, topnStr)
				}
			}
			wboxParam["data"] = data
		} else if user, ok := wboxParam["user"].(map[string]any); ok {
			wboxParam["user"] = t.append(user, topnStr)
		}
		card["wboxParam"] = wboxParam
	}

	// 处理双列卡片的特殊字段
	specialFields := []struct {
		field    string
		subField string
	}{
		{"left_element", "mblog"},
		{"right_element", "mblog"},
		{"left_element", "itemid"},
		{"right_element", "itemid"},
	}

	for _, sf := range specialFields {
		if elem, ok := card[sf.field].(map[string]any); ok {
			if sf.subField == "itemid" {
				card[sf.field] = t.append(elem, topnStr)
			} else if subVal, ok := elem[sf.subField].(map[string]any); ok {
				elem[sf.subField] = t.append(subVal, topnStr)
				card[sf.field] = elem
			}
		}
	}

	return card
}

// isCalculatePos 是否需要计算位置
func (t *topn) isCalculatePos(card map[string]any, pos int, isFromGroup bool) bool {
	cardType := getCardType(card)
	noTopnPos := getNoTopnPos(card)

	// 143，197 等作为纯分割线和间距使用，这类忽略位置计算
	// group 组合卡片内，存在的 title 卡片，也忽略位置计算
	if pos <= 0 &&
		!t.containsInt(t.spanCardTypes, cardType) &&
		!t.containsInt(t.titleCardTypes, cardType) &&
		noTopnPos == 0 {
		return true
	}

	return false
}

// isAppend 是否要在日志中追加 topn 口径位置字段
func (t *topn) isAppend(card map[string]any) bool {
	cardType := getCardType(card)
	noTopnPos := getNoTopnPos(card)

	if !t.containsInt(t.spanCardTypes, cardType) &&
		!t.containsInt(t.titleCardTypes, cardType) &&
		(noTopnPos == 0 || noTopnPos == 2) {
		return true
	}

	return false
}

// append 追加 topn 口径位置字段
func (t *topn) append(card map[string]any, topnStr string) map[string]any {
	card = t.feedExpo(card, topnStr)
	return t.actionlog(card, topnStr)
}

// feedExpo 卡片真实阅读字段追加 topn_pos 字段
func (t *topn) feedExpo(card map[string]any, topnStr string) map[string]any {
	readtimetype := ""
	if rt, ok := card["readtimetype"].(string); ok {
		readtimetype = rt
	}

	if readtimetype == "" || readtimetype == "card" {
		if itemid, ok := card["itemid"].(string); ok && itemid != "" {
			card["itemid"] = itemid + topnStr
			// 需要更新 item_id
			if updatePath, ok := card["topn_update_itemid_path"].(string); ok && updatePath != "" {
				t.updateItemid[updatePath] = card["itemid"].(string)
				delete(card, "topn_update_itemid_path")
			}
		}

		// cell 动态化卡片是 itemId
		if itemId, ok := card["itemId"].(string); ok && itemId != "" {
			card["itemId"] = itemId + topnStr
		}
	}

	if analysisExtra, ok := card["analysis_extra"].(string); ok && analysisExtra != "" {
		card["analysis_extra"] = analysisExtra + topnStr
	}

	if mblog, ok := card["mblog"].(map[string]any); ok {
		if analysisExtra, ok := mblog["analysis_extra"].(string); ok && analysisExtra != "" {
			mblog["analysis_extra"] = analysisExtra + topnStr
			card["mblog"] = mblog
		}
	}

	return card
}

// actionlog 卡片点击日志追加 topn_pos 字段
func (t *topn) actionlog(card map[string]any, topnStr string) map[string]any {
	if actionlog, ok := card["actionlog"].(map[string]any); ok {
		if ext, ok := actionlog["ext"].(string); ok && ext != "" {
			actionlog["ext"] = ext + topnStr
			card["actionlog"] = actionlog
			return card
		}
	}

	if actionLog, ok := card["action_log"].(map[string]any); ok {
		if ext, ok := actionLog["ext"].(string); ok && ext != "" {
			actionLog["ext"] = ext + topnStr
			card["action_log"] = actionLog
			return card
		}
	}

	if userActionlog, ok := card["user_actionlog"].(map[string]any); ok {
		if ext, ok := userActionlog["ext"].(string); ok && ext != "" {
			userActionlog["ext"] = ext + topnStr
			card["user_actionlog"] = userActionlog
			return card
		}
		if cardid, ok := userActionlog["cardid"].(string); ok && cardid != "" {
			userActionlog["cardid"] = cardid + topnStr
			card["user_actionlog"] = userActionlog
			return card
		}
	}

	if click, ok := card["click"].(map[string]any); ok {
		if actionlog, ok := click["actionlog"].(map[string]any); ok {
			if ext, ok := actionlog["ext"].(string); ok && ext != "" {
				actionlog["ext"] = ext + topnStr
				click["actionlog"] = actionlog
				card["click"] = click
				return card
			}
		}
	}

	// 处理 flow 类型卡片
	if cardType, ok := card["card_type"].(string); ok && cardType == "flow" {
		if click, ok := card["click"].(map[string]any); ok {
			if modules, ok := click["modules"].([]any); ok {
				for key, module := range modules {
					if moduleMap, ok := module.(map[string]any); ok {
						if actionlog, ok := moduleMap["actionlog"].(map[string]any); ok {
							if ext, ok := actionlog["ext"].(string); ok && ext != "" {
								actionlog["ext"] = ext + topnStr
								moduleMap["actionlog"] = actionlog
								modules[key] = moduleMap
							}
						}
					}
				}
				click["modules"] = modules
				card["click"] = click
			}
		}
	}

	return card
}

// updateNestedValues 更新嵌套数组中的 itemid 值
func (t *topn) updateNestedValues(result map[string]any) map[string]any {
	if len(t.updateItemid) == 0 {
		return result
	}

	for ukey, uvalue := range t.updateItemid {
		pathParts := strings.Split(ukey, ".")
		current := result

		pathExists := true
		for i := 0; i < len(pathParts)-1; i++ {
			if next, ok := current[pathParts[i]].(map[string]any); ok {
				current = next
			} else {
				pathExists = false
				break
			}
		}

		if pathExists {
			lastKey := pathParts[len(pathParts)-1]
			if _, ok := current[lastKey]; ok {
				current[lastKey] = uvalue
			}
		}
	}

	return result
}

// containsInt 检查切片是否包含指定整数
func (t *topn) containsInt(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// getCardType 获取卡片类型
func getCardType(card map[string]any) int {
	switch v := card["card_type"].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	return 0
}

// getShowType 获取展示类型
func getShowType(card map[string]any) int {
	switch v := card["show_type"].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	return 0
}

// getNoTopnPos 获取 no_topn_pos 值
func getNoTopnPos(card map[string]any) int {
	// 先检查 card 本身的 no_topn_pos
	if noTopnPos, ok := card["no_topn_pos"].(int); ok {
		return noTopnPos
	}
	if noTopnPos, ok := card["no_topn_pos"].(float64); ok {
		return int(noTopnPos)
	}

	// 再检查 wboxParam 中的 no_topn_pos
	if wboxParam, ok := card["wboxParam"].(map[string]any); ok {
		if noTopnPos, ok := wboxParam["no_topn_pos"].(int); ok {
			return noTopnPos
		}
		if noTopnPos, ok := wboxParam["no_topn_pos"].(float64); ok {
			return int(noTopnPos)
		}
	}

	return 0
}
