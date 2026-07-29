package card

// C192 视频集合卡片
// 用于视频集合展示
type C192 struct {
	BaseItem
	TopContent   *C192TopContent `json:"top_content,omitempty"`   // 顶部内容
	TopDesc      string          `json:"top_desc,omitempty"`      // 卡片顶部标题文案
	Type         int             `json:"type,omitempty"`          // 卡片类型
	ReadtimeType string          `json:"readtimetype,omitempty"`  // 阅读时间类型
	HighPriority string          `json:"high_priority,omitempty"` // 高优先级标识
	SectionId    string          `json:"section_id,omitempty"`    // 分区id
	Videos       []*C192Video    `json:"videos,omitempty"`        // 视频集合
	OpenUrl      string          `json:"openurl,omitempty"`       // 视频滑动最后跳转scheme
}

// C192TopContent 顶部内容
type C192TopContent struct {
	ImgUrl            string `json:"img_url,omitempty"`              // 第一的网络图片
	Title             string `json:"title,omitempty"`                // 视频名称
	VideoCountString  string `json:"video_count_string,omitempty"`   // 视频个数
	PlayCountString   string `json:"play_count_string,omitempty"`    // 播放热度文案
	PlayCountIcon     string `json:"play_count_icon,omitempty"`      // 播放热度icon
	PlayCountIconDark string `json:"play_count_icon_dark,omitempty"` // iOS暗黑模式播放热度icon
	FlagUrl           string `json:"flag_url,omitempty"`             // 标识url
	Scheme            string `json:"scheme,omitempty"`               // 点击跳转scheme
}

// C192Video 视频
type C192Video struct {
	Status           any            `json:"status,omitempty"`             // 视频结构
	ExposureLog      map[string]any `json:"exposure_log,omitempty"`       // 曝光日志
	Mid              string         `json:"mid,omitempty"`                // 微博mid
	Protocol         string         `json:"protocol,omitempty"`           // 播放协议
	MediaId          string         `json:"media_id,omitempty"`           // 播放id
	MediaType        string         `json:"media_type,omitempty"`         // 类型 video 视频 pic 图片
	VideoTitle       string         `json:"video_title,omitempty"`        // 标题
	AutoPlay         bool           `json:"auto_play,omitempty"`          // 是否自动播放
	CenterCrop       bool           `json:"center_crop,omitempty"`        // 是否裁切
	Tag              string         `json:"tag,omitempty"`                // 标签
	Oid              string         `json:"oid,omitempty"`                // 播放oid
	ItemRadius       int            `json:"item_radius,omitempty"`        // 圆角
	Scheme           string         `json:"scheme,omitempty"`             // 跳转scheme
	Cover            string         `json:"cover,omitempty"`              // 封面图
	BigCard          bool           `json:"big_card,omitempty"`           // 是否是大卡
	SubInfoDuration  string         `json:"sub_info_duration,omitempty"`  // 起播以后副标题隐藏时间
	MainInfoDuration string         `json:"main_info_duration,omitempty"` // 起播以后主标题隐藏时间
	ActionLog        map[string]any `json:"actionlog,omitempty"`          // 点击日志结构
	AnalysisExtra    string         `json:"analysis_extra,omitempty"`     // 统计分析extra
	VideoOrientation string         `json:"video_orientation,omitempty"`  // 视频方向 horizontal/vertical
	Index            int            `json:"index,omitempty"`              // 索引
	HighPriority     bool           `json:"high_priority,omitempty"`      // 是否高优先级
	Type             string         `json:"type,omitempty"`               // 子项类型，如 card_fourth
	User             map[string]any `json:"user,omitempty"`               // 用户信息
	MblogMenusNew    []any          `json:"mblog_menus_new,omitempty"`    // 负反馈菜单
}

// NewCard192 创建Card192实例
func NewCard192() *C192 {
	return &C192{BaseItem: BaseItem{Base: Base{CardType: 192}}}
}
