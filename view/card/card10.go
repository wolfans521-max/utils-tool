package card

// C10 用户卡片
// 支持内嵌标题
// 支持2段或3段文字和CommonButton
// 支持来源文字可配
type C10 struct {
	BaseItem
	Title          string          `json:"title,omitempty"`            // 标题
	TitleExtraText string          `json:"title_extra_text,omitempty"` // 副标题
	ShowTitleArrow int8            `json:"show_title_arrow,omitempty"` // 标题右侧箭头
	DisplayArrow   int8            `json:"display_arrow,omitempty"`    // card右侧箭头
	User           *C10User        `json:"user,omitempty"`             // 用户结构
	Buttons        []*Button       `json:"buttons,omitempty"`          // CommonButton
	Desc1          string          `json:"desc1,omitempty"`            // 描述1
	Desc2          string          `json:"desc2,omitempty"`            // 描述2
	Desc2Struct    []*C10Desc2Item `json:"desc2_struct,omitempty"`     // 来源文字可配
	OpenUrl        string          `json:"openurl,omitempty"`          // 打开链接
}

// C10User 用户结构
type C10User struct {
	Id              int64  `json:"id,omitempty"`
	Idstr           string `json:"idstr,omitempty"`
	ScreenName      string `json:"screen_name,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
	AvatarLarge     string `json:"avatar_large,omitempty"`
	AvatarHd        string `json:"avatar_hd,omitempty"`
	Verified        bool   `json:"verified,omitempty"`
	VerifiedType    int8   `json:"verified_type,omitempty"`
	VerifiedReason  string `json:"verified_reason,omitempty"`
	Description     string `json:"description,omitempty"`
}

// C10Desc2Item 来源文字配置项
type C10Desc2Item struct {
	Name   string `json:"name,omitempty"`   // 文字名称
	Scheme string `json:"scheme,omitempty"` // 跳转链接
}

// NewCard10 创建Card10实例
func NewCard10() *C10 {
	return &C10{BaseItem: BaseItem{Base: Base{CardType: 10}}}
}
