package card

// C108 视频专辑卡片
type C108 struct {
	BaseItem
	Playlist     *C108Playlist `json:"playlist,omitempty"`      // 专辑结构（必须下发）
	NeedFeedback int8          `json:"need_feedback,omitempty"` // 右上角是否需要显示负反馈关闭按钮；1表示显示，0表示不显示
	FeedBackExt  string        `json:"feed_back_ext,omitempty"` // 负反馈接口请求参数
}

// C108Playlist 专辑结构
type C108Playlist struct {
	Id           int64         `json:"id,omitempty"`            // 专辑id
	Name         string        `json:"name,omitempty"`          // 专辑名称
	VideoCount   int           `json:"video_count,omitempty"`   // 专辑视频数
	IsSubscribed int8          `json:"is_subscribed,omitempty"` // 是否订阅了当前专辑
	Author       *C108Author   `json:"author,omitempty"`        // 专辑作者信息
	Statuses     []*C108Status `json:"statuses,omitempty"`      // 专辑内的微博数组：目前只支持展示一条微博
	Keywords     []string      `json:"keywords,omitempty"`
	RecomCode    string        `json:"recom_code,omitempty"`
	FeedBackExt  string        `json:"feed_back_ext,omitempty"`
	Index        int           `json:"index,omitempty"`
}

// C108Author 专辑作者信息
type C108Author struct {
	Id              int64  `json:"id,omitempty"`
	ScreenName      string `json:"screen_name,omitempty"`
	Name            string `json:"name,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
	CoverImage      string `json:"cover_image,omitempty"`
	Following       bool   `json:"following,omitempty"`
	Verified        bool   `json:"verified,omitempty"`
	VerifiedType    int8   `json:"verified_type,omitempty"`
	AvatarLarge     string `json:"avatar_large,omitempty"`
	AvatarHd        string `json:"avatar_hd,omitempty"`
	VerifiedReason  string `json:"verified_reason,omitempty"`
	VerifiedLevel   int8   `json:"verified_level,omitempty"`
	VerifiedTypeExt int8   `json:"verified_type_ext,omitempty"`
	FollowMe        bool   `json:"follow_me,omitempty"`
	Level           int8   `json:"level,omitempty"`
	FollowersCount  int64  `json:"followers_count,omitempty"`
}

// C108Status 微博状态
type C108Status struct {
	CreatedAt      string           `json:"created_at,omitempty"`
	Id             int64            `json:"id,omitempty"`
	Idstr          string           `json:"idstr,omitempty"`
	Mid            string           `json:"mid,omitempty"`
	Text           string           `json:"text,omitempty"`
	TextLength     int              `json:"textLength,omitempty"`
	Source         string           `json:"source,omitempty"`
	Favorited      bool             `json:"favorited,omitempty"`
	Truncated      bool             `json:"truncated,omitempty"`
	PicIds         []string         `json:"pic_ids,omitempty"`
	User           *C108StatusUser  `json:"user,omitempty"`
	RepostsCount   int              `json:"reposts_count,omitempty"`
	CommentsCount  int              `json:"comments_count,omitempty"`
	AttitudesCount int              `json:"attitudes_count,omitempty"`
	Visible        *C108Visible     `json:"visible,omitempty"`
	MblogId        string           `json:"mblogid,omitempty"`
	Scheme         string           `json:"scheme,omitempty"`
	PageInfo       *C108PageInfo    `json:"page_info,omitempty"`
	UrlStruct      []*C108UrlStruct `json:"url_struct,omitempty"`
	ObjExt         string           `json:"obj_ext,omitempty"` // 观看次数
	IsLongText     bool             `json:"isLongText,omitempty"`
	MultiAttitude  []*C108Attitude  `json:"multi_attitude,omitempty"`
}

// C108StatusUser 微博用户
type C108StatusUser struct {
	Id              int64  `json:"id,omitempty"`
	Idstr           string `json:"idstr,omitempty"`
	ScreenName      string `json:"screen_name,omitempty"`
	Name            string `json:"name,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
	AvatarLarge     string `json:"avatar_large,omitempty"`
	AvatarHd        string `json:"avatar_hd,omitempty"`
	Verified        bool   `json:"verified,omitempty"`
	VerifiedType    int8   `json:"verified_type,omitempty"`
	VerifiedReason  string `json:"verified_reason,omitempty"`
	Following       bool   `json:"following,omitempty"`
	FollowMe        bool   `json:"follow_me,omitempty"`
	FollowersCount  int64  `json:"followers_count,omitempty"`
	FriendsCount    int64  `json:"friends_count,omitempty"`
	StatusesCount   int64  `json:"statuses_count,omitempty"`
	Description     string `json:"description,omitempty"`
	Gender          string `json:"gender,omitempty"`
	Location        string `json:"location,omitempty"`
}

