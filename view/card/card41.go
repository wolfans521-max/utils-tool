package card

// C41 多功能信息卡片
// 显示简介、会员、认证、设置备注、好友关系展示、作品集展示、标签展示、徽章展示
type C41 struct {
	BaseItem
	ItemName         string      `json:"item_name,omitempty"`          // 项目名称
	ItemContent      string      `json:"item_content,omitempty"`       // 项目内容
	ItemType         string      `json:"item_type,omitempty"`          // 项目类型
	ContentColorType int8        `json:"content_color_type,omitempty"` // item_content颜色 0默认 1蓝色
	Uid              string      `json:"uid,omitempty"`                // 用户ID
	DisplayArrow     string      `json:"display_arrow,omitempty"`      // 显示箭头
	Mbtype           int8        `json:"mbtype,omitempty"`             // 会员类型
	Mbrank           int8        `json:"mbrank,omitempty"`             // 会员等级
	Badges           []*C41Badge `json:"badges,omitempty"`             // 徽章列表
	Tags             []*C41Tag   `json:"tags,omitempty"`               // 标签列表
	ItemPic          *C41ItemPic `json:"item_pic,omitempty"`           // 作品集
	ItemUsers        []*C41User  `json:"item_users,omitempty"`         // 好友关系用户列表
	OpenUrl          string      `json:"openurl,omitempty"`            // 打开链接
}

// C41Badge 徽章
type C41Badge struct {
	PicUrl string `json:"pic_url,omitempty"` // 徽章图片URL
}

// C41Tag 标签
type C41Tag struct {
	DisplayName string `json:"display_name,omitempty"`  // 显示名称
	ShowEditTag int8   `json:"show_edit_tag,omitempty"` // 显示编辑标签
	TagScheme   string `json:"tag_scheme,omitempty"`    // 标签跳转
	BorderColor string `json:"border_color,omitempty"`  // 边框颜色
	TextColor   string `json:"text_color,omitempty"`    // 文字颜色
	BgColor     string `json:"bg_color,omitempty"`      // 背景颜色
}

// C41ItemPic 作品集
type C41ItemPic struct {
	Pics      []*C41Pic `json:"pics,omitempty"`       // 图片列表
	PicSum    string    `json:"pic_sum,omitempty"`    // 图片数量
	ShowLayer int8      `json:"show_layer,omitempty"` // 显示蒙层
}

// C41Pic 图片
type C41Pic struct {
	Pic   string `json:"pic,omitempty"`   // 图片URL
	Desc1 string `json:"desc1,omitempty"` // 描述
}

// C41User 用户
type C41User struct {
	Id    int64  `json:"id,omitempty"`
	Idstr string `json:"idstr,omitempty"`
	Class int8   `json:"class,omitempty"`
}

// NewCard41 创建Card41实例
func NewCard41() *C41 {
	return &C41{BaseItem: BaseItem{Base: Base{CardType: 41}}}
}
