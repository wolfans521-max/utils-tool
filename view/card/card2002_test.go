package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard2002 校验工厂方法是否正确设置 card_type
func TestNewCard2002(t *testing.T) {
	c := NewCard2002()
	if c == nil {
		t.Fatalf("NewCard2002 returned nil")
	}
	if c.CardType != 2002 {
		t.Fatalf("CardType expected 2002, got %d", c.CardType)
	}
}

// TestC2002JSONStructure 测试 C2002 的 JSON 序列化
func TestC2002JSONStructure(t *testing.T) {
	card := NewCard2002()
	card.Title = "距离下一个活动"
	card.EndTitle = "结束了"
	card.CountType = 1 // 1:正计时 (使用非零值以便测试)
	card.StartTime = 1640000000
	card.TitleColor = "#333333"
	card.NumberColor = "#FF0000"
	card.NumberWarnColor = "#FF6600"
	card.UnitColor = "#666666"
	card.UnitWarnColor = "#FF6600"

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C2002 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	if v, ok := m["card_type"]; !ok || v != float64(2002) {
		t.Fatalf("card_type expect 2002, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title"]; !ok || v != "距离下一个活动" {
		t.Fatalf("title expect 距离下一个活动, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["end_title"]; !ok || v != "结束了" {
		t.Fatalf("end_title expect 结束了, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["count_type"]; !ok || v != float64(1) {
		t.Fatalf("count_type expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["start_time"]; !ok || v != float64(1640000000) {
		t.Fatalf("start_time expect 1640000000, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_color"]; !ok || v != "#333333" {
		t.Fatalf("title_color expect #333333, got %v (ok=%v)", v, ok)
	}
}
