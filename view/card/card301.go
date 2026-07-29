package card

// C301ClickObj 神评卡片点击对象（评论区/博文区）
type C301ClickObj struct {
	Scheme    string         `json:"scheme"`
	ActionLog map[string]any `json:"action_log,omitempty"`
}

// C301Image 神评卡片图片对象（序号图/头像图）
type C301Image struct {
	Url             string  `json:"url,omitempty"`
	UrlDark         string  `json:"url_dark,omitempty"`
	Width           float64 `json:"width,omitempty"`
	Height          float64 `json:"height,omitempty"`
	RadiusCorner    float64 `json:"radius_corner,omitempty"`
	BorderWidth     float64 `json:"border_width,omitempty"`
	BorderColor     string  `json:"border_color,omitempty"`
	BorderColorDark string  `json:"border_color_dark,omitempty"`
}

// C301Text 神评卡片文本对象（用户名）
type C301Text struct {
	Text          string `json:"text,omitempty"`
	TextSize      int    `json:"text_size,omitempty"`
	TextColor     string `json:"text_color,omitempty"`
	TextColorDark string `json:"text_color_dark,omitempty"`
}

// C301UserObj 神评卡片用户对象（序号图 + 头像 + 用户名）
type C301UserObj struct {
	IndexImage *C301Image `json:"index_image,omitempty"` // 序号图
	AvataImage *C301Image `json:"avata_image,omitempty"` // 用户头像
	Name       *C301Text  `json:"name,omitempty"`        // 用户名
}

// C301 神评评论卡片
// 用于发现页神评流，展示单条神评内容（用户信息 + 关联博文）。
type C301 struct {
	BaseItem
	UserObj         *C301UserObj   `json:"user_obj,omitempty"`          // 用户信息（序号图/头像/用户名）
	CommentClickObj *C301ClickObj  `json:"comment_click_obj,omitempty"` // 评论区点击
	BlogClickObj    *C301ClickObj  `json:"blog_click_obj,omitempty"`    // 博文区点击
	Comment         map[string]any `json:"comment,omitempty"`           // 评论数据(comments_show_batch 结果), 其 status 字段为博文 showbatch 结果
	ItemExt         *ItemExt       `json:"itemExt,omitempty"`           // item 扩展（markInterval 等，用于 cursor 翻页计数）
}

// NewCard301 创建Card301实例。
func NewCard301() *C301 {
	return &C301{BaseItem: BaseItem{Base: Base{CardType: 301}}}
}
