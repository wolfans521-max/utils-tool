package card

// C197 发现页标题卡片
// 用于发现页标题展示
type C197 struct {
	BaseItem
	BackgroundColor     string  `json:"backgroundColor,omitempty"`     // 背景颜色
	BackgroundDarkColor string  `json:"backgroundDarkColor,omitempty"` // 深色模式背景颜色
	Padding             []int   `json:"padding,omitempty"`             // 间隔 [左，上，下，右]
	Height              float64 `json:"height,omitempty"`              // 高度
	FloatHeight         float64 `json:"floatHeight,omitempty"`         // 高度
	OpenUrl             string  `json:"openurl,omitempty"`             // 打开链接
}

// NewCard197 创建Card197实例
func NewCard197() *C197 {
	return &C197{BaseItem: BaseItem{Base: Base{CardType: 197}}}
}
