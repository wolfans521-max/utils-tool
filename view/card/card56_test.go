package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard56 校验工厂方法是否正确设置 card_type
func TestNewCard56(t *testing.T) {
	c := NewCard56()
	if c == nil {
		t.Fatalf("NewCard56 returned nil")
	}
	if c.CardType != 56 {
		t.Fatalf("CardType expected 56, got %d", c.CardType)
	}
}

// TestC56JSONStructure 构造一个完整的 C56，校验关键字段的 JSON 结构
func TestC56JSONStructure(t *testing.T) {
	card := NewCard56()

	card.Title = "你支持哪一方？"
	card.PositiveSide = &C56Side{
		Desc:  "支持",
		Count: "1000人",
		Button: &Button{
			Type:     "follow",
			Name:     "支持按钮",
			BtnTitle: "支持",
			Scheme:   "sinaweibo://vote?action=positive",
		},
		ActionLog: map[string]any{
			"act_code": 1,
			"ext":      "positive_ext",
			"uicode":   "uicode",
			"fid":      "fid",
			"luicode":  "luicode",
			"lfid":     "lfid",
		},
		Pic: "https://img.test/positive.png",
	}
	card.NegativeSide = &C56Side{
		Desc:  "反对",
		Count: "800人",
		Button: &Button{
			Type:     "follow",
			Name:     "反对按钮",
			BtnTitle: "反对",
			Scheme:   "sinaweibo://vote?action=negative",
		},
		ActionLog: map[string]any{
			"act_code": 2,
			"ext":      "negative_ext",
			"uicode":   "uicode",
			"fid":      "fid",
			"luicode":  "luicode",
			"lfid":     "lfid",
		},
		Pic: "https://img.test/negative.png",
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C56 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(56) {
		t.Fatalf("card_type expect 56, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title"]; !ok || v != "你支持哪一方？" {
		t.Fatalf("title expect '你支持哪一方？', got %v (ok=%v)", v, ok)
	}

	// positive_side 字段
	positiveSide, ok := m["positive_side"].(map[string]any)
	if !ok {
		t.Fatalf("positive_side expect object, got %T", m["positive_side"])
	}
	if v := positiveSide["desc"]; v != "支持" {
		t.Fatalf("positive_side.desc unexpected: %v", v)
	}
	if v := positiveSide["count"]; v != "1000人" {
		t.Fatalf("positive_side.count unexpected: %v", v)
	}
	if v := positiveSide["pic"]; v != "https://img.test/positive.png" {
		t.Fatalf("positive_side.pic unexpected: %v", v)
	}
	if v := positiveSide["button"]; v == nil {
		t.Fatalf("positive_side.button should not be nil")
	}
	if v := positiveSide["action_log"]; v == nil {
		t.Fatalf("positive_side.action_log should not be nil")
	}

	// negative_side 字段
	negativeSide, ok := m["negative_side"].(map[string]any)
	if !ok {
		t.Fatalf("negative_side expect object, got %T", m["negative_side"])
	}
	if v := negativeSide["desc"]; v != "反对" {
		t.Fatalf("negative_side.desc unexpected: %v", v)
	}
	if v := negativeSide["count"]; v != "800人" {
		t.Fatalf("negative_side.count unexpected: %v", v)
	}
	if v := negativeSide["pic"]; v != "https://img.test/negative.png" {
		t.Fatalf("negative_side.pic unexpected: %v", v)
	}
	if v := negativeSide["button"]; v == nil {
		t.Fatalf("negative_side.button should not be nil")
	}
	if v := negativeSide["action_log"]; v == nil {
		t.Fatalf("negative_side.action_log should not be nil")
	}
}

// TestC56EmptyFields 测试空字段时 omitempty 是否生效
func TestC56EmptyFields(t *testing.T) {
	card := NewCard56()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C56 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(56) {
		t.Fatalf("card_type expect 56, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["title"]; ok {
		t.Fatalf("title should be omitted when empty")
	}
	if _, ok := m["positive_side"]; ok {
		t.Fatalf("positive_side should be omitted when nil")
	}
	if _, ok := m["negative_side"]; ok {
		t.Fatalf("negative_side should be omitted when nil")
	}
}

// TestC56PartialFields 测试部分字段填充
func TestC56PartialFields(t *testing.T) {
	card := NewCard56()
	card.Title = "测试标题"
	card.PositiveSide = &C56Side{
		Desc:  "支持理由",
		Count: "500人",
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C56 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["title"]; !ok || v != "测试标题" {
		t.Fatalf("title expect '测试标题', got %v (ok=%v)", v, ok)
	}

	positiveSide, ok := m["positive_side"].(map[string]any)
	if !ok {
		t.Fatalf("positive_side expect object, got %T", m["positive_side"])
	}
	if v := positiveSide["desc"]; v != "支持理由" {
		t.Fatalf("positive_side.desc unexpected: %v", v)
	}
	if v := positiveSide["count"]; v != "500人" {
		t.Fatalf("positive_side.count unexpected: %v", v)
	}

	// 验证 omitempty 生效的字段
	if _, ok := positiveSide["button"]; ok {
		t.Fatalf("positive_side.button should be omitted when nil")
	}
	if _, ok := positiveSide["action_log"]; ok {
		t.Fatalf("positive_side.action_log should be omitted when nil")
	}
	if _, ok := positiveSide["pic"]; ok {
		t.Fatalf("positive_side.pic should be omitted when empty")
	}

	// negative_side 应该被忽略
	if _, ok := m["negative_side"]; ok {
		t.Fatalf("negative_side should be omitted when nil")
	}
}
