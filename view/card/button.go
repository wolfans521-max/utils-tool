package card

type Button struct {
	Type                     string             `json:"type"`
	Name                     string             `json:"name,omitempty"`
	BtnStyle                 int8               `json:"btn_style,omitempty"`
	BtnTitle                 string             `json:"btn_title,omitempty"`
	ButtonStyle              string             `json:"button_style,omitempty"`
	BorderColor              string             `json:"border_color,omitempty"`
	TitleColor               string             `json:"title_color,omitempty"`
	TitleColorDark           string             `json:"title_color_dark,omitempty"`
	NormalBgColor            string             `json:"normal_bg_color,omitempty"`
	NormalBgColorDark        string             `json:"normal_bg_color_dark,omitempty"`
	StartColor               string             `json:"start_color,omitempty"`
	StartColorDark           string             `json:"start_color_dark,omitempty"`
	EndColor                 string             `json:"end_color,omitempty"`
	EndColorDark             string             `json:"end_color_dark,omitempty"`
	Pic                      string             `json:"pic,omitempty"`
	PicDark                  string             `json:"pic_dark,omitempty"`
	IsExpand                 int8               `json:"is_expand,omitempty"`
	SubType                  int8               `json:"sub_type"`
	ShowLoading              int8               `json:"show_loading,omitempty"`
	SkipFormat               int8               `json:"skip_format,omitempty"`
	Scheme                   string             `json:"scheme,omitempty"`
	RedirectScheme           string             `json:"redirect_scheme,omitempty"`
	BtnLeftIconUrl           string             `json:"btn_left_icon_url,omitempty"`
	Params                   map[string]any     `json:"params,omitempty"`
	ActionLog                map[string]any     `json:"actionlog,omitempty"`
	ExtUid                   string             `json:"ext_uid,omitempty"` // 同步签到状态，值需要和一键签到ext_button_uid_list、底导我的超话中contanerid一致，并且badge_img_type为1
	IsShieldAdLog            bool               `json:"is_shield_ad_log,omitempty"`
	CanFollow                int8               `json:"can_follow,omitempty"`
	CanUnfollow              *int8              `json:"can_unfollow,omitempty"`
	Relationship             int8               `json:"relationship,omitempty"`
	Disable                  bool               `json:"disable,omitempty"`
	JumpType                 string             `json:"jump_type,omitempty"`
	FollowWithoutSelectGroup bool               `json:"follow_without_select_group,omitempty"` // 关注后不出选择分组弹窗
	FollowRes                *ButtonFollowRes   `json:"follow_res,omitempty"`
	UnfollowRes              *ButtonUnfollowRes `json:"unfollow_res,omitempty"`
}

type ButtonFollowRes struct {
	Title      string `json:"title"`
	TitleColor string `json:"title_color"`
	Pic        string `json:"pic"`
	PicDark    string `json:"pic_dark"`
}

type ButtonUnfollowRes struct {
	Title string `json:"title"`
}
