package card

// C146 左图右文卡片
// 通用小图的左图右文的样式
type C146 struct {
	BaseItem
	Image                     string      `json:"image,omitempty"`                        // 左图url
	CoverInfo                 string      `json:"cover_info,omitempty"`                   // 左图底部的文本
	CoverDecorText            string      `json:"cover_decor_text,omitempty"`             // 左图右上角角标文字
	CoverDecorBackgroundColor string      `json:"cover_decor_background_color,omitempty"` // 左图右上角角标背景色
	MainTitle                 string      `json:"main_title,omitempty"`                   // 右文：标题
	SubTitle                  string      `json:"sub_title,omitempty"`                    // 右文：子标题
	Mid                       string      `json:"mid,omitempty"`                          // 对应视频的mid
	Oid                       string      `json:"oid,omitempty"`                          // 视频的oid
	ShowMoreMenu              int8        `json:"show_more_menu,omitempty"`               // 是否显示更多菜单
	MoreMenuContent           []*C146Menu `json:"more_menu_content,omitempty"`            // 菜单数组
	OpenUrl                   string      `json:"openurl,omitempty"`                      // 打开链接
}

// C146Menu 菜单
type C146Menu struct {
	Type     string `json:"type,omitempty"`     // 区分不同菜单点击行为
	Name     string `json:"name,omitempty"`     // 菜单文本
	Scheme   string `json:"scheme,omitempty"`   // 菜单使用的跳转scheme
	Editable int8   `json:"editable,omitempty"` // 视频编辑菜单，是否支持编辑
}

// NewCard146 创建Card146实例
func NewCard146() *C146 {
	return &C146{BaseItem: BaseItem{Base: Base{CardType: 146}}}
}
