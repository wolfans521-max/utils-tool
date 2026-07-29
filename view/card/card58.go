package card

// C58 间隔卡片（支持type 0-5多种样式）
type C58 struct {
	BaseItem
	Type                              int8   `json:"type"`                                            // 控制间隔card的样式，可支持0-5
	Name                              string `json:"name,omitempty"`                                  // 标题
	BackgroundColor                   string `json:"background_color,omitempty"`                      // 卡片背景浅色
	BackgroundColorDark               string `json:"background_color_dark,omitempty"`                 // 卡片背景深色
	NameSize                          int    `json:"name_size,omitempty"`                             // 标题的文字大小
	NameColor                         string `json:"name_color,omitempty"`                            // 默认模式下文本的颜色
	NameColorDark                     string `json:"name_color_dark,omitempty"`                       // 深色模式下文本的颜色
	NameBold                          bool   `json:"name_bold,omitempty"`                             // 字体是否加粗
	Height                            int    `json:"height,omitempty"`                                // 整个card的高度
	RichText                          string `json:"rich_text,omitempty"`                             // 富文本文字
	RichTextColor                     string `json:"rich_text_color,omitempty"`                       // 富文本文字浅色
	RichTextDarkColor                 string `json:"rich_text_dark_color,omitempty"`                  // 富文本文字暗色
	RichTextScheme                    string `json:"rich_text_scheme,omitempty"`                      // 富文本链接，containerID
	RichArrow                         bool   `json:"rich_arrow,omitempty"`                            // 是否展示右侧箭头
	RichArrowImage                    string `json:"rich_arrow_image,omitempty"`                      // 右侧箭头浅色图片url
	RichArrowDarkImage                string `json:"rich_arrow_dark_image,omitempty"`                 // 右侧箭头深色图片url
	DottedLine                        bool   `json:"dotted_line,omitempty"`                           // 是否展示虚线
	ButtonBackgroundColor             string `json:"button_background_color,omitempty"`               // 按钮浅色背景色
	ButtonDarkBackgroundColor         string `json:"button_dark_background_color,omitempty"`          // 按钮深色背景色
	ButtonSelectedBackgroundColor     string `json:"button_selected_background_color,omitempty"`      // 按钮选中态浅色背景色
	ButtonSelectedDarkBackgroundColor string `json:"button_selected_dark_background_color,omitempty"` // 按钮选中态深色背景色
	ButtonSelectedTextColor           string `json:"button_selected_text_color,omitempty"`            // 按钮文字选中态浅色色值
	ButtonSelectedDarkTextColor       string `json:"button_selected_dark_text_color,omitempty"`       // 按钮文字选中态深色色值
	ButtonScheme                      string `json:"button_scheme,omitempty"`                         // 按钮跳转链接，containerID
	UserIconImage                     string `json:"user_icon_image,omitempty"`                       // 按钮内浅色图片url
	UserIconImageDark                 string `json:"user_icon_image_dark,omitempty"`                  // 按钮内深色图片url
	Openurl                           string `json:"openurl,omitempty"`
}

// NewCard58 创建Card58实例
func NewCard58() *C58 {
	return &C58{BaseItem: BaseItem{Base: Base{CardType: 58}}}
}
