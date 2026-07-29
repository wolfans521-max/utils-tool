package card

// C143 灰条卡片
// 提供和cardlistInfo showStyle=1相同的能力，即一个灰条
type C143 struct {
	BaseItem
	BgColor     string `json:"bg_color,omitempty"`      // 背景色
	BgColorDark string `json:"bg_color_dark,omitempty"` // 暗黑背景色，默认值#151515
	Height      int    `json:"height,omitempty"`        // Card高度
	TopLine     bool   `json:"top_line,omitempty"`      // 顶部1px线（仅安卓支持）
	BottomLine  bool   `json:"bottom_line,omitempty"`   // 底部1px线（仅安卓支持）
	OpenUrl     string `json:"openurl,omitempty"`       // 打开链接
}

// NewCard143 创建Card143实例
func NewCard143() *C143 {
	return &C143{BaseItem: BaseItem{Base: Base{CardType: 143}}}
}
