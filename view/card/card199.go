package card

// C199 实时讨论卡片（比赛弹幕）
// 分为content、match、flow_comments、comment_bar四部分
type C199 struct {
	BaseItem
	BarrageSpeed      int                `json:"barrage_speed,omitempty"`        // 弹幕速度
	BgColor           string             `json:"bgcolor,omitempty"`              // 背景颜色
	BgColorDark       string             `json:"bgcolor_dark,omitempty"`         // 背景颜色暗黑
	BgTopImg          string             `json:"bgtop_img,omitempty"`            // 上背景图
	BgBottomImg       string             `json:"bgbottom_img,omitempty"`         // 下背景图
	BgTopImgDark      string             `json:"bgtop_img_dark,omitempty"`       // 上背景图暗黑
	BgBottomImgDark   string             `json:"bgbottom_img_dark,omitempty"`    // 下背景图暗黑
	MarginLeft        int                `json:"margin_left,omitempty"`          // card内部左间距
	MarginRight       int                `json:"margin_right,omitempty"`         // card内部右间距
	MarginTop         int                `json:"margin_top,omitempty"`           // card内部上间距
	MarginBottom      int                `json:"margin_bottom,omitempty"`        // card内部下间距
	BgBorderColor     string             `json:"bg_border_color,omitempty"`      // card描边颜色
	BgBorderColorDark string             `json:"bg_border_color_dark,omitempty"` // card描边颜色暗黑
	BgBorderRadius    int                `json:"bg_border_radius,omitempty"`     // card描边圆角大小
	TitleInfo         *C199TitleInfo     `json:"title_info,omitempty"`           // 标题部分
	Match             *C199Match         `json:"match,omitempty"`                // 比分部分
	CommentBar        *C199CommentBar    `json:"comment_bar,omitempty"`          // 评论条部分
	FlowComments      []*C199FlowComment `json:"flow_comments,omitempty"`        // 评论弹幕部分
	ContentInfo       any                `json:"content_info,omitempty"`         // content部分
	Component         *C199Component     `json:"component,omitempty"`            // 中部组件（许愿池）
	OpenUrl           string             `json:"openurl,omitempty"`              // 打开链接
}

// C199TitleInfo 标题信息
type C199TitleInfo struct {
	LeftImg       string `json:"leftimg,omitempty"`        // 左侧图片
	LeftImgWidth  int    `json:"leftimg_width,omitempty"`  // 左侧图片宽度
	LeftImgHeight int    `json:"leftimg_height,omitempty"` // 左侧图片高度
	RightTxt      any    `json:"right_txt,omitempty"`      // 右侧文字（富文本数组）
	RightScheme   string `json:"right_scheme,omitempty"`   // 右侧跳转scheme
}

// C199Match 比分信息
type C199Match struct {
	LeftInfo     *C199MatchTeam `json:"left_info,omitempty"`     // 左侧队伍信息
	StateInfo    *C199StateInfo `json:"state_info,omitempty"`    // 状态信息
	RightInfo    *C199MatchTeam `json:"right_info,omitempty"`    // 右侧队伍信息
	CanRequest   int            `json:"can_request,omitempty"`   // 是否开启比分轮询
	Path         string         `json:"path,omitempty"`          // 比分轮询的path
	TimeInterval int            `json:"time_interval,omitempty"` // 比分轮询间隔
}

// C199MatchTeam 队伍信息
type C199MatchTeam struct {
	Image           string  `json:"image,omitempty"`             // 图片
	ImageDark       string  `json:"image_dark,omitempty"`        // 图片暗黑
	Text            string  `json:"text,omitempty"`              // 文本
	TextColor       string  `json:"text_color,omitempty"`        // 文本颜色
	TextColorDark   string  `json:"text_color_dark,omitempty"`   // 文本颜色暗黑
	TextSize        int     `json:"text_size,omitempty"`         // 文本大小
	ImageCorner     int     `json:"image_corner,omitempty"`      // 图片圆角
	BorderColor     string  `json:"border_color,omitempty"`      // 边框颜色
	BorderDarkColor string  `json:"border_dark_color,omitempty"` // 边框颜色暗黑
	BorderWidth     float64 `json:"border_width,omitempty"`      // 边框宽度
	ImageWidth      int     `json:"image_width,omitempty"`       // 图片宽度
	ImageHeight     int     `json:"image_height,omitempty"`      // 图片高度
}

// C199StateInfo 状态信息
type C199StateInfo struct {
	TopText                   string           `json:"top_text,omitempty"`                     // 顶部文本
	TopTextColor              string           `json:"top_text_color,omitempty"`               // 顶部文本颜色
	TopTextColorDark          string           `json:"top_text_color_dark,omitempty"`          // 顶部文本颜色暗黑
	TopTextSize               int              `json:"top_text_size,omitempty"`                // 顶部文本大小
	Center                    *C199ScoreCenter `json:"center,omitempty"`                       // 中间比分
	BottomText                string           `json:"bottom_text,omitempty"`                  // 底部文本
	BottomTextColor           string           `json:"bottom_text_color,omitempty"`            // 底部文本颜色
	BottomTextColorDark       string           `json:"bottom_text_color_dark,omitempty"`       // 底部文本颜色暗黑
	BottomTextSize            int              `json:"bottom_text_size,omitempty"`             // 底部文本大小
	BottomBackgroundColor     string           `json:"bottom_background_color,omitempty"`      // 底部背景颜色
	BottomBackgroundColorDark string           `json:"bottom_background_color_dark,omitempty"` // 底部背景颜色暗黑
	BottomCorner              int              `json:"bottom_corner,omitempty"`                // 底部圆角
}

