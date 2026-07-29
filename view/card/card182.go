package card

// C182 横滑展示卡片
// 支持横滑展示内容的Card，可以向其中添加不同种类的内部小卡片
type C182 struct {
	BaseItem
	Padding  []int          `json:"padding,omitempty"`   // 上，左，下，右 间隔
	Spacing  int            `json:"spacing,omitempty"`   // 小卡片中间的间隔
	SubItems []*C182SubItem `json:"sub_items,omitempty"` // 小卡片结构数组
	OpenUrl  string         `json:"openurl,omitempty"`   // 打开链接
}

// C182SubItem 小卡片
type C182SubItem struct {
	Width       int              `json:"width,omitempty"`         // 宽度
	Height      int              `json:"height,omitempty"`        // 高度
	BgImage     string           `json:"bg_image,omitempty"`      // 背景图
	BgImageDark string           `json:"bg_image_dark,omitempty"` // 深色模式背景图
	Header      *C182Header      `json:"header,omitempty"`        // 头部
	Components  []*C182Component `json:"components,omitempty"`    // Components 容器
}

// C182Header 头部
type C182Header struct {
	Type      string          `json:"type,omitempty"`       // 类型 richTextHeader
	Icon      string          `json:"icon,omitempty"`       // header 左侧图标
	Title     any             `json:"title,omitempty"`      // header 标题（可以是字符串或富文本数组）
	SubTitle  []*C182RichText `json:"subTitle,omitempty"`   // 副标题
	Content   string          `json:"content,omitempty"`    // header右侧文字
	ShowRight bool            `json:"show_right,omitempty"` // header是否展示右侧箭头
	Scheme    string          `json:"scheme,omitempty"`     // header 跳转链接
	ActionLog map[string]any  `json:"actionlog,omitempty"`  // header 点击日志
	Cleaned   bool            `json:"cleaned,omitempty"`    // 是否清理
}

// C182Component 组件
type C182Component struct {
	Type       int              `json:"type,omitempty"`        // 组件类型 1:热议 2:人物/直播 3:超话 4:潮流 5:赛程 6:超话活动聚合
	Title      []*C182RichText  `json:"title,omitempty"`       // 主标题（支持图文混排）
	SubTitle   []*C182RichText  `json:"sub_title,omitempty"`   // 副标题（支持图文混排）
	Image      string           `json:"image,omitempty"`       // 右侧图片
	Icon       string           `json:"icon,omitempty"`        // 左侧图片
	Scheme     string           `json:"scheme,omitempty"`      // 点击跳转链接
	ActionLog  map[string]any   `json:"actionlog,omitempty"`   // 点击日志
	User       any              `json:"user,omitempty"`        // 用户信息（type=2时使用）
	Button     *C182Button      `json:"button,omitempty"`      // 右侧按钮
	Items      []*C182TrendItem `json:"items,omitempty"`       // 潮流内容数组（type=4时使用）
	LeftItem   *C182MatchItem   `json:"leftItem,omitempty"`    // 左边队伍（type=5时使用）
	RightItem  *C182MatchItem   `json:"rightItem,omitempty"`   // 右边队伍（type=5时使用）
	RightTitle []*C182RichText  `json:"right_title,omitempty"` // 右侧标题（type=6时使用）
	Cleaned    bool             `json:"cleaned,omitempty"`     // 是否清理
}

// C182RichText 富文本
type C182RichText struct {
	Type    string         `json:"type,omitempty"`    // 类型 text/icon
	Style   *C182TextStyle `json:"style,omitempty"`   // 样式
	Content string         `json:"content,omitempty"` // 内容
	IconUrl string         `json:"iconUrl,omitempty"` // 图标URL
}

// C182TextStyle 文本样式
type C182TextStyle struct {
	Width         int    `json:"width,omitempty"`         // 宽度
	Height        int    `json:"height,omitempty"`        // 高度
	DarkMode      string `json:"darkMode,omitempty"`      // 暗黑模式
	TextColor     string `json:"textColor,omitempty"`     // 文字颜色
	TextColorKey  string `json:"textColorKey,omitempty"`  // 文字颜色Key
	TextColorDark string `json:"textColorDark,omitempty"` // 文字颜色暗黑
	TextSize      int    `json:"textSize,omitempty"`      // 文字大小
}

// C182Button 按钮
type C182Button struct {
	Type            string          `json:"type,omitempty"`              // 按钮类型
	SubType         int             `json:"sub_type,omitempty"`          // 子类型
	Name            string          `json:"name,omitempty"`              // 按钮名称
	NameAfter       string          `json:"name_after,omitempty"`        // 点击后名称
	TitleAfterClick string          `json:"title_after_click,omitempty"` // 点击后标题
	DtStart         int64           `json:"dt_start,omitempty"`          // 开始时间
	Title           any             `json:"title,omitempty"`             // 标题（可以是字符串或富文本数组）
	SkipFormat      int             `json:"skip_format,omitempty"`       // 跳过格式化
	ShowLoading     int             `json:"show_loading,omitempty"`      // 显示加载
	Scheme          string          `json:"scheme,omitempty"`            // 跳转链接
	Params          any             `json:"params,omitempty"`            // 参数
	ActionLog       map[string]any  `json:"actionlog,omitempty"`         // 点击日志
	Image           string          `json:"image,omitempty"`             // 左边图片icon
	ImageDark       string          `json:"imageDark,omitempty"`         // 左边图片icon暗黑
	Border          *C182Border     `json:"border,omitempty"`            // 边框
	Background      *C182Background `json:"background,omitempty"`        // 背景
	Radius          int             `json:"radius,omitempty"`            // 圆角
}

// C182Border 边框
type C182Border struct {
	Color     string `json:"color,omitempty"`     // 边框颜色
	ColorDark string `json:"colorDark,omitempty"` // 边框暗黑颜色
	Width     int    `json:"width,omitempty"`     // 边框宽度
}

// C182Background 背景
type C182Background struct {
	HightLightColor     string `json:"hightLightColor,omitempty"`     // 选中背景色
	HightLightDarkColor string `json:"hightLightDarkColor,omitempty"` // 选中暗黑背景色
}

// C182TrendItem 潮流项
type C182TrendItem struct {
	Title     []*C182RichText `json:"title,omitempty"`     // 底部文字
	Icon      string          `json:"icon,omitempty"`      // 主图片
	Scheme    string          `json:"scheme,omitempty"`    // 跳转链接
	ActionLog map[string]any  `json:"actionlog,omitempty"` // 点击日志
}

// C182MatchItem 赛程队伍项
type C182MatchItem struct {
	Pic     string          `json:"pic,omitempty"`     // 队伍图标
	PicDark string          `json:"picDark,omitempty"` // 队伍图标暗黑
	Name    []*C182RichText `json:"name,omitempty"`    // 队名
}

// NewCard182 创建Card182实例
func NewCard182() *C182 {
	return &C182{BaseItem: BaseItem{Base: Base{CardType: 182}}}
}
