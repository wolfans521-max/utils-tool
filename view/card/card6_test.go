package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard6 校验工厂方法是否正确设置 card_type
func TestNewCard6(t *testing.T) {
	c := NewCard6()
	if c == nil {
		t.Fatalf("NewCard6 returned nil")
	}
	if c.CardType != 6 {
		t.Fatalf("CardType expected 6, got %d", c.CardType)
	}
}

// TestC6JSONStructure 构造一个完整的 C6，校验关键字段的 JSON 结构
func TestC6JSONStructure(t *testing.T) {
	card := NewCard6()

	card.ItemId = "card6_test_item"
	card.Scheme = "sinaweibo://detail?id=123"
	card.Desc = "这是一个文本卡片内容，支持多种颜色显示"
	card.TitleColor = "#FF6600"
	card.TitleColorDark = "#FF9933"
	card.ShowType = 1
	card.OpenUrl = "sinaweibo://browser?url=https://example.com"

	card.Buttons = []*Button{
		{
			Type: "link",
			Name: "查看详情",
			Pic:  "http://example.com/btn.png",
			Params: map[string]any{
				"url": "https://example.com",
			},
			ActionLog: map[string]any{
				"act_code": 100,
			},
		},
		{
			Type: "follow",
			Name: "关注",
			Params: map[string]any{
				"uid": 123456789,
			},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C6 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(6) {
		t.Fatalf("card_type expect 6, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc"]; !ok || v != "这是一个文本卡片内容，支持多种颜色显示" {
		t.Fatalf("desc unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_color"]; !ok || v != "#FF6600" {
		t.Fatalf("title_color unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_color_dark"]; !ok || v != "#FF9933" {
		t.Fatalf("title_color_dark unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["show_type"]; !ok || v != float64(1) {
		t.Fatalf("show_type unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["openurl"]; !ok || v != "sinaweibo://browser?url=https://example.com" {
		t.Fatalf("openurl unexpected: %v (ok=%v)", v, ok)
	}

	// buttons 数组
	buttons, ok := m["buttons"].([]any)
	if !ok || len(buttons) != 2 {
		t.Fatalf("buttons expect array of 2, got %T len=%d", m["buttons"], len(buttons))
	}
	btn0, ok := buttons[0].(map[string]any)
	if !ok {
		t.Fatalf("buttons[0] expect object, got %T", buttons[0])
	}
	if v := btn0["type"]; v != "link" {
		t.Fatalf("buttons[0].type unexpected: %v", v)
	}
	if v := btn0["name"]; v != "查看详情" {
		t.Fatalf("buttons[0].name unexpected: %v", v)
	}

	btn1, ok := buttons[1].(map[string]any)
	if !ok {
		t.Fatalf("buttons[1] expect object, got %T", buttons[1])
	}
	if v := btn1["type"]; v != "follow" {
		t.Fatalf("buttons[1].type unexpected: %v", v)
	}
}

// TestC6EmptyFields 测试空字段时 omitempty 是否生效
func TestC6EmptyFields(t *testing.T) {
	card := NewCard6()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C6 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(6) {
		t.Fatalf("card_type expect 6, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["desc"]; ok {
		t.Fatalf("desc should be omitted when empty")
	}
	if _, ok := m["title_color"]; ok {
		t.Fatalf("title_color should be omitted when empty")
	}
	if _, ok := m["buttons"]; ok {
		t.Fatalf("buttons should be omitted when nil")
	}
}

// TestC6PartialFields 测试部分字段填充
func TestC6PartialFields(t *testing.T) {
	card := NewCard6()
	card.Desc = "简单文本内容"
	card.ShowType = 2

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C6 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["desc"]; !ok || v != "简单文本内容" {
		t.Fatalf("desc expect '简单文本内容', got %v (ok=%v)", v, ok)
	}
	if v, ok := m["show_type"]; !ok || v != float64(2) {
		t.Fatalf("show_type expect 2, got %v (ok=%v)", v, ok)
	}

	// 验证不存在的字段
	if _, ok := m["title_color"]; ok {
		t.Fatalf("title_color should be omitted when empty")
	}
	if _, ok := m["buttons"]; ok {
		t.Fatalf("buttons should be omitted when nil")
	}
}

// TestC6ShowTypeColors 测试不同 show_type 值
func TestC6ShowTypeColors(t *testing.T) {
	tests := []struct {
		name     string
		showType int8
		desc     string
	}{
		{"默认颜色带箭头", 0, "默认样式"},
		{"绿色", 1, "绿色文本"},
		{"红色", 2, "红色文本"},
		{"蓝色", 3, "蓝色文本"},
		{"高亮白色", 4, "高亮白色文本"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := NewCard6()
			card.Desc = tt.desc
			card.ShowType = tt.showType

			b, err := json.Marshal(card)
			if err != nil {
				t.Fatalf("marshal C6 failed: %v", err)
			}

			var m map[string]any
			if err := json.Unmarshal(b, &m); err != nil {
				t.Fatalf("unmarshal to map failed: %v", err)
			}

			if v, ok := m["show_type"]; tt.showType != 0 && (!ok || v != float64(tt.showType)) {
				t.Fatalf("show_type expect %d, got %v (ok=%v)", tt.showType, v, ok)
			}
		})
	}
}

// TestC6WithButtons 测试带按钮的卡片
func TestC6WithButtons(t *testing.T) {
	card := NewCard6()
	card.Desc = "带按钮的文本卡片"
	card.Buttons = []*Button{
		{
			Type: "default",
			Name: "按钮1",
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C6 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	buttons, ok := m["buttons"].([]any)
	if !ok || len(buttons) != 1 {
		t.Fatalf("buttons expect array of 1")
	}

	btn, ok := buttons[0].(map[string]any)
	if !ok {
		t.Fatalf("buttons[0] expect object")
	}

	if v := btn["name"]; v != "按钮1" {
		t.Fatalf("buttons[0].name unexpected: %v", v)
	}
}
