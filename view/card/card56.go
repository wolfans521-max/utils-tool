package card

// C56 投票对战卡片
type C56 struct {
	BaseItem
	Title        string   `json:"title,omitempty"`
	PositiveSide *C56Side `json:"positive_side,omitempty"`
	NegativeSide *C56Side `json:"negative_side,omitempty"`
}

// C56Side 投票方数据结构（支持/反对）
type C56Side struct {
	Desc      string         `json:"desc,omitempty"`       // 支持/反对理由
	Count     string         `json:"count,omitempty"`      // 支持/反对人数
	Button    *Button        `json:"button,omitempty"`     // commonButton通用数据
	ActionLog map[string]any `json:"action_log,omitempty"` // 打点数据
	Pic       string         `json:"pic,omitempty"`        // 支持/反对图片
}

// NewCard56 创建Card56实例
func NewCard56() *C56 {
	return &C56{BaseItem: BaseItem{Base: Base{CardType: 56}}}
}
