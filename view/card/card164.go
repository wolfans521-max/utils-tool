package card

// C164 潮流双卡片
// 采用159数据结构进行拼接的双排card
type C164 struct {
	BaseItem
	LeftElement  *C159  `json:"left_element,omitempty"`  // 左边元素，card159数据结构
	RightElement *C159  `json:"right_element,omitempty"` // 右边元素，card159数据结构
	OpenUrl      string `json:"openurl,omitempty"`       // 打开链接
}

// NewCard164 创建Card164实例
func NewCard164() *C164 {
	return &C164{BaseItem: BaseItem{Base: Base{CardType: 164}}}
}
