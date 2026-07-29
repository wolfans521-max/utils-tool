package card

// C145 视频通用标题卡片
// 用于视频通用标题展示
type C145 struct {
	BaseItem
	BottomDivider      int8           `json:"bottom_divider,omitempty"`        // 1有底部分割线，0无底部分割线
	TopDivider         int8           `json:"top_divider,omitempty"`           // 1有顶部分割线，0无顶部分割线
	ShowArrow          int8           `json:"show_arrow,omitempty"`            // 1显示箭头，0不显示箭头
	TitleText          string         `json:"title_text,omitempty"`            // 标题文本
	TitleColor         string         `json:"title_color,omitempty"`           // 标题色值
	TitleColorDark     string         `json:"title_color_dark,omitempty"`      // 深色模式标题色值
	LeftTagImg         string         `json:"left_tag_img,omitempty"`          // 左侧图标url
	LeftTagImgDark     string         `json:"left_tag_img_dark,omitempty"`     // 深色模式左侧图标url
	ArrowDesc          string         `json:"arrow_desc,omitempty"`            // 右侧箭头文本
	ArrowDescColor     string         `json:"arrow_desc_color,omitempty"`      // 右侧箭头文本色值
	ArrowDescColorDark string         `json:"arrow_desc_color_dark,omitempty"` // 右侧箭头文本深色模式色值
	ArrowImg           string         `json:"arrow_img,omitempty"`             // 右侧箭头url
	ArrowImgDark       string         `json:"arrow_img_dark,omitempty"`        // 右侧箭头深色模式url
	ExtraDesc          string         `json:"extra_desc,omitempty"`            // 左侧昵称
	ArrowDescScheme    string         `json:"arrow_desc_scheme,omitempty"`     // 右侧昵称点击跳转地址
	ArrowDescActionLog map[string]any `json:"arrow_desc_actionlog,omitempty"`  // 右侧昵称点击actionlog
	ExtraDescScheme    string         `json:"extra_desc_scheme,omitempty"`     // 左侧昵称点击跳转地址
	ExtraDescActionLog map[string]any `json:"extra_desc_actionlog,omitempty"`  // 左侧昵称点击actionlog
	DescScheme         string         `json:"desc_scheme,omitempty"`           // 描述跳转链接
	DecorateColor      string         `json:"decorate_color,omitempty"`        // 装饰颜色
	DecorateColorDark  string         `json:"decorate_color_dark,omitempty"`   // 深色模式装饰颜色
	OpenUrl            string         `json:"openurl,omitempty"`               // 打开链接
}

// NewCard145 创建Card145实例
func NewCard145() *C145 {
	return &C145{BaseItem: BaseItem{Base: Base{CardType: 145}}}
}
