package card

// C21 卡片类型21
// 具体功能参考文档
type C21 struct {
	BaseItem
	CardTypeName string `json:"card_type_name,omitempty"` // 卡片类型名称
	Title        string `json:"title,omitempty"`          // 标题
	Desc         string `json:"desc,omitempty"`           // 描述
	Pic          string `json:"pic,omitempty"`            // 图片
	OpenUrl      string `json:"openurl,omitempty"`        // 打开链接
}

// NewCard21 创建Card21实例
func NewCard21() *C21 {
	return &C21{BaseItem: BaseItem{Base: Base{CardType: 21}}}
}
