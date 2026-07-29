package card

// C4 文本卡片
// 支持右侧配图标或圆头像
// 支持左侧文案、右侧配置可点击链接
// 11.11.0版本新增点击变灰功能
type C4 struct {
	BaseItem
	ImageRightPadding int8            `json:"image_right_padding,omitempty"` // 左侧图标与右侧标题文案的间距，单位dp
	ContentStyle      int8            `json:"content_style,omitempty"`       // 内容样式
	Desc              string          `json:"desc,omitempty"`                // 左侧标题
	RightDesc         string          `json:"right_desc,omitempty"`          // 右侧文案
	Height            int             `json:"height,omitempty"`              // 支持高度可配
	DisableBottomLine int8            `json:"disable_bottom_line,omitempty"` // 支持组合卡片时不显示下划线
	Bold              int8            `json:"bold,omitempty"`                // 标题加粗
	IsHtmlText        bool            `json:"is_html_text,omitempty"`        // 标题是否支持富文本
	DescExtr          string          `json:"desc_extr,omitempty"`           // 描述
	Pic               string          `json:"pic,omitempty"`                 // 左图
	DisplayArrow      int8            `json:"display_arrow,omitempty"`       // 右侧箭头
	OpenUrl           string          `json:"openurl,omitempty"`             // 打开链接
	AvatarUrl         string          `json:"avatar_url,omitempty"`          // 圆头像（和icon同时拥有时显示圆头像）
	Icon              string          `json:"icon,omitempty"`                // 右侧icon，支持gif
	IconWidth         int             `json:"icon_width,omitempty"`          // 右侧icon宽度
	IconHeight        int             `json:"icon_height,omitempty"`         // 右侧icon高度
	IconUrl           string          `json:"icon_url,omitempty"`            // 点击右侧icon后，删掉当前card，且请求接口
	ReadConfigs       []*C4ReadConfig `json:"read_configs,omitempty"`        // 点击后对应的时间段内显示的颜色
}

// C4ReadConfig 已读配置
type C4ReadConfig struct {
	StartTime       int    `json:"start_time,omitempty"`        // 开始时间
	EndTime         int    `json:"end_time,omitempty"`          // 结束时间，单位s
	ContentColorKey string `json:"content_color_key,omitempty"` // 标题已读颜色key(ios使用)
	ContentColor    string `json:"content_color,omitempty"`     // 标题已读颜色(安卓使用)
}

// NewCard4 创建Card4实例
func NewCard4() *C4 {
	return &C4{BaseItem: BaseItem{Base: Base{CardType: 4}}}
}
