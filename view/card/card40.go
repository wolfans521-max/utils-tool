package card

// C40 多位置文字卡片
// 支持多位置文字
// 支持右下文字配色
type C40 struct {
	BaseItem
	PicUrl     string `json:"pic_url,omitempty"`     // 图片URL
	PicScheme  string `json:"pic_scheme,omitempty"`  // 图片跳转
	TitleSub   string `json:"title_sub,omitempty"`   // 标题(左上)
	Desc1      string `json:"desc1,omitempty"`       // 左下文字
	Desc2      string `json:"desc2,omitempty"`       // 右下文字
	Desc3      string `json:"desc3,omitempty"`       // 右上文字
	Desc2Color string `json:"desc2_color,omitempty"` // 右下文字配色（0绿色 1红色 2灰色）
}

// NewCard40 创建Card40实例
func NewCard40() *C40 {
	return &C40{BaseItem: BaseItem{Base: Base{CardType: 40}}}
}
