package card

// C3 图片网格卡片
// 样式一：展示图片网格布局，支持标题文字显示
// 样式二：展示用户头像网格布局
// 样式三：展示图片网格布局，带名称和副标题
type C3 struct {
	BaseItem
	CardTypeName   string       `json:"card_type_name,omitempty"`   // 卡片类型名称
	AvatarSize     int          `json:"avatar_size,omitempty"`      // 头像大小 BB1生效
	TitleSize      int          `json:"title_size,omitempty"`       // 标题字体大小 BB1生效
	Title          string       `json:"title,omitempty"`            // 标题
	TitleExtraText string       `json:"title_extra_text,omitempty"` // 副标题 (791 起生效)
	ShowTitleArrow int8         `json:"show_title_arrow,omitempty"` // 标题左侧箭头 (791 起生效)
	ShowAvatar     int8         `json:"show_avatar,omitempty"`      // 是否显示头像 (791 起生效)
	RoundedCorner  int8         `json:"roundedcorner,omitempty"`    // 圆角
	DisplayArrow   int8         `json:"display_arrow,omitempty"`    // 显示箭头
	MaxItemCount   int          `json:"max_item_count,omitempty"`   // 默认是4个item，可以配置
	PicSum         string       `json:"pic_sum,omitempty"`          // 最后一张图显示数字或文字
	ShowLayer      int8         `json:"show_layer,omitempty"`       // 最后一张图是否显示蒙层
	Pics           []*C3Pic     `json:"pics,omitempty"`             // 图片列表
	Users          []*C3User    `json:"users,omitempty"`            // 用户列表 (791 起生效)
	Elements       []*C3Element `json:"elements,omitempty"`         // 显示列表 (791 起生效)
	FlagPic        string       `json:"flag_pic,omitempty"`         // 标记图片
	OpenUrl        string       `json:"open_url,omitempty"`         // 打开链接
}

// C3Pic 图片
type C3Pic struct {
	Pic       string         `json:"pic,omitempty"`       // 图片URL
	Desc1     string         `json:"desc1,omitempty"`     // 标题
	Desc2     string         `json:"desc2,omitempty"`     // 描述
	ActionLog map[string]any `json:"actionlog,omitempty"` // 打点日志
}

// C3User 用户
type C3User struct {
	Id    int64  `json:"id,omitempty"`
	Idstr string `json:"idstr,omitempty"`
}

// C3Element 元素
type C3Element struct {
	Uid        string         `json:"uid,omitempty"`         // 用户ID
	Desc1      string         `json:"desc1,omitempty"`       // 用户名
	Desc2      string         `json:"desc2,omitempty"`       // 描述
	Scheme     string         `json:"scheme,omitempty"`      // 跳转链接
	ActionLog  map[string]any `json:"actionlog,omitempty"`   // 打点日志
	Highlight  *C3Highlight   `json:"highlight,omitempty"`   // 高亮配置
	ShowBorder bool           `json:"show_border,omitempty"` // 头像是否显示内边框 BB1生效
}

// C3Highlight 高亮配置
type C3Highlight struct {
	DescFont int     `json:"desc_font,omitempty"` // 高亮字体大小 BB1生效
	DescEm   [][]int `json:"desc_em,omitempty"`   // 高亮位置
}

// NewCard3 创建Card3实例
func NewCard3() *C3 {
	return &C3{BaseItem: BaseItem{Base: Base{CardType: 3}}}
}
