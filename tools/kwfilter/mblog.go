package kwfilter

import (
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
	"github.com/spf13/cast"
)

// MblogFilterConfig 微博过滤配置
type MblogFilterConfig struct {
	AllowedStatuses  []int // 允许的微博状态列表
	BlockedUserTypes []int // 屏蔽的用户类型列表
}

// 默认允许的微博状态
var defaultAllowedStatuses = []int{0, 1, 39, 37, 40, 47, 44}

// 默认屏蔽的用户类型
var defaultBlockedUserTypes = []int{
	6, 7, 8, 9, 11, 12, 13, 14, 15, 16, 19, 701, 801, 1201, 1301, 1302, 1303, 1304, 1701, 10000,
}

// GetDefaultMblogFilterConfig 获取默认的微博过滤配置
func GetDefaultMblogFilterConfig() *MblogFilterConfig {
	return &MblogFilterConfig{
		AllowedStatuses:  defaultAllowedStatuses,
		BlockedUserTypes: defaultBlockedUserTypes,
	}
}

// FilterMblogDeleted 微博删除状态过滤
func FilterMblogDeleted(weibo map[string]any) bool {
	deleted := cast.ToInt(utils.ArrayValue(weibo, "deleted", 0))
	return deleted > 0
}

// FilterMblogStatus 微博状态过滤
func FilterMblogStatus(weibo map[string]any, allowedStatuses []int) bool {
	status := cast.ToInt(utils.ArrayValue(weibo, "status", 0))

	for _, allowed := range allowedStatuses {
		if status == allowed {
			return false // 状态合法，不过滤
		}
	}
	return true // 状态不在允许列表，过滤
}

// FilterUserType 用户类型过滤
func FilterUserType(weibo map[string]any, blockedUserTypes []int) bool {
	user := utils.ArrayValue(weibo, "user", nil)
	if user == nil {
		return false
	}

	userType := cast.ToInt(utils.ArrayValue(user, "type", 0))

	for _, blocked := range blockedUserTypes {
		if userType == blocked {
			return true // 用户类型在屏蔽列表，过滤
		}
	}
	return false // 用户类型不在屏蔽列表，不过滤
}
