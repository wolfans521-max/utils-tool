package card

// C280TopContent 内部卡片上半部分内容
type C280TopContent struct {
	Pic            string `json:"pic,omitempty"`              // 内部卡片图片
	PicDark        string `json:"pic_dark,omitempty"`         // 暗黑模式图片
	PlayIcon       string `json:"play_icon,omitempty"`        // 播放按钮图片
	PlayIconDark   string `json:"play_icon_dark,omitempty"`   // 暗黑模式播放按钮
	PlayIconWidth  int    `json:"play_icon_width,omitempty"`  // 播放按钮宽
	PlayIconHeight int    `json:"play_icon_height,omitempty"` // 播放按钮高
	RadiusCorner   []int  `json:"radius_corner,omitempty"`    // 上半部分圆角 [左上,右上,左下,右下]
}

// C280BottomContent 内部卡片下半部分内容
type C280BottomContent struct {
	Text                 string `json:"text,omitempty"`                   // 文字
	TextAlignLeft        bool   `json:"text_align_left,omitempty"`        // 是否居左
	TextSize             int    `json:"text_size,omitempty"`              // 文字大小
	TextColor            string `json:"text_color,omitempty"`             // 文字颜色
	TextColorDark        string `json:"text_color_dark,omitempty"`        // 暗黑模式文字颜色
	TextMargin           []int  `json:"text_margin,omitempty"`            // 文字左右边距 [左,0,右,0]
	BottomHeight         int    `json:"bottom_height,omitempty"`          // 文字部分高度
	BottomBackground     string `json:"bottom_background,omitempty"`      // 文字部分背景色
	BottomBackgroundDark string `json:"bottom_background_dark,omitempty"` // 暗黑模式背景色
	RadiusCorner         []int  `json:"radius_corner,omitempty"`          // 文字部分圆角 [左上,右上,左下,右下]
}

// C280Item 内部卡片项
type C280Item struct {
	Scheme            string             `json:"scheme,omitempty"`             // 跳转scheme
	TopContent        *C280TopContent    `json:"top_content,omitempty"`        // 上半部分内容
	BottomContent     *C280BottomContent `json:"bottom_content,omitempty"`     // 下半部分内容
	ExposureActionLog any                `json:"exposure_actionlog,omitempty"` // 曝光日志
	ClickActionLog    any                `json:"click_actionlog,omitempty"`    // 点击日志
	Promotion         any                `json:"promotion,omitempty"`          // 三方监控日志
}

// C280 横滑卡片
// 支持内部卡片横滑展示，可配置图片、文字、圆角等
type C280 struct {
	BaseItem
	ItemInterval         int         `json:"item_interval,omitempty"`           // 内部卡片间距
	PicWidth             int         `json:"pic_width,omitempty"`               // 内部卡片宽度
	PicWidthRatio        int         `json:"pic_width_ratio,omitempty"`         // 内部卡片宽比例
	PicHeightRatio       int         `json:"pic_height_ratio,omitempty"`        // 内部卡片高比例
	PageLeftRightPadding int         `json:"page_left_right_padding,omitempty"` // 横滑card左右边距
	CopyMargin           []int       `json:"copy_margin,omitempty"`             // 外边距 iOS专用
	CopyPadding          []int       `json:"copy_padding,omitempty"`            // 内边距 iOS专用
	Items                []*C280Item `json:"items,omitempty"`                   // 内部卡片列表
}

// NewCard280 创建Card280实例
func NewCard280() *C280 {
	return &C280{BaseItem: BaseItem{Base: Base{CardType: 280}}}
}
