package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard4 校验工厂方法是否正确设置 card_type
func TestNewCard4(t *testing.T) {
	c := NewCard4()
	if c == nil {
		t.Fatalf("NewCard4 returned nil")
	}
	if c.CardType != 4 {
		t.Fatalf("CardType expected 4, got %d", c.CardType)
	}
}

// TestC4JSONStructure 构造一个完整的 C4，校验关键字段的 JSON 结构
func TestC4JSONStructure(t *testing.T) {
	card := NewCard4()

	card.ItemId = "card4_test_item"
	card.Scheme = "sinaweibo://detail?id=123"
	card.ImageRightPadding = 8
	card.ContentStyle = 1
	card.Desc = "这是一个文本卡片标题"
	card.RightDesc = "查看详情"
	card.Height = 60
	card.DisableBottomLine = 1
	card.Bold = 1
	card.IsHtmlText = true
	card.DescExtr = "附加描述信息"
	card.Pic = "http://example.com/left.png"
	card.DisplayArrow = 1
	card.OpenUrl = "sinaweibo://browser?url=https://example.com"
	card.AvatarUrl = "http://example.com/avatar.jpg"
	card.Icon = "http://example.com/icon.gif"
	card.IconUrl = "http://api.example.com/dismiss"

	card.ReadConfigs = []*C4ReadConfig{
		{
			StartTime:       0,
			EndTime:         86400,
			ContentColorKey: "gray_text",
			ContentColor:    "#999999",
		},
		{
			StartTime:       86400,
			EndTime:         172800,
			ContentColorKey: "light_gray",
			ContentColor:    "#CCCCCC",
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C4 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(4) {
		t.Fatalf("card_type expect 4, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc"]; !ok || v != "这是一个文本卡片标题" {
		t.Fatalf("desc unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["right_desc"]; !ok || v != "查看详情" {
		t.Fatalf("right_desc unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["height"]; !ok || v != float64(60) {
		t.Fatalf("height unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["bold"]; !ok || v != float64(1) {
		t.Fatalf("bold unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["is_html_text"]; !ok || v != true {
		t.Fatalf("is_html_text unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["pic"]; !ok || v != "http://example.com/left.png" {
		t.Fatalf("pic unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["avatar_url"]; !ok || v != "http://example.com/avatar.jpg" {
		t.Fatalf("avatar_url unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["icon"]; !ok || v != "http://example.com/icon.gif" {
		t.Fatalf("icon unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["display_arrow"]; !ok || v != float64(1) {
		t.Fatalf("display_arrow unexpected: %v (ok=%v)", v, ok)
	}

	// read_configs 数组
	readConfigs, ok := m["read_configs"].([]any)
	if !ok || len(readConfigs) != 2 {
		t.Fatalf("read_configs expect array of 2, got %T len=%d", m["read_configs"], len(readConfigs))
	}
	config0, ok := readConfigs[0].(map[string]any)
	if !ok {
		t.Fatalf("read_configs[0] expect object, got %T", readConfigs[0])
	}
	if v := config0["end_time"]; v != float64(86400) {
		t.Fatalf("read_configs[0].end_time unexpected: %v", v)
	}
	if v := config0["content_color_key"]; v != "gray_text" {
		t.Fatalf("read_configs[0].content_color_key unexpected: %v", v)
	}
	if v := config0["content_color"]; v != "#999999" {
		t.Fatalf("read_configs[0].content_color unexpected: %v", v)
	}
}

// TestC4EmptyFields 测试空字段时 omitempty 是否生效
func TestC4EmptyFields(t *testing.T) {
	card := NewCard4()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C4 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(4) {
		t.Fatalf("card_type expect 4, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["desc"]; ok {
		t.Fatalf("desc should be omitted when empty")
	}
	if _, ok := m["right_desc"]; ok {
		t.Fatalf("right_desc should be omitted when empty")
	}
	if _, ok := m["read_configs"]; ok {
		t.Fatalf("read_configs should be omitted when nil")
	}
	if _, ok := m["is_html_text"]; ok {
		t.Fatalf("is_html_text should be omitted when false")
	}
}

// TestC4PartialFields 测试部分字段填充
func TestC4PartialFields(t *testing.T) {
	card := NewCard4()
	card.Desc = "简单标题"
	card.DisplayArrow = 1
	card.OpenUrl = "sinaweibo://browser?url=https://example.com"

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C4 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["desc"]; !ok || v != "简单标题" {
		t.Fatalf("desc expect '简单标题', got %v (ok=%v)", v, ok)
	}
	if v, ok := m["display_arrow"]; !ok || v != float64(1) {
		t.Fatalf("display_arrow expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["openurl"]; !ok || v != "sinaweibo://browser?url=https://example.com" {
		t.Fatalf("openurl unexpected: %v (ok=%v)", v, ok)
	}

	// 验证不存在的字段
	if _, ok := m["right_desc"]; ok {
		t.Fatalf("right_desc should be omitted when empty")
	}
	if _, ok := m["read_configs"]; ok {
		t.Fatalf("read_configs should be omitted when nil")
	}
}

// TestC4ReadConfigStructure 测试已读配置结构
func TestC4ReadConfigStructure(t *testing.T) {
	card := NewCard4()
	card.Desc = "带已读配置的卡片"
	card.ReadConfigs = []*C4ReadConfig{
		{
			StartTime:       100,
			EndTime:         3600,
			ContentColorKey: "read_color",
			ContentColor:    "#888888",
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C4 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	readConfigs, ok := m["read_configs"].([]any)
	if !ok || len(readConfigs) != 1 {
		t.Fatalf("read_configs expect array of 1")
	}

	config, ok := readConfigs[0].(map[string]any)
	if !ok {
		t.Fatalf("read_configs[0] expect object")
	}

	// 注意：start_time 为 0 时会被 omitempty 忽略，所以使用非零值测试
	if v := config["start_time"]; v != float64(100) {
		t.Fatalf("start_time unexpected: %v", v)
	}
	if v := config["end_time"]; v != float64(3600) {
		t.Fatalf("end_time unexpected: %v", v)
	}
}
