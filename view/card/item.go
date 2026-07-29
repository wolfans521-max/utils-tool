package card

type Category string

const (
	CategoryCell    Category = "cell"
	CategoryGroup   Category = "group"
	CategoryCard    Category = "card"
	CategoryFeed    Category = "feed"
	CategoryDynamic Category = "dynamic"
)

type ItemExt struct {
	AnchorId     string `json:"anchorId,omitempty"`    // 详情页评论定位用ID
	ExtraItemId  string `json:"extraItemId,omitempty"` // 记录日志用的item_id
	BizType      string `json:"bizType,omitempty"`
	FilterType   string `json:"filterType,omitempty"`
	HasWindow    int    `json:"has_window,omitempty"`
	MarkInterval int    `json:"markInterval,omitempty"` // 标记间隔
}

type ItemBase struct {
	Category Category `json:"category"`
}
