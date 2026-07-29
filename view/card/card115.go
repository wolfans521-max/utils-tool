package card

// C115 视频运营活动节目顶部卡片
// 用于视频运营活动节目顶部展示
type C115 struct {
	BaseItem
	BgPic          string `json:"bg_pic,omitempty"`          // 背景1:1大图
	TitleTextColor string `json:"title_textcolor,omitempty"` // 标题文字颜色，不下发时走默认颜色白色
	Mblog          any    `json:"mblog,omitempty"`           // 微博数据
	OpenUrl        string `json:"openurl,omitempty"`         // 打开链接
}

// NewCard115 创建Card115实例
func NewCard115() *C115 {
	return &C115{BaseItem: BaseItem{Base: Base{CardType: 115}}}
}
