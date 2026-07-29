package card

// C87 横滑卡片项
// 作为item，用于card86的横滑上，可用于展示文字、图片、视频（文字支持两行）
type C87 struct {
	BaseItem
	CardCommTitle string         `json:"card_comm_title,omitempty"` // 卡片通用标题
	TitleSub      string         `json:"title_sub,omitempty"`       // 标题
	Pic           string         `json:"pic,omitempty"`             // 图片
	PicUrl        string         `json:"pic_url,omitempty"`         // 图片URL
	UnreadId      string         `json:"unread_id,omitempty"`       // 未读ID
	ObjectId      string         `json:"object_id,omitempty"`       // 对象ID
	ObjectType    string         `json:"object_type,omitempty"`     // 对象类型
	ActStatus     int8           `json:"act_status,omitempty"`      // 状态
	MediaInfo     map[string]any `json:"media_info,omitempty"`      // 媒体信息（下游控制，使用map类型）
	OpenUrl       string         `json:"openurl,omitempty"`         // 打开链接
}

// NewCard87 创建Card87实例
func NewCard87() *C87 {
	return &C87{BaseItem: BaseItem{Base: Base{CardType: 87}}}
}
