package card

// C35 宫格图片卡片
// 支持多行3列的宫格图片
// 支持图片2行描述
// 支持card标题
type C35 struct {
	BaseItem
	Title   string    `json:"title,omitempty"`   // 标题
	Cols    int       `json:"cols,omitempty"`    // 列数
	Pics    []*C35Pic `json:"pics,omitempty"`    // 宫格图片
	OpenUrl string    `json:"openurl,omitempty"` // 打开链接
}

// C35Pic 图片项
type C35Pic struct {
	PicSmall string `json:"pic_small,omitempty"` // 小图
	PicBig   string `json:"pic_big,omitempty"`   // 大图
	Desc1    string `json:"desc1,omitempty"`     // 第一段描述（支持2行）
	Desc2    string `json:"desc2,omitempty"`     // 第二段描述（支持1行）
	Scheme   string `json:"scheme,omitempty"`    // 跳转链接
}

// NewCard35 创建Card35实例
func NewCard35() *C35 {
	return &C35{BaseItem: BaseItem{Base: Base{CardType: 35}}}
}
