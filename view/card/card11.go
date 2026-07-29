package card

// C11 组合卡片
// 支持竖向拼接多个Card
// card_header支持吸顶
type C11 struct {
	BaseItem
	ShowType   int8           `json:"show_type,omitempty"`   // 设置card_group中的分割线样式，3没有线
	CardHeader any            `json:"card_header,omitempty"` // 支持所有类型的card
	CardGroup  []any          `json:"card_group,omitempty"`  // 可以包含所有类型的card
	OpenUrl    string         `json:"openurl,omitempty"`     // 打开链接
	GroupStyle map[string]any `json:"group_style,omitempty"` // 组样式，如 padding、margin
}

// NewCard11 创建Card11实例
func NewCard11() *C11 {
	return &C11{BaseItem: BaseItem{Base: Base{CardType: 11}}}
}
