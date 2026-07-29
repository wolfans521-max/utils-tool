package card

// C174 搜索视频背景和图片背景预约卡片
// 用于搜索视频背景和图片背景预约展示
type C174 struct {
	BaseItem
	CoverUrl           string             `json:"coverurl,omitempty"`           // 封面背景图
	Mblog              any                `json:"mblog,omitempty"`              // 封面视频
	Style              *C174Style         `json:"style,omitempty"`              // 样式
	Click              *C174Click         `json:"click,omitempty"`              // 按钮点击封装
	Title              *C174Title         `json:"title,omitempty"`              // 封面标题
	ReservationEndTime int64              `json:"reservationEndTime,omitempty"` // 开始时间
	CalendarInfo       *C174CalendarInfo  `json:"calendarInfo,omitempty"`       // 日历预约信息
	Reservation        []*C174Reservation `json:"reservation,omitempty"`        // 预约按钮列表
	Entrance           []*C174Entrance    `json:"entrance,omitempty"`           // 入口列表
	OpenUrl            string             `json:"openurl,omitempty"`            // 打开链接
}

// C174Style 样式
type C174Style struct {
	Ratio           float64 `json:"ratio,omitempty"`           // 比例（高/宽）
	MaskEndColor    string  `json:"maskEndColor,omitempty"`    // 蒙层渐变色底色
	MaskEndColorKey string  `json:"maskEndColorKey,omitempty"` // 蒙层渐变色底色暗黑模式
}

// C174Click 点击
type C174Click struct {
	Type      string         `json:"type,omitempty"`      // 类型
	Scheme    string         `json:"scheme,omitempty"`    // 跳转链接
	ActionLog map[string]any `json:"actionlog,omitempty"` // 打点日志
}

// C174Title 标题
type C174Title struct {
	Style   *C174TitleStyle `json:"style,omitempty"`   // 样式
	Content string          `json:"content,omitempty"` // 内容
}

// C174TitleStyle 标题样式
type C174TitleStyle struct {
	TextColor    string `json:"textColor,omitempty"`    // 文字颜色
	TextColorKey string `json:"textColorKey,omitempty"` // 文字颜色暗黑模式
	TextSize     int    `json:"textSize,omitempty"`     // 文字大小
	MaxLines     int    `json:"maxlines,omitempty"`     // 最大行数
}

// C174CalendarInfo 日历信息
type C174CalendarInfo struct {
	Type      string `json:"type,omitempty"`       // 类型
	Title     string `json:"title,omitempty"`      // 存储到日历中的标题
	Des       string `json:"des,omitempty"`        // 存储到日历中的描述
	Scheme    string `json:"scheme,omitempty"`     // 存储到日历中的URL
	DtStart   string `json:"dt_start,omitempty"`   // 存储到日历中的开始时间
	AlarmList []int  `json:"alarm_list,omitempty"` // 存储到日历中的提醒时间
}

// C174Reservation 预约
type C174Reservation struct {
	Style   *C174ReservationStyle `json:"style,omitempty"`   // 样式
	Content string                `json:"content,omitempty"` // 内容
	Icon    string                `json:"icon,omitempty"`    // 图标
	Click   *C174Click            `json:"click,omitempty"`   // 点击
}

// C174ReservationStyle 预约样式
type C174ReservationStyle struct {
	BorderColor    string `json:"borderColor,omitempty"`    // 边框颜色
	BorderColorKey string `json:"borderColorKey,omitempty"` // 边框颜色暗黑模式
	BorderWidth    int    `json:"borderWidth,omitempty"`    // 边框宽度
	TextColor      string `json:"textColor,omitempty"`      // 文字颜色
	TextColorKey   string `json:"textColorKey,omitempty"`   // 文字颜色暗黑模式
	TextSize       int    `json:"textSize,omitempty"`       // 文字大小
	Radius         int    `json:"radius,omitempty"`         // 圆角
	Width          int    `json:"width,omitempty"`          // 宽度
	Height         int    `json:"height,omitempty"`         // 高度
}

// C174Entrance 入口
type C174Entrance struct {
	Content string          `json:"content,omitempty"` // 内容
	Style   *C174TitleStyle `json:"style,omitempty"`   // 样式
	Click   *C174Click      `json:"click,omitempty"`   // 点击
}

// NewCard174 创建Card174实例
func NewCard174() *C174 {
	return &C174{BaseItem: BaseItem{Base: Base{CardType: 174}}}
}
