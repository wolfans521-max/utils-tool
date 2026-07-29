package card

// C116 今日热点事件卡片
type C116 struct {
	BaseItem
	PageSize     int               `json:"page_size,omitempty"`     // 每次显示的事件数量
	Events       []*C116Event      `json:"events,omitempty"`        // 事件列表
	ActionButton *C116ActionButton `json:"action_button,omitempty"` // 换一换按钮，如果不下发则不显示
}

// C116Event 事件
type C116Event struct {
	Id        string         `json:"id,omitempty"`
	Title     string         `json:"title,omitempty"`
	Scheme    string         `json:"scheme,omitempty"`
	ActionLog map[string]any `json:"actionlog,omitempty"`
	Color     string         `json:"color,omitempty"`
}

// C116ActionButton 换一换按钮
type C116ActionButton struct {
	Text   string `json:"text,omitempty"`
	Scheme string `json:"scheme,omitempty"`
}

// NewCard116 创建Card116实例
func NewCard116() *C116 {
	return &C116{BaseItem: BaseItem{Base: Base{CardType: 116}}}
}