// C108Visible 可见性
type C108Visible struct {
	Type   int8 `json:"type,omitempty"`
	ListId int  `json:"list_id,omitempty"`
}

// C108PageInfo 页面信息
type C108PageInfo struct {
	Type       string         `json:"type,omitempty"`
	PageId     string         `json:"page_id,omitempty"`
	ObjectType string         `json:"object_type,omitempty"`
	Oid        int64          `json:"oid,omitempty"`
	PageTitle  string         `json:"page_title,omitempty"`
	PagePic    string         `json:"page_pic,omitempty"`
	TypeIcon   string         `json:"type_icon,omitempty"`
	PageUrl    string         `json:"page_url,omitempty"`
	ObjectId   string         `json:"object_id,omitempty"`
	MediaInfo  *C108MediaInfo `json:"media_info,omitempty"`
	AuthorId   int64          `json:"author_id,omitempty"`
	Authorid   int64          `json:"authorid,omitempty"`
	Cards      []*C108Card    `json:"cards,omitempty"`
	ActionLog  map[string]any `json:"actionlog,omitempty"`
}

// C108MediaInfo 媒体信息
type C108MediaInfo struct {
	Name                  string                      `json:"name,omitempty"`
	Duration              int                         `json:"duration,omitempty"`
	StreamUrl             string                      `json:"stream_url,omitempty"`
	StreamUrlHd           string                      `json:"stream_url_hd,omitempty"`
	H5Url                 string                      `json:"h5_url,omitempty"`
	Mp4SdUrl              string                      `json:"mp4_sd_url,omitempty"`
	Mp4HdUrl              string                      `json:"mp4_hd_url,omitempty"`
	PrefetchType          int8                        `json:"prefetch_type,omitempty"`
	PrefetchSize          int                         `json:"prefetch_size,omitempty"`
	ActStatus             int8                        `json:"act_status,omitempty"`
	MediaId               string                      `json:"media_id,omitempty"`
	NextTitle             string                      `json:"next_title,omitempty"`
	VideoDetails          []*C108VideoDetail          `json:"video_details,omitempty"`
	PlayCompletionActions []*C108PlayCompletionAction `json:"play_completion_actions,omitempty"`
	VideoPublishTime      int64                       `json:"video_publish_time,omitempty"`
	PlayLoopType          int8                        `json:"play_loop_type,omitempty"`
	Titles                []*C108Title                `json:"titles,omitempty"`
	AuthorName            string                      `json:"author_name,omitempty"`
	PlaylistId            int64                       `json:"playlist_id,omitempty"`
	IsPlaylist            int8                        `json:"is_playlist,omitempty"`
	GetPlaylistId         int64                       `json:"get_playlist_id,omitempty"`
	ExtraInfo             *C108ExtraInfo              `json:"extra_info,omitempty"`
	HasRecommendVideo     int8                        `json:"has_recommend_video,omitempty"`
	BackPasterInfo        *C108BackPasterInfo         `json:"back_paster_info,omitempty"`
	AuthorVerifiedType    int8                        `json:"author_verified_type,omitempty"`
	VideoUnmute           int8                        `json:"video_unmute,omitempty"`
	OnlineUsers           string                      `json:"online_users,omitempty"`
	OnlineUsersNumber     int                         `json:"online_users_number,omitempty"`
	IsKeepCurrentMblog    int8                        `json:"is_keep_current_mblog,omitempty"`
	VideoTags             []any                       `json:"video_tags,omitempty"`
}

// C108VideoDetail 视频详情
type C108VideoDetail struct {
	Size    int    `json:"size,omitempty"`
	Bitrate int    `json:"bitrate,omitempty"`
	Label   string `json:"label,omitempty"`
}

