package card

// C162 媒体直播主播说卡片
// 用于媒体直播主播说展示
type C162 struct {
	BaseItem
	LineConfig *C162LineConfig `json:"line_config,omitempty"` // 线配置
	Timestamp  *C162Timestamp  `json:"timestamp,omitempty"`   // 时间戳
	Text       *C162Text       `json:"text,omitempty"`        // 文本
	OpenUrl    string          `json:"openurl,omitempty"`     // 打开链接
}

// C162LineConfig 线配置
type C162LineConfig struct {
	LineType   string `json:"line_type,omitempty"`   // 左侧线的类型：top顶端，normal中间，bottom底部，single单条
	PointColor string `json:"point_color,omitempty"` // 点的颜色
	LineColor  string `json:"line_color,omitempty"`  // 左侧线的颜色
}

// C162Timestamp 时间戳
type C162Timestamp struct {
	Color string `json:"color,omitempty"` // 发布时间的颜色，默认黑色
	Value int64  `json:"value,omitempty"` // 发布的时间戳
}

// C162Text 文本
type C162Text struct {
	Color string `json:"color,omitempty"` // 发布内容的颜色
	Value string `json:"value,omitempty"` // 发布的内容
}

// NewCard162 创建Card162实例
func NewCard162() *C162 {
	return &C162{BaseItem: BaseItem{Base: Base{CardType: 162}}}
}
