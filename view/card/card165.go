package card

// C165 一起看视频卡片
// 用于一起看视频展示
type C165 struct {
	BaseItem
	SinceId                    string           `json:"since_id,omitempty"`                      // 通过下方icons打开半屏页时传递的参数
	Mblog                      any              `json:"mblog,omitempty"`                         // 微博结构
	Desc                       string           `json:"desc,omitempty"`                          // 用户认证信息
	Desc2                      string           `json:"desc2,omitempty"`                         // 视频标题下方标签文本
	DescScheme                 string           `json:"desc_scheme,omitempty"`                   // 视频标题下方标签点击跳转scheme
	CoverDecorIcon             *C165DecorIcon   `json:"cover_decor_icon,omitempty"`              // 角标对象
	AvatarStack                *C165AvatarStack `json:"avatar_stack,omitempty"`                  // 头像堆叠
	Icons                      []*C165Icon      `json:"icons,omitempty"`                         // 右下角的icon数组
	FlowCommentPosition        string           `json:"flow_comment_position,omitempty"`         // 滚动热评默认的位置
	FlowCommentDefaultDisplay  int8             `json:"flow_comment_default_display,omitempty"`  // 滚动热评默认是否显示
	RightBottomButtonsDuration int              `json:"right_bottom_buttons_duration,omitempty"` // 右下角按钮展示时间
	HotComment                 any              `json:"hot_comment,omitempty"`                   // 评论对象
	LeftIcon                   *C165LeftIcon    `json:"left_icon,omitempty"`                     // 左侧动画icon
	Strategy                   *C165Strategy    `json:"strategy,omitempty"`                      // 策略配置
	CommentGuide               *C165Guide       `json:"comment_guide,omitempty"`                 // 评论引导
	OpenUrl                    string           `json:"openurl,omitempty"`                       // 打开链接
}

// C165DecorIcon 角标
type C165DecorIcon struct {
	Img string `json:"img,omitempty"` // 角标图片地址
}

// C165AvatarStack 头像堆叠
type C165AvatarStack struct {
	Avatars    []*C165Avatar `json:"avatars,omitempty"`     // 头像列表
	TotalCount int           `json:"total_count,omitempty"` // 总数
	Name       string        `json:"name,omitempty"`        // 名称
	TagCount   string        `json:"tag_count,omitempty"`   // 标签数
}

// C165Avatar 头像
type C165Avatar struct {
	Name string `json:"name,omitempty"` // 名称
	Url  string `json:"url,omitempty"`  // URL
}

// C165Icon 图标
type C165Icon struct {
	Title           string `json:"title,omitempty"`             // 右下角的icon的文本
	UrlNormal       string `json:"url_normal,omitempty"`        // 右下角的icon普通态url
	UrlSelected     string `json:"url_selected,omitempty"`      // 右下角的icon的激活态url
	UrlNormalDark   string `json:"url_normal_dark,omitempty"`   // 右下角的icon普通态深色模式url
	UrlSelectedDark string `json:"url_selected_dark,omitempty"` // 右下角的icon的激活态深色模式url
	IconType        int8   `json:"icon_type,omitempty"`         // 图标类型
	Scheme          string `json:"scheme,omitempty"`            // 跳转链接
	AttitudesStatus int8   `json:"attitudes_status,omitempty"`  // icon是否激活
}

// C165LeftIcon 左侧图标
type C165LeftIcon struct {
	Text              string         `json:"text,omitempty"`                // 滚动文本内容
	Icon              string         `json:"icon,omitempty"`                // 左侧icon默认模式url地址
	IconDark          string         `json:"icon_dark,omitempty"`           // 左侧icon深色模式url
	Scheme            string         `json:"scheme,omitempty"`              // 左侧icon和文本点击跳转scheme
	ActionLog         map[string]any `json:"actionlog,omitempty"`           // 左侧按钮点击后上报的行为码
	TitleShowDuration int            `json:"title_show_duration,omitempty"` // 展示时间
}

// C165Strategy 策略
type C165Strategy struct {
	VoiceButtonEnable       int8 `json:"voice_button_enable,omitempty"`        // 声音开关
	ContentTextEnable       int8 `json:"content_text_enable,omitempty"`        // 展示正文
	TimeShowEnable          int8 `json:"time_show_enable,omitempty"`           // 显示时间
	FullScreenButtonEnable  int8 `json:"full_screen_button_enable,omitempty"`  // 全屏在下
	FlowCommentLoopInterval int  `json:"flow_comment_loop_interval,omitempty"` // 热评循环展示前的时间间隔
	SeekBarShowTime         int  `json:"seek_bar_show_time,omitempty"`         // 进度条自动显示时间
	VoiceMuteEnable         int8 `json:"voice_mute_enable,omitempty"`          // 是否静音
}

// C165Guide 引导
type C165Guide struct {
	ShowTime int `json:"show_time,omitempty"` // 展示开始时间
}

// NewCard165 创建Card165实例
func NewCard165() *C165 {
	return &C165{BaseItem: BaseItem{Base: Base{CardType: 165}}}
}
