package card

// C61 用户推荐卡片
type C61 struct {
	BaseItem
	CardCommTitle      string           `json:"card_comm_title,omitempty"`
	CardBackgroudcolor string           `json:"card_backgroudcolor,omitempty"`
	IsShowBorder       string           `json:"is_show_border,omitempty"`
	Desc1              string           `json:"desc1,omitempty"`
	DeleteAction       *C61DeleteAction `json:"delete_action,omitempty"`
	Buttons            []*Button        `json:"buttons,omitempty"`
	User               *C61User         `json:"user,omitempty"`
}

// C61DeleteAction 删除操作
type C61DeleteAction struct {
	CanDelete int8           `json:"can_delete,omitempty"`
	ActionLog map[string]any `json:"actionlog,omitempty"`
}

// C61User 用户信息
type C61User struct {
	Id              int64  `json:"id,omitempty"`
	ScreenName      string `json:"screen_name,omitempty"`
	Description     string `json:"description,omitempty"`
	Gender          string `json:"gender,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
	AvatarLarge     string `json:"avatar_large,omitempty"`
	AvatarHd        string `json:"avatar_hd,omitempty"`
	Verified        bool   `json:"verified,omitempty"`
	VerifiedType    int8   `json:"verified_type,omitempty"`
	FollowersCount  int64  `json:"followers_count,omitempty"`
	FriendsCount    int64  `json:"friends_count,omitempty"`
	VerifiedTypeExt int8   `json:"verified_type_ext,omitempty"`
	VerifiedReason  string `json:"verified_reason,omitempty"`
	Level           int8   `json:"level,omitempty"`
	Ptype           int8   `json:"ptype,omitempty"`
}

// NewCard61 创建Card61实例
func NewCard61() *C61 {
	return &C61{BaseItem: BaseItem{Base: Base{CardType: 61}}}
}
