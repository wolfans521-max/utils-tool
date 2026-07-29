package card

// C122 投票卡片
type C122 struct {
	BaseItem
	VoteObject *C122VoteObject `json:"vote_object,omitempty"`
}

// C122VoteObject 投票对象
type C122VoteObject struct {
	Id                  string          `json:"id,omitempty"`
	VoteType            int8            `json:"vote_type,omitempty"`             // 投票类型 0 文字 1 图片+文字
	Content             string          `json:"content,omitempty"`               // 投票标题/描述
	Parted              int8            `json:"parted,omitempty"`                // 是否已经参与 0-未参与 1-已参与
	PartInfo            string          `json:"part_info,omitempty"`             // 投票人数信息说明 如"100人参与"
	State               int8            `json:"state,omitempty"`                 // 投票状态 0-结束 1-进行中
	ShareScheme         string          `json:"share_scheme,omitempty"`          // 分享scheme
	UserId              string          `json:"user_id,omitempty"`               // 发起投票的用户ID
	UserNick            string          `json:"user_nick,omitempty"`             // 发起投票的用户昵称
	ExpireDate          int64           `json:"expire_date,omitempty"`           // 结束时间 单位为毫秒
	VoteList            []*C122VoteItem `json:"vote_list,omitempty"`             // 投票情况列表
	ChoiceMoreScheme    string          `json:"choice_more_scheme,omitempty"`    // D31添加，下发后配置choice_count字段可以控制投票展示个数
	ChoiceMoreActionLog map[string]any  `json:"choice_more_actionlog,omitempty"` // D31添加，点击choice_more_scheme跳转的日志
	ChoiceCount         int             `json:"choice_count,omitempty"`          // 投票展示个数
}

// C122VoteItem 投票选项
type C122VoteItem struct {
	Id        string  `json:"id,omitempty"`
	Content   string  `json:"content,omitempty"`    // 投票Title
	Selected  int8    `json:"selected,omitempty"`   // 是否为当前选中的投票 1-选中 0-未选中
	PartNum   string  `json:"part_num,omitempty"`   // 投票人数
	PartRatio float64 `json:"part_ratio,omitempty"` // 投票占比，保留两位小数
	Pic       string  `json:"pic,omitempty"`        // 图片（vote_type=1时使用）
}

// NewCard122 创建Card122实例
func NewCard122() *C122 {
	return &C122{BaseItem: BaseItem{Base: Base{CardType: 122}}}
}
