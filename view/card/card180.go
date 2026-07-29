package card

// C180 排行榜超话卡片
// 用于排行榜超话展示
type C180 struct {
	BaseItem
	LayoutStyle int8         `json:"layout_style,omitempty"` // 区分样式 1:专辑卡片新样式 2:可动态配置子数目的新样式
	Group       []*C180Group `json:"group,omitempty"`        // 排行分组
	OpenUrl     string       `json:"openurl,omitempty"`      // 打开链接
}

// C180Group 排行分组
type C180Group struct {
	BackgroundUrl  string         `json:"background_url,omitempty"`   // 单组排行榜背景
	Title          string         `json:"title,omitempty"`            // 排行榜分组名称
	TitleColor     string         `json:"title_color,omitempty"`      // 排行分组名称颜色
	TitleColorDark string         `json:"title_color_dark,omitempty"` // 排行分组名称颜色(暗黑)
	TitleIcon      string         `json:"title_icon,omitempty"`       // 排行榜分组名称icon
	Scheme         string         `json:"scheme,omitempty"`           // 排行榜分类scheme
	SubItems       []*C180SubItem `json:"sub_items,omitempty"`        // 排行超话
	ActionLog      map[string]any `json:"actionlog,omitempty"`        // 跳转scheme actionlog
	BackColors     []string       `json:"back_colors,omitempty"`      // 背景颜色渐变
	BackColorsDark []string       `json:"back_colors_dark,omitempty"` // 背景颜色渐变(暗黑)
	HeadBgUrl      string         `json:"head_bg_url,omitempty"`      // 头部区域背景图
	HeadBgUrlDark  string         `json:"head_bg_url_dark,omitempty"` // 头部区域背景图(暗黑)
	Bottom         *C180Bottom    `json:"bottom,omitempty"`           // 底部按钮
}

// C180SubItem 排行超话项
type C180SubItem struct {
	Scheme         string         `json:"scheme,omitempty"`           // 超话scheme
	AvatarUrl      string         `json:"avatar_url,omitempty"`       // 超话头像
	AvatarCorner   int            `json:"avatar_corner,omitempty"`    // 头像圆角
	AvatarStroke   string         `json:"avatar_stroke,omitempty"`    // 头像描边背景
	Rank           string         `json:"rank,omitempty"`             // 排名
	Title          string         `json:"title,omitempty"`            // 超话名称
	TitleColor     string         `json:"title_color,omitempty"`      // 超话名称颜色
	TitleColorDark string         `json:"title_color_dark,omitempty"` // 超话名称暗黑颜色
	Des            string         `json:"des,omitempty"`              // 超话描述
	DesColor       string         `json:"des_color,omitempty"`        // 超话描述颜色
	DesColorDark   string         `json:"des_color_dark,omitempty"`   // 超话描述暗黑颜色
	Button         *Button        `json:"button,omitempty"`           // 关注按钮
	ActionLog      map[string]any `json:"actionlog,omitempty"`        // 跳转scheme actionlog
}

// C180Bottom 底部按钮
type C180Bottom struct {
	Title     string         `json:"title,omitempty"`     // 按钮文案
	Pic       string         `json:"pic,omitempty"`       // 按钮图标地址
	PicDark   string         `json:"pic_dark,omitempty"`  // 按钮图标地址(暗黑)
	Scheme    string         `json:"scheme,omitempty"`    // 按钮行为
	ActionLog map[string]any `json:"actionlog,omitempty"` // 打点日志
}

// NewCard180 创建Card180实例
func NewCard180() *C180 {
	return &C180{BaseItem: BaseItem{Base: Base{CardType: 180}}}
}
