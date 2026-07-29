package card

// C132 落地页标题卡片
// 用于落地页标题展示
type C132 struct {
	BaseItem
	Title       string `json:"title,omitempty"`        // 标题文字
	InfoTitle   string `json:"info_title,omitempty"`   // 点击icon后显示信息对话框的标题，为空不显示入口
	InfoContent string `json:"info_content,omitempty"` // 点击icon后显示信息对话框的文本内容，为空不显示入口
	OpenUrl     string `json:"openurl,omitempty"`      // 打开链接
}

// NewCard132 创建Card132实例
func NewCard132() *C132 {
	return &C132{BaseItem: BaseItem{Base: Base{CardType: 132}}}
}
