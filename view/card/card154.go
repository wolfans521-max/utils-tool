package card

// C154 创作中心小标题卡片
// 用于创作中心小标题展示
type C154 struct {
	BaseItem
	Img       string `json:"img,omitempty"`        // 图片
	ArrowImg  string `json:"arrow_img,omitempty"`  // 箭头图片
	MainTitle string `json:"main_title,omitempty"` // 主标题
	OpenUrl   string `json:"openurl,omitempty"`    // 打开链接
}

// NewCard154 创建Card154实例
func NewCard154() *C154 {
	return &C154{BaseItem: BaseItem{Base: Base{CardType: 154}}}
}
