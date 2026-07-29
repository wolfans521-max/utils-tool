package card

type C19Group struct {
	TitleSub  string         `json:"title_sub"`
	Pic       string         `json:"pic"`
	Scheme    string         `json:"scheme"`
	ItemId    string         `json:"itemid,omitempty"`
	ActionLog map[string]any `json:"action_log,omitempty"`
	UnfoldPic string         `json:"unfold_pic,omitempty"`
}

type C19 struct {
	BaseItem
	DefaultRows          int            `json:"default_rows"`
	DividerColor         string         `json:"divider_color,omitempty"`
	CardBgColor          string         `json:"card_bg_color,omitempty"`
	CardBgColorDark      string         `json:"card_bg_color_dark,omitempty"`
	IsNewsquareUistyle   int            `json:"is_newsquare_uistyle,omitempty"`
	NewSquareStyle       int            `json:"new_square_style,omitempty"`
	Mode                 int            `json:"mode,omitempty"`
	PosId                string         `json:"posid,omitempty"`
	IsRefactorStyle      int            `json:"isRefactorStyle,omitempty"`
	Col                  int            `json:"col,omitempty"`
	Group                []C19Group     `json:"group,omitempty"`
	CardPadding          map[string]int `json:"card_padding,omitempty"`
	CardBgImage          string         `json:"card_bg_image,omitempty"`
	CardBgImageDark      string         `json:"card_bg_image_dark,omitempty"`
	CardBgImageNinePatch string         `json:"card_bg_image_nine_patch,omitempty"`
	UnfoldCol            int            `json:"unfold_col,omitempty"`
	MorePic              string         `json:"more_pic,omitempty"`
	DisplayAnimType      int            `json:"display_anim_type,omitempty"`
}

func NewCard19() *C19 {
	return &C19{BaseItem: BaseItem{Base: Base{CardType: 19}}}
}

func (c *C19) SetDisplayAnimType(value int) {
	c.DisplayAnimType = value
}
