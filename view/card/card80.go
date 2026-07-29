package card

// C80 大图卡片
// 16:9大图+人物信息+描述+JsonButton
type C80 struct {
	BaseItem
	TitleSub  string    `json:"title_sub,omitempty"`  // 标题
	TypeIcon  string    `json:"type_icon,omitempty"`  // 类型图标
	Buttons   []*Button `json:"buttons,omitempty"`    // 按钮列表
	Pic       string    `json:"pic,omitempty"`        // 图片
	PicWidth  int       `json:"pic_width,omitempty"`  // 图片宽度比例
	PicHeight int       `json:"pic_height,omitempty"` // 图片高度比例
	OpenUrl   string    `json:"openurl,omitempty"`    // 打开链接
}

// NewCard80 创建Card80实例
func NewCard80() *C80 {
	return &C80{BaseItem: BaseItem{Base: Base{CardType: 80}}}
}
