package card

// C88 双列直播卡片
// 左右两侧各展示一个直播内容
type C88 struct {
	BaseItem
	LeftElement  *C88Element `json:"left_element,omitempty"`  // 左边元素
	RightElement *C88Element `json:"right_element,omitempty"` // 右边元素
	UnreadId     string      `json:"unread_id,omitempty"`     // 未读ID
	OpenUrl      string      `json:"openurl,omitempty"`       // 打开链接
}

// C88Element 元素
type C88Element struct {
	Mblog          any            `json:"mblog,omitempty"`            // 微博数据
	Recommend      string         `json:"recommend,omitempty"`        // 推荐标签，如"星主播"
	Position       string         `json:"position,omitempty"`         // 位置信息
	IsIconShow     int8           `json:"is_icon_show,omitempty"`     // 为1时显示播放按钮icon
	LiveBottomInfo *C88BottomInfo `json:"live_bottom_info,omitempty"` // 直播底部信息
}

// C88BottomInfo 直播底部信息
type C88BottomInfo struct {
	Title  string `json:"title,omitempty"`  // 标题
	Scheme string `json:"scheme,omitempty"` // 跳转链接
}

// NewCard88 创建Card88实例
func NewCard88() *C88 {
	return &C88{BaseItem: BaseItem{Base: Base{CardType: 88}}}
}
