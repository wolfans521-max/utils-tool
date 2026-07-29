package card

// C99 运动周报卡片
// 主要用于在微博运动与私信中显示周报数据
type C99 struct {
	BaseItem
	DisplayArrow        int8              `json:"display_arrow,omitempty"`        // 显示箭头
	Title               string            `json:"title,omitempty"`                // 标题
	MultimediaActionLog string            `json:"multimedia_actionlog,omitempty"` // 多媒体打点
	IsAsyn              int8              `json:"is_asyn,omitempty"`              // 是否异步
	HighlightState      string            `json:"highlight_state,omitempty"`      // 高亮状态
	User                any               `json:"user,omitempty"`                 // 用户信息
	TitleStruct         []*C99TitleStruct `json:"title_struct,omitempty"`         // 标题结构
	Data                []string          `json:"data,omitempty"`                 // 数据
	RightText           string            `json:"right_text,omitempty"`           // 右侧文本
	SubRightText        string            `json:"sub_right_text,omitempty"`       // 右侧副文本
	Desc                string            `json:"desc,omitempty"`                 // 描述
	Timestamp           int64             `json:"timestamp,omitempty"`            // 时间戳
	OpenUrl             string            `json:"openurl,omitempty"`              // 打开链接
}

// C99TitleStruct 标题结构
type C99TitleStruct struct {
	Name      string `json:"name,omitempty"`       // 名称
	Scheme    string `json:"scheme,omitempty"`     // 跳转链接
	Color     string `json:"color,omitempty"`      // 颜色
	ColorSkin string `json:"color_skin,omitempty"` // 皮肤颜色
}

// NewCard99 创建Card99实例
func NewCard99() *C99 {
	return &C99{BaseItem: BaseItem{Base: Base{CardType: 99}}}
}
