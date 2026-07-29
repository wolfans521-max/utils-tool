package card

// C27 标题+图+描述+按钮卡片
// 支持趋势
type C27 struct {
	BaseItem
	TitleSub string    `json:"title_sub,omitempty"` // 标题
	Desc     string    `json:"desc,omitempty"`      // 描述
	Pic      string    `json:"pic,omitempty"`       // 图片
	Buttons  []*Button `json:"buttons,omitempty"`   // 按钮列表
	OpenUrl  string    `json:"openurl,omitempty"`   // 打开链接
}

// NewCard27 创建Card27实例
func NewCard27() *C27 {
	return &C27{BaseItem: BaseItem{Base: Base{CardType: 27}}}
}
