package card

// C217 热搜推送通知卡片
// 发现页热搜模块上方新增热搜推送通知card
type C217 struct {
	BaseItem
	CloseTime    int               `json:"closeTime,omitempty"`    // 自动关闭时间单位秒
	CardHeight   int               `json:"cardHeight,omitempty"`   // card的高度
	CloseRequest *C217CloseRequest `json:"closeRequest,omitempty"` // 手动删除card后通知服务器的参数
	BgView       *C217BgView       `json:"bgView,omitempty"`       // 背景视图
	LeftIcon     *C217Icon         `json:"leftIcon,omitempty"`     // 左侧图标
	Text         *C217Text         `json:"text,omitempty"`         // 文本
	RightIcon    *C217Icon         `json:"rightIcon,omitempty"`    // 右侧图标
	OpenUrl      string            `json:"openurl,omitempty"`      // 打开链接
	RelativeMid  string            `json:"relative_mid,omitempty"`
}

// C217CloseRequest 关闭请求
type C217CloseRequest struct {
	Path   string `json:"path,omitempty"`   // 请求路径
	Params any    `json:"params,omitempty"` // 请求参数
}

// C217BgView 背景视图
type C217BgView struct {
	Scheme      string      `json:"scheme,omitempty"`      // 跳转scheme
	Margin      *C217Margin `json:"margin,omitempty"`      // 边距
	BgColor     string      `json:"bgColor,omitempty"`     // 背景颜色
	BgColorDark string      `json:"bgColorDark,omitempty"` // 背景颜色暗黑
	Corner      int         `json:"corner,omitempty"`      // 圆角
}

// C217Margin 边距
type C217Margin struct {
	Left   int `json:"left,omitempty"`   // 左边距
	Top    int `json:"top,omitempty"`    // 上边距
	Right  int `json:"right,omitempty"`  // 右边距
	Bottom int `json:"bottom,omitempty"` // 下边距
}

// C217Icon 图标
type C217Icon struct {
	IsIgnoreClick bool           `json:"isIgnoreClick,omitempty"` // 不响应点击事件
	Margin        *C217Margin    `json:"margin,omitempty"`        // 边距
	Width         int            `json:"width,omitempty"`         // 宽度
	Height        int            `json:"height,omitempty"`        // 高度
	Url           string         `json:"url,omitempty"`           // 图片URL
	UrlDark       string         `json:"urlDark,omitempty"`       // 图片URL暗黑
	ActionLog     map[string]any `json:"actionlog,omitempty"`     // 点击日志
}

// C217Text 文本
type C217Text struct {
	ChangeTime int             `json:"changeTime,omitempty"` // 自动轮播时间，单位秒
	Margin     *C217Margin     `json:"margin,omitempty"`     // 边距
	Richtexts  []*C217RichText `json:"richtexts,omitempty"`  // 轮播的内容
}

// C217RichText 富文本
type C217RichText struct {
	Contents any `json:"contents,omitempty"` // 富文本内容数组
}

// NewCard217 创建Card217实例
func NewCard217() *C217 {
	return &C217{BaseItem: BaseItem{Base: Base{CardType: 217}}}
}
