package card

// C159 瀑布流潮流卡片
// 用于瀑布流潮流展示
type C159 struct {
	BaseItem
	RightTagIcon       string     `json:"right_tag_icon,omitempty"`       // 右上角图片
	RightTagIconDark   string     `json:"right_tag_icon_dark,omitempty"`  // 右上角图片暗黑模式
	RightTagBackground string     `json:"right_tag_background,omitempty"` // 右上角背景
	RightTagText       string     `json:"right_tag_text,omitempty"`       // 右上角图片右侧文字
	PicWidth           int        `json:"pic_width,omitempty"`            // 背景图宽
	PicHeight          int        `json:"pic_height,omitempty"`           // 背景图高
	RightBottomText    string     `json:"right_bottom_text,omitempty"`    // 背景图右下角文案
	AdjustHeight       bool       `json:"adjust_height,omitempty"`        // 图片是否自适应宽高
	CoverImage         *C159Cover `json:"cover_image,omitempty"`          // 背景图
	PicScheme          string     `json:"pic_scheme,omitempty"`           // 点击背景图跳转的scheme
	Mblog              any        `json:"mblog,omitempty"`                // Status对象
	TextLines          int        `json:"text_lines,omitempty"`           // 标题强制显示多少行
	TextMaxLine        int        `json:"text_max_line,omitempty"`        // 标题最多显示多少行
	TagIcons           []string   `json:"tag_icons,omitempty"`            // 标签图标
	Text               string     `json:"text,omitempty"`                 // 文本
	TextScheme         string     `json:"text_scheme,omitempty"`          // 文本跳转链接
	Buttons            []*Button  `json:"buttons,omitempty"`              // 整个card右下角评论/点赞可配
	Mid                string     `json:"mid,omitempty"`                  // 微博MID
	User               any        `json:"user,omitempty"`                 // 用户信息
	AttitudesCount     int        `json:"attitudes_count,omitempty"`      // 点赞数
	ItemType           int8       `json:"item_type,omitempty"`            // 项目类型
	PageInfo           any        `json:"page_info,omitempty"`            // 页面信息
	PicRecommendIndex  int        `json:"pic_recommend_index,omitempty"`  // 图片推荐索引
	AdBrandStyle       bool       `json:"ad_brand_style,omitempty"`       // 广告业务
	PicList            []string   `json:"pic_list,omitempty"`             // 广告业务图片列表
	OpenUrl            string     `json:"openurl,omitempty"`              // 打开链接
}

// C159Cover 封面
type C159Cover struct {
	Url     string `json:"url,omitempty"`      // URL
	CutType int8   `json:"cut_type,omitempty"` // 裁剪类型
}

// NewCard159 创建Card159实例
func NewCard159() *C159 {
	return &C159{BaseItem: BaseItem{Base: Base{CardType: 159}}}
}