// C108PlayCompletionAction 播放完成动作
type C108PlayCompletionAction struct {
	Type         string         `json:"type,omitempty"`
	Icon         string         `json:"icon,omitempty"`
	Text         string         `json:"text,omitempty"`
	Link         string         `json:"link,omitempty"`
	Scheme       string         `json:"scheme,omitempty"`
	BtnCode      int            `json:"btn_code,omitempty"`
	ShowPosition int            `json:"show_position,omitempty"`
	ActionLog    map[string]any `json:"actionlog,omitempty"`
}

// C108Title 标题
type C108Title struct {
	Default bool   `json:"default,omitempty"`
	Title   string `json:"title,omitempty"`
}

// C108ExtraInfo 额外信息
type C108ExtraInfo struct {
	Sceneid          string `json:"sceneid,omitempty"`
	PlaylistDesc     string `json:"playlist_desc,omitempty"`
	PlaylistItemDesc string `json:"playlist_item_desc,omitempty"`
}

// C108BackPasterInfo 后贴片信息
type C108BackPasterInfo struct {
	HasBackPaster int            `json:"has_back_paster,omitempty"`
	RequestParam  map[string]any `json:"request_param,omitempty"`
}

// C108Card 卡片
type C108Card struct {
	Type       string         `json:"type,omitempty"`
	PageId     string         `json:"page_id,omitempty"`
	ObjectType string         `json:"object_type,omitempty"`
	ObjectId   string         `json:"object_id,omitempty"`
	Content1   string         `json:"content1,omitempty"`
	Content2   string         `json:"content2,omitempty"`
	ActStatus  int8           `json:"act_status,omitempty"`
	MediaInfo  *C108MediaInfo `json:"media_info,omitempty"`
	PagePic    string         `json:"page_pic,omitempty"`
	PageTitle  string         `json:"page_title,omitempty"`
	PageUrl    string         `json:"page_url,omitempty"`
	PageDesc   string         `json:"page_desc,omitempty"`
	PicInfo    *C108PicInfo   `json:"pic_info,omitempty"`
	Oid        int64          `json:"oid,omitempty"`
	TypeIcon   string         `json:"type_icon,omitempty"`
	AuthorId   int64          `json:"author_id,omitempty"`
	Authorid   int64          `json:"authorid,omitempty"`
	Warn       string         `json:"warn,omitempty"`
	ActionLog  map[string]any `json:"actionlog,omitempty"`
}

// C108PicInfo 图片信息
type C108PicInfo struct {
	PicBig    *C108Pic `json:"pic_big,omitempty"`
	PicSmall  *C108Pic `json:"pic_small,omitempty"`
	PicMiddle *C108Pic `json:"pic_middle,omitempty"`
}

// C108Pic 图片
type C108Pic struct {
	Height string `json:"height,omitempty"`
	Width  string `json:"width,omitempty"`
	Url    string `json:"url,omitempty"`
}

// C108UrlStruct URL结构
type C108UrlStruct struct {
	UrlTitle    string         `json:"url_title,omitempty"`
	UrlTypePic  string         `json:"url_type_pic,omitempty"`
	OriUrl      string         `json:"ori_url,omitempty"`
	PageId      string         `json:"page_id,omitempty"`
	ShortUrl    string         `json:"short_url,omitempty"`
	UrlType     int            `json:"url_type,omitempty"`
	Result      bool           `json:"result,omitempty"`
	ActionLog   map[string]any `json:"actionlog,omitempty"`
	StorageType string         `json:"storage_type,omitempty"`
	Hide        int8           `json:"hide,omitempty"`
	ObjectType  string         `json:"object_type,omitempty"`
	NeedSaveObj int8           `json:"need_save_obj,omitempty"`
	Log         string         `json:"log,omitempty"`
}

// C108Attitude 态度
type C108Attitude struct {
	Type  int8 `json:"type,omitempty"`
	Count int  `json:"count,omitempty"`
}

// NewCard108 创建Card108实例
func NewCard108() *C108 {
	return &C108{BaseItem: BaseItem{Base: Base{CardType: 108}}}
}
