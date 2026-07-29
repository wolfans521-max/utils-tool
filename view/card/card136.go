package card

// C136 视频社区热门频道入口卡片
// 用于视频社区热门频道入口展示
type C136 struct {
	BaseItem
	CardConfig *C136CardConfig `json:"card_config,omitempty"` // 卡片配置
	Entrance   []*C136Entrance `json:"entrance,omitempty"`    // 入口列表
	OpenUrl    string          `json:"openurl,omitempty"`     // 打开链接
}

// C136CardConfig 卡片配置
type C136CardConfig struct {
	RefreshRemove int8 `json:"refresh_remove,omitempty"` // 刷新移除
}

// C136Entrance 入口
type C136Entrance struct {
	Title       string         `json:"title,omitempty"`        // 标题
	Version     int64          `json:"version,omitempty"`      // 版本
	Icon        string         `json:"icon,omitempty"`         // 图标
	ContainerId int64          `json:"container_id,omitempty"` // 容器ID
	Scheme      string         `json:"scheme,omitempty"`       // 跳转链接
	ActionLog   map[string]any `json:"actionlog,omitempty"`    // 打点日志
}

// NewCard136 创建Card136实例
func NewCard136() *C136 {
	return &C136{BaseItem: BaseItem{Base: Base{CardType: 136}}}
}
