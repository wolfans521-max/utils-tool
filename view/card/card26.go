package card

// C26 左图右文卡片（与Card25类似）
// 支持内嵌标题
// 支持左图右两行文字
// 支持多媒体播放
// 支持标题角标
type C26 struct {
	BaseItem
	Title          string        `json:"title,omitempty"`            // 标题
	TitleExtraText string        `json:"title_extra_text,omitempty"` // 副标题
	ShowTitleArrow int8          `json:"show_title_arrow,omitempty"` // 标题右侧箭头
	Score          string        `json:"score,omitempty"`            // 评分
	Grade          string        `json:"grade,omitempty"`            // 等级
	DisplayArrow   int8          `json:"display_arrow,omitempty"`    // 显示箭头
	TitleSub       string        `json:"title_sub,omitempty"`        // 副标题
	Desc           string        `json:"desc,omitempty"`             // 描述
	Pic            string        `json:"pic,omitempty"`              // 左图
	Buttons        []*Button     `json:"buttons,omitempty"`          // 按钮列表
	MediaInfo      *C26MediaInfo `json:"media_info,omitempty"`       // 多媒体信息
	TitleFlagPic   string        `json:"title_flag_pic,omitempty"`   // 左上角角标
	FlagPic        string        `json:"flag_pic,omitempty"`         // 标题角标
	IsCollection   bool          `json:"is_collection,omitempty"`    // 是否显示合集背景
	BorderRadius   int           `json:"borderRadius,omitempty"`     // 图片圆角大小
	PicWidth       int           `json:"picWidth,omitempty"`         // 图片宽度
	DescContent    string        `json:"desc_content,omitempty"`     // 第二行文案描述字段
	MaxLines       int           `json:"Max_lines,omitempty"`        // 第二行最大行数
	TagInfo        *C26TagInfo   `json:"tag_info,omitempty"`         // 标签信息
	OpenUrl        string        `json:"openurl,omitempty"`          // 打开链接
}

// C26MediaInfo 多媒体信息
type C26MediaInfo struct {
	Name     string `json:"name,omitempty"`     // 名称
	Duration int    `json:"duration,omitempty"` // 时长
}

// C26TagInfo 标签信息
type C26TagInfo struct {
	Text      string         `json:"text,omitempty"`      // 标签内容
	TextColor string         `json:"textcolor,omitempty"` // 标签内容颜色
	BgColor   string         `json:"bgcolor,omitempty"`   // 标签背景颜色
	Scheme    string         `json:"scheme,omitempty"`    // 跳转url
	ActionLog map[string]any `json:"actionlog,omitempty"` // 埋点数据
}

// NewCard26 创建Card26实例
func NewCard26() *C26 {
	return &C26{BaseItem: BaseItem{Base: Base{CardType: 26}}}
}
