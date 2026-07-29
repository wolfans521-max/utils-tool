package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard265 校验工厂方法是否正确设置 card_type
func TestNewCard265(t *testing.T) {
	c := NewCard265()
	if c == nil {
		t.Fatalf("NewCard265 returned nil")
	}
	if c.CardType != 265 {
		t.Fatalf("CardType expected 265, got %d", c.CardType)
	}
}

// TestC265JSONStructure 测试 C265 的 JSON 序列化
func TestC265JSONStructure(t *testing.T) {
	card := NewCard265()
	card.ItemInterval = 7
	card.PicWidth = 142
	card.PicWidthRatio = 16
	card.PicHeightRatio = 9
	card.PageLeftRightPadding = 14

	card.Items = []*C265Item{
		{
			Scheme: "http://test.com/scheme",
			TopContent: &C265TopContent{
				Pic:            "http://img.test/pic.jpg",
				PlayIconWidth:  18,
				PlayIconHeight: 21,
				RadiusCorner:   []int{6, 6, 6, 6},
			},
			BottomContent: &C265BottomContent{
				Text:          "Apple 贺岁片",
				TextAlignLeft: true,
				TextSize:      12,
				TextColor:     "#333333",
				BottomHeight:  25,
			},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C265 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	if v, ok := m["card_type"]; !ok || v != float64(265) {
		t.Fatalf("card_type expect 265, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["item_interval"]; !ok || v != float64(7) {
		t.Fatalf("item_interval expect 7, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["pic_width"]; !ok || v != float64(142) {
		t.Fatalf("pic_width expect 142, got %v (ok=%v)", v, ok)
	}

	// 验证 items 数组
	items, ok := m["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("items expect array with 1 element, got %T", m["items"])
	}

	item := items[0].(map[string]any)
	if v := item["scheme"]; v != "http://test.com/scheme" {
		t.Fatalf("item.scheme unexpected: %v", v)
	}

	// 验证 top_content
	tc, ok := item["top_content"].(map[string]any)
	if !ok {
		t.Fatalf("top_content expect object, got %T", item["top_content"])
	}
	if v := tc["pic"]; v != "http://img.test/pic.jpg" {
		t.Fatalf("top_content.pic unexpected: %v", v)
	}

	// 验证 bottom_content
	bc, ok := item["bottom_content"].(map[string]any)
	if !ok {
		t.Fatalf("bottom_content expect object, got %T", item["bottom_content"])
	}
	if v := bc["text"]; v != "Apple 贺岁片" {
		t.Fatalf("bottom_content.text unexpected: %v", v)
	}
}
