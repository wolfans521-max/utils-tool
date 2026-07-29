package card

// C127 Profile视频签视频卡片
// 用于Profile页视频展示
type C127 struct {
	BaseItem
	Visible     *C127Visible   `json:"visible,omitempty"`       // 可见性
	HideFlag    int8           `json:"hide_flag,omitempty"`     // 隐藏标记
	TextMaxLine int            `json:"text_max_line,omitempty"` // 文本最大行数
	Mid         string         `json:"mid,omitempty"`           // 微博MID
	Text        string         `json:"text,omitempty"`          // 正文文案
	ItemType    int8           `json:"item_type,omitempty"`     // 项目类型
	VideoInfo   *C127VideoInfo `json:"video_info,omitempty"`    // 视频信息
	User        any            `json:"user,omitempty"`          // 用户信息
	Icon        string         `json:"icon,omitempty"`          // 图标
	OpenUrl     string         `json:"openurl,omitempty"`       // 打开链接
}

// C127Visible 可见性
type C127Visible struct {
	Type int8 `json:"type,omitempty"` // 类型
}

// C127VideoInfo 视频信息
type C127VideoInfo struct {
	MediaId     int64             `json:"media_id,omitempty"`     // 媒体ID
	MediaIdStr  string            `json:"media_idStr,omitempty"`  // 媒体ID字符串
	Cover       string            `json:"cover,omitempty"`        // 封面
	CoverWidth  int               `json:"cover_width,omitempty"`  // 封面宽度
	CoverHeight int               `json:"cover_height,omitempty"` // 封面高度
	CoverInfo   *C127CoverInfo    `json:"cover_info,omitempty"`   // 封面信息
	Duration    string            `json:"duration,omitempty"`     // 时长
	IsLiked     int8              `json:"is_liked,omitempty"`     // 是否点赞
	LikeCount   int               `json:"like_count,omitempty"`   // 点赞数
	PlayCount   int               `json:"play_count,omitempty"`   // 播放次数
	Orientation string            `json:"orientation,omitempty"`  // 方向
	ShowTitle   string            `json:"show_title,omitempty"`   // 显示标题
	Urls        map[string]string `json:"urls,omitempty"`         // 播放流地址
	Protocol    string            `json:"protocol,omitempty"`     // 协议
	VideoType   int8              `json:"video_type,omitempty"`   // 视频类型
}

// C127CoverInfo 封面信息
type C127CoverInfo struct {
	Url         string `json:"url,omitempty"`           // URL
	Pid         string `json:"pid,omitempty"`           // 图片ID
	Width       int    `json:"width,omitempty"`         // 宽度
	Height      int    `json:"height,omitempty"`        // 高度
	HasFace     bool   `json:"has_face,omitempty"`      // 是否有人脸
	IsSelfCover int8   `json:"is_self_cover,omitempty"` // 是否自定义封面
	CoverType   string `json:"cover_type,omitempty"`    // 封面类型
}

// NewCard127 创建Card127实例
func NewCard127() *C127 {
	return &C127{BaseItem: BaseItem{Base: Base{CardType: 127}}}
}
