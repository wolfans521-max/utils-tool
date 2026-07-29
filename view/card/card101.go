package card

type C101 struct {
	BaseItem
	Title              string `json:"title"`
	SubTitle           string `json:"sub_title,omitempty"`
	Desc               string `json:"desc"`
	FontSize           int8   `json:"font_size"`
	TitleIsBold        int8   `json:"title_is_bold,omitempty"`
	DisplayArrow       int8   `json:"display_arrow,omitempty"`
	NoHighlight        int8   `json:"no_highlight,omitempty"`
	NewTitleIsBold     int8   `json:"new_title_is_bold,omitempty"`
	TitleColor         string `json:"title_color"`
	TitleColorDark     string `json:"title_color_dark,omitempty"`
	DescColor          string `json:"desc_color,omitempty"`
	DescColorDark      string `json:"desc_color_dark,omitempty"`
	TitleDark          string `json:"title_dark,omitempty"`
	LeftTagImg         string `json:"left_tag_img"`
	LeftTagImgHeight   int    `json:"left_tag_img_height,omitempty"`
	LeftTagImgWidth    int    `json:"left_tag_img_width,omitempty"`
	IsShowArrow        int8   `json:"is_show_arrow"`
	BottomLine         int8   `json:"bottom_line,omitempty"`
	TopTagImgPadding   int8   `json:"top_tag_img_padding"`
	TagImg             string `json:"tag_img,omitempty"`
	TagImgDark         string `json:"tag_img_dark,omitempty"`
	TagImgWidth        int    `json:"tag_img_width,omitempty"`
	TagImgHeight       int    `json:"tag_img_height,omitempty"`
	DescIcon           string `json:"desc_icon,omitempty"`
	DescIconDark       string `json:"desc_icon_dark,omitempty"`
	DescIconWidth      int    `json:"desc_icon_width,omitempty"`
	DescIconHeight     int    `json:"desc_icon_height,omitempty"`
	TopPadding         int8   `json:"top_padding,omitempty"`
	BottomPadding      int8   `json:"bottom_padding,omitempty"`
	UseJsonColor       bool   `json:"use_json_color"`
	SubTitleColor      string `json:"sub_title_color,omitempty"`
	SubTitleColorDark  string `json:"sub_title_color_dark,omitempty"`
	Height             int8   `json:"height,omitempty"`
	ContentAlignBottom int8   `json:"content_align_bottom,omitempty"`
	PicTagStyle        int8   `json:"pic_tag_style"`
	LeftTagImgPadding  int8   `json:"left_tag_img_padding"`
	IsCenterVertical   int    `json:"is_center_vertical,omitempty"`
	LeftTagImgDark     string `json:"left_tag_img_dark"`
	Scheme             string `json:"scheme,omitempty"`
	RightTagImgPadding int    `json:"right_tag_img_padding,omitempty"`
	RelativeMid        string `json:"relative_mid,omitempty"`
}

func NewCard101() *C101 {
	return &C101{BaseItem: BaseItem{Base: Base{CardType: 101}}}
}
