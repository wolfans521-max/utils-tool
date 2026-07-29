package card

// C9 微博卡片
type C9 struct {
	BaseItem
	Rating             string          `json:"rating,omitempty"`     // 显示星级评分（791 起生效）
	WeiboNeed          string          `json:"weibo_need,omitempty"` // 如果设置为mblog，则只要在mblog字段中带mid字段即可
	Mblog              map[string]any  `json:"mblog,omitempty"`      // 微博结构
	ShowType           int8            `json:"show_type,omitempty"`  // 0简单数据，1详细信息，2热门微博（无图），3热门微博（有图），4新鲜事
	Hidebtns           int8            `json:"hidebtns,omitempty"`   // 客户端是否显示微博下面的三个按钮
	Openurl            string          `json:"openurl,omitempty"`
	MblogButtons       []*MblogButtons `json:"mblog_buttons,omitempty"`
	ContainerColor     string          `json:"container_color,omitempty"`
	ContainerColorDark string          `json:"container_color_dark,omitempty"`
	DisplayArrow       int8            `json:"display_arrow,omitempty"`
	Scheme             string          `json:"scheme,omitempty"`
}

type MblogButtons struct {
	Type      string `json:"type,omitempty"`
	Name      string `json:"name,omitempty"`
	Pic       string `json:"pic,omitempty"`
	ActionLog string `json:"actionlog,omitempty"`
}

// NewCard9 创建Card9实例
func NewCard9() *C9 {
	return &C9{BaseItem: BaseItem{Base: Base{CardType: 9}}}
}
