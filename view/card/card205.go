package card

// C205 赛程卡片
// 用于赛程展示
type C205 struct {
	BaseItem
	LeftItem  *C205TeamItem `json:"leftItem,omitempty"`  // 左边队伍
	RightItem *C205TeamItem `json:"rightItem,omitempty"` // 右边队伍
	Title     any           `json:"title,omitempty"`     // 标题（富文本数组）
	Desc      any           `json:"desc,omitempty"`      // 副标题（富文本数组）
	SubTitle  any           `json:"subTitle,omitempty"`  // 比分（富文本数组）
	Button    *C205Button   `json:"button,omitempty"`    // 按钮
	OpenUrl   string        `json:"openurl,omitempty"`   // 打开链接
}

// C205TeamItem 队伍项
type C205TeamItem struct {
	Pic     string `json:"pic,omitempty"`     // 国旗、标志图片地址
	PicDark string `json:"picDark,omitempty"` // 国旗、标志图片地址暗黑
	Name    any    `json:"name,omitempty"`    // 队名（富文本数组）
}

// C205Button 按钮
type C205Button struct {
	Image      string           `json:"image,omitempty"`      // 按钮左边图片
	ImageDark  string           `json:"imageDark,omitempty"`  // 按钮左边图片暗黑
	Title      *C205ButtonTitle `json:"title,omitempty"`      // 按钮标题
	Border     *C205Border      `json:"border,omitempty"`     // 边框
	Background *C205Background  `json:"background,omitempty"` // 背景
	Params     any              `json:"params,omitempty"`     // 参数
	Radius     int              `json:"radius,omitempty"`     // 按钮圆角
	Scheme     string           `json:"scheme,omitempty"`     // 按钮scheme地址
	ActionLog  map[string]any   `json:"actionlog,omitempty"`  // 日志
}

// C205ButtonTitle 按钮标题
type C205ButtonTitle struct {
	Style   *C205TextStyle `json:"style,omitempty"`   // 样式
	Content string         `json:"content,omitempty"` // 文本内容
}

// C205TextStyle 文本样式
type C205TextStyle struct {
	TextColor    string `json:"textColor,omitempty"`    // 文字颜色
	TextColorKey string `json:"textColorKey,omitempty"` // 文字颜色key
	TextSize     int    `json:"textSize,omitempty"`     // 文字大小
}

// C205Border 边框
type C205Border struct {
	Color    string `json:"color,omitempty"`    // 边框色值
	ColorKey string `json:"colorKey,omitempty"` // 边框色值key
	Width    int    `json:"width,omitempty"`    // 边框宽度
}

// C205Background 背景
type C205Background struct {
	HightLightColor     string `json:"hightLightColor,omitempty"`     // 选中背景色
	HightLightDarkColor string `json:"hightLightDarkColor,omitempty"` // 选中暗黑背景色
}

// NewCard205 创建Card205实例
func NewCard205() *C205 {
	return &C205{BaseItem: BaseItem{Base: Base{CardType: 205}}}
}
