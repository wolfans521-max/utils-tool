package card

// C292Image 神评卡片图片对象
type C292Image struct {
	Url             string  `json:"url,omitempty"`
	UrlDark         string  `json:"url_dark,omitempty"`
	Width           float64 `json:"width,omitempty"`
	Height          float64 `json:"height,omitempty"`
	RadiusCorner    float64 `json:"radius_corner,omitempty"`
	BorderColor     string  `json:"border_color,omitempty"`
	BorderColorDark string  `json:"border_color_dark,omitempty"`
	BorderWidth     float64 `json:"border_width,omitempty"`
	ContentMode     string  `json:"content_mode,omitempty"`
}

// C292Text 神评卡片文本对象
type C292Text struct {
	Text          string  `json:"text,omitempty"`
	TextSize      int     `json:"text_size,omitempty"`
	TextColor     string  `json:"text_color,omitempty"`
	TextColorDark string  `json:"text_color_dark,omitempty"`
	Width         float64 `json:"width,omitempty"`
	Height        float64 `json:"height,omitempty"`
}

// C292UserObj 神评卡片用户对象
type C292UserObj struct {
	AvataImage *C292Image `json:"avata_image,omitempty"`
	Name       *C292Text  `json:"name,omitempty"`
}

// C292CommentObj 神评卡片评论对象
type C292CommentObj struct {
	Text          string `json:"text,omitempty"`
	TextSize      int    `json:"text_size,omitempty"`
	TextColor     string `json:"text_color,omitempty"`
	TextColorDark string `json:"text_color_dark,omitempty"`
}

// C292CommentItem 神评卡片评论项
type C292CommentItem struct {
	UserObj         *C292UserObj    `json:"user_obj,omitempty"`
	CommentObj      *C292CommentObj `json:"comment_obj,omitempty"`
	Scheme          string          `json:"scheme,omitempty"`
	ExposeActionLog any             `json:"expose_actionlog,omitempty"`
	ClickActionLog  any             `json:"click_actionlog,omitempty"`
}

// C292ControlObj 神评卡片折叠控制对象
type C292ControlObj struct {
	NeedFold         bool                 `json:"need_fold,omitempty"`
	FoldSeconds      int                  `json:"fold_seconds,omitempty"`
	FoldHeightAdjust int                  `json:"fold_height_adjust,omitempty"`
	FoldToDic        map[string]string    `json:"fold_to_dic,omitempty"`
}

// C292Mask 神评卡片遮罩对象
type C292Mask struct {
	Url     string  `json:"url,omitempty"`
	UrlDark string  `json:"url_dark,omitempty"`
	Height  float64 `json:"height,omitempty"`
}

// C292Left 神评卡片左侧动画对象
type C292Left struct {
	PaddingVertical   int          `json:"padding_vertical,omitempty"`
	PaddingHorizontal int          `json:"padding_horizontal,omitempty"`
	AnmationInterval  float64      `json:"anmation_interval,omitempty"`
	Images            []*C292Image `json:"images,omitempty"`
}

// C292Right 神评卡片右侧箭头对象
type C292Right struct {
	PaddingVertical   int        `json:"padding_vertical,omitempty"`
	PaddingHorizontal int        `json:"padding_horizontal,omitempty"`
	ImageObj          *C292Image `json:"image_obj,omitempty"`
}

// C292 神评卡片
// 用于发现页神评模块，支持评论轮播、左侧动画、折叠控制等。
type C292 struct {
	BaseItem
	ClickActionLog        any                `json:"click_actionlog,omitempty"`
	Height                int                `json:"height,omitempty"`
	CommentArray          []*C292CommentItem `json:"comment_array,omitempty"`
	CommentScrollInterval int                `json:"comment_scroll_interval,omitempty"`
	PaddingVertical       int                `json:"padding_vertical,omitempty"`
	BackgroundPic         string             `json:"background_pic,omitempty"`
	BackgroundPicDark     string             `json:"background_pic_dark,omitempty"`
	ControlObj            *C292ControlObj    `json:"control_obj,omitempty"`
	CommentTopmask        *C292Mask          `json:"comment_topmask,omitempty"`
	CommentBottommask     *C292Mask          `json:"comment_bottommask,omitempty"`
	DisplayAnimType       int                `json:"display_anim_type,omitempty"`
	BorderColor           string             `json:"border_color,omitempty"`
	BorderColorDark       string             `json:"border_color_dark,omitempty"`
	Left                  *C292Left          `json:"left,omitempty"`
	Right                 *C292Right         `json:"right,omitempty"`
}

// SetDisplayAnimType 设置展示动画类型。
func (c *C292) SetDisplayAnimType(value int) {
	c.DisplayAnimType = value
}

// NewCard292 创建Card292实例。
func NewCard292() *C292 {
	return &C292{BaseItem: BaseItem{Base: Base{CardType: 292}}}
}
