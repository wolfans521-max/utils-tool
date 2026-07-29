package card

// C7 文本卡片
// 样式一：纯文本，支持单行标题、多行标题、内容文本和来源、文本高亮
// 样式二：Card标题，支持单行、多行
type C7 struct {
	BaseItem
	CardTypeName             string       `json:"card_type_name,omitempty"`              // 卡片类型名称
	Desc                     string       `json:"desc,omitempty"`                        // 内容文本，超过三行将显示更多按钮
	DescAlignment            int8         `json:"desc_alignment,omitempty"`              // desc内容对齐方式，0-左对齐，1-居中，2-右对齐
	DescMaxLine              int          `json:"desc_max_line,omitempty"`               // -1:不限制行数；0：默认三行；>0：设置自定义行数
	TitleFontStyle           int8         `json:"title_font_stytle,omitempty"`           // 0或不传表示默认字体，1表示12dp大小字号字体
	ContentTextColor         string       `json:"content_text_color,omitempty"`          // desc内容正常状态下字体颜色
	ContentFontSize          int8         `json:"content_font_size,omitempty"`           // desc内容字体，0或不传表示默认字体
	HideContentBottomPadding int8         `json:"hide_content_bottom_padding,omitempty"` // 0或者不传表示默认间距，1表示隐藏内容下边距
	HideContentTopPadding    int8         `json:"hide_content_top_padding,omitempty"`    // 0或者不传表示默认间距，1表示隐藏内容上边距
	HideLine                 int8         `json:"hide_line,omitempty"`                   // 0或者不传表示有分隔线，1表示隐藏分隔线
	Highlight                *C7Highlight `json:"highlight,omitempty"`                   // 高亮配置
	Title                    string       `json:"title,omitempty"`                       // 标题
	Source                   string       `json:"source,omitempty"`                      // 来源
	HighlightColorType       int8         `json:"highlight_color_type,omitempty"`        // source高亮颜色 0默认，1蓝色
	DisplayArrow             int8         `json:"display_arrow,omitempty"`               // 标题右侧箭头
	ShowType                 string       `json:"show_type,omitempty"`                   // 只在Android发挥作用，1表示显示2行，0显示偶同文本
	ShowMoreStyle            int8         `json:"show_more_style,omitempty"`             // 0:默认 1:强制不显示尾部更多
}

// C7Highlight 高亮配置
type C7Highlight struct {
	TitleEm  [][]any `json:"title_em,omitempty"`  // title标题高亮字体范围（黄色）
	DescEm   [][]any `json:"desc_em,omitempty"`   // desc内容高亮字体范围（黄色）
	SourceEm [][]any `json:"source_em,omitempty"` // source来源，字体可点击范围（蓝色）
}

// NewCard7 创建Card7实例
func NewCard7() *C7 {
	return &C7{BaseItem: BaseItem{Base: Base{CardType: 7}}}
}
