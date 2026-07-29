package card

// C156 Profile页小视频tab已赞按钮卡片
// 可横滑
type C156 struct {
	BaseItem
	CommendInfo []*C156Commend `json:"commend_info,omitempty"` // 推荐信息列表
	OpenUrl     string         `json:"openurl,omitempty"`      // 打开链接
}

// C156Commend 推荐信息
type C156Commend struct {
	Icon        string         `json:"icon,omitempty"`         // 显示图标
	Text        string         `json:"text,omitempty"`         // 显示文案
	AccessRight int8           `json:"access_right,omitempty"` // 查看权限 1：所有人 2:我关注的人 3:我的粉丝 4:仅自己
	Scheme      string         `json:"scheme,omitempty"`       // 跳转链接
	ActionLog   map[string]any `json:"actionlog,omitempty"`    // 打点日志
}

// NewCard156 创建Card156实例
func NewCard156() *C156 {
	return &C156{BaseItem: BaseItem{Base: Base{CardType: 156}}}
}
