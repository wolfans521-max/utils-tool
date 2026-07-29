package card

// C242 用户/超话信息卡片
// 用于用户或超话信息展示
type C242 struct {
	BaseItem
	Intro                   string           `json:"intro,omitempty"`                      // 简介
	Fans                    string           `json:"fans,omitempty"`                       // 粉丝数
	Topic                   string           `json:"topic,omitempty"`                      // 超话标题
	TopicIcon               string           `json:"topic_icon,omitempty"`                 // 超话icon
	TopicAvatar             string           `json:"topic_avatar,omitempty"`               // 超话头像
	MaskHeight              int              `json:"mask_height,omitempty"`                // 顶部渐变遮罩高度
	MaskStartColor          string           `json:"mask_start_color,omitempty"`           // 顶部渐变遮罩起始颜色
	MaskEndColor            string           `json:"mask_end_color,omitempty"`             // 顶部渐变遮罩结束颜色
	UserBackgroundColor     string           `json:"user_background_color,omitempty"`      // 用户信息View背景色
	UserBackgroundDarkColor string           `json:"user_background_dark_color,omitempty"` // 用户信息View暗黑背景色
	CornerMarkData          *C242CornerMark  `json:"corner_mark_data,omitempty"`           // 右上角广告标
	BgPicUrl                string           `json:"bg_pic_url,omitempty"`                 // 背景图
	BrandWindow             *C242BrandWindow `json:"brand_window,omitempty"`               // 品牌橱窗
	User                    any              `json:"user,omitempty"`                       // 用户信息
	FollowButton            *Button          `json:"follow_button,omitempty"`              // 关注按钮
	OpenUrl                 string           `json:"openurl,omitempty"`                    // 打开链接
}

// C242CornerMark 角标
type C242CornerMark struct {
	Title           string `json:"title,omitempty"`            // 标题
	TitleSize       int    `json:"title_size,omitempty"`       // 标题大小
	TitleColor      string `json:"title_color,omitempty"`      // 标题颜色
	BorderColor     string `json:"border_color,omitempty"`     // 边框颜色
	BorderWidth     int    `json:"border_width,omitempty"`     // 边框宽度
	Padding         []int  `json:"padding,omitempty"`          // 内边距 [左, 上, 右, 下]
	Radius          []int  `json:"radius,omitempty"`           // 圆角 [左上, 右上, 左下, 右下]
	BackgroundColor string `json:"background_color,omitempty"` // 背景颜色
}

// C242BrandWindow 品牌橱窗
type C242BrandWindow struct {
	Padding []int            `json:"padding,omitempty"` // 内边距 [左, 上, 右, 下]
	Items   []*C242BrandItem `json:"items,omitempty"`   // 品牌橱窗数组
}

// C242BrandItem 品牌项
type C242BrandItem struct {
	Icon                string         `json:"icon,omitempty"`                  // 图标
	IconDark            string         `json:"icon_dark,omitempty"`             // 图标暗黑
	DisplayName         string         `json:"display_name,omitempty"`          // 显示名称
	TextColor           string         `json:"text_color,omitempty"`            // 文字颜色
	TextDarkColor       string         `json:"text_dark_color,omitempty"`       // 文字暗黑颜色
	TextSize            int            `json:"text_size,omitempty"`             // 文字大小
	BackgroundColor     string         `json:"background_color,omitempty"`      // 背景颜色
	BackgroundDarkColor string         `json:"background_dark_color,omitempty"` // 背景暗黑颜色
	Scheme              string         `json:"scheme,omitempty"`                // 跳转scheme
	ActionLog           map[string]any `json:"actionlog,omitempty"`             // 点击日志
}

// NewCard242 创建Card242实例
func NewCard242() *C242 {
	return &C242{BaseItem: BaseItem{Base: Base{CardType: 242}}}
}
