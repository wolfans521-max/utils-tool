package card

// C60 图文卡片
// 支持标题、描述、图片和按钮
type C60 struct {
	BaseItem
	CardCommTitle string    `json:"card_comm_title,omitempty"` // 卡片通用标题
	TitleSub      string    `json:"title_sub,omitempty"`       // 第一行文案
	Desc1         string    `json:"desc1,omitempty"`           // 第二行文案（支持星星或者半星）
	Desc2         string    `json:"desc2,omitempty"`           // 第三行文案（支持星星或者半星）
	Buttons       []*Button `json:"buttons,omitempty"`         // 右侧button
	Pic           string    `json:"pic,omitempty"`             // 图片
	PicWidth      int       `json:"pic_width,omitempty"`       // 图片宽度比例
	PicHeight     int       `json:"pic_height,omitempty"`      // 图片高度比例
	UnreadId      string    `json:"unread_id,omitempty"`       // 未读ID
	OpenUrl       string    `json:"openurl,omitempty"`         // 打开链接
	CornerRadius  int       `json:"corner_radius,omitempty"`   // 图片圆角度数
}

// NewCard60 创建Card60实例
func NewCard60() *C60 {
	return &C60{BaseItem: BaseItem{Base: Base{CardType: 60}}}
}
