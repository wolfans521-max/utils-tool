package card

// C246BannerStyle banner样式配置
type C246BannerStyle struct {
	Ratio        float64 `json:"ratio,omitempty"`         // 比例
	LeftMargin   int     `json:"left_margin,omitempty"`   // 左边距
	RightMargin  int     `json:"right_margin,omitempty"`  // 右边距
	TopMargin    int     `json:"top_margin,omitempty"`    // 上边距
	BottomMargin int     `json:"bottom_margin,omitempty"` // 下边距
}

// C246MoreFrameStyle 多帧样式配置
type C246MoreFrameStyle struct {
	Ratio        float64 `json:"ratio,omitempty"`         // 比例
	LeftMargin   int     `json:"left_margin,omitempty"`   // 左边距
	RightMargin  int     `json:"right_margin,omitempty"`  // 右边距
	TopMargin    int     `json:"top_margin,omitempty"`    // 上边距
	BottomMargin int     `json:"bottom_margin,omitempty"` // 下边距
}

// C246MoreFrame 多帧内容
type C246MoreFrame struct {
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

// C246Banner banner项
type C246Banner struct {
	Pic       string `json:"pic,omitempty"`       // 图片
	Scheme    string `json:"scheme,omitempty"`    // 跳转链接
	ActionLog any    `json:"actionlog,omitempty"` // 行为日志
	Cleaned   bool   `json:"cleaned,omitempty"`   // 是否已清理
}

// C246Grid 宫格项
type C246Grid struct {
	Pic       string `json:"pic,omitempty"`       // 图片
	Text      string `json:"text,omitempty"`      // 文本
	Scheme    string `json:"scheme,omitempty"`    // 跳转链接
	ActionLog any    `json:"actionlog,omitempty"` // 行为日志
	Cleaned   bool   `json:"cleaned,omitempty"`   // 是否已清理
}

// C246 搜索业务广告Card
// 基于196改造，支持一卡多帧，支持图片轮播和一卡多帧的margin配置
type C246 struct {
	BaseItem
	PaddingTop        int                 `json:"padding_top,omitempty"`         // 上内边距
	PaddingBottom     int                 `json:"padding_bottom,omitempty"`      // 下内边距
	BannerBorderColor string              `json:"banner_border_color,omitempty"` // banner边框颜色
	TitlePicWidth     int                 `json:"title_pic_width,omitempty"`     // 标题图片宽度
	BackgroundPic     string              `json:"background_pic,omitempty"`      // 背景图片
	TitlePic          string              `json:"title_pic,omitempty"`           // 标题图片
	LogoPic           string              `json:"logo_pic,omitempty"`            // logo图片
	BannerStyle       *C246BannerStyle    `json:"banner_style,omitempty"`        // banner样式
	MoreFrameStyle    *C246MoreFrameStyle `json:"more_frame_style,omitempty"`    // 多帧样式
	MoreFrame         *C246MoreFrame      `json:"more_frame,omitempty"`          // 多帧内容
	Banner            []*C246Banner       `json:"banner,omitempty"`              // banner列表
	GridColumn        int                 `json:"grid_column,omitempty"`         // 宫格列数
	Grid              []*C246Grid         `json:"grid,omitempty"`                // 宫格列表
	OpenUrl           string              `json:"openurl,omitempty"`             // 打开链接
}

// NewCard246 创建Card246实例
func NewCard246() *C246 {
	return &C246{BaseItem: BaseItem{Base: Base{CardType: 246}}}
}
