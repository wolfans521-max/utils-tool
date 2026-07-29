package card

// C30 左图右文卡片
// 支持左图右2行文字
// 支持右侧可配按钮
// 支持配底色
// 支持右侧显示图片
type C30 struct {
	BaseItem
	DisplayArrow    int8      `json:"display_arrow,omitempty"`    // 显示箭头
	User            *C30User  `json:"user,omitempty"`             // 用户信息
	UserColor       string    `json:"user_color,omitempty"`       // 用户名颜色
	Desc1           string    `json:"desc1,omitempty"`            // 描述
	BackgroundColor int8      `json:"background_color,omitempty"` // 底色（只支持配1，显示黄色）
	RightPic        string    `json:"right_pic,omitempty"`        // 右侧显示图片
	Buttons         []*Button `json:"buttons,omitempty"`          // 可配按钮
	OpenUrl         string    `json:"openurl,omitempty"`          // 打开链接
}

// C30User 用户信息
type C30User struct {
	Id              int64  `json:"id,omitempty"`
	Idstr           string `json:"idstr,omitempty"`
	ScreenName      string `json:"screen_name,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
	AvatarLarge     string `json:"avatar_large,omitempty"`
	Verified        bool   `json:"verified,omitempty"`
	VerifiedType    int8   `json:"verified_type,omitempty"`
}

// NewCard30 创建Card30实例
func NewCard30() *C30 {
	return &C30{BaseItem: BaseItem{Base: Base{CardType: 30}}}
}
