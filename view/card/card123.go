package card

// C123 聊天记录卡片
// 支持上下两行标题，time字段为时间
// 出现场景：消息箱-搜索聊天记录-群聊中的单人聊天记录
type C123 struct {
	BaseItem
	TitleSub  string         `json:"title_sub,omitempty"` // 第一行标题
	Desc      string         `json:"desc,omitempty"`      // 第二行标题
	Pic       string         `json:"pic,omitempty"`       // 图片
	Time      string         `json:"time,omitempty"`      // 时间
	Highlight *C123Highlight `json:"highlight,omitempty"` // 高亮配置
	OpenUrl   string         `json:"openurl,omitempty"`   // 打开链接
}

// C123Highlight 高亮配置
type C123Highlight struct {
	DescEm [][]int `json:"desc_em,omitempty"` // 高亮位置
}

// NewCard123 创建Card123实例
func NewCard123() *C123 {
	return &C123{BaseItem: BaseItem{Base: Base{CardType: 123}}}
}
