package card

// C261Content 播报card内容项
type C261Content struct {
	Type    string         `json:"type,omitempty"`    // 内容类型: text/icon
	Style   *C261TextStyle `json:"style,omitempty"`   // 样式
	Content string         `json:"content,omitempty"` // 文本内容（type=text时使用）
}

// C261TextStyle 播报card文本样式
type C261TextStyle struct {
	TextColor     string `json:"textColor,omitempty"`     // 文本颜色
	TextColorDark string `json:"textColorDark,omitempty"` // 暗黑模式文本颜色
	TextSize      int    `json:"textSize,omitempty"`      // 文本大小
	Bold          int    `json:"bold,omitempty"`          // 是否加粗 1=加粗
	Width         int    `json:"width,omitempty"`         // 宽度（type=icon时使用）
}

// C261 夺金播报card
// 用于展示夺金播报广告模块，支持动画、背景图、富文本内容
type C261 struct {
	BaseItem
	CardHeight      int            `json:"card_height,omitempty"`        // card高度
	ShowAnimation   int            `json:"show_animation,omitempty"`     // 是否展示动画 1=展示
	LeftImg         string         `json:"left_img,omitempty"`           // 左侧图片
	LeftImgDark     string         `json:"left_img_dark,omitempty"`      // 暗黑模式左侧图片
	RightLottie     string         `json:"right_lottie,omitempty"`       // 右侧lottie动画
	RightLottieDark string         `json:"right_lottie_dark,omitempty"`  // 暗黑模式右侧lottie
	LottieBkImg     string         `json:"lottie_bk_img,omitempty"`      // lottie背景图
	LottieBkImgDark string         `json:"lottie_bk_img_dark,omitempty"` // 暗黑模式lottie背景图
	BgImg           string         `json:"bg_img,omitempty"`             // 背景图
	BgImgDark       string         `json:"bg_img_dark,omitempty"`        // 暗黑模式背景图
	Contents        []*C261Content `json:"contents,omitempty"`           // 富文本内容列表
}

// NewCard261 创建Card261实例
func NewCard261() *C261 {
	return &C261{BaseItem: BaseItem{Base: Base{CardType: 261}}}
}
