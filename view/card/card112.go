package card

// C112 视频节目顶部标题卡片
// 最左侧支持下发两种大小的图片，title、desc展示方向支持可配
type C112 struct {
	BaseItem
	Icon      string      `json:"icon,omitempty"`       // 图片地址
	IconType  int8        `json:"icon_type,omitempty"`  // 0表示小头像，1表示大头像
	TitleType int8        `json:"title_type,omitempty"` // 0表示横向，1表示纵向，2表示提示文案(蓝色字体)
	Title     string      `json:"title,omitempty"`      // 大标题文案
	Desc      string      `json:"desc,omitempty"`       // 描述
	Button    *C112Button `json:"button,omitempty"`     // 按钮
	OpenUrl   string      `json:"openurl,omitempty"`    // 打开链接
}

// C112Button 按钮
type C112Button struct {
	ButtonType int8           `json:"button_type,omitempty"`  // 0表示无按钮，1表示文字按钮，2表示倒三角按钮，3表示右向箭头
	ButtonText string         `json:"button_text,omitempty"`  // 按钮文案，button_type为1时才生效
	FeedBackEx string         `json:"feed_back_ex,omitempty"` // 用于负反馈接口透传参数，button_type为2时才生效
	ActionLog  map[string]any `json:"actionlog,omitempty"`    // 打点日志
}

// NewCard112 创建Card112实例
func NewCard112() *C112 {
	return &C112{BaseItem: BaseItem{Base: Base{CardType: 112}}}
}
