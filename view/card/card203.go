package card

// C203 超话新用户引导卡片
// 分为新手指引、活动模块、感兴趣的超话三个模块
type C203 struct {
	BaseItem
	CloseButton      *Button               `json:"close_button,omitempty"`      // 关闭按钮
	NewerGuide       *C203NewerGuide       `json:"newer_guide,omitempty"`       // 新手指引模块
	Activities       *C203Activities       `json:"activities,omitempty"`        // 活动模块
	InterestedTopics *C203InterestedTopics `json:"interested_topics,omitempty"` // 感兴趣的超话模块
	OpenUrl          string                `json:"openurl,omitempty"`           // 打开链接
}

// C203NewerGuide 新手指引
type C203NewerGuide struct {
	Scheme        string         `json:"scheme,omitempty"`          // 跳转链接
	ActionLog     map[string]any `json:"actionlog,omitempty"`       // actionlog
	GuidePic      string         `json:"guide_pic,omitempty"`       // 指引图片
	GuideText     string         `json:"guide_text,omitempty"`      // 指引文案
	GuideTextDark string         `json:"guide_text_dark,omitempty"` // 指引文案(暗黑)
	ArrowPic      string         `json:"arrow_pic,omitempty"`       // 指引文案后的右箭头
}

// C203Activities 活动模块
type C203Activities struct {
	Title                 string         `json:"title,omitempty"`                     // title
	TitlePre              string         `json:"title_pre,omitempty"`                 // title前半部分
	TitleTopic            string         `json:"title_topic,omitempty"`               // title中的超话词
	TitleNext             string         `json:"title_next,omitempty"`                // title后半部分
	TitleColor            string         `json:"title_color,omitempty"`               // title颜色
	TitleColorDark        string         `json:"title_color_dark,omitempty"`          // title颜色(暗黑)
	Detail                string         `json:"detail,omitempty"`                    // 副标题
	DetailColor           string         `json:"detail_color,omitempty"`              // 副标题颜色
	DetailColorDark       string         `json:"detail_color_dark,omitempty"`         // 副标题颜色(暗黑)
	Scheme                string         `json:"scheme,omitempty"`                    // 跳转链接
	ActionLog             map[string]any `json:"actionlog,omitempty"`                 // actionlog
	BgPic                 string         `json:"bg_pic,omitempty"`                    // 背景图
	BgPicStretchPx        int            `json:"bg_pic_stretch_px,omitempty"`         // 背景图片拉伸位置
	RightText             string         `json:"right_text,omitempty"`                // 右侧按钮文案
	RightButtonWidth      int            `json:"right_button_width,omitempty"`        // 右边按钮宽度
	RightTextColor        string         `json:"right_text_color,omitempty"`          // 右侧按钮文案颜色
	RightTextColorDark    string         `json:"right_text_color_dark,omitempty"`     // 右侧按钮文案颜色(暗黑)
	RightTextBg           string         `json:"right_text_bg,omitempty"`             // 右侧按钮背景颜色
	RightTextBgDark       string         `json:"right_text_bg_dark,omitempty"`        // 右侧按钮背景颜色(暗黑)
	RightTextBgColors     []string       `json:"right_text_bg_colors,omitempty"`      // 右侧按钮背景渐变颜色
	RightTextBgDarkColors []string       `json:"right_text_bg_dark_colors,omitempty"` // 右侧按钮背景渐变颜色(暗黑)
	LeftIcon              string         `json:"left_icon,omitempty"`                 // 左侧图片
	LeftIconWidth         int            `json:"left_icon_width,omitempty"`           // 左侧图片宽
	LeftIconHeight        int            `json:"left_icon_height,omitempty"`          // 左侧图片高
	RichTitle             any            `json:"rich_title,omitempty"`                // 富文本标题
}

// C203InterestedTopics 感兴趣的超话
type C203InterestedTopics struct {
	Title              string           `json:"title,omitempty"`                 // title
	TitleColor         string           `json:"title_color,omitempty"`           // title颜色
	TitleColorDark     string           `json:"title_color_dark,omitempty"`      // title颜色(暗黑)
	ArrowPic           string           `json:"arrow_pic,omitempty"`             // 右箭头
	ArrowText          string           `json:"arrow_text,omitempty"`            // arrow文案
	ArrowTextColor     string           `json:"arrow_text_color,omitempty"`      // arrow文案颜色
	ArrowTextColorDark string           `json:"arrow_text_color_dark,omitempty"` // arrow文案颜色(暗黑)
	Scheme             string           `json:"scheme,omitempty"`                // 跳转链接
	ActionLog          map[string]any   `json:"actionlog,omitempty"`             // actionlog
	Items              []*C203TopicItem `json:"items,omitempty"`                 // 感兴趣的超话list
}

// C203TopicItem 超话项
type C203TopicItem struct {
	Title               string         `json:"title,omitempty"`                  // title
	TitleColor          string         `json:"title_color,omitempty"`            // title颜色
	TitleColorDark      string         `json:"title_color_dark,omitempty"`       // title颜色(暗黑)
	Pic                 string         `json:"pic,omitempty"`                    // 超话头像
	Friends             []string       `json:"friends,omitempty"`                // 关注该超话的好友头像列表
	FollowText          string         `json:"follow_text,omitempty"`            // 提示文案(已关注)
	FollowTextColor     string         `json:"follow_text_color,omitempty"`      // 提示文案颜色
	FollowTextColorDark string         `json:"follow_text_color_dark,omitempty"` // 提示文案颜色(暗黑)
	FollowButton        *Button        `json:"follow_button,omitempty"`          // 关注按钮
	ActionLog           map[string]any `json:"actionlog,omitempty"`              // actionlog
	Scheme              string         `json:"scheme,omitempty"`                 // 跳转链接
	ItemId              string         `json:"item_id,omitempty"`                // 真实曝光itemid
	ActCode             string         `json:"act_code,omitempty"`               // 真实曝光act_code
}

// NewCard203 创建Card203实例
func NewCard203() *C203 {
	return &C203{BaseItem: BaseItem{Base: Base{CardType: 203}}}
}
