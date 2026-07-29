// Package types 提供 render_show_batch 响应体的强类型 struct 定义
// 替换 map[string]any 以消除 GC 压力与 cast 反射开销
package types

// Feed render_show_batch 返回的单条博文
type Feed struct {
	// 基础标识字段
	IdStr string `json:"idstr,omitempty"`
	Mid   string `json:"mid,omitempty"`
	Id    string `json:"id,omitempty"`
	Text  string `json:"text,omitempty"`

	// 媒体相关
	PicNum      int    `json:"pic_num,omitempty"`
	PicIds      []any  `json:"pic_ids,omitempty"`
	OriginalPic string `json:"original_pic,omitempty"`
	Scheme      string `json:"scheme,omitempty"`

	// 来源与区域
	Source     string `json:"source,omitempty"`
	RegionName string `json:"region_name,omitempty"`
	Itemid     string `json:"itemid,omitempty"`

	// 插入广告标记（类型是 any，可能是 bool/struct）
	InsertAdFeed    any    `json:"insert_ad_feed,omitempty"`
	InsertHotAdFeed any    `json:"insert_hot_ad_feed,omitempty"`
	ItemCategory    string `json:"item_category,omitempty"`
	CandType        string `json:"cand_type,omitempty"`

	// 热点字段
	HotFeedPositiveFB bool   `json:"hot_feed_positive_feedback,omitempty"`
	HotSpotName       string `json:"hot_spot_name,omitempty"`
	HotSpotMultiName  string `json:"hot_spot_multi_name,omitempty"`
	HotSpotTag        any    `json:"hot_spot_tag,omitempty"`
	RegionSpotName    string `json:"region_spot_tag_name,omitempty"`
	RegionSpotTag     any    `json:"region_spot_tag,omitempty"`
	FeedType          string `json:"feed_type,omitempty"`
	Hotspot           int    `json:"hotspot,omitempty"`
	MblogType         int    `json:"mblogtype,omitempty"`

	// 正反馈标记位（由下游 set）
	PositiveRecomType string `json:"positive_recom_type,omitempty"`

	// --- 嵌套结构（仅用到的字段） ---
	User         *FeedUser         `json:"user,omitempty"`
	PageInfo     *FeedPageInfo     `json:"page_info,omitempty"`
	ExtraInfo    *FeedExtraInfo    `json:"extra_info,omitempty"`
	MixMediaInfo *FeedMixMediaInfo `json:"mix_media_info,omitempty"`

	// --- 以下字段由 WrapFeedData 等下游填充（json 中不存在） ---
	AnalysisExtra           string `json:"-"`
	StreamEntryId           string `json:"-"`
	RecommendSource         string `json:"-"`
	FromCateId              string `json:"-"`
	SourceMid               string `json:"-"`
	EnableCommentGuide      bool   `json:"-"`
	CommentGuidePlaceholder string `json:"-"`
	SearchDebugInfo         string `json:"-"`
	MblogMenuNewStyle       int    `json:"-"`
	Buttons                 any    `json:"-"`
	MblogButtons            any    `json:"-"`
	MblogFeedBackMenus      any    `json:"-"`
	ExtraButtonInfo         any    `json:"-"`
	Actionlog               any    `json:"-"`
	ClickActionlog          any    `json:"-"`
	InsertItems             any    `json:"-"`
}

// FeedUser 用户信息（仅用到的字段）
type FeedUser struct {
	Id              any    `json:"id,omitempty"`
	IdStr           string `json:"idstr,omitempty"`
	ScreenName      string `json:"screen_name,omitempty"`
	Name            string `json:"name,omitempty"`
	Following       bool   `json:"following,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
	AvatarLarge     string `json:"avatar_large,omitempty"`
	Verified        any    `json:"verified,omitempty"`
	VerifiedType    int    `json:"verified_type,omitempty"`
	VerifiedReason  string `json:"verified_reason,omitempty"`
	Description     string `json:"description,omitempty"`
	Level           any    `json:"level,omitempty"`
	CardidHash      string `json:"cardid_secret,omitempty"`
}

// FeedPageInfo 页面信息
type FeedPageInfo struct {
	ObjectType       string         `json:"object_type,omitempty"`
	OriginObjectType string         `json:"origin_object_type,omitempty"`
	ObjectId         string         `json:"object_id,omitempty"`
	PagePic          any            `json:"page_pic,omitempty"`
	MediaInfo        *FeedMediaInfo `json:"media_info,omitempty"`
}

// FeedMediaInfo 媒体信息
type FeedMediaInfo struct {
	MediaId string `json:"media_id,omitempty"`
}

// FeedExtraInfo 额外信息
type FeedExtraInfo struct {
	RecommendSource     string                 `json:"recommend_source,omitempty"`
	GroupId             string                 `json:"group_id,omitempty"`
	Category            string                 `json:"category,omitempty"`
	RecommendReason     string                 `json:"recommend_reason,omitempty"`
	RecommendDictionary *FeedDictionary        `json:"recommend_dictionary,omitempty"`
	Dictionary          *FeedDictionary        `json:"dictionary,omitempty"`
	NegativeTagsDetail  []*FeedNegativeTag     `json:"negative_tags_detail,omitempty"`
	BusinessSupport     int                    `json:"business_support,omitempty"`
	AttitudeDynamicAd   *FeedAttitudeDynamicAd `json:"attitude_dynamic_ad,omitempty"`
	BgCard              *FeedBgCard            `json:"bg_card,omitempty"`
}

// FeedDictionary 字典结构
type FeedDictionary struct {
	Pairs []*FeedDictPair `json:"pairs,omitempty"`
}

// FeedDictPair 字典键值对
type FeedDictPair struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// FeedNegativeTag 负反馈标签
type FeedNegativeTag struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	Tag  string `json:"tag,omitempty"`
}

// FeedAttitudeDynamicAd 态度动态广告
type FeedAttitudeDynamicAd struct {
	AdType string `json:"adType,omitempty"`
	Adid   string `json:"adid,omitempty"`
}

// FeedBgCard 背景卡片
type FeedBgCard struct {
	Adid     string `json:"adid,omitempty"`
	Keywords any    `json:"keywords,omitempty"`
}

// FeedMixMediaInfo 混合媒体信息
type FeedMixMediaInfo struct {
	Items []*FeedMixMediaItem `json:"items,omitempty"`
}

// FeedMixMediaItem 混合媒体项
type FeedMixMediaItem struct {
	Type string            `json:"type,omitempty"`
	Data *FeedMixMediaData `json:"data,omitempty"`
}

// FeedMixMediaData 混合媒体数据
type FeedMixMediaData struct {
	ObjectType string `json:"object_type,omitempty"`
}
