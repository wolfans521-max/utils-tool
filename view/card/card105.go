package card

// C105 功能型复访入口卡片
// 支持上图下文样式和数字样式
type C105 struct {
	BaseItem
	ItemSizeType int8        `json:"item_size_type,omitempty"` // 由小到大分为0,1,2
	Items        []*C105Item `json:"items,omitempty"`          // 入口列表
	OpenUrl      string      `json:"openurl,omitempty"`        // 打开链接
}

// C105Item 入口项
type C105Item struct {
	ItemType   string         `json:"item_type,omitempty"`   // 0表示上图下文样式，1表示上面是文字的样式
	Pic        string         `json:"pic,omitempty"`         // 图片
	Title      string         `json:"title,omitempty"`       // 数字标题
	TitleColor string         `json:"title_color,omitempty"` // 标题颜色
	DescMain   string         `json:"desc_main,omitempty"`   // 主描述
	DescSub    string         `json:"desc_sub,omitempty"`    // 副描述
	Scheme     string         `json:"scheme,omitempty"`      // 跳转链接
	UnreadId   string         `json:"unread_id,omitempty"`   // 未读ID
	ActionLog  map[string]any `json:"action_log,omitempty"`  // 打点日志
}

// NewCard105 创建Card105实例
func NewCard105() *C105 {
	return &C105{BaseItem: BaseItem{Base: Base{CardType: 105}}}
}
