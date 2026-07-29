package card

// C119TextItem 视窗文本项（用于query_config）
type C119TextItem struct {
	ActionLog map[string]any `json:"actionlog,omitempty"`  // 点击日志
	LeftIcon  string         `json:"left_icon,omitempty"`  // 左侧图标
	RightIcon string         `json:"right_icon,omitempty"` // 右侧图标（箭头）
	Text      string         `json:"text,omitempty"`       // 文本内容
	TextSize  int            `json:"text_size,omitempty"`  // 文本大小
	TextColor string         `json:"text_color,omitempty"` // 文本颜色
	Scheme    string         `json:"scheme,omitempty"`     // 跳转链接
}

type C119SubItem struct {
	Scheme            string                     `json:"scheme,omitempty"`
	ItemId            string                     `json:"itemid,omitempty"`
	ContentType       int                        `json:"content_type"`
	Pic               string                     `json:"pic,omitempty"`
	NewStyle          int                        `json:"newStyle,omitempty"`
	Title             string                     `json:"title,omitempty"`
	IsBigPic          int                        `json:"is_big_pic,omitempty"`
	IconUrl           string                     `json:"icon_url"`
	Desc              string                     `json:"desc"`
	ChannelBackground string                     `json:"channel_background,omitempty"`
	ActionLog         map[string]any             `json:"actionlog,omitempty"`
	Promotion         map[string]any             `json:"promotion,omitempty"`
	CornerMarkData    *C119SubItemCornerMarkData `json:"corner_mark_data,omitempty"`
	AdVideoInfo       map[string]any             `json:"ad_videoinfo,omitempty"`
	IsLiveContext     bool                       `json:"is_live_context,omitempty"`
	RightType         string                     `json:"right_type,omitempty"`
	JumpUrl           string                     `json:"jumpUrl,omitempty"`
	LiveInfo          map[string]any             `json:"live_info,omitempty"`
	TextItems         []C119TextItem             `json:"text_items,omitempty"`     // 文本项列表（query_config）
	IsGroupStyle      bool                       `json:"is_group_style,omitempty"` // 三合一样式标识
	GroupInfo         *C119GroupInfo             `json:"group_info,omitempty"`     // 三合一布局信息
}

// C119GroupInfo 三合一布局信息
type C119GroupInfo struct {
	BackgroundColor     string        `json:"background_color,omitempty"`
	BackgroundColorDark string        `json:"background_color_dark,omitempty"`
	HorizontalSpaceVal  int           `json:"horizontal_space_val,omitempty"`
	LeftObj             *C119LeftObj  `json:"left_obj,omitempty"`
	RightObj            *C119RightObj `json:"right_obj,omitempty"`
}

// C119LeftObj 三合一左侧对象
type C119LeftObj struct {
	Width       int              `json:"width,omitempty"`
	Height      int              `json:"height,omitempty"`
	ContentType string           `json:"content_type,omitempty"` // pic 或 video
	PicImageUrl string           `json:"pic_image_url,omitempty"`
	Scheme      string           `json:"scheme,omitempty"`
	TagImageObj *C119TagImageObj `json:"tag_image_obj,omitempty"`
}

// C119TagImageObj 左上角标识
type C119TagImageObj struct {
	Url         string `json:"url,omitempty"`
	UrlDark     string `json:"url_dark,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	PaddingTop  int    `json:"padding_top,omitempty"`
	PaddingLeft int    `json:"padding_left,omitempty"`
}

// C119RightObj 三合一右侧对象
type C119RightObj struct {
	VerticalSpaceVal int             `json:"vertical_space_val,omitempty"`
	Items            []C119RightItem `json:"items,omitempty"`
}

// C119RightItem 三合一右侧项
type C119RightItem struct {
	ImageUrl  string         `json:"image_url,omitempty"`
	Scheme    string         `json:"scheme,omitempty"`
	ActionLog map[string]any `json:"actionlog,omitempty"`
}

type C119SubItemCornerMarkData struct {
	BorderWidth     int    `json:"border_width"`
	TitleSize       int    `json:"title_size,omitempty"`
	Title           string `json:"title,omitempty"`
	TitleColor      string `json:"title_color,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
	BorderColor     string `json:"border_color,omitempty"`
	Padding         []int  `json:"padding,omitempty"`
	Radius          []int  `json:"radius,omitempty"`
}

type C119 struct {
	BaseItem
	SubItem []C119SubItem `json:"sub_item,omitempty"`
}

func NewCard119() *C119 {
	return &C119{BaseItem: BaseItem{Base: Base{CardType: 119}}}
}
