package card

// C262BannerStyle banner样式配置
type C262BannerStyle struct {
	Ratio        float64 `json:"ratio,omitempty"`         // 比例
	LeftMargin   int     `json:"left_margin,omitempty"`   // 左边距
	RightMargin  int     `json:"right_margin,omitempty"`  // 右边距
	TopMargin    int     `json:"top_margin,omitempty"`    // 上边距
	BottomMargin int     `json:"bottom_margin,omitempty"` // 下边距
}

// C262MoreFrameStyle 多帧样式配置
type C262MoreFrameStyle struct {
	Ratio        float64 `json:"ratio,omitempty"`         // 比例
	LeftMargin   int     `json:"left_margin,omitempty"`   // 左边距
	RightMargin  int     `json:"right_margin,omitempty"`  // 右边距
	TopMargin    int     `json:"top_margin,omitempty"`    // 上边距
	BottomMargin int     `json:"bottom_margin,omitempty"` // 下边距
}

// C262MoreFrame 多帧内容
type C262MoreFrame struct {
	ItemId        string `json:"itemid,omitempty"`       // 项目ID
	LeftPadding   int    `json:"left_padding,omitempty"` // 左内边距
	Items         []any  `json:"items,omitempty"`        // 帧项目列表
	AspectRatio   string `json:"aspect_ratio,omitempty"` // 宽高比
	BottomPadding int    `json:"bottom_padding,omitempty"`
	TopPadding    int    `json:"top_padding,omitempty"`
	RightPadding  int    `json:"right_padding,omitempty"`
	CardType      int    `json:"card_type,omitempty"`
	ShowStyle     int    `json:"show_style,omitempty"`
	HideBottom    int    `json:"hide_bottom,omitempty"`
	HideTopPad    int    `json:"hide_top_padding,omitempty"`
	VideoMargin   []int  `json:"video_margin,omitempty"`
	VideoPadding  []int  `json:"video_padding,omitempty"`
	CardRadius    int    `json:"card_radius,omitempty"`
	VideoCorner   int    `json:"video_corner,omitempty"`
}

// C262Banner banner项
type C262Banner struct {
	Pic       string `json:"pic,omitempty"`       // 图片
	Scheme    string `json:"scheme,omitempty"`    // 跳转链接
	ActionLog any    `json:"actionlog,omitempty"` // 行为日志
	Cleaned   bool   `json:"cleaned,omitempty"`   // 是否已清理
}

// C262Grid 宫格项
type C262Grid struct {
	Pic       string `json:"pic,omitempty"`       // 图片
	Text      string `json:"text,omitempty"`      // 文本
	Scheme    string `json:"scheme,omitempty"`    // 跳转链接
	ActionLog any    `json:"actionlog,omitempty"` // 行为日志
	Cleaned   bool   `json:"cleaned,omitempty"`   // 是否已清理
}

// C262 搜索业务广告Card
// 基于196改造，支持一卡多帧，支持图片轮播和一卡多帧的margin配置
type C262 struct {
	BaseItem
	PaddingTop        int                 `json:"padding_top,omitempty"`         // 上内边距
	PaddingBottom     int                 `json:"padding_bottom,omitempty"`      // 下内边距
	BannerBorderColor string              `json:"banner_border_color,omitempty"` // banner边框颜色
	TitlePicWidth     int                 `json:"title_pic_width,omitempty"`     // 标题图片宽度
	BackgroundPic     string              `json:"background_pic,omitempty"`      // 背景图片
	TitlePic          string              `json:"title_pic,omitempty"`           // 标题图片
	LogoPic           string              `json:"logo_pic,omitempty"`            // logo图片
	BannerStyle       *C262BannerStyle    `json:"banner_style,omitempty"`        // banner样式
	MoreFrameStyle    *C262MoreFrameStyle `json:"more_frame_style,omitempty"`    // 多帧样式
	MoreFrame         *C262MoreFrame      `json:"more_frame,omitempty"`          // 多帧内容
	Banner            []*C262Banner       `json:"banner,omitempty"`              // banner列表
	GridColumn        int                 `json:"grid_column,omitempty"`         // 宫格列数
	Grid              []*C262Grid         `json:"grid,omitempty"`                // 宫格列表
	OpenUrl           string              `json:"openurl,omitempty"`             // 打开链接
}

// NewCard262 创建Card262实例
func NewCard262() *C262 {
	return &C262{BaseItem: BaseItem{Base: Base{CardType: 262}}}
}
