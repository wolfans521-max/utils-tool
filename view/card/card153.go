package card

// C153 创作中心数据展示卡片
// 用于创作中心数据展示
type C153 struct {
	BaseItem
	Col      int         `json:"col,omitempty"`       // 呈现为N列展示，默认1列
	DataList []*C153Data `json:"data_list,omitempty"` // 需要显示的各大指标
	OpenUrl  string      `json:"openurl,omitempty"`   // 打开链接
}

// C153Data 数据
type C153Data struct {
	Scheme  string `json:"scheme,omitempty"`  // 跳转链接
	Content string `json:"content,omitempty"` // 指标内容
	Label   string `json:"label,omitempty"`   // 指标标题
}

// NewCard153 创建Card153实例
func NewCard153() *C153 {
	return &C153{BaseItem: BaseItem{Base: Base{CardType: 153}}}
}
