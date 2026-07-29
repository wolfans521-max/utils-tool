package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard17 校验工厂方法是否正确设置 card_type
func TestNewCard17(t *testing.T) {
	c := NewCard17()
	if c == nil {
		t.Fatalf("NewCard17 returned nil")
	}
	if c.CardType != 17 {
		t.Fatalf("CardType expected 17, got %d", c.CardType)
	}
}

// TestC17JSONStructure 构造一个完整的 C17，校验关键字段的 JSON 结构
// 使用工厂方法 NewCard17() 创建，以符合标准使用方式
func TestC17JSONStructure(t *testing.T) {
	card := NewCard17()

	card.Col = 2
	card.Group = []*C17Texts{
		{
			TitleSub:      "sub-title",
			Icon:          "https://img.test/icon.png",
			Scheme:        "sinaweibo://cardlist?containerid=1001001",
			TitleSubColor: "#333333",
			ActionLog: map[string]any{
				"act_code": 1,
				"ext":      "ext",
				"uicode":   "uicode",
				"fid":      "fid",
				"luicode":  "luicode",
				"lfid":     "lfid",
			},
			IconWidth:  16,
			IconHeight: 16,
		},
	}
	card.Bottom = &C17Bottom{
		IconSize:     18,
		TopMargin:    8,
		BottomMargin: 8,
		ActLog: map[string]any{
			"act_code": 2,
		},
		Title:      "bottom-title",
		Icon:       "https://img.test/bottom-icon.png",
		TitleColor: "#999999",
		Scheme:     "sinaweibo://cardlist?containerid=1001002",
		TitleSize:  12,
	}
	card.CoverHeight = 120
	card.CoverContent = "cover-content"
	card.Top = &C17Top{
		RightScheme: "sinaweibo://cardlist?containerid=1001003",
		TitleSize:   16,
		Title:       "top-title",
		TopMargin:   4,
		Icon:        "https://img.test/top-icon.png",
		ActLog: map[string]any{
			"act_code": 3,
		},
		BottomMargin:          4,
		CloseStatus:           1,
		IconSize:              20,
		ItemHeight:            40,
		TitleLeftMargin:       10,
		Desc:                  "desc",
		Refresh:               &C17TopRefresh{Icon: "https://img.test/refresh.png", Title: "刷新", Scheme: "sinaweibo://refresh"},
		HideTopLine:           true,
		TitleLeftMarginOffset: 2,
	}
	card.TitleBold = 1
	card.TitleColor = "#000000"
	card.TitleColorDark = "#ffffff"
	card.GroupItemStartOffset = 1

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C17 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(17) {
		t.Fatalf("card_type expect 17, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["col"]; !ok || v != float64(2) {
		t.Fatalf("col expect 2, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["cover_height"]; !ok || v != float64(120) {
		t.Fatalf("cover_height expect 120, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["cover_content"]; !ok || v != "cover-content" {
		t.Fatalf("cover_content expect cover-content, got %v (ok=%v)", v, ok)
	}

	// 标题相关字段
	if v, ok := m["title_bold"]; !ok || v != float64(1) {
		t.Fatalf("title_bold expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_color"]; !ok || v != "#000000" {
		t.Fatalf("title_color expect #000000, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_color_dark"]; !ok || v != "#ffffff" {
		t.Fatalf("title_color_dark expect #ffffff, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["group_item_start_offset"]; !ok || v != float64(1) {
		t.Fatalf("group_item_start_offset expect 1, got %v (ok=%v)", v, ok)
	}

	// group 数组
	group, ok := m["group"].([]any)
	if !ok || len(group) == 0 {
		t.Fatalf("group expect non-empty array, got %T len=%d", m["group"], len(group))
	}
	g0, ok := group[0].(map[string]any)
	if !ok {
		t.Fatalf("group[0] expect object, got %T", group[0])
	}
	if v := g0["title_sub"]; v != "sub-title" {
		t.Fatalf("group[0].title_sub unexpected: %v", v)
	}
	if v := g0["icon"]; v != "https://img.test/icon.png" {
		t.Fatalf("group[0].icon unexpected: %v", v)
	}
	if v := g0["scheme"]; v != "sinaweibo://cardlist?containerid=1001001" {
		t.Fatalf("group[0].scheme unexpected: %v", v)
	}
	if v := g0["title_sub_color"]; v != "#333333" {
		t.Fatalf("group[0].title_sub_color unexpected: %v", v)
	}
	if v := g0["action_log"]; v == nil {
		t.Fatalf("group[0].action_log should not be nil")
	}

	// top 字段
	top, ok := m["top"].(map[string]any)
	if !ok {
		t.Fatalf("top expect object, got %T", m["top"])
	}
	if v := top["title"]; v != "top-title" {
		t.Fatalf("top.title unexpected: %v", v)
	}
	if v := top["right_scheme"]; v != "sinaweibo://cardlist?containerid=1001003" {
		t.Fatalf("top.right_scheme unexpected: %v", v)
	}
	if v := top["actlog"]; v == nil {
		t.Fatalf("top.actlog should not be nil")
	}
	if v := top["refresh"]; v == nil {
		t.Fatalf("top.refresh should not be nil")
	}

	// bottom 字段
	bottom, ok := m["bottom"].(map[string]any)
	if !ok {
		t.Fatalf("bottom expect object, got %T", m["bottom"])
	}
	if v := bottom["title"]; v != "bottom-title" {
		t.Fatalf("bottom.title unexpected: %v", v)
	}
	if v := bottom["scheme"]; v != "sinaweibo://cardlist?containerid=1001002" {
		t.Fatalf("bottom.scheme unexpected: %v", v)
	}
	if v := bottom["actlog"]; v == nil {
		t.Fatalf("bottom.actlog should not be nil")
	}
}
