package card

// C24 多功能卡片
// type=0: 支持显示最多7个头像
// type=1: 支持显示多标签
// type=2: 支持显示profile多行标签
// type=3: 搜索历史数据
// type=4: 带左标题的5个头像，支持展开显示推荐用户
type C24 struct {
	BaseItem
	Type           int8              `json:"type,omitempty"`             // 0为7头像，1为标签card，2为多行标签，3为搜索历史，4为带左标题头像
	CardTypeName   string            `json:"card_type_name,omitempty"`   // 卡片类型名称
	Title          string            `json:"title,omitempty"`            // 标题
	TitleExtraText string            `json:"title_extra_text,omitempty"` // 副标题
	ShowTitleArrow int8              `json:"show_title_arrow,omitempty"` // 标题右箭头
	ShowType       int8              `json:"show_type,omitempty"`        // 显示类型
	RoundedCorner  int8              `json:"roundedcorner,omitempty"`    // 圆角
	WeiboNeed      string            `json:"weibo_need,omitempty"`       // 微博需要
	ProfileTag     *C24ProfileTag    `json:"profile_tag,omitempty"`      // profile多行标签
	ItemsFeature   *C24ItemsFeature  `json:"items_feature,omitempty"`    // 单行可滑动标签
	Elements       []*C24Element     `json:"elements,omitempty"`         // 头像元素列表
	Users          []*C24User        `json:"users,omitempty"`            // 用户列表
	Items          []string          `json:"items,omitempty"`            // 搜索历史数据
	DataFrom       int8              `json:"data_from,omitempty"`        // 数据来源
	Margin         []int             `json:"margin,omitempty"`           // 控制card左右边距
	ShowTopPadding int8              `json:"show_top_padding,omitempty"` // 上边距
	ItemName       string            `json:"item_name,omitempty"`        // 左标题字段
	Uid            string            `json:"uid,omitempty"`              // 展开显示推荐用户列表的请求字段
	ArrowStruct    map[string]string `json:"arrow_struct,omitempty"`     // 请求推荐用户数据的请求字段
	OpenUrl        string            `json:"openurl,omitempty"`          // 打开链接
}

// C24ProfileTag profile标签
type C24ProfileTag struct {
	Title            string    `json:"title,omitempty"`              // 标题
	FoldLineCount    int       `json:"fold_line_count,omitempty"`    // 折叠行数
	UnfoldLineCount  int       `json:"unfold_line_count,omitempty"`  // 展开行数
	TopBottomPadding int       `json:"top_bottom_padding,omitempty"` // 上下间距
	Tags             []*C24Tag `json:"tags,omitempty"`               // 标签列表
	MoreTag          *C24Tag   `json:"moreTag,omitempty"`            // 更多标签
	TagUrl           string    `json:"tag_url,omitempty"`            // tag点击跳转地址
}

// C24Tag 标签
type C24Tag struct {
	DisplayName string `json:"display_name,omitempty"`  // 显示名称
	ShowEditTag int8   `json:"show_edit_tag,omitempty"` // 显示编辑标签
	TagScheme   string `json:"tag_scheme,omitempty"`    // 标签跳转
}

// C24ItemsFeature 单行可滑动标签
type C24ItemsFeature struct {
	Items    []*C24LabelItem `json:"items,omitempty"`    // 胶囊数据集合
	Padding  []int           `json:"padding,omitempty"`  // 控制胶囊滑动的左右边距
	Distance int             `json:"distance,omitempty"` // 控制胶囊之间的间距
}

// C24LabelItem 胶囊数据
type C24LabelItem struct {
	Type        int8           `json:"type,omitempty"`         // 类型
	Icon        string         `json:"icon,omitempty"`         // 胶囊左边的图标地址
	IconDark    string         `json:"icon_dark,omitempty"`    // 暗黑模式图标
	DisplayName string         `json:"display_name,omitempty"` // 胶囊文案
	Scheme      string         `json:"scheme,omitempty"`       // 胶囊跳转
	ActionLog   map[string]any `json:"actionlog,omitempty"`    // 打点日志
}

// C24Element 头像元素
type C24Element struct {
	Uid       string         `json:"uid,omitempty"`       // 用户ID
	Scheme    string         `json:"scheme,omitempty"`    // 跳转链接
	ActionLog map[string]any `json:"actionlog,omitempty"` // 打点日志
}

// C24User 用户信息
type C24User struct {
	Id              int64  `json:"id,omitempty"`
	ScreenName      string `json:"screen_name,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
	AvatarLarge     string `json:"avatar_large,omitempty"`
	Verified        bool   `json:"verified,omitempty"`
	VerifiedType    int8   `json:"verified_type,omitempty"`
	VerifiedTypeExt int8   `json:"verified_type_ext,omitempty"`
	Level           int8   `json:"level,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowersCount  int64  `json:"followers_count,omitempty"`
	FriendsCount    int64  `json:"friends_count,omitempty"`
	Mbtype          int8   `json:"mbtype,omitempty"`
	Mbrank          int8   `json:"mbrank,omitempty"`
	Remark          string `json:"remark,omitempty"`
	Following       bool   `json:"following,omitempty"`
}

// NewCard24 创建Card24实例
func NewCard24() *C24 {
	return &C24{BaseItem: BaseItem{Base: Base{CardType: 24}}}
}
