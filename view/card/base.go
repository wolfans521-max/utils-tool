package card

import (
	"regexp"
)

// 预编译正则表达式，避免每次调用都编译
var (
	// cateIdRegex 用于从字符串中匹配 cate_id
	cateIdRegex = regexp.MustCompile(`cate=(\d+)`)
)

type NoFormat struct{}

// ForceGreySetter 蒙层设置接口
// 所有嵌入Base的card类型都会自动实现此接口
type ForceGreySetter interface {
	SetForceGrey(value int8)
}

type Base struct {
	NoFormat
	Scheme        string         `json:"scheme,omitempty"`
	CardType      uint16         `json:"card_type,omitempty"`
	CardTypeName  string         `json:"card_type_name,omitempty"`
	ItemId        string         `json:"itemid,omitempty"`
	CardId        string         `json:"cardid,omitempty"`
	Title         string         `json:"title,omitempty"`
	ActionLog     map[string]any `json:"actionlog,omitempty"`
	RoundedCorner int            `json:"roundedcorner,omitempty"`
	ActStatus     any            `json:"act_status,omitempty"`
	ObjectId      string         `json:"object_id,omitempty"`
	ObjectType    string         `json:"object_type,omitempty"`
	MediaInfo     any            `json:"media_info,omitempty"`
	Promotion     any            `json:"promotion,omitempty"`
	Viewer        any            `json:"viewer,omitempty"`
	OpenUrl       string         `json:"openurl,omitempty"`
	AnalysisExtra string         `json:"analysis_extra,omitempty"`
	ReadtimeType  string         `json:"readtimetype,omitempty"`
	CateId        string         `json:"cate_id,omitempty"`
	IsForceGrey   int8           `json:"is_force_grey,omitempty"` // 蒙层标记
}

// SetForceGrey 实现ForceGreySetter接口
func (b *Base) SetForceGrey(value int8) {
	b.IsForceGrey = value
}

// SetActionLog 设置 ActionLog 并自动处理 AnalysisExtra 和 CateId
func (b *Base) SetActionLog(actionLog map[string]any) {
	b.ActionLog = actionLog
	b.ProcessAnalysisExtra()
}

// ProcessAnalysisExtra 处理 AnalysisExtra 和 CateId
// 从 ActionLog.ext 中提取 analysis_extra 和 cate_id
func (b *Base) ProcessAnalysisExtra() {
	if b.ActionLog == nil {
		return
	}
	if ext, ok := b.ActionLog["ext"].(string); ok && ext != "" {
		if b.AnalysisExtra == "" {
			b.AnalysisExtra = ext
		}
		if b.CateId == "" {
			b.CateId = MatchCateId(ext)
		}
	}
}

type BaseItem struct {
	Base
	ItemStyle *Style `json:"item_style,omitempty"`
}

// MatchCateId 从字符串中匹配 cate_id
func MatchCateId(str string) string {
	if str == "" {
		return ""
	}

	// 使用预编译的正则表达式
	matches := cateIdRegex.FindStringSubmatch(str)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
