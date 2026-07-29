package card

// C59 双列微博对比卡片
// 左右两侧各展示一条微博内容
type C59 struct {
	BaseItem
	LeftElement  *C59Element `json:"left_element,omitempty"`  // 左边card
	RightElement *C59Element `json:"right_element,omitempty"` // 右边card
	OpenUrl      string      `json:"openurl,omitempty"`       // 打开链接
}

// C59Element 元素
type C59Element struct {
	Mblog   *C59Mblog `json:"mblog,omitempty"`   // 微博数据
	Desc1   string    `json:"desc1,omitempty"`   // 描述1
	Desc2   string    `json:"desc2,omitempty"`   // 描述2
	Buttons []*Button `json:"buttons,omitempty"` // 底部可配按钮
}

// C59Mblog 微博数据
type C59Mblog struct {
	CreatedAt        string   `json:"created_at,omitempty"`        // 创建时间
	Id               int64    `json:"id,omitempty"`                // 微博ID
	Mid              string   `json:"mid,omitempty"`               // 微博MID
	Idstr            string   `json:"idstr,omitempty"`             // 微博ID字符串
	Text             string   `json:"text,omitempty"`              // 微博文本
	SourceAllowClick int8     `json:"source_allowclick,omitempty"` // 来源是否可点击
	SourceType       int8     `json:"source_type,omitempty"`       // 来源类型
	Source           string   `json:"source,omitempty"`            // 来源
	Favorited        bool     `json:"favorited,omitempty"`         // 是否收藏
	PicIds           []string `json:"pic_ids,omitempty"`           // 图片ID列表
	MblogId          string   `json:"mblogid,omitempty"`           // 微博ID
	Scheme           string   `json:"scheme,omitempty"`            // 跳转链接
	MblogTypeName    string   `json:"mblogtypename,omitempty"`     // 微博类型名称
	AttitudesStatus  int8     `json:"attitudes_status,omitempty"`  // 点赞状态
	RecomState       int8     `json:"recom_state,omitempty"`       // 推荐状态
}

// NewCard59 创建Card59实例
func NewCard59() *C59 {
	return &C59{BaseItem: BaseItem{Base: Base{CardType: 59}}}
}
