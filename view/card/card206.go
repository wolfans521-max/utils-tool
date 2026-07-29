package card

// C206 一卡多帧视频卡片
// 展示多个视频，自动轮流播放
type C206 struct {
	BaseItem
	AspectRatio   float64     `json:"aspect_ratio,omitempty"`   // 长宽比例，默认16:9（1.778）
	LeftPadding   int         `json:"left_padding,omitempty"`   // 左边距
	RightPadding  int         `json:"right_padding,omitempty"`  // 右边距
	TopPadding    int         `json:"top_padding,omitempty"`    // 上边距
	BottomPadding int         `json:"bottom_padding,omitempty"` // 下边距
	Items         []*C206Item `json:"items,omitempty"`          // 视频项列表
	OpenUrl       string      `json:"openurl,omitempty"`        // 打开链接
}

// C206Item 视频项
type C206Item struct {
	MediaInfo *C206MediaInfo `json:"media_info,omitempty"` // 媒体信息
	ObjectId  string         `json:"object_id,omitempty"`  // 对象ID
	PagePic   string         `json:"page_pic,omitempty"`   // 视频封面
	PageUrl   string         `json:"page_url,omitempty"`   // 页面URL
	ActionLog map[string]any `json:"actionlog,omitempty"`  // 点击日志
}

// C206MediaInfo 媒体信息
type C206MediaInfo struct {
	VideoDownloadStrategy any                     `json:"video_download_strategy,omitempty"` // 视频下载策略
	HasRecommendVideo     int                     `json:"has_recommend_video,omitempty"`     // 是否有推荐视频
	PlayLoopType          int                     `json:"play_loop_type,omitempty"`          // 播放循环类型
	SearchScheme          string                  `json:"search_scheme,omitempty"`           // 搜索scheme
	VideoSchemeH5         bool                    `json:"video_scheme_h5,omitempty"`         // 视频scheme H5
	ExtraInfo             any                     `json:"extra_info,omitempty"`              // 额外信息
	IsShortVideo          int                     `json:"is_short_video,omitempty"`          // 是否短视频
	ForwardStrategy       int                     `json:"forward_strategy,omitempty"`        // 转发策略
	Titles                []*C206Title            `json:"titles,omitempty"`                  // 标题列表
	NextTitle             string                  `json:"next_title,omitempty"`              // 下一个标题
	TitlesDisplayTime     string                  `json:"titles_display_time,omitempty"`     // 标题显示时间
	H5Url                 string                  `json:"h5_url,omitempty"`                  // H5 URL
	VideoDetails          []*C206VideoDetail      `json:"video_details,omitempty"`           // 视频详情
	Mp4SdUrl              string                  `json:"mp4_sd_url,omitempty"`              // MP4标清URL
	PrefetchSize          int                     `json:"prefetch_size,omitempty"`           // 预取大小
	HevcMp4Hd             string                  `json:"hevc_mp4_hd,omitempty"`             // HEVC MP4高清
	Name                  string                  `json:"name,omitempty"`                    // 名称
	Mp4HdUrl              string                  `json:"mp4_hd_url,omitempty"`              // MP4高清URL
	ExtInfo               any                     `json:"ext_info,omitempty"`                // 扩展信息
	VideoOrientation      string                  `json:"video_orientation,omitempty"`       // 视频方向
	StreamUrl             string                  `json:"stream_url,omitempty"`              // 流URL
	AuthorName            string                  `json:"author_name,omitempty"`             // 作者名称
	PlayCompletionActions []*C206CompletionAction `json:"play_completion_actions,omitempty"` // 播放完成动作
	ActStatus             int                     `json:"act_status,omitempty"`              // 动作状态
	HevcMp4Ld             string                  `json:"hevc_mp4_ld,omitempty"`             // HEVC MP4低清
	StreamUrlHd           string                  `json:"stream_url_hd,omitempty"`           // 流URL高清
	KolTitle              string                  `json:"kol_title,omitempty"`               // KOL标题
	OpenScheme            string                  `json:"open_scheme,omitempty"`             // 打开scheme
	MediaId               string                  `json:"media_id,omitempty"`                // 媒体ID
	Protocol              string                  `json:"protocol,omitempty"`                // 协议
	Autoplay              int                     `json:"autoplay,omitempty"`                // 自动播放
	VoteIsShow            int                     `json:"vote_is_show,omitempty"`            // 投票是否显示
	VideoPublishTime      int64                   `json:"video_publish_time,omitempty"`      // 视频发布时间
	PrefetchType          int                     `json:"prefetch_type,omitempty"`           // 预取类型
	Duration              int                     `json:"duration,omitempty"`                // 时长
	AdType                int                     `json:"ad_type,omitempty"`                 // 广告类型
	BelongCollection      int                     `json:"belong_collection,omitempty"`       // 所属集合
	IsKeepCurrentMblog    int                     `json:"is_keep_current_mblog,omitempty"`   // 是否保持当前微博
}

// C206Title 标题
type C206Title struct {
	Title   string `json:"title,omitempty"`   // 标题
	Default bool   `json:"default,omitempty"` // 是否默认
}

// C206VideoDetail 视频详情
type C206VideoDetail struct {
	PrefetchSize int    `json:"prefetch_size,omitempty"` // 预取大小
	Bitrate      int    `json:"bitrate,omitempty"`       // 比特率
	Label        string `json:"label,omitempty"`         // 标签
	Size         int    `json:"size,omitempty"`          // 大小
}

// C206CompletionAction 播放完成动作
type C206CompletionAction struct {
	Scheme                   string         `json:"scheme,omitempty"`                      // scheme
	Type                     string         `json:"type,omitempty"`                        // 类型
	ActionLog                map[string]any `json:"actionlog,omitempty"`                   // 点击日志
	BtnCode                  int            `json:"btn_code,omitempty"`                    // 按钮代码
	Text                     string         `json:"text,omitempty"`                        // 文本
	Link                     string         `json:"link,omitempty"`                        // 链接
	Icon                     string         `json:"icon,omitempty"`                        // 图标
	ThirdPartySchemeProtocol string         `json:"third_party_scheme_protocol,omitempty"` // 第三方scheme协议
	ShowPosition             int            `json:"show_position,omitempty"`               // 显示位置
	Cleaned                  bool           `json:"cleaned,omitempty"`                     // 是否清理
}

// NewCard206 创建Card206实例
func NewCard206() *C206 {
	return &C206{BaseItem: BaseItem{Base: Base{CardType: 206}}}
}
