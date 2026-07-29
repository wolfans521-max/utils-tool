package card

// C209 实况讨论/深度议题入口卡片
// 发现页新增实况讨论/深度议题入口card（与Card208类似）
type C209 struct {
	BaseItem
	CardPadding         *C209Padding   `json:"card_padding,omitempty"`          // card上下左右内边距
	TimeInterval        int            `json:"time_interval,omitempty"`         // 每个item执行动画时长，-1表示时长无限大
	Direction           int            `json:"direction,omitempty"`             // 动画方向，0从上面滚出，1从下面滚出
	Pic                 string         `json:"pic,omitempty"`                   // 左侧图片
	PicDark             string         `json:"pic_dark,omitempty"`              // 左侧图片暗黑
	BackgroundPic       string         `json:"background_pic,omitempty"`        // 左侧背景图
	BackgroundPicDark   string         `json:"background_pic_dark,omitempty"`   // 左侧背景图暗黑
	BackgroundColor     string         `json:"background_color,omitempty"`      // 左侧背景颜色
	BackgroundColorDark string         `json:"background_color_dark,omitempty"` // 左侧背景颜色暗黑
	BackgroundCorner    int            `json:"background_corner,omitempty"`     // 背景圆角
	TitleInfo           *C209TitleInfo `json:"title_info,omitempty"`            // 标题
	Groups              []*C209Group   `json:"groups,omitempty"`                // 标题下面滚动区域
	RightMore           *C209RightMore `json:"right_more,omitempty"`            // 右侧区域
	OpenUrl             string         `json:"openurl,omitempty"`               // 打开链接
}

// C209Padding 内边距
type C209Padding struct {
	Left   int `json:"left,omitempty"`   // 左边距
	Right  int `json:"right,omitempty"`  // 右边距
	Bottom int `json:"bottom,omitempty"` // 下边距
	Top    int `json:"top,omitempty"`    // 上边距
}

// C209TitleInfo 标题信息
type C209TitleInfo struct {
	TextColor     string `json:"text_color,omitempty"`      // 文字颜色
	TextColorDark string `json:"text_color_dark,omitempty"` // 文字颜色暗黑
	Content       string `json:"content,omitempty"`         // 内容
	LeftIcon      string `json:"left_icon,omitempty"`       // 标题左侧icon
	LeftIconDark  string `json:"left_icon_dark,omitempty"`  // 标题左侧icon暗黑
	RightIcon     string `json:"right_icon,omitempty"`      // 标题右侧icon
	RightIconDark string `json:"right_icon_dark,omitempty"` // 标题右侧icon暗黑
}

// C209Group 滚动区域组
type C209Group struct {
	Icons           []string `json:"icons,omitempty"`             // 图标列表
	Users           []string `json:"users,omitempty"`             // 用户头像列表
	Type            string   `json:"type,omitempty"`              // 类型 circle圆形
	Content         string   `json:"content,omitempty"`           // 内容
	Border          int      `json:"border,omitempty"`            // 边框
	BorderColor     string   `json:"border_color,omitempty"`      // 边框颜色
	BorderColorDark string   `json:"border_color_dark,omitempty"` // 边框颜色暗黑
	TextColor       string   `json:"text_color,omitempty"`        // 文字颜色
	TextColorDark   string   `json:"text_color_drak,omitempty"`   // 文字颜色暗黑
}

// C209RightMore 右侧更多
type C209RightMore struct {
	BackgroundPic       string         `json:"background_pic,omitempty"`        // 背景图
	BackgroundPicDark   string         `json:"background_pic_dark,omitempty"`   // 背景图暗黑
	BackgroundColor     string         `json:"background_color,omitempty"`      // 背景颜色
	BackgroundColorDark string         `json:"background_color_dark,omitempty"` // 背景颜色暗黑
	BackgroundCorner    int            `json:"background_corner,omitempty"`     // 背景圆角
	Content             string         `json:"content,omitempty"`               // 内容
	TextColor           string         `json:"text_color,omitempty"`            // 文字颜色
	TextColorDark       string         `json:"text_color_drak,omitempty"`       // 文字颜色暗黑
	RightIcon           string         `json:"right_icon,omitempty"`            // 右侧icon
	RightIconDark       string         `json:"right_icon_dark,omitempty"`       // 右侧icon暗黑
	Scheme              string         `json:"scheme,omitempty"`                // 跳转scheme
	ActionLog           map[string]any `json:"actionlog,omitempty"`             // 点击日志
}

// NewCard209 创建Card209实例
func NewCard209() *C209 {
	return &C209{BaseItem: BaseItem{Base: Base{CardType: 209}}}
}
