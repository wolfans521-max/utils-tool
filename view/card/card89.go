package card

// C89 视频微博卡片
type C89 struct {
	BaseItem
	Mblog         *C89Mblog     `json:"mblog,omitempty"`
	HideBottom    int8          `json:"hide_bottom,omitempty"` // 控制非视频区域显示逻辑。1:只显示视频区域 0:全显示
	HeaderTag     *C89HeaderTag `json:"header_tag,omitempty"`
	HotDiscussion string        `json:"hot_discussion,omitempty"` // 视频区域视频标题前方图标icon地址
	NeedFeedback  int8          `json:"need_feedback,omitempty"`
}

// C89HeaderTag 头部标签
type C89HeaderTag struct {
	Icon  string `json:"icon,omitempty"`
	Title string `json:"title,omitempty"`
}

// C89Mblog 微博内容
type C89Mblog struct {
	User     *C89User     `json:"user,omitempty"`
	PageInfo *C89PageInfo `json:"page_info,omitempty"`
}

// C89User 用户信息
type C89User struct {
	Id             int64  `json:"id,omitempty"`
	Idstr          string `json:"idstr,omitempty"`
	ScreenName     string `json:"screen_name,omitempty"`
	AvatarLarge    string `json:"avatar_large,omitempty"`
	VerifiedReason string `json:"verified_reason,omitempty"`
}

// C89PageInfo 页面信息
type C89PageInfo struct {
	Type       string         `json:"type,omitempty"`
	PageId     string         `json:"page_id,omitempty"`
	ObjectType string         `json:"object_type,omitempty"`
	ObjectId   string         `json:"object_id,omitempty"`
	Content1   string         `json:"content1,omitempty"`
	Content2   string         `json:"content2,omitempty"`
	ActStatus  int8           `json:"act_status,omitempty"`
	MediaInfo  *C89MediaInfo  `json:"media_info,omitempty"`
	PagePic    string         `json:"page_pic,omitempty"`
	PageTitle  string         `json:"page_title,omitempty"`
	PageUrl    string         `json:"page_url,omitempty"`
	PicInfo    *C89PicInfo    `json:"pic_info,omitempty"`
	Oid        string         `json:"oid,omitempty"`
	TypeIcon   string         `json:"type_icon,omitempty"`
	AuthorId   string         `json:"author_id,omitempty"`
	Authorid   string         `json:"authorid,omitempty"`
	Warn       string         `json:"warn,omitempty"`
	ActionLog  map[string]any `json:"actionlog,omitempty"`
}

// C89MediaInfo 媒体信息
type C89MediaInfo struct {
	Name                  string                     `json:"name,omitempty"`
	Autoplay              int8                       `json:"autoplay,omitempty"` // 0 默认(按用户设置) 1 （wifi下）强制自动播放 2 强制不自动播放
	Duration              int                        `json:"duration,omitempty"`
	StreamUrl             string                     `json:"stream_url,omitempty"`
	StreamUrlHd           string                     `json:"stream_url_hd,omitempty"`
	H5Url                 string                     `json:"h5_url,omitempty"`
	Mp4SdUrl              string                     `json:"mp4_sd_url,omitempty"`
	Mp4HdUrl              string                     `json:"mp4_hd_url,omitempty"`
	H265Mp4Hd             string                     `json:"h265_mp4_hd,omitempty"`
	H265Mp4Ld             string                     `json:"h265_mp4_ld,omitempty"`
	Inch4Mp4Hd            string                     `json:"inch_4_mp4_hd,omitempty"`
	Inch5Mp4Hd            string                     `json:"inch_5_mp4_hd,omitempty"`
	Inch55Mp4Hd           string                     `json:"inch_5_5_mp4_hd,omitempty"`
	Mp4720pMp4            string                     `json:"mp4_720p_mp4,omitempty"`
	HevcMp4720p           string                     `json:"hevc_mp4_720p,omitempty"`
	HevcMp4Hd             string                     `json:"hevc_mp4_hd,omitempty"`
	PrefetchType          int8                       `json:"prefetch_type,omitempty"`
	PrefetchSize          int                        `json:"prefetch_size,omitempty"`
	ActStatus             int8                       `json:"act_status,omitempty"`
	NextTitle             string                     `json:"next_title,omitempty"`
	PlayCompletionActions []*C89PlayCompletionAction `json:"play_completion_actions,omitempty"`
	VideoPublishTime      int64                      `json:"video_publish_time,omitempty"`
	PlayLoopType          int8                       `json:"play_loop_type,omitempty"`
	ScrubberPreview       *C89ScrubberPreview        `json:"scrubber_preview,omitempty"`
	Titles                []*C89Title                `json:"titles,omitempty"`
	VideoTags             []*C89VideoTag             `json:"video_tags,omitempty"`
	AuthorMid             string                     `json:"author_mid,omitempty"`
	AuthorName            string                     `json:"author_name,omitempty"`
	HasRecommendVideo     int8                       `json:"has_recommend_video,omitempty"`
	VideoFeedShowCustomBg int8                       `json:"video_feed_show_custom_bg,omitempty"`
	OnlineUsers           string                     `json:"online_users,omitempty"`
	OnlineUsersNumber     int                        `json:"online_users_number,omitempty"`
	Ttl                   int                        `json:"ttl,omitempty"`
	StorageType           string                     `json:"storage_type,omitempty"`
	IsKeepCurrentMblog    int8                       `json:"is_keep_current_mblog,omitempty"`
}

// C89PlayCompletionAction 播放完成动作
type C89PlayCompletionAction struct {
	Type         string         `json:"type,omitempty"`
	Icon         string         `json:"icon,omitempty"`
	Text         string         `json:"text,omitempty"`
	Link         string         `json:"link,omitempty"`
	Scheme       string         `json:"scheme,omitempty"`
	BtnCode      int            `json:"btn_code,omitempty"`
	ShowPosition int            `json:"show_position,omitempty"`
	ActionLog    map[string]any `json:"actionlog,omitempty"`
}

// C89ScrubberPreview 进度条预览
type C89ScrubberPreview struct {
	Previews []*C89Preview `json:"previews,omitempty"`
}

// C89Preview 预览图
type C89Preview struct {
	Width    int      `json:"width,omitempty"`
	Height   int      `json:"height,omitempty"`
	Row      int      `json:"row,omitempty"`
	Col      int      `json:"col,omitempty"`
	Interval int      `json:"interval,omitempty"`
	Label    string   `json:"label,omitempty"`
	Urls     []string `json:"urls,omitempty"`
}

// C89Title 标题
type C89Title struct {
	Title   string `json:"title,omitempty"`
	Default bool   `json:"default,omitempty"`
}

// C89VideoTag 视频标签
type C89VideoTag struct {
	Name      string         `json:"name,omitempty"`
	Scheme    string         `json:"scheme,omitempty"`
	Link      string         `json:"link,omitempty"`
	ActionLog map[string]any `json:"actionlog,omitempty"`
}

// C89PicInfo 图片信息
type C89PicInfo struct {
	PicBig    *C89Pic `json:"pic_big,omitempty"`
	PicSmall  *C89Pic `json:"pic_small,omitempty"`
	PicMiddle *C89Pic `json:"pic_middle,omitempty"`
}

// C89Pic 图片
type C89Pic struct {
	Height string `json:"height,omitempty"`
	Width  string `json:"width,omitempty"`
	Url    string `json:"url,omitempty"`
}

// NewCard89 创建Card89实例
func NewCard89() *C89 {
	return &C89{BaseItem: BaseItem{Base: Base{CardType: 89}}}
}
