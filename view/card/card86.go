package card

// C86 横滑卡片容器
// 用于展示横滑轮播图
type C86 struct {
	BaseItem
	PagePadding        int    `json:"page_padding,omitempty"`       // 轮播图内部左右侧间距，单位dp
	NotSyncFollowState int8   `json:"notSyncFollowState,omitempty"` // 接受关注通知是否同步删除card (ios使用)
	SubCards           []any  `json:"sub_cards,omitempty"`          // 子卡片列表
	OpenUrl            string `json:"openurl,omitempty"`            // 打开链接
}

// NewCard86 创建Card86实例
func NewCard86() *C86 {
	return &C86{BaseItem: BaseItem{Base: Base{CardType: 86}}}
}
