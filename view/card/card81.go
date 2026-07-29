package card

// C81 直播卡片
// 作为card容器，仅在profile页精选第一条、和部分搜索结果使用
type C81 struct {
	BaseItem
	HideTopProfile int8           `json:"hide_top_profile,omitempty"` // 隐藏上方头像+文字+发布时间+来源
	DisplayArrow   int8           `json:"display_arrow,omitempty"`    // 显示箭头
	CardCommTitle  string         `json:"card_comm_title,omitempty"`  // 卡片通用标题
	HideBtns       string         `json:"hidebtns,omitempty"`         // 隐藏按钮
	Title          string         `json:"title,omitempty"`            // 标题
	VideoCorner    int            `json:"video_corner,omitempty"`     // 播放器圆角，单位dp
	Mblog          any            `json:"mblog,omitempty"`            // 微博数据
	UnreadId       string         `json:"unread_id,omitempty"`        // 未读ID
	LiveBottomInfo *C81BottomInfo `json:"live_bottom_info,omitempty"` // 直播底部信息
	OpenUrl        string         `json:"openurl,omitempty"`          // 打开链接
}

// C81BottomInfo 直播底部信息
type C81BottomInfo struct {
	Scheme string `json:"scheme,omitempty"` // 点击title跳转地址
	Title  string `json:"title,omitempty"`  // title字段
}

// NewCard81 创建Card81实例
func NewCard81() *C81 {
	return &C81{BaseItem: BaseItem{Base: Base{CardType: 81}}}
}
