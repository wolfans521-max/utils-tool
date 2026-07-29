package card

// C170 商品评分卡片
// 用于商品评分展示
type C170 struct {
	BaseItem
	Buttons  []*C170Button `json:"buttons,omitempty"`   // 按钮列表
	StarInfo *C170StarInfo `json:"star_info,omitempty"` // 评分信息
	Style    *C170Style    `json:"style,omitempty"`     // 样式
	TagInfo  *C170TagInfo  `json:"tag_info,omitempty"`  // 标签信息
	OpenUrl  string        `json:"openurl,omitempty"`   // 打开链接
}

// C170Button 按钮
type C170Button struct {
	SkipFormat int8              `json:"skip_format,omitempty"` // 跳过格式化
	Type       string            `json:"type,omitempty"`        // 类型
	Name       string            `json:"name,omitempty"`        // 名称
	Params     map[string]string `json:"params,omitempty"`      // 参数
	Pic        string            `json:"pic,omitempty"`         // 图片
	Style      *C170ButtonStyle  `json:"style,omitempty"`       // 样式
	ExtTitle   string            `json:"ext_title,omitempty"`   // 扩展标题
	ActionLog  map[string]any    `json:"actionlog,omitempty"`   // 打点日志
}

// C170ButtonStyle 按钮样式
type C170ButtonStyle struct {
	FontSize    int    `json:"font_size,omitempty"`    // 字体大小
	SubSize     int    `json:"sub_size,omitempty"`     // 副字体大小
	FontColor   string `json:"font_color,omitempty"`   // 字体颜色
	BorderColor string `json:"border_color,omitempty"` // 边框颜色
}

// C170StarInfo 评分信息
type C170StarInfo struct {
	Scheme    string         `json:"scheme,omitempty"`    // 跳转链接
	ActionLog map[string]any `json:"actionlog,omitempty"` // 打点日志
	Desc      *C170Desc      `json:"desc,omitempty"`      // 描述
	Star      *C170Star      `json:"star,omitempty"`      // 星星
}

// C170Desc 描述
type C170Desc struct {
	Title string         `json:"title,omitempty"` // 标题
	Style *C170DescStyle `json:"style,omitempty"` // 样式
}

// C170DescStyle 描述样式
type C170DescStyle struct {
	FontSize  int    `json:"font_size,omitempty"`  // 字体大小
	FontColor string `json:"font_color,omitempty"` // 字体颜色
}

// C170Star 星星
type C170Star struct {
	Icon          string         `json:"icon,omitempty"`            // 图标
	IconBaseColor string         `json:"icon_base_color,omitempty"` // 图标基础颜色
	IconSize      int            `json:"icon_size,omitempty"`       // 图标大小
	Space         int            `json:"space,omitempty"`           // 间距
	Style         *C170DescStyle `json:"style,omitempty"`           // 样式
	Rate          string         `json:"rate,omitempty"`            // 评分
	StarCount     int            `json:"star_count,omitempty"`      // 星星数量
}

// C170Style 样式
type C170Style struct {
	BackgroundColor string   `json:"background_color,omitempty"` // 背景颜色
	BorderColors    []string `json:"border_colors,omitempty"`    // 边框颜色
	BorderWidth     int      `json:"border_width,omitempty"`     // 边框宽度
}

// C170TagInfo 标签信息
type C170TagInfo struct {
	MaxLines int        `json:"max_lines,omitempty"` // 最大行数
	Tags     []*C170Tag `json:"tags,omitempty"`      // 标签列表
}

// C170Tag 标签
type C170Tag struct {
	Title     string            `json:"title,omitempty"`     // 标题
	Style     *C170TagStyle     `json:"style,omitempty"`     // 样式
	Params    map[string]string `json:"params,omitempty"`    // 参数
	ActionLog map[string]any    `json:"actionlog,omitempty"` // 打点日志
}

// C170TagStyle 标签样式
type C170TagStyle struct {
	FontSize        int    `json:"font_size,omitempty"`        // 字体大小
	FontColor       string `json:"font_color,omitempty"`       // 字体颜色
	BackgroundColor string `json:"background_color,omitempty"` // 背景颜色
}

// NewCard170 创建Card170实例
func NewCard170() *C170 {
	return &C170{BaseItem: BaseItem{Base: Base{CardType: 170}}}
}
