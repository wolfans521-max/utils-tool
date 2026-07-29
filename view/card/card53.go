package card

// C53 红包卡片
type C53 struct {
	BaseItem
	User           *C53User     `json:"user,omitempty"`
	BgUrl          string       `json:"bg_url,omitempty"`           // 背景图片
	TxtTitle       string       `json:"txt_title,omitempty"`        // 标题
	Desc1          string       `json:"desc1,omitempty"`            // 描述文案
	TitleColor     string       `json:"title_color,omitempty"`      // 标题颜色
	TitleSize      int          `json:"title_size,omitempty"`       // 标题字号
	TitleColorSkin string       `json:"title_color_skin,omitempty"` // 有皮肤时的标题颜色
	Desc1Color     string       `json:"desc1_color,omitempty"`
	Desc1Size      int          `json:"desc1_size,omitempty"`
	Desc1ColorSkin string       `json:"desc1_color_skin,omitempty"`
	Price          *C53Price    `json:"price,omitempty"`       // 已领取的价格
	ButtonText     string       `json:"button_text,omitempty"` // 按钮文案
	Buttons        []*C53Button `json:"buttons,omitempty"`     // 按钮点击配置
}

// C53User 用户信息
type C53User struct {
	Id          int64  `json:"id,omitempty"`
	ScreenName  string `json:"screen_name,omitempty"`
	AvatarLarge string `json:"avatar_large,omitempty"`
}

// C53Price 已领取的价格
type C53Price struct {
	TextColor     string          `json:"text_color,omitempty"`
	TextSize      int             `json:"text_size,omitempty"`
	TextColorSkin string          `json:"text_color_skin,omitempty"`
	Text          string          `json:"text,omitempty"`
	Number        *C53PriceNumber `json:"number,omitempty"`
}

// C53PriceNumber 价格数字
type C53PriceNumber struct {
	Num           string `json:"num,omitempty"`
	TextColor     string `json:"text_color,omitempty"`
	TextColorSkin string `json:"text_color_skin,omitempty"`
}

// C53Button 红包按钮
type C53Button struct {
	Type      string           `json:"type,omitempty"`
	Name      string           `json:"name,omitempty"`
	Pic       string           `json:"pic,omitempty"`
	Params    *C53ButtonParams `json:"params,omitempty"`
	ActionLog map[string]any   `json:"actionlog,omitempty"`
}

// C53ButtonParams 按钮参数
type C53ButtonParams struct {
	Scheme       string `json:"scheme,omitempty"`
	TxtBg        string `json:"txt_bg,omitempty"`
	TxtColor     string `json:"txt_color,omitempty"`
	TxtColorSkin string `json:"txt_color_skin,omitempty"`
}

// NewCard53 创建Card53实例
func NewCard53() *C53 {
	return &C53{BaseItem: BaseItem{Base: Base{CardType: 53}}}
}
