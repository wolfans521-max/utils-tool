package card

// C16 双按钮卡片
// 支持双按钮
type C16 struct {
	BaseItem
	CardTypeName string          `json:"card_type_name,omitempty"` // 卡片类型名称
	Group        []*C16GroupItem `json:"group,omitempty"`          // 双按钮
	OpenUrl      string          `json:"openurl,omitempty"`        // 打开链接
}

// C16GroupItem 按钮项
type C16GroupItem struct {
	TitleSub string `json:"title_sub,omitempty"` // 按钮标题
	Pic      string `json:"pic,omitempty"`       // 按钮图标
	Type     string `json:"type,omitempty"`      // 类型
	Scheme   string `json:"scheme,omitempty"`    // 跳转链接
}

// NewCard16 创建Card16实例
func NewCard16() *C16 {
	return &C16{BaseItem: BaseItem{Base: Base{CardType: 16}}}
}
