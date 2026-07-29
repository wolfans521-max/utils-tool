package card

// C77 预约卡片
// 主要用于比赛预约等场景
type C77 struct {
	BaseItem
	SportName       string          `json:"sport_name,omitempty"`        // 比赛名称
	Title           string          `json:"title,omitempty"`             // 比赛标题
	TitleLeftColor  string          `json:"title_left_color,omitempty"`  // 左侧颜色条颜色
	Type            int8            `json:"type,omitempty"`              // 显示类型，0表示单行描述样式，1表示国旗选手VS样式
	Desc            string          `json:"desc,omitempty"`              // 单行描述样式的单行描述内容
	StartTime       string          `json:"start_time,omitempty"`        // 左侧开始时间
	LeftFlagIcon    string          `json:"left_flag_icon,omitempty"`    // 国旗选手VS样式左侧国旗图标
	LeftPlayerName  string          `json:"left_player_name,omitempty"`  // 国旗选手VS样式左侧选手名称
	RightFlagIcon   string          `json:"right_flag_icon,omitempty"`   // 国旗选手VS样式右侧国旗图标
	RightPlayerName string          `json:"right_player_name,omitempty"` // 国旗选手VS样式右侧选手名称
	OrderButton     *C77OrderButton `json:"order_button,omitempty"`      // 预约按钮结构
	OpenUrl         string          `json:"openurl,omitempty"`           // 打开链接
}

// C77OrderButton 预约按钮
type C77OrderButton struct {
	Type      string `json:"type,omitempty"`       // 按钮类型
	Name      string `json:"name,omitempty"`       // 按钮名称
	Title     string `json:"title,omitempty"`      // 存储到日历中的标题
	Des       string `json:"des,omitempty"`        // 存储到日历中的描述
	Scheme    string `json:"scheme,omitempty"`     // 存储到日历中的URL
	DtStart   string `json:"dt_start,omitempty"`   // 存储到日历中的开始时间
	AlarmList []int  `json:"alarm_list,omitempty"` // 存储到日历中的提醒时间，可设置多个，表示几分钟前开始提醒
}

// NewCard77 创建Card77实例
func NewCard77() *C77 {
	return &C77{BaseItem: BaseItem{Base: Base{CardType: 77}}}
}
