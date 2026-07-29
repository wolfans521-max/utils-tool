package card

// C184 自动换行可折叠卡片
// item自动换行、可折叠 card
type C184 struct {
	BaseItem
	Fold       bool        `json:"fold,omitempty"`        // 是否折叠
	MinRows    int         `json:"min_rows,omitempty"`    // 折叠后显示的行数
	FoldText   string      `json:"fold_text,omitempty"`   // 折叠按钮文案
	UnfoldText string      `json:"unfold_text,omitempty"` // 展开按钮文案
	Items      []*C184Item `json:"items,omitempty"`       // item数组
	OpenUrl    string      `json:"openurl,omitempty"`     // 打开链接
}

// C184Item 项目
type C184Item struct {
	Content         string         `json:"content,omitempty"`           // item文案
	Scheme          string         `json:"scheme,omitempty"`            // item scheme
	BgColor         string         `json:"bg_color,omitempty"`          // item背景颜色
	BgColorDark     string         `json:"bg_color_dark,omitempty"`     // item背景颜色_深色模式
	TextColor       string         `json:"text_color,omitempty"`        // item字体颜色
	TextColorDark   string         `json:"text_color_dark,omitempty"`   // item字体颜色_深色模式
	RightIcon       string         `json:"right_icon,omitempty"`        // item右侧icon
	RightIconDark   string         `json:"right_icon_dark,omitempty"`   // item右侧icon_深色模式
	BorderColor     string         `json:"border_color,omitempty"`      // 边框色
	BorderColorDark string         `json:"border_color_dark,omitempty"` // 边框色(暗黑)
	ActionType      int            `json:"action_type,omitempty"`       // 事件类型(1表示刷新当前list)
	ActionLog       map[string]any `json:"actionlog,omitempty"`         // action log
}

// NewCard184 创建Card184实例
func NewCard184() *C184 {
	return &C184{BaseItem: BaseItem{Base: Base{CardType: 184}}}
}
