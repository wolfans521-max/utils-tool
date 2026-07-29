package card

// C211 富文本和背景图片卡片
// 用于富文本和背景图片展示
type C211 struct {
	BaseItem
	LeftMargin            int                  `json:"leftMargin,omitempty"`            // 左边距
	RightMargin           int                  `json:"rightMargin,omitempty"`           // 右边距
	Height                int                  `json:"height,omitempty"`                // card高度
	RichTexts             []*C211RichText      `json:"richTexts,omitempty"`             // 富文本数组
	BackgroundImage       *C211BackgroundImage `json:"backgroundImage,omitempty"`       // 整体背景图
	RightImage            *C211RightImage      `json:"rightImage,omitempty"`            // 右侧背景图
	RightRichText         *C211RichText        `json:"rightRichText,omitempty"`         // 右侧富文本
	RightRichTextSelected *C211RichText        `json:"rightRichTextSelected,omitempty"` // 右侧富文本选中
	RightTextWidth        int                  `json:"rightTextWidth,omitempty"`        // 右侧富文本预留宽度
	OpenUrl               string               `json:"openurl,omitempty"`               // 打开链接
	Cleaned               bool                 `json:"cleaned,omitempty"`               // 是否清理
}

// C211RichText 富文本
type C211RichText struct {
	Content   any            `json:"content,omitempty"`   // 富文本内容数组
	Scheme    string         `json:"scheme,omitempty"`    // 跳转scheme
	ActionLog map[string]any `json:"actionlog,omitempty"` // 点击日志
	Cleaned   bool           `json:"cleaned,omitempty"`   // 是否清理
}

// C211BackgroundImage 背景图
type C211BackgroundImage struct {
	Url         string `json:"url,omitempty"`         // 图片URL
	UrlDark     string `json:"url_dark,omitempty"`    // 暗黑地址
	ContentMode int    `json:"contentMode,omitempty"` // 填充模式 1.WBScaleAspectBottomFit 2.WBScaleToFill else.WBScaleAspectTopFit
}

// C211RightImage 右侧图片
type C211RightImage struct {
	Url         string `json:"url,omitempty"`         // 图片URL
	UrlDark     string `json:"url_dark,omitempty"`    // 暗黑地址
	Width       int    `json:"width,omitempty"`       // 宽度
	Height      int    `json:"height,omitempty"`      // 高度
	RightMargin int    `json:"rightMargin,omitempty"` // 右边距
	TopMargin   int    `json:"topMargin,omitempty"`   // 上边距
}

// NewCard211 创建Card211实例
func NewCard211() *C211 {
	return &C211{BaseItem: BaseItem{Base: Base{CardType: 211}}}
}
