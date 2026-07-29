package card

// C150 视频排行榜横滑卡片
// 用于视频排行榜展示
type C150 struct {
	BaseItem
	Index   int         `json:"index,omitempty"`   // 在页面上各card中的排序位置
	Items   []*C150Item `json:"items,omitempty"`   // 包含多个横滑的item数据
	OpenUrl string      `json:"openurl,omitempty"` // 打开链接
}

// C150Item 项目
type C150Item struct {
	IndexIcon    string            `json:"index_icon,omitempty"`    // 左上角排行角标icon的url
	Image        string            `json:"image,omitempty"`         // 每个item的封面图的url
	ActionLog    map[string]any    `json:"actionlog,omitempty"`     // 服务端下发日志字段
	Mid          string            `json:"mid,omitempty"`           // 该视频对应的mid
	UnifiedParam *C150UnifiedParam `json:"unified_param,omitempty"` // 视频对应的参数
}

// C150UnifiedParam 统一参数
type C150UnifiedParam struct {
	BizType    string `json:"biz_type,omitempty"`    // 业务类型
	BizId      string `json:"biz_id,omitempty"`      // 业务ID
	NextCursor int    `json:"next_cursor,omitempty"` // 下一页游标
}

// NewCard150 创建Card150实例
func NewCard150() *C150 {
	return &C150{BaseItem: BaseItem{Base: Base{CardType: 150}}}
}
