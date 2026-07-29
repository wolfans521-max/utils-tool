package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard19 校验工厂方法是否正确设置 card_type
func TestNewCard19(t *testing.T) {
	c := NewCard19()
	if c == nil {
		t.Fatalf("NewCard19 returned nil")
	}
	if c.CardType != 19 {
		t.Fatalf("CardType expected 19, got %d", c.CardType)
	}
}

// TestC19JSONStructure 构造一个 C19，校验关键字段的 JSON 结构
// 使用工厂方法 NewCard19() 创建，以符合标准使用方式
func TestC19JSONStructure(t *testing.T) {
	card := NewCard19()

	card.DefaultRows = 2
	card.DividerColor = "#eeeeee"
	card.CardBgColor = "#ffffff"
	card.CardBgColorDark = "#000000"
	card.IsNewsquareUistyle = 1
	card.Mode = 3
	card.PosId = "posid-1"
	card.IsRefactorStyle = 1
	card.CardPadding = map[string]int{
		"left":   -2,
		"top":    0,
		"right":  2,
		"bottom": 2,
	}
	card.Col = 3
	card.Group = []C19Group{
		{
			TitleSub: "sub1",
			Pic:      "http://img.test/p1.png",
			Scheme:   "sinaweibo://cardlist?containerid=1002001",
			ItemId:   "item1",
			ActionLog: map[string]any{
				"act_code": 1,
			},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C19 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(19) {
		t.Fatalf("card_type expect 19, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["default_rows"]; !ok || v != float64(2) {
		t.Fatalf("default_rows expect 2, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["divider_color"]; !ok || v != "#eeeeee" {
		t.Fatalf("divider_color expect #eeeeee, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["card_bg_color"]; !ok || v != "#ffffff" {
		t.Fatalf("card_bg_color expect #ffffff, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["card_bg_color_dark"]; !ok || v != "#000000" {
		t.Fatalf("card_bg_color_dark expect #000000, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["is_newsquare_uistyle"]; !ok || v != float64(1) {
		t.Fatalf("is_newsquare_uistyle expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["mode"]; !ok || v != float64(3) {
		t.Fatalf("mode expect 3, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["posid"]; !ok || v != "posid-1" {
		t.Fatalf("posid expect posid-1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["isRefactorStyle"]; !ok || v != float64(1) {
		t.Fatalf("isRefactorStyle expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["col"]; !ok || v != float64(3) {
		t.Fatalf("col expect 3, got %v (ok=%v)", v, ok)
	}

	// group 列表
	group, ok := m["group"].([]any)
	if !ok || len(group) == 0 {
		t.Fatalf("group expect non-empty array, got %T len=%d", m["group"], len(group))
	}
	g0, ok := group[0].(map[string]any)
	if !ok {
		t.Fatalf("group[0] expect object, got %T", group[0])
	}
	if v := g0["title_sub"]; v != "sub1" {
		t.Fatalf("group[0].title_sub unexpected: %v", v)
	}
	if v := g0["pic"]; v != "http://img.test/p1.png" {
		t.Fatalf("group[0].pic unexpected: %v", v)
	}
	if v := g0["scheme"]; v != "sinaweibo://cardlist?containerid=1002001" {
		t.Fatalf("group[0].scheme unexpected: %v", v)
	}
	if v := g0["itemid"]; v != "item1" {
		t.Fatalf("group[0].itemid unexpected: %v", v)
	}
	if v := g0["action_log"]; v == nil {
		t.Fatalf("group[0].action_log should not be nil")
	}
}
