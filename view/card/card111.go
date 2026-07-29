package card

// C111 三图卡片容器
// card110容器，固定展示3个card110数据，不可左右横滑图片
type C111 struct {
	BaseItem
	Spacing      int    `json:"spacing,omitempty"`       // 图片横向间距
	Padding      int    `json:"padding,omitempty"`       // 图片左右边距
	MarginBottom int    `json:"margin_bottom,omitempty"` // 图片竖向间距
	ShowAllPics  bool   `json:"show_all_pics,omitempty"` // 点击item进入大图页是否传入cardlist中所有图片
	SubCards     []any  `json:"sub_cards,omitempty"`     // Card110数组
	OpenUrl      string `json:"openurl,omitempty"`       // 打开链接
}

// NewCard111 创建Card111实例
func NewCard111() *C111 {
	return &C111{BaseItem: BaseItem{Base: Base{CardType: 111}}}
}
