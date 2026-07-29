package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard246 校验工厂方法是否正确设置 card_type
func TestNewCard246(t *testing.T) {
	c := NewCard246()
	if c == nil {
		t.Fatalf("NewCard246 returned nil")
	}
	if c.CardType != 246 {
		t.Fatalf("CardType expected 246, got %d", c.CardType)
	}
}

// TestC246JSONStructure 测试 C246 的 JSON 序列化
func TestC246JSONStructure(t *testing.T) {
	card := NewCard246()
	card.PaddingTop = 10
	card.PaddingBottom = 12
	card.BackgroundPic = "http://img.test/bg.jpg"
	card.TitlePic = "http://img.test/title.png"
	card.GridColumn = 5

	card.BannerStyle = &C246BannerStyle{
		Ratio:       0.5,
		LeftMargin:  14,
		RightMargin: 14,
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C246 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	if v, ok := m["card_type"]; !ok || v != float64(246) {
		t.Fatalf("card_type expect 246, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["padding_top"]; !ok || v != float64(10) {
		t.Fatalf("padding_top expect 10, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["grid_column"]; !ok || v != float64(5) {
		t.Fatalf("grid_column expect 5, got %v (ok=%v)", v, ok)
	}

	// 验证嵌套结构
	bs, ok := m["banner_style"].(map[string]any)
	if !ok {
		t.Fatalf("banner_style expect object, got %T", m["banner_style"])
	}
	if v := bs["ratio"]; v != float64(0.5) {
		t.Fatalf("banner_style.ratio unexpected: %v", v)
	}
}
