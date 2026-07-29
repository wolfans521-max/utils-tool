package card

// C142 星球号推荐位理由卡片
// 星球号推荐位理由，以及可能的更多操作
type C142 struct {
	BaseItem
	RecommendDesc    string `json:"recommend_desc,omitempty"`     // 左侧文字
	OperTxt          string `json:"oper_txt,omitempty"`           // 非空展示
	NeedBottomBorder bool   `json:"need_bottom_border,omitempty"` // 是否添加底线
	OpenUrl          string `json:"openurl,omitempty"`            // 打开链接
}

// NewCard142 创建Card142实例
func NewCard142() *C142 {
	return &C142{BaseItem: BaseItem{Base: Base{CardType: 142}}}
}
