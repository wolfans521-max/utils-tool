package card

// C240 榜单页热搜推送卡片
// 榜单页热搜推送，点击查看更多
type C240 struct {
	BaseItem
	Fold       *C240FoldState `json:"fold,omitempty"`        // 折叠状态
	Loading    *C240FoldState `json:"loading,omitempty"`     // 加载状态
	Unfold     *C240FoldState `json:"unfold,omitempty"`      // 展开状态
	RequestUrl string         `json:"request_url,omitempty"` // 请求URL
	OpenUrl    string         `json:"openurl,omitempty"`     // 打开链接
}

// C240FoldState 折叠状态
type C240FoldState struct {
	Icon           string `json:"icon,omitempty"`             // 图标
	IconDark       string `json:"icon_dark,omitempty"`        // 图标暗黑
	TitleFont      int    `json:"title_font,omitempty"`       // 标题字体大小
	Title          string `json:"title,omitempty"`            // 标题
	TitleColor     string `json:"title_color,omitempty"`      // 标题颜色
	TitleColorDark string `json:"title_color_dark,omitempty"` // 标题颜色暗黑
}

// NewCard240 创建Card240实例
func NewCard240() *C240 {
	return &C240{BaseItem: BaseItem{Base: Base{CardType: 240}}}
}
