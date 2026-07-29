package card

// C31 搜索卡片
// 支持显示默认背景文字
// 支持右侧按钮
type C31 struct {
	BaseItem
	Desc     string     `json:"desc,omitempty"`      // 搜索提示文字
	ShowType int8       `json:"show_type,omitempty"` // UI样式的类型
	Struct   *C31Struct `json:"struct,omitempty"`    // 右侧按钮信息
	OpenUrl  string     `json:"openurl,omitempty"`   // 打开链接
}

// C31Struct 右侧按钮信息
type C31Struct struct {
	Pic            string         `json:"pic,omitempty"`              // 按钮图片
	PicDesc        string         `json:"pic_desc,omitempty"`         // 按钮文字
	Scheme         string         `json:"scheme,omitempty"`           // 按钮scheme
	ActionLog      map[string]any `json:"actionlog,omitempty"`        // 点击按钮时的actionlog
	BgImg          string         `json:"bg_img,omitempty"`           // 按钮背景
	TitleColor     string         `json:"title_color,omitempty"`      // 按钮文字颜色
	TitleColorDark string         `json:"title_color_dark,omitempty"` // 按钮文字暗黑颜色
}

// NewCard31 创建Card31实例
func NewCard31() *C31 {
	return &C31{BaseItem: BaseItem{Base: Base{CardType: 31}}}
}
