package card

// C124 视频榜单轮播卡片
// 出现场景：美食频道
type C124 struct {
	BaseItem
	Items               string      `json:"items,omitempty"`                 // itemid必须有
	Lists               []*C124List `json:"lists,omitempty"`                 // 榜单列表
	ThemeColor          string      `json:"theme_color,omitempty"`           // 主题色值
	ForegroundColor     string      `json:"foregroud_color,omitempty"`       // 各标题文字色值
	TopImg              string      `json:"top_img,omitempty"`               // 广告区，顶部图片区
	ArrowDown           string      `json:"arrow_down,omitempty"`            // 隐藏图标列表的图标
	ArrowUp             string      `json:"arrow_up,omitempty"`              // 展开列表的图标
	ButtonBackground    string      `json:"button_background,omitempty"`     // 左侧按钮和弹幕开关背景色
	ButtonTextColor     string      `json:"button_textcolor,omitempty"`      // 左侧按钮文字颜色
	PlayingChannelIndex int         `json:"playing_channel_index,omitempty"` // 播放第几个频道
	PlayingVideoIndex   int         `json:"playing_video_index,omitempty"`   // 播放第几个视频
	IconKeyboard        string      `json:"icon_keyboard,omitempty"`         // 输入法顶部悬浮输入框右侧图标，键盘图标
	IconEmoji           string      `json:"icon_emoji,omitempty"`            // 输入法顶部悬浮输入框右侧图标，表情图标
	DanmakuOn           string      `json:"danmaku_on,omitempty"`            // 弹幕打开图标
	DanmakuOff          string      `json:"danmaku_off,omitempty"`           // 弹幕关闭图标
	NeedTitle           int8        `json:"need_title,omitempty"`            // 是否需要标题
	OpenUrl             string      `json:"openurl,omitempty"`               // 打开链接
}

// C124List 榜单
type C124List struct {
	Title    string `json:"title,omitempty"`    // 标题
	Order    int    `json:"order,omitempty"`    // 排序
	Statuses []any  `json:"statuses,omitempty"` // 微博数据列表
}

// NewCard124 创建Card124实例
func NewCard124() *C124 {
	return &C124{BaseItem: BaseItem{Base: Base{CardType: 124}}}
}
