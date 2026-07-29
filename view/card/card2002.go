package card

// C2002 倒计时/正计时卡片
// 用于展示活动倒计时或正计时
type C2002 struct {
	BaseItem
	Title           string `json:"title,omitempty"`                // 标题
	EndTitle        string `json:"end_title,omitempty"`            // 结束标题
	CountType       int    `json:"count_type,omitempty"`           // 计时类型 0:倒计时 1:正计时
	StartTime       int64  `json:"start_time,omitempty"`           // 开始时间
	TitleColor      string `json:"title_color,omitempty"`          // 标题颜色
	NumberColor     string `json:"number_color,omitempty"`         // 数字颜色
	NumberWarnColor string `json:"number_warning_color,omitempty"` // 数字警告颜色
	UnitColor       string `json:"unit_color,omitempty"`           // 单位颜色
	UnitWarnColor   string `json:"unit_warning_color,omitempty"`   // 单位警告颜色
}

// NewCard2002 创建Card2002实例
func NewCard2002() *C2002 {
	return &C2002{BaseItem: BaseItem{Base: Base{CardType: 2002}}}
}
