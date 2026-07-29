package card

// C102 搜索页面多图文章卡片
// 用于搜索页面展示多图文章
type C102 struct {
	BaseItem
	CardTypeName string     `json:"card_type_name,omitempty"` // 卡片类型名称
	Title        string     `json:"title,omitempty"`          // Card要显示的内容
	MaxLines     string     `json:"maxLines,omitempty"`       // 文本显示的最大行数
	Source       string     `json:"source,omitempty"`         // 新闻来源文本
	Forward      string     `json:"forward,omitempty"`        // 转发描述
	Time         string     `json:"time,omitempty"`           // 发布时间描述
	Pics         []*C102Pic `json:"pics,omitempty"`           // 图片列表
	OpenUrl      string     `json:"openurl,omitempty"`        // 打开链接
}

// C102Pic 图片
type C102Pic struct {
	Url    string `json:"url,omitempty"`    // 图片链接
	Width  int    `json:"width,omitempty"`  // 宽度比例
	Height int    `json:"height,omitempty"` // 高度比例
}

// NewCard102 创建Card102实例
func NewCard102() *C102 {
	return &C102{BaseItem: BaseItem{Base: Base{CardType: 102}}}
}
