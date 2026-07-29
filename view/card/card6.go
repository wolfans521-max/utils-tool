package card

// C6 文本卡片
// 文本支持5种颜色
// 支持CommonButton
// 支持箭头颜色可配
type C6 struct {
	BaseItem
	Desc           string    `json:"desc,omitempty"`             // 文本内容
	TitleColor     string    `json:"title_color,omitempty"`      // 文字颜色
	TitleColorDark string    `json:"title_color_dark,omitempty"` // 文字颜色暗黑
	ShowType       int8      `json:"show_type,omitempty"`        // 0现有默认颜色，1绿色，2红色，3蓝色，4高亮白色。当0时显示箭头
	Buttons        []*Button `json:"buttons,omitempty"`          // 支持CommonButton
	OpenUrl        string    `json:"openurl,omitempty"`          // 打开链接
	DisplayArrow   int       `json:"display_arrow,omitempty"`    // 更多箭头
}

// NewCard6 创建Card6实例
func NewCard6() *C6 {
	return &C6{BaseItem: BaseItem{Base: Base{CardType: 6}}}
}
