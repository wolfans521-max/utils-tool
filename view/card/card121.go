package card

// C121 过滤器卡片
// page页过滤器
type C121 struct {
	BaseItem
	Type             int8               `json:"type,omitempty"`               // 类型
	Default          string             `json:"default,omitempty"`            // 选中第几项，默认第0项；都不选中下发-1
	Align            string             `json:"align,omitempty"`              // 居左下发"left"，默认居中
	SelectedFontBold bool               `json:"selected_font_bold,omitempty"` // 选中项字体加粗，默认不加粗
	FilterGroup      []*C121FilterGroup `json:"filter_group,omitempty"`       // 各个子标签
	OpenUrl          string             `json:"openurl,omitempty"`            // 打开链接
}

// C121FilterGroup 过滤器组
type C121FilterGroup struct {
	Name                string `json:"name,omitempty"`                   // 标签标题
	ContainerId         string `json:"containerid,omitempty"`            // 标签的id
	FontSize            int    `json:"font_size,omitempty"`              // 标签字体大小，默认12
	ColorNormal         string `json:"color_normal,omitempty"`           // 未选中文字颜色
	ColorSelected       string `json:"color_selected,omitempty"`         // 选中文字颜色
	ColorNormalDark     string `json:"color_normal_dark,omitempty"`      // 暗黑未选中文字颜色
	ColorSelectedDark   string `json:"color_selected_dark,omitempty"`    // 暗黑选中文字颜色
	BgColorNormal       string `json:"bg_color_normal,omitempty"`        // 未选中背景颜色
	BgColorSelected     string `json:"bg_color_selected,omitempty"`      // 选中背景颜色
	BgColorNormalDark   string `json:"bg_color_normal_dark,omitempty"`   // 暗黑未选中背景颜色
	BgColorSelectedDark string `json:"bg_color_selected_dark,omitempty"` // 暗黑选中背景颜色
	RedDot              bool   `json:"red_dot,omitempty"`                // 未选中时是否展示小红点
	BusinessName        string `json:"business_name,omitempty"`          // 业务名称
}

// NewCard121 创建Card121实例
func NewCard121() *C121 {
	return &C121{BaseItem: BaseItem{Base: Base{CardType: 121}}}
}
