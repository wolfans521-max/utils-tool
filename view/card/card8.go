package card

// C8 左图右文卡片
// 支持内嵌标题
// 支持左图右三行文字
// 支持右侧CommonButton
// 支持标题扩展图片
type C8 struct {
	BaseItem
	CardTypeName       string         `json:"card_type_name,omitempty"`        // 卡片类型名称
	Title              string         `json:"title,omitempty"`                 // 大标题
	TitleSize          string         `json:"title_size,omitempty"`            // 大标题字体
	TitleExtraText     string         `json:"title_extra_text,omitempty"`      // 副标题
	ShowTitleArrow     int8           `json:"show_title_arrow,omitempty"`      // 标题右侧箭头
	DisplayArrow       int8           `json:"display_arrow,omitempty"`         // 显示箭头
	Pic                string         `json:"pic,omitempty"`                   // 左图
	TitleSub           string         `json:"title_sub,omitempty"`             // 标题
	Desc1              string         `json:"desc1,omitempty"`                 // 正文内容1
	Desc2              string         `json:"desc2,omitempty"`                 // 正文内容2
	CardDisplayType    int8           `json:"card_display_type,omitempty"`     // 卡片显示类型
	Buttons            []*Button      `json:"buttons,omitempty"`               // 按钮列表
	OpenUrl            string         `json:"openurl,omitempty"`               // 打开链接
	DescPicUrl         string         `json:"desc_pic_url,omitempty"`          // 描述右边显示"app专享icon"
	TitleFlagPicList   []string       `json:"title_flag_pic_list,omitempty"`   // 标题多icon列表
	RightPannel        *C8RightPannel `json:"right_pannel,omitempty"`          // 右侧面板
	TopPadding         int            `json:"top_padding,omitempty"`           // 上边距
	BottomPadding      int            `json:"bottom_padding,omitempty"`        // 下边距
	CardHeightStyle    int8           `json:"card_height_style,omitempty"`     // 高度样式 0:默认 1：创作者中心
	IsInside           bool           `json:"is_inside,omitempty"`             // 图片是否撑满
	TopMarkPic         string         `json:"top_mark_pic,omitempty"`          // pic上角标
	TopMarkPicTop      int            `json:"top_mark_pic_top,omitempty"`      // 角标到pic顶部距离
	TopMarkPicLeft     int            `json:"top_mark_pic_left,omitempty"`     // 角标到pic左边距离
	TopMarkPicWidth    int            `json:"top_mark_pic_width,omitempty"`    // 角标宽度
	TopMarkPicHeight   int            `json:"top_mark_pic_height,omitempty"`   // 角标高度
	PicBorderColor     string         `json:"pic_border_color,omitempty"`      // 头像描边颜色
	PicBorderColorDark string         `json:"pic_border_color_dark,omitempty"` // 头像描边颜色（暗黑模式）
	PicBorderWidth     float64        `json:"pic_border_width,omitempty"`      // 头像描边宽度
	PicCornerRadius    int            `json:"pic_corner_radius,omitempty"`     // 图片圆角
	PicSize            int            `json:"pic_size,omitempty"`              // 头像大小
	NameFontSize       int            `json:"name_font_size,omitempty"`        // 标题字体大小
}

// C8RightPannel 右侧面板
type C8RightPannel struct {
	Type      int8     `json:"type,omitempty"`       // 类型
	FlagImgs  []string `json:"flag_imgs,omitempty"`  // 标记图片列表
	Desc      string   `json:"desc,omitempty"`       // 描述
	MiddleImg string   `json:"middle_img,omitempty"` // 中间图片
}

// NewCard8 创建Card8实例
func NewCard8() *C8 {
	return &C8{BaseItem: BaseItem{Base: Base{CardType: 8}}}
}
