package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard243 校验工厂方法是否正确设置 card_type
func TestNewCard243(t *testing.T) {
	c := NewCard243()
	if c == nil {
		t.Fatalf("NewCard243 returned nil")
	}
	if c.CardType != 243 {
		t.Fatalf("CardType expected 243, got %d", c.CardType)
	}
}

// TestC243JSONStructure 测试 C243 的 JSON 序列化
func TestC243JSONStructure(t *testing.T) {
	card := NewCard243()
	card.Text = "测试博文"
	card.MaxLine = 3
	card.ExtraInfo = "左下角文本"
	card.Pic = "http://img.test/pic.jpg"
	card.Duration = "03:45"

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C243 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	if v, ok := m["card_type"]; !ok || v != float64(243) {
		t.Fatalf("card_type expect 243, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["text"]; !ok || v != "测试博文" {
		t.Fatalf("text expect 测试博文, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["max_line"]; !ok || v != float64(3) {
		t.Fatalf("max_line expect 3, got %v (ok=%v)", v, ok)
	}
}
