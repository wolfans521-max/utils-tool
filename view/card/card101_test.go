package card

import (
	"encoding/json"
	"testing"
)

// TestNew101 校验工厂方法是否正确设置 card_type
func TestNew101(t *testing.T) {
	c := NewCard101()
	if c == nil {
		t.Fatalf("New101 returned nil")
	}
	if c.CardType != 101 {
		t.Fatalf("CardType expected 101, got %d", c.CardType)
	}
}

// TestC101JSONStructure 构造一个 C101，校验关键字段的 JSON 结构
// 使用工厂方法 NewCard101() 创建，以符合标准使用方式
func TestC101JSONStructure(t *testing.T) {
	card := NewCard101()

	card.Title = "title"
	card.SubTitle = "sub"
	card.Desc = "desc"
	card.FontSize = 14
	card.TitleIsBold = 1
	card.DisplayArrow = 1
	card.NoHighlight = 1
	card.NewTitleIsBold = 1
	card.TitleColor = "#111111"
	card.TitleColorDark = "#222222"
	card.DescColor = "#333333"
	card.DescColorDark = "#444444"
	card.TitleDark = "#555555"
	card.LeftTagImg = "http://img.test/tag.png"
	card.LeftTagImgHeight = 20
	card.IsShowArrow = 1
	card.BottomLine = 1
	card.TopTagImgPadding = 2
	card.TagImg = "http://img.test/tag2.png"
	card.TopPadding = 3
	card.BottomPadding = 4
	card.UseJsonColor = true
	card.SubTitleColor = "#666666"
	card.SubTitleColorDark = "#777777"
	card.Height = 10
	card.ContentAlignBottom = 1
	card.PicTagStyle = 2
	card.LeftTagImgPadding = 5
	card.IsCenterVertical = 1
	card.LeftTagImgDark = "http://img.test/tag_dark.png"
	card.LeftTagImgWidth = 8
	card.Scheme = "sinaweibo://cardlist?containerid=1003001"

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C101 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(101) {
		t.Fatalf("card_type expect 101, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title"]; !ok || v != "title" {
		t.Fatalf("title expect title, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["sub_title"]; !ok || v != "sub" {
		t.Fatalf("sub_title expect sub, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc"]; !ok || v != "desc" {
		t.Fatalf("desc expect desc, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["font_size"]; !ok || v != float64(14) {
		t.Fatalf("font_size expect 14, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_is_bold"]; !ok || v != float64(1) {
		t.Fatalf("title_is_bold expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["display_arrow"]; !ok || v != float64(1) {
		t.Fatalf("display_arrow expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["no_highlight"]; !ok || v != float64(1) {
		t.Fatalf("no_highlight expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_color"]; !ok || v != "#111111" {
		t.Fatalf("title_color expect #111111, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_color_dark"]; !ok || v != "#222222" {
		t.Fatalf("title_color_dark expect #222222, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc_color"]; !ok || v != "#333333" {
		t.Fatalf("desc_color expect #333333, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc_color_dark"]; !ok || v != "#444444" {
		t.Fatalf("desc_color_dark expect #444444, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["is_show_arrow"]; !ok || v != float64(1) {
		t.Fatalf("is_show_arrow expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["use_json_color"]; !ok || v != true {
		t.Fatalf("use_json_color expect true, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["scheme"]; !ok || v != "sinaweibo://cardlist?containerid=1003001" {
		t.Fatalf("scheme expect sinaweibo://cardlist?containerid=1003001, got %v (ok=%v)", v, ok)
	}
}

