package card

// C181 奖牌榜卡片
// 用于奖牌榜展示（如奥运会奖牌榜）
type C181 struct {
	BaseItem
	LeftImg   *C181LeftImg     `json:"left_img,omitempty"`   // 左侧图片
	RightData []*C181RightData `json:"right_data,omitempty"` // 右侧数据
	OpenUrl   string           `json:"openurl,omitempty"`    // 打开链接
}

// C181LeftImg 左侧图片
type C181LeftImg struct {
	CoverImage string  `json:"cover_image,omitempty"` // 封面图片
	RatioWH    float64 `json:"ratio_w_h,omitempty"`   // 宽高比
	Scheme     string  `json:"scheme,omitempty"`      // 跳转scheme
}

// C181RightData 右侧数据
type C181RightData struct {
	NumberPic    string `json:"number_pic,omitempty"`    // 数字图片
	Title        string `json:"title,omitempty"`         // 标题（国家名称）
	FirstNumber  string `json:"first_number,omitempty"`  // 金牌数
	FirstPic     string `json:"first_pic,omitempty"`     // 金牌图片
	SecondNumber string `json:"second_number,omitempty"` // 银牌数
	SecondPic    string `json:"second_pic,omitempty"`    // 银牌图片
	ThirdNumber  string `json:"third_number,omitempty"`  // 铜牌数
	ThirdPic     string `json:"third_pic,omitempty"`     // 铜牌图片
	Scheme       string `json:"scheme,omitempty"`        // 跳转scheme
}

// NewCard181 创建Card181实例
func NewCard181() *C181 {
	return &C181{BaseItem: BaseItem{Base: Base{CardType: 181}}}
}
