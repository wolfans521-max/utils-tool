package card

// C179 我的超话卡片
// 用于我的超话展示
type C179 struct {
	BaseItem
	BgImg      string      `json:"bg_img,omitempty"`      // 背景图
	Header     *C179Header `json:"header,omitempty"`      // 头部
	ItemsStyle int8        `json:"items_style,omitempty"` // 1表示支持副标题；0表示不支持
	Items      []*C179Item `json:"items,omitempty"`       // 底部图片数组
	OpenUrl    string      `json:"openurl,omitempty"`     // 打开链接
}

// C179Header 头部
type C179Header struct {
	Scheme         string         `json:"scheme,omitempty"`           // 跳转scheme
	Title          string         `json:"title,omitempty"`            // 标题
	TitleColor     string         `json:"title_color,omitempty"`      // 标题颜色
	TitleColorDark string         `json:"title_color_dark,omitempty"` // 标题颜色-暗黑
	TitleIcon      string         `json:"title_icon,omitempty"`       // 标题后icon
	Desc           string         `json:"desc,omitempty"`             // 右侧描述
	DescColor      string         `json:"desc_color,omitempty"`       // 右侧描述颜色
	DescColorDark  string         `json:"desc_color_dark,omitempty"`  // 右侧描述颜色-暗黑
	DescScheme     string         `json:"desc_scheme,omitempty"`      // 右侧文案点击scheme
	ArrowImg       string         `json:"arrow_img,omitempty"`        // 右侧箭头图片
	ActionLog      map[string]any `json:"actionlog,omitempty"`        // 打点日志
}

// C179Item 项目
type C179Item struct {
	Scheme          string         `json:"scheme,omitempty"`            // 跳转scheme
	Image           string         `json:"image,omitempty"`             // 超话头像
	BadgeImg        string         `json:"badge_img,omitempty"`         // 角标图片
	BadgeImgTop     string         `json:"badge_img_top,omitempty"`     // 图片右上角icon
	Title           string         `json:"title,omitempty"`             // 超话名称
	TitleColor      string         `json:"title_color,omitempty"`       // 超话名称颜色
	TitleColorDark  string         `json:"title_color_dark,omitempty"`  // 超话名称颜色-暗黑
	Detail          string         `json:"detail,omitempty"`            // 详情
	DetailColor     string         `json:"detail_color,omitempty"`      // 详情颜色
	DetailColorDark string         `json:"detail_color_dark,omitempty"` // 详情颜色-暗黑
	Desc            string         `json:"desc,omitempty"`              // 右侧描述
	DescColor       string         `json:"desc_color,omitempty"`        // 右侧描述颜色
	DescColorDark   string         `json:"desc_color_dark,omitempty"`   // 右侧描述颜色-暗黑
	ArrowImg        string         `json:"arrow_img,omitempty"`         // 右侧箭头图片
	BadgeImgType    int8           `json:"badge_img_type,omitempty"`    // 角标类型 1 待签
	ContainerId     string         `json:"containerid,omitempty"`       // 超话ID
	ActionLog       map[string]any `json:"actionlog,omitempty"`         // 打点日志
}

// NewCard179 创建Card179实例
func NewCard179() *C179 {
	return &C179{BaseItem: BaseItem{Base: Base{CardType: 179}}}
}
