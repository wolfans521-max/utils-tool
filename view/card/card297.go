package card

// C297Color 渐变背景色
type C297Color struct {
	StartColor     string `json:"start_color,omitempty"`
	EndColor       string `json:"end_color,omitempty"`
	StartColorKey  string `json:"start_color_key,omitempty"`
	EndColorKey    string `json:"end_color_key,omitempty"`
	StartColorDark string `json:"start_color_dark,omitempty"`
	EndColorDark   string `json:"end_color_dark,omitempty"`
}

// C297Capsule 胶囊（引导词）
type C297Capsule struct {
	Text                string         `json:"text,omitempty"`
	TextColor           string         `json:"text_color,omitempty"`
	TextColorKey        string         `json:"text_color_key,omitempty"`
	TextColorDark       string         `json:"text_color_dark,omitempty"`
	BackgroundColor     string         `json:"background_color,omitempty"`
	BackgroundColorKey  string         `json:"background_color_key,omitempty"`
	BackgroundColorDark string         `json:"background_color_dark,omitempty"`
	BorderColor         string         `json:"border_color,omitempty"`
	BorderColorKey      string         `json:"border_color_key,omitempty"`
	BorderColorDark     string         `json:"border_color_dark,omitempty"`
	Scheme              string         `json:"scheme,omitempty"`
	ActionLog           map[string]any `json:"actionlog,omitempty"`
}

// C297 搜索引导词胶囊卡片（客户端原生卡片）
type C297 struct {
	BaseItem
	BackgroundColor *C297Color     `json:"background_color,omitempty"`
	SlideActionLog  map[string]any `json:"slide_action_log,omitempty"`
	Capsules        []*C297Capsule `json:"capsules,omitempty"`
	// TopPadding 顶部内边距，bigday 场景下发（值为 0），非 bigday 不下发该 key
	TopPadding *int `json:"top_padding,omitempty"`
}

func NewCard297() *C297 {
	return &C297{BaseItem: BaseItem{Base: Base{CardType: 297}}}
}
