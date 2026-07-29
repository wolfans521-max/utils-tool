package card

// C90 评论卡片
// 用于展示评论内容
type C90 struct {
	BaseItem
	Comment *C90Comment `json:"comment,omitempty"` // 评论通用结构
	OpenUrl string      `json:"openurl,omitempty"` // 打开链接
}

// C90Comment 评论
type C90Comment struct {
	CreatedAt        string `json:"created_at,omitempty"`        // 创建时间
	Id               int64  `json:"id,omitempty"`                // 评论ID
	RootId           int64  `json:"rootid,omitempty"`            // 根评论ID
	FloorNumber      int    `json:"floor_number,omitempty"`      // 楼层号
	Text             string `json:"text,omitempty"`              // 评论文本
	SourceAllowClick int8   `json:"source_allowclick,omitempty"` // 来源是否可点击
	SourceType       int8   `json:"source_type,omitempty"`       // 来源类型
	Source           string `json:"source,omitempty"`            // 来源
	User             any    `json:"user,omitempty"`              // 用户信息
	Mid              string `json:"mid,omitempty"`               // 微博MID
	Idstr            string `json:"idstr,omitempty"`             // ID字符串
	Liked            bool   `json:"liked,omitempty"`             // 是否已点赞
	LikeCounts       int    `json:"like_counts,omitempty"`       // 点赞数
}

// NewCard90 创建Card90实例
func NewCard90() *C90 {
	return &C90{BaseItem: BaseItem{Base: Base{CardType: 90}}}
}
