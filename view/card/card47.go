package card

// C47 图片网格卡片（支持标签）
// 主要用在查看图片及视频card
// 使用场景：发博-图片-图片及视频查看页
type C47 struct {
	BaseItem
	CardTypeName string    `json:"card_type_name,omitempty"` // 卡片类型名称
	Title        string    `json:"title,omitempty"`          // 标题信息
	Spacing      int       `json:"spacing,omitempty"`        // 图片之间间隔
	Pics         []*C47Pic `json:"pics,omitempty"`           // 图片墙
	OpenUrl      string    `json:"openurl,omitempty"`        // 打开链接
}

// C47Pic 图片
type C47Pic struct {
	PicSmall  string         `json:"pic_small,omitempty"`  // 小图
	PicMiddle string         `json:"pic_middle,omitempty"` // 中图
	PicBig    string         `json:"pic_big,omitempty"`    // 大图
	LayerText string         `json:"layer_text,omitempty"` // 蒙层文字
	Mblog     *C47Mblog      `json:"mblog,omitempty"`      // 微博数据
	ObjectId  string         `json:"object_id,omitempty"`  // 对象ID
	Scheme    string         `json:"scheme,omitempty"`     // 跳转链接
	Corner    int            `json:"corner,omitempty"`     // 图片圆角
	ActionLog map[string]any `json:"actionlog,omitempty"`  // 打点日志
}

// C47Mblog 微博数据
type C47Mblog struct {
	Id     string `json:"id,omitempty"`     // 微博ID
	Mid    string `json:"mid,omitempty"`    // 微博MID
	Text   string `json:"text,omitempty"`   // 微博文本
	Scheme string `json:"scheme,omitempty"` // 跳转链接
}

// NewCard47 创建Card47实例
func NewCard47() *C47 {
	return &C47{BaseItem: BaseItem{Base: Base{CardType: 47}}}
}
