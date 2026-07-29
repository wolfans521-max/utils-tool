package card

// C226 Card226 个性化开关卡片
type C226 struct {
	BaseItem
	CardPadding         map[string]int    `json:"card_padding,omitempty"`
	SwitchSubtype       int               `json:"switch_subtype,omitempty"`
	Path                string            `json:"path,omitempty"`
	BackgroundColor     string            `json:"backgroundColor,omitempty"`
	BackgroundColorDark string            `json:"backgroundColorDark,omitempty"`
	CardHeight          int               `json:"cardHeight,omitempty"`
	SwitchLayout        *C226SwitchLayout `json:"switchLayout,omitempty"`
}

// C226SwitchLayout 开关布局
type C226SwitchLayout struct {
	SwitchLayoutBg     string            `json:"switchLayoutBg,omitempty"`
	SwitchLayoutBgDark string            `json:"switchLayoutBgDark,omitempty"`
	Height             int               `json:"height,omitempty"`
	Width              int               `json:"width,omitempty"`
	Radius             int               `json:"radius,omitempty"`
	LeftButton         *C226SwitchButton `json:"leftButton,omitempty"`
	RightButton        *C226SwitchButton `json:"rightButton,omitempty"`
}

// C226SwitchButton 开关按钮
type C226SwitchButton struct {
	SelectedBgColor       string         `json:"selectedBgColor,omitempty"`
	SelectedBgColorDark   string         `json:"selectedBgColorDark,omitempty"`
	SelectedTextColor     string         `json:"selectedTextColor,omitempty"`
	SelectedTextColorDark string         `json:"selectedTextColorDark,omitempty"`
	DefaultBgColor        string         `json:"defaultBgColor,omitempty"`
	DefaultBgColorDark    string         `json:"defaultBgColorDark,omitempty"`
	DefaultTextColor      string         `json:"defaultTextColor,omitempty"`
	DefaultTextColorDark  string         `json:"defaultTextColorDark,omitempty"`
	SelectedTextSize      int            `json:"selectedTextSize,omitempty"`
	DefaultTextSize       int            `json:"defaultTextSize,omitempty"`
	SwitchBgColor         string         `json:"switchBgColor,omitempty"`
	SwitchBgColorDark     string         `json:"switchBgColorDark,omitempty"`
	ActionLog             map[string]any `json:"actionlog,omitempty"`
	Text                  string         `json:"text,omitempty"`
	Selected              bool           `json:"selected,omitempty"`
	SortType              string         `json:"sortType,omitempty"`
	Padding               map[string]int `json:"padding,omitempty"`
	TitleText             *C226TitleText `json:"titleText,omitempty"`
	FeedTopToast          int            `json:"feed_top_toast,omitempty"`
}

// C226TitleText 标题文本
type C226TitleText struct {
	Text          string         `json:"text,omitempty"`
	TextColor     string         `json:"textColor,omitempty"`
	TextColorDark string         `json:"textColorDark,omitempty"`
	TextSize      int            `json:"textSize,omitempty"`
	Margin        map[string]int `json:"margin,omitempty"`
}

// NewCard226 创建Card226实例
func NewCard226() *C226 {
	return &C226{BaseItem: BaseItem{Base: Base{CardType: 226}}}
}
