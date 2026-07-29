package card

// C171 商品评分卡片
// 用于商品评分展示（与card170类似）
type C171 struct {
	BaseItem
	Buttons  []*C170Button `json:"buttons,omitempty"`   // 按钮列表
	StarInfo *C170StarInfo `json:"star_info,omitempty"` // 评分信息
	Style    *C170Style    `json:"style,omitempty"`     // 样式
	TagInfo  *C170TagInfo  `json:"tag_info,omitempty"`  // 标签信息
	OpenUrl  string        `json:"openurl,omitempty"`   // 打开链接
}

// NewCard171 创建Card171实例
func NewCard171() *C171 {
	return &C171{BaseItem: BaseItem{Base: Base{CardType: 171}}}
}
