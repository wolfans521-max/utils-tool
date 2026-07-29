package card

// C243 原创微博卡片
// 用于原创微博展示，支持纯文本、文本+图片、文本+视频封面
type C243 struct {
	BaseItem
	User          any    `json:"user,omitempty"`           // 用户信息
	CreatedAt     string `json:"created_at,omitempty"`     // 博文创建时间
	RightTagImg   string `json:"right_tag_img,omitempty"`  // 右上角图片
	ProfileScheme string `json:"profile_scheme,omitempty"` // profile页scheme
	Text          string `json:"text,omitempty"`           // 博文
	MaxLine       int    `json:"max_line,omitempty"`       // 文本最大行数
	ExtraInfo     string `json:"extra_info,omitempty"`     // 左下角文本
	Pic           string `json:"pic,omitempty"`            // 图片
	PlayIcon      string `json:"play_icon,omitempty"`      // 播放图标
	Duration      string `json:"duration,omitempty"`       // 视频时长
	OpenUrl       string `json:"openurl,omitempty"`        // 打开链接
}

// NewCard243 创建Card243实例
func NewCard243() *C243 {
	return &C243{BaseItem: BaseItem{Base: Base{CardType: 243}}}
}
