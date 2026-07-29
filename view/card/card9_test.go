package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard9 校验工厂方法是否正确设置 card_type
func TestNewCard9(t *testing.T) {
	c := NewCard9()
	if c == nil {
		t.Fatalf("NewCard9 returned nil")
	}
	if c.CardType != 9 {
		t.Fatalf("CardType expected 9, got %d", c.CardType)
	}
}

// TestC9EmptyFields 测试空字段时 omitempty 是否生效
func TestC9EmptyFields(t *testing.T) {
	card := NewCard9()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C9 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(9) {
		t.Fatalf("card_type expect 9, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["mblog"]; ok {
		t.Fatalf("mblog should be omitted when nil")
	}
	if _, ok := m["rating"]; ok {
		t.Fatalf("rating should be omitted when empty")
	}
}
