package card

import (
	"encoding/json"
	"testing"
)

// TestNew118 校验工厂方法是否正确设置 card_type
func TestNew118(t *testing.T) {
	c := NewCard118()
	if c == nil {
		t.Fatalf("New118 returned nil")
	}
	if c.CardType != 118 {
		t.Fatalf("CardType expected 118, got %d", c.CardType)
	}
}

// TestC118JSONStructure 构造一个 C118，校验关键字段的 JSON 结构
// 使用工厂方法 NewCard118() 创建，以符合标准使用方式
func TestC118JSONStructure(t *testing.T) {
	card := NewCard118()

	card.LeftPadding = 10
	card.RightPadding = 20
	card.LoopInterval = 5
	card.NewStyle = 1
	card.CardPadding = C118CardPadding{
		Left:   1,
		Top:    2,
		Right:  3,
		Bottom: 4,
	}
	card.Items = []C119{
		{
			BaseItem: BaseItem{Base: Base{ItemId: "item-1", CardType: 119}},
			SubItem: []C119SubItem{
				{
					Scheme:            "sinaweibo://cardlist?containerid=1004001",
					ItemId:            "subitem-1",
					ContentType:       1,
					Pic:               "http://img.test/p1.png",
					NewStyle:          1,
					Title:             "title1",
					IsBigPic:          1,
					IconUrl:           "http://img.test/icon.png",
					Desc:              "desc1",
					ChannelBackground: "#ffffff",
					ActionLog: map[string]any{
						"act_code": 1,
					},
					Promotion: map[string]any{
						"key": "val",
					},
				},
			},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C118 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(118) {
		t.Fatalf("card_type expect 118, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["left_padding"]; !ok || v != float64(10) {
		t.Fatalf("left_padding expect 10, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["right_padding"]; !ok || v != float64(20) {
		t.Fatalf("right_padding expect 20, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["loop_interval"]; !ok || v != float64(5) {
		t.Fatalf("loop_interval expect 5, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["new_style"]; !ok || v != float64(1) {
		t.Fatalf("new_style expect 1, got %v (ok=%v)", v, ok)
	}

	// card_padding
	cp, ok := m["card_padding"].(map[string]any)
	if !ok {
		t.Fatalf("card_padding expect object, got %T", m["card_padding"])
	}
	if v := cp["left"]; v != float64(1) {
		t.Fatalf("card_padding.left unexpected: %v", v)
	}
	if v := cp["top"]; v != float64(2) {
		t.Fatalf("card_padding.top unexpected: %v", v)
	}
	if v := cp["right"]; v != float64(3) {
		t.Fatalf("card_padding.right unexpected: %v", v)
	}
	if v := cp["bottom"]; v != float64(4) {
		t.Fatalf("card_padding.bottom unexpected: %v", v)
	}

	// items 列表
	items, ok := m["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatalf("items expect non-empty array, got %T len=%d", m["items"], len(items))
	}
	it0, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("items[0] expect object, got %T", items[0])
	}
	if v := it0["itemid"]; v != "item-1" {
		t.Fatalf("items[0].itemid unexpected: %v", v)
	}
	if v := it0["card_type"]; v != float64(119) {
		t.Fatalf("items[0].card_type unexpected: %v", v)
	}

	subs, ok := it0["sub_item"].([]any)
	if !ok || len(subs) == 0 {
		t.Fatalf("sub_item expect non-empty array, got %T len=%d", it0["sub_item"], len(subs))
	}
	s0, ok := subs[0].(map[string]any)
	if !ok {
		t.Fatalf("sub_item[0] expect object, got %T", subs[0])
	}
	if v := s0["scheme"]; v != "sinaweibo://cardlist?containerid=1004001" {
		t.Fatalf("sub_item[0].scheme unexpected: %v", v)
	}
	if v := s0["itemid"]; v != "subitem-1" {
		t.Fatalf("sub_item[0].itemid unexpected: %v", v)
	}
	if v := s0["content_type"]; v != float64(1) {
		t.Fatalf("sub_item[0].content_type unexpected: %v", v)
	}
	if v := s0["pic"]; v != "http://img.test/p1.png" {
		t.Fatalf("sub_item[0].pic unexpected: %v", v)
	}
	if v := s0["newStyle"]; v != float64(1) {
		t.Fatalf("sub_item[0].newStyle unexpected: %v", v)
	}
	if v := s0["title"]; v != "title1" {
		t.Fatalf("sub_item[0].title unexpected: %v", v)
	}
	if v := s0["is_big_pic"]; v != float64(1) {
		t.Fatalf("sub_item[0].is_big_pic unexpected: %v", v)
	}
	if v := s0["icon_url"]; v != "http://img.test/icon.png" {
		t.Fatalf("sub_item[0].icon_url unexpected: %v", v)
	}
	if v := s0["desc"]; v != "desc1" {
		t.Fatalf("sub_item[0].desc unexpected: %v", v)
	}
	if v := s0["channel_background"]; v != "#ffffff" {
		t.Fatalf("sub_item[0].channel_background unexpected: %v", v)
	}
	if v := s0["actionlog"]; v == nil {
		t.Fatalf("sub_item[0].actionlog should not be nil")
	}
	if v := s0["promotion"]; v == nil {
		t.Fatalf("sub_item[0].promotion should not be nil")
	}
}
