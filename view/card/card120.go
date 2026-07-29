package card

// C120 滚动卡片容器
// 可以指定任意类型card执行上下滚动
type C120 struct {
	BaseItem
	TimeInterval int        `json:"time_interval,omitempty"` // 滚动时间间隔
	Direction    int8       `json:"direction,omitempty"`     // 默认0是向上，1是向下
	SubcardType  int        `json:"subcard_type,omitempty"`  // 标记需要滚动的card类型
	Groups       []any      `json:"groups,omitempty"`        // 需要滚动的cards
	StillDict    *C120Still `json:"still_dict,omitempty"`    // 静态区域
	OpenUrl      string     `json:"openurl,omitempty"`       // 打开链接
}

// C120Still 静态区域
type C120Still struct {
	Icon string `json:"icon,omitempty"` // 图标
}

// NewCard120 创建Card120实例
func NewCard120() *C120 {
	return &C120{BaseItem: BaseItem{Base: Base{CardType: 120}}}
}
