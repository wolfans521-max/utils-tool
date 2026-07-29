package card

// C148 催更卡片
// 用于视频催更展示
type C148 struct {
	BaseItem
	AuthorUid      string      `json:"author_uid,omitempty"`      // 视频作者uid
	PromoterEnable int8        `json:"promoter_enable,omitempty"` // 当前是否可以催更，1为可催更
	Content1       string      `json:"content1,omitempty"`        // 第一行文本
	Content1Color  string      `json:"content1_color,omitempty"`  // 第一行文本字色
	Content1Size   int         `json:"content1_size,omitempty"`   // 第一行文本字体
	Content2       string      `json:"content2,omitempty"`        // 第二行文本
	Content2Color  string      `json:"content2_color,omitempty"`  // 第二行文本字色
	Content2Size   int         `json:"content2_size,omitempty"`   // 第二行文本字体
	BackgroundType int8        `json:"background_type,omitempty"` // 背景样式 0：默认 1：样式一
	IconArray      []string    `json:"icon_array,omitempty"`      // 催更需要展示的头像icon数组
	Button         *C148Button `json:"button,omitempty"`          // 按钮
	OpenUrl        string      `json:"openurl,omitempty"`         // 打开链接
}

// C148Button 按钮
type C148Button struct {
	Enable      int8           `json:"enable,omitempty"`       // 是否可以点击，1是可点击
	Style       int8           `json:"style,omitempty"`        // 1是蓝色，2是黄色，默认蓝色
	EnableText  string         `json:"enable_text,omitempty"`  // 可点击时按钮文案
	DisableText string         `json:"disable_text,omitempty"` // 不可点击时按钮文案
	Operation   *C148Operation `json:"operation,omitempty"`    // 按钮的点击行为
}

// C148Operation 操作
type C148Operation struct {
	Scheme        string            `json:"scheme,omitempty"`         // 点击按钮后执行scheme
	Path          string            `json:"path,omitempty"`           // 点击按钮请求接口的path
	ReturnDisable int8              `json:"return_disable,omitempty"` // 按钮交互成功后是否需要置灰
	Params        map[string]string `json:"params,omitempty"`         // path接口请求参数
	SuccessToast  string            `json:"success_toast,omitempty"`  // 按钮点击后成功提示toast文案
	FailToast     string            `json:"fail_toast,omitempty"`     // 按钮点击后失败提示文案
}

// NewCard148 创建Card148实例
func NewCard148() *C148 {
	return &C148{BaseItem: BaseItem{Base: Base{CardType: 148}}}
}
