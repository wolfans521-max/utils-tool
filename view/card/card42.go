package card

// C42 菜单卡片
// 功能：
// 42_1 显示菜单类型
// 42_2 显示无分割线类型
// 42_3 支持配色
// 42_4 支持文本里插入icon
type C42 struct {
	BaseItem
	Style          int8             `json:"style,omitempty"`            // 控制标题左侧icon显示隐藏，0显示，1隐藏
	Pic            string           `json:"pic,omitempty"`              // 图片
	Desc           string           `json:"desc,omitempty"`             // 标题
	TitleExtraText string           `json:"title_extra_text,omitempty"` // 副标题
	DisplayArrow   int8             `json:"display_arrow,omitempty"`    // 显示右侧箭头（只适用于无分割线类型）
	DisplayType    int8             `json:"display_type,omitempty"`     // 0表示默认类型、1表示无分割线类型、2表示菜单类型
	DecorateColor  string           `json:"decorate_color,omitempty"`   // 配色（左边条颜色和右箭头底色）
	Menus          []*Button        `json:"menus,omitempty"`            // 菜单（只适用于显示菜单类型）
	DescStruct     []*C42DescStruct `json:"desc_struct,omitempty"`      // 描述结构
	OpenUrl        string           `json:"openurl,omitempty"`          // 打开链接
}

// C42DescStruct 描述结构
type C42DescStruct struct {
	Name     string `json:"name,omitempty"`      // 显示icon BB1生效
	Icon     string `json:"icon,omitempty"`      // 图标
	IconDark string `json:"icon_dark,omitempty"` // 暗黑模式图标
}

// NewCard42 创建Card42实例
func NewCard42() *C42 {
	return &C42{BaseItem: BaseItem{Base: Base{CardType: 42}}}
}
