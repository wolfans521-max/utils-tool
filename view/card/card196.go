package card

// C196 大型活动卡片（如冬奥会、世界杯）
// 分为banner、宫格、奖牌榜三部分，可以添加头部图片
type C196 struct {
	BaseItem
	PaddingTop        int              `json:"padding_top,omitempty"`         // card 内部上间距
	PaddingBottom     int              `json:"padding_bottom,omitempty"`      // card 内部下间距
	BackgroundPic     string           `json:"background_pic,omitempty"`      // 背景图
	TitlePic          string           `json:"title_pic,omitempty"`           // 左上标题图
	TitlePicWidth     int              `json:"title_pic_width,omitempty"`     // 左上标题图宽度
	LogoPic           string           `json:"logo_pic,omitempty"`            // 右上logo图
	BannerBorderColor string           `json:"banner_border_color,omitempty"` // banner描边颜色
	Banner            []*C196Banner    `json:"banner,omitempty"`              // 轮播
	GridColumn        int              `json:"grid_column,omitempty"`         // 宫格列数
	Grid              []*C196Grid      `json:"grid,omitempty"`                // 宫格部分
	Rank              *C196Rank        `json:"rank,omitempty"`                // 奖牌榜
	FlutterWbox       *C196FlutterWbox `json:"flutterWbox,omitempty"`         // 小程序card数据结构
	OpenUrl           string           `json:"openurl,omitempty"`             // 打开链接
}

// C196Banner 轮播
type C196Banner struct {
	Pic       string         `json:"pic,omitempty"`       // 图片
	Scheme    string         `json:"scheme,omitempty"`    // 跳转scheme
	ActionLog map[string]any `json:"actionlog,omitempty"` // 点击日志
	Cleaned   bool           `json:"cleaned,omitempty"`   // 是否清理
}

// C196Grid 宫格
type C196Grid struct {
	Pic       string         `json:"pic,omitempty"`       // 图片
	Text      string         `json:"text,omitempty"`      // 文字
	Scheme    string         `json:"scheme,omitempty"`    // 跳转scheme
	ActionLog map[string]any `json:"actionlog,omitempty"` // 点击日志
	Cleaned   bool           `json:"cleaned,omitempty"`   // 是否清理
}

// C196Rank 奖牌榜
type C196Rank struct {
	BackgroundPic string          `json:"background_pic,omitempty"` // 背景图
	LeftImg       string          `json:"left_img,omitempty"`       // 左侧图片
	Scheme        string          `json:"scheme,omitempty"`         // 跳转scheme
	ActionLog     map[string]any  `json:"actionlog,omitempty"`      // 点击日志
	Title         *C196RankTitle  `json:"title,omitempty"`          // 标题
	RightData     []*C196RankData `json:"right_data,omitempty"`     // 右侧数据
	Cleaned       bool            `json:"cleaned,omitempty"`        // 是否清理
}

// C196RankTitle 奖牌榜标题
type C196RankTitle struct {
	LeftPic   string         `json:"left_pic,omitempty"`   // 左侧图片
	RightText string         `json:"right_text,omitempty"` // 右侧文字
	Scheme    string         `json:"scheme,omitempty"`     // 跳转scheme
	ArrowPic  string         `json:"arrow_pic,omitempty"`  // 箭头图片
	ActionLog map[string]any `json:"actionlog,omitempty"`  // 点击日志
	Cleaned   bool           `json:"cleaned,omitempty"`    // 是否清理
}

// C196RankData 奖牌榜数据
type C196RankData struct {
	NumberPic    string `json:"number_pic,omitempty"`    // 排名图片
	Title        string `json:"title,omitempty"`         // 国家名称
	FirstNumber  string `json:"first_number,omitempty"`  // 金牌数
	FirstPic     string `json:"first_pic,omitempty"`     // 金牌图片
	SecondNumber string `json:"second_number,omitempty"` // 银牌数
	SecondPic    string `json:"second_pic,omitempty"`    // 银牌图片
	ThirdNumber  string `json:"third_number,omitempty"`  // 铜牌数
	ThirdPic     string `json:"third_pic,omitempty"`     // 铜牌图片
}

// C196FlutterWbox 小程序数据
type C196FlutterWbox struct {
	WboxCacheId   string `json:"wboxCacheId,omitempty"`   // 缓存ID
	WboxScheme    string `json:"wboxScheme,omitempty"`    // 小程序scheme
	WboxParams    any    `json:"wboxParams,omitempty"`    // 小程序参数
	Height        any    `json:"height,omitempty"`        // 高度
	BackgroundUrl string `json:"backgroundUrl,omitempty"` // 背景URL
}

// NewCard196 创建Card196实例
func NewCard196() *C196 {
	return &C196{BaseItem: BaseItem{Base: Base{CardType: 196}}}
}
