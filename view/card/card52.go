package card

// C52 双列商品卡片
// 功能：
// 52_1 仅有描述信息
// 52_2 既有描述信息，也有价格信息
type C52 struct {
	BaseItem
	Type    int8       `json:"type,omitempty"`    // 1表示只有描述信息、2表示有描述信息和价格信息
	Width   int        `json:"width,omitempty"`   // 宽高比-宽度
	Height  int        `json:"height,omitempty"`  // 宽高比-高度
	Title   string     `json:"title,omitempty"`   // 标题
	Items   []*C52Item `json:"items,omitempty"`   // 商品列表
	OpenUrl string     `json:"openurl,omitempty"` // 打开链接
}

// C52Item 商品项
type C52Item struct {
	Pic           string         `json:"pic,omitempty"`            // 图片
	PicBig        string         `json:"pic_big,omitempty"`        // 大图
	PicScheme     string         `json:"pic_scheme,omitempty"`     // 图片跳转链接
	Scheme        string         `json:"scheme,omitempty"`         // 跳转链接
	Title         string         `json:"title,omitempty"`          // 标题
	Desc1         string         `json:"desc1,omitempty"`          // 描述信息
	Price2        string         `json:"price2,omitempty"`         // 价格信息
	Desc2         string         `json:"desc2,omitempty"`          // type=2 有价格时的描述信息
	ObjectId      string         `json:"object_id,omitempty"`      // 对象ID
	ActStatus     int8           `json:"act_status,omitempty"`     // 状态
	ObjectType    string         `json:"object_type,omitempty"`    // 对象类型
	OnlineUsers   int            `json:"online_users,omitempty"`   // 观看次数
	VideoDuration int            `json:"video_duration,omitempty"` // 视频时长，单位秒
	ActionLog     map[string]any `json:"actionlog,omitempty"`      // 打点日志
}

// NewCard52 创建Card52实例
func NewCard52() *C52 {
	return &C52{BaseItem: BaseItem{Base: Base{CardType: 52}}}
}
