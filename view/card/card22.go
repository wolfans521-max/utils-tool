package card

type C22PicItems struct {
	PicDark            string         `json:"pic_dark,omitempty"`
	Pic                string         `json:"pic,omitempty"`
	Scheme             string         `json:"scheme,omitempty"`
	PicUnenableClick   bool           `json:"pic_unenable_click,omitempty"`
	UserImageView      string         `json:"user_imageview,omitempty"`
	ActionLog          map[string]any `json:"actionlog,omitempty"`
	NeedClearBackGround int8          `json:"need_clear_backGround,omitempty"`
}

type C22 struct {
	BaseItem
	Pic                 string         `json:"pic,omitempty"`
	PicDark             string         `json:"pic_dark,omitempty"`
	Width               int32          `json:"width,omitempty"`
	Height              int32          `json:"height,omitempty"`
	FlowGap             int            `json:"flow_gap,omitempty"`
	AutoFlow            int8           `json:"auto_flow,omitempty"`
	BottomPadding       int            `json:"bottom_padding,omitempty"`
	TopPadding          int            `json:"top_padding,omitempty"`
	LeftRightPadding    int            `json:"left_right_padding,omitempty"`
	PicItems            []*C22PicItems `json:"pic_items,omitempty"`
	PicHWScale          float64        `json:"pic_h_w_scale,omitempty"`
	PicContentMode      int8           `json:"pic_content_mode,omitempty"`
	IsShowCornerRadius  int8           `json:"is_show_corner_radius,omitempty"`
	CardAdStyle         int8           `json:"card_ad_style,omitempty"`
	PicBgColorType      int8           `json:"pic_bgcolor_type,omitempty"`
	CornerRadius        int8           `json:"corner_radius,omitempty"`
	HiddenSeperator     int8           `json:"hidden_seperator,omitempty"`
	PicUnenableClick    bool           `json:"pic_unenable_click"`
	PicBig              string         `json:"pic_big,omitempty"`
	PicLarge            string         `json:"pic_large,omitempty"`
	FoldWidth           string         `json:"foldWidth,omitempty"`
	CommonButton        *Button        `json:"common_button,omitempty"`
	ForceNeedSeperator  int            `json:"force_need_seperator,omitempty"`
	Cleaned             bool           `json:"cleaned,omitempty"`
	Scheme              string         `json:"scheme,omitempty"`
	ActionLog           map[string]any `json:"actionlog,omitempty"`
	NeedClearBackGround int8           `json:"need_clear_backGround,omitempty"`
	CardBackgroundColor string         `json:"card_background_color,omitempty"`
	DisplayAnimType     int            `json:"display_anim_type,omitempty"`
}

func NewCard22() *C22 {
	return &C22{BaseItem: BaseItem{Base: Base{CardType: 22}}}
}

func (c *C22) SetDisplayAnimType(value int) {
	c.DisplayAnimType = value
}