// C199ScoreCenter 比分中心
type C199ScoreCenter struct {
	LeftText            string `json:"left_text,omitempty"`              // 左侧比分
	LeftTextColor       string `json:"left_text_color,omitempty"`        // 左侧比分颜色
	LeftTextColorDark   string `json:"left_text_color_dark,omitempty"`   // 左侧比分颜色暗黑
	LeftTextSize        int    `json:"left_text_size,omitempty"`         // 左侧比分大小
	MiddleText          string `json:"middle_text,omitempty"`            // 中间文本
	MiddleTextColor     string `json:"middle_text_color,omitempty"`      // 中间文本颜色
	MiddleTextColorDark string `json:"middle_text_color_dark,omitempty"` // 中间文本颜色暗黑
	MiddleTextSize      int    `json:"middle_text_size,omitempty"`       // 中间文本大小
	RightText           string `json:"right_text,omitempty"`             // 右侧比分
	RightTextColor      string `json:"right_text_color,omitempty"`       // 右侧比分颜色
	RightTextColorDark  string `json:"right_text_color_dark,omitempty"`  // 右侧比分颜色暗黑
	RightTextSize       int    `json:"right_text_size,omitempty"`        // 右侧比分大小
}

// C199CommentBar 评论条
type C199CommentBar struct {
	User            any            `json:"user,omitempty"`             // 用户信息
	BgColor         string         `json:"bgcolor,omitempty"`          // 背景颜色
	BgColorDark     string         `json:"bgcolor_dark,omitempty"`     // 背景颜色暗黑
	BorderColor     string         `json:"bordercolor,omitempty"`      // 边框颜色
	BorderColorDark string         `json:"bordercolor_dark,omitempty"` // 边框颜色暗黑
	Hint            string         `json:"hint,omitempty"`             // 提示文字
	Scheme          string         `json:"scheme,omitempty"`           // 跳转scheme
	ActionLog       map[string]any `json:"actionlog,omitempty"`        // 点击日志
}

// C199FlowComment 弹幕评论
type C199FlowComment struct {
	BgColor                 string          `json:"bgcolor,omitempty"`                    // 背景颜色
	BgColorDark             string          `json:"bgcolor_dark,omitempty"`               // 背景颜色暗黑
	User                    any             `json:"user,omitempty"`                       // 用户信息
	Text                    string          `json:"text,omitempty"`                       // 文本
	TextColor               string          `json:"text_color,omitempty"`                 // 文本颜色
	Scheme                  string          `json:"scheme,omitempty"`                     // 跳转scheme
	ActionLog               map[string]any  `json:"actionlog,omitempty"`                  // 点击日志
	GradientBorderColor     []string        `json:"gradient_border_color,omitempty"`      // 弹幕边框渐变色
	GradientBorderColorDark []string        `json:"gradient_border_color_dark,omitempty"` // 弹幕边框暗黑渐变色
	GradientBgColor         []string        `json:"gradient_bg_color,omitempty"`          // 弹幕背景渐变色
	GradientBgColorDark     []string        `json:"gradient_bg_color_dark,omitempty"`     // 弹幕背景暗黑渐变色
	RightLabel              *C199RightLabel `json:"right_label,omitempty"`                // 弹幕右上角角标图片
	BorderWidth             float64         `json:"border_width,omitempty"`               // 弹幕边框宽度
	Cleaned                 bool            `json:"cleaned,omitempty"`                    // 是否清理
}

// C199RightLabel 右上角标签
type C199RightLabel struct {
	Url    string `json:"url,omitempty"`    // 图片URL
	Width  int    `json:"width,omitempty"`  // 宽度
	Height int    `json:"height,omitempty"` // 高度
}

// C199Component 组件（许愿池）
type C199Component struct {
	Type          int              `json:"type,omitempty"`            // 类型
	LeftTopImg    *C199RightLabel  `json:"left_top_img,omitempty"`    // 左上图片
	LeftBottomImg *C199RightLabel  `json:"left_bottom_img,omitempty"` // 左下图片
	Desc1         any              `json:"desc1,omitempty"`           // 描述1
	Desc2         any              `json:"desc2,omitempty"`           // 描述2
	Desc3         any              `json:"desc3,omitempty"`           // 描述3
	RightButton   *C199RightButton `json:"right_button,omitempty"`    // 右侧按钮
}

// C199RightButton 右侧按钮
type C199RightButton struct {
	Text          string  `json:"text,omitempty"`            // 文本
	TextColor     string  `json:"text_color,omitempty"`      // 文本颜色
	TextColorDark string  `json:"text_color_dark,omitempty"` // 文本颜色暗黑
	BgColor       string  `json:"bg_color,omitempty"`        // 背景颜色
	BgColorDark   string  `json:"bg_color_dark,omitempty"`   // 背景颜色暗黑
	Scheme2       string  `json:"scheme2,omitempty"`         // 跳转scheme
	CommonButton  *Button `json:"common_button,omitempty"`   // 通用按钮
}

// NewCard199 创建Card199实例
func NewCard199() *C199 {
	return &C199{BaseItem: BaseItem{Base: Base{CardType: 199}}}
}
