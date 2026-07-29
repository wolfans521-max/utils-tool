package card

// C229 横滑标签卡片
// 横滑标签card
type C229 struct {
	BaseItem
	CardPadding *C229Padding `json:"card_padding,omitempty"` // card内边距
	Tags        []*C229Tag   `json:"tags,omitempty"`         // 标签列表
	OpenUrl     string       `json:"openurl,omitempty"`      // 打开链接
}

// C229Padding 内边距
type C229Padding struct {
	Top    int `json:"top,omitempty"`    // 上边距
	Bottom int `json:"bottom,omitempty"` // 下边距
	Left   int `json:"left,omitempty"`   // 左边距
	Right  int `json:"right,omitempty"`  // 右边距
}

// C229Tag 标签
type C229Tag struct {
	Title     string         `json:"title,omitempty"`     // 标题
	Style     *C229TagStyle  `json:"style,omitempty"`     // 样式
	Params    *C229TagParams `json:"params,omitempty"`    // 参数
	ActionLog map[string]any `json:"actionlog,omitempty"` // 点击日志
}

// C229TagStyle 标签样式
type C229TagStyle struct {
	FontSize             int    `json:"font_size,omitempty"`              // 字体大小
	FontColor            string `json:"font_color,omitempty"`             // 字体颜色
	FontColorLight       string `json:"font_color_light,omitempty"`       // 字体颜色亮色
	FontColorDark        string `json:"font_color_dark,omitempty"`        // 字体颜色暗黑
	BackgroundColor      string `json:"background_color,omitempty"`       // 背景颜色
	BackgroundColorLight string `json:"background_color_light,omitempty"` // 背景颜色亮色
	BackgroundColorDark  string `json:"background_color_dark,omitempty"`  // 背景颜色暗黑
}

// C229TagParams 标签参数
type C229TagParams struct {
	Scheme  string `json:"scheme,omitempty"`  // 跳转scheme
	Cleaned bool   `json:"cleaned,omitempty"` // 是否清理
}

// NewCard229 创建Card229实例
func NewCard229() *C229 {
	return &C229{BaseItem: BaseItem{Base: Base{CardType: 229}}}
}
