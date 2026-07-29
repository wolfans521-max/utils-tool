package card

type C17 struct {
	BaseItem
	Col                  int8        `json:"col,omitempty"`
	Group                []*C17Texts `json:"group,omitempty"`
	Bottom               *C17Bottom  `json:"bottom,omitempty"`
	CoverHeight          int         `json:"cover_height,omitempty"`
	CoverContent         string      `json:"cover_content,omitempty"`
	Top                  *C17Top     `json:"top,omitempty"`
	TitleBold            int8        `json:"title_bold,omitempty"`
	TitleColor           string      `json:"title_color,omitempty"`
	TitleColorDark       string      `json:"title_color_dark,omitempty"`
	GroupItemStartOffset int8        `json:"group_item_start_offset,omitempty"`
}

type C17Texts struct {
	TitleSub      string         `json:"title_sub,omitempty"`
	Icon          string         `json:"icon,omitempty"`
	Scheme        string         `json:"scheme,omitempty"`
	TitleSubColor string         `json:"title_sub_color,omitempty"`
	ActionLog     map[string]any `json:"action_log,omitempty"`
	IconWidth     int            `json:"icon_width,omitempty"`
	IconHeight    int            `json:"icon_height,omitempty"`
}

type C17Top struct {
	RightScheme           string         `json:"right_scheme,omitempty"`
	TitleSize             int            `json:"title_size,omitempty"`
	Title                 string         `json:"title,omitempty"`
	TopMargin             int            `json:"top_margin,omitempty"`
	Icon                  string         `json:"icon,omitempty"`
	ActLog                map[string]any `json:"actlog,omitempty"`
	BottomMargin          int            `json:"bottom_margin,omitempty"`
	CloseStatus           int            `json:"close_status,omitempty"`
	IconSize              int            `json:"icon_size,omitempty"`
	ItemHeight            int            `json:"item_height,omitempty"`
	TitleLeftMargin       int            `json:"title_left_margin,omitempty"`
	Desc                  string         `json:"desc,omitempty"`
	Refresh               *C17TopRefresh `json:"refresh,omitempty"`
	HideTopLine           bool           `json:"hide_top_line,omitempty"`
	TitleLeftMarginOffset int            `json:"title_left_margin_offset,omitempty"`
}

type C17TopRefresh struct {
	Icon   string `json:"icon,omitempty"`
	Title  string `json:"title,omitempty"`
	Scheme string `json:"scheme,omitempty"`
	ActLog string `json:"actlog,omitempty"`
}

type C17Bottom struct {
	IconSize     int            `json:"icon_size,omitempty"`
	TopMargin    int            `json:"top_margin,omitempty"`
	BottomMargin int            `json:"bottom_margin,omitempty"`
	ActLog       map[string]any `json:"actlog,omitempty"`
	Title        string         `json:"title,omitempty"`
	Icon         string         `json:"icon,omitempty"`
	TitleColor   string         `json:"title_color,omitempty"`
	Scheme       string         `json:"scheme,omitempty"`
	TitleSize    int            `json:"title_size,omitempty"`
}

func NewCard17() *C17 {
	return &C17{BaseItem: BaseItem{Base: Base{CardType: 17}}}
}
