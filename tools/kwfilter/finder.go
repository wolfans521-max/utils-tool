package kwfilter

import (
	"encoding/json"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/log"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
	"github.com/spf13/cast"
)

// FilterLog 过滤日志
type FilterLog struct {
	Mid           string `json:"mid"`
	Uid           string `json:"uid"`
	InterveneType string `json:"intervene_type"`
	InterveneInfo string `json:"intervene_info"`
	Timestamp     int64  `json:"time"`
}

// 热点流允许的微博状态
var hotFeedAllowedStatuses = []int{0, 1, 39, 37, 40, 47, 44}

// 热点流屏蔽的用户类型
var hotFeedBlockedUserTypes = []int{
	6, 7, 8, 9, 11, 12, 13, 14, 15, 16, 19, 701, 801, 1201, 1301, 1302, 1303, 1304, 1701, 10000,
}

// isGovernmentOrMedia 判断是否为政府或媒体账号
func isGovernmentOrMedia(weibo map[string]any) bool {
	user := utils.ArrayValue(weibo, "user", nil)
	if user == nil {
		return false
	}

	verified := cast.ToBool(utils.ArrayValue(user, "verified", false))
	verifiedType := cast.ToInt(utils.ArrayValue(user, "verified_type", 0))

	return verified && (verifiedType == 1 || verifiedType == 3)
}

// filterHotFeedSpecific 发现页特有过滤（谣言、safe_tags、urisk）
func filterHotFeedSpecific(weibo map[string]any) (bool, struct{ typ, msg string }) {
	// 1. 谣言、举报类型过滤
	mlevel := cast.ToInt(utils.ArrayValue(weibo, "mlevel", 0))
	ifYaoyan := (mlevel & 0x01) != 0
	ifJubao := (mlevel & 0x80) != 0

	if ifYaoyan || ifJubao {
		return true, struct{ typ, msg string }{"mid_yaoyan_jubao", "sens_mid_yaoyan_jubao"}
	}

	// 2. safe_tags 过滤
	if safeTags := utils.ArrayValue(weibo, "safe_tags", nil); safeTags != nil {
		safeTagsValue := cast.ToInt(safeTags)
		ifNosafe := (safeTagsValue&0x2) != 0 || (safeTagsValue&0x800) != 0
		if ifNosafe {
			return true, struct{ typ, msg string }{"mid_no_safe", "sens_mid_no_safe"}
		}
	}

	// 3. urisk 过滤
	user := utils.ArrayValue(weibo, "user", nil)
	if user != nil {
		if urisk := utils.ArrayValue(user, "urisk", nil); urisk != nil {
			uriskValue := cast.ToInt(urisk)
			ifNosafeUser := (uriskValue & 0x20) != 0
			if ifNosafeUser {
				return true, struct{ typ, msg string }{"no_safe_user", "sens_no_safe_user"}
			}
		}
	}

	return false, struct{ typ, msg string }{"", ""}
}

// createFilterLog 创建过滤日志
func createFilterLog(weibo map[string]any, filterType, filterMsg string) *FilterLog {
	mid := cast.ToString(utils.ArrayValue(weibo, "idstr", ""))
	uid := ""

	user := utils.ArrayValue(weibo, "user", nil)
	if user != nil {
		uid = cast.ToString(utils.ArrayValue(user, "idstr", ""))
	}

	return &FilterLog{
		Mid:           mid,
		Uid:           uid,
		InterveneType: filterType,
		InterveneInfo: filterMsg,
		Timestamp:     time.Now().Unix(),
	}
}

// hotFeedRiskManageFilter 发现页热门流风控策略过滤
func hotFeedRiskManageFilter(weibo map[string]any) (bool, *FilterLog) {
	// 1. 删除状态过滤
	if FilterMblogDeleted(weibo) {
		return true, createFilterLog(weibo, "mid_delete", "sens_mid_delete")
	}

	// 2. 微博状态过滤
	if FilterMblogStatus(weibo, hotFeedAllowedStatuses) {
		return true, createFilterLog(weibo, "mblog_status", "sens_mblog_status")
	}

	// 3. 用户类型过滤
	if FilterUserType(weibo, hotFeedBlockedUserTypes) {
		return true, createFilterLog(weibo, "user_status", "sens_user_status")
	}

	// 4. 发现页特有过滤（政府对媒体豁免）
	if !isGovernmentOrMedia(weibo) {
		if filtered, reason := filterHotFeedSpecific(weibo); filtered {
			return true, createFilterLog(weibo, reason.typ, reason.msg)
		}
	}

	return false, nil
}

// DiscoverHotFeed 发现页热门流风控过滤
// 参数:
//   - ctx: 请求上下文（可选，为 nil 时使用默认的 uid/seqId）
//   - weibo: 微博数据（map[string]any 结构）
//
// 返回:
//   - int: 过滤类型，0=未过滤，4=风控敏感过滤
//   - *FilterLog: 过滤日志（如果被过滤）
func DiscoverHotFeed(ctx *context.Context, weibo map[string]any) (int, *FilterLog) {
	if weibo == nil {
		return 0, nil
	}

	// 执行风控过滤
	if filtered, filterLog := hotFeedRiskManageFilter(weibo); filtered {
		logMsg, _ := json.Marshal(filterLog)
		log.Write("finder_hot_feed_sensitive", logMsg, ctx)

		return 4, filterLog
	}

	return 0, nil
}

// DiscoverHotFeedList 批量过滤发现页热门流微博
// 参数:
//   - ctx: 请求上下文
//   - statuses: 微博内容列表
//
// 返回:
//   - 过滤后的微博列表
func DiscoverHotFeedList(ctx *context.Context, statuses any) any {
	statusesSlice, ok := statuses.([]any)
	if !ok {
		return statuses
	}

	result := make([]any, 0, len(statusesSlice))

	for _, item := range statusesSlice {
		weibo, ok := item.(map[string]any)
		if !ok {
			result = append(result, item)
			continue
		}

		filterType, _ := DiscoverHotFeed(ctx, weibo)
		if filterType == 0 {
			result = append(result, item)
		}
	}

	return result
}
