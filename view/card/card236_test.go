package card

import (
	"encoding/json"
	"testing"
)

// TestNew236 校验工厂方法是否正确设置 card_type
func TestNew236(t *testing.T) {
	c := NewCard236()
	if c == nil {
		t.Fatalf("New236 returned nil")
	}
	if c.CardType != 236 {
		t.Fatalf("CardType expected 236, got %d", c.CardType)
	}
}

// TestC236JSONStructure 构造一个 C236，校验嵌套 display / displayBackup 以及 wbox 字段的 JSON 结构
// 使用工厂方法 NewCard236() 创建，以符合标准使用方式
func TestC236JSONStructure(t *testing.T) {
	card := NewCard236()

	card.Display.MaxHeight = 300
	card.Display.CoverImageUrl = "http://img.test/cover.png"
	card.Display.CoverImageUrlDark = "http://img.test/cover_dark.png"
	card.Display.DisplayHeight = 200
	card.Display.DisplayRatio = 1.5
	card.Display.MeasureHeight = &MeasureHeight{
		Expression:        "min(200, screenHeight*0.5)",
		BundleVersionCode: "1000023",
	}

	card.DisplayBackup.DisplayHeight = 150
	card.DisplayBackup.DisplayRatio = 1.2
	card.DisplayBackup.CoverImageUrl = "http://img.test/cover_backup.png"
	card.DisplayBackup.CoverImageUrlDark = "http://img.test/cover_backup_dark.png"

	card.Scheme = "sinaweibo://cardlist?containerid=1005001"
	card.WboxScheme = "sinaweibo://wbox?app=foo"
	card.WboxMinRuntimeVersion = 10
	card.WboxMinSDKVersion = 20
	card.WboxParam = map[string]any{
		"foo": "bar",
		"num": 1,
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C236 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(236) {
		t.Fatalf("card_type expect 236, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["scheme"]; !ok || v != "sinaweibo://cardlist?containerid=1005001" {
		t.Fatalf("scheme expect sinaweibo://cardlist?containerid=1005001, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["wboxScheme"]; !ok || v != "sinaweibo://wbox?app=foo" {
		t.Fatalf("wboxScheme expect sinaweibo://wbox?app=foo, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["wboxMinRuntimeVersion"]; !ok || v != float64(10) {
		t.Fatalf("wboxMinRuntimeVersion expect 10, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["wboxMinSDKVersion"]; !ok || v != float64(20) {
		t.Fatalf("wboxMinSDKVersion expect 20, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["wboxParam"]; !ok || v == nil {
		t.Fatalf("wboxParam expect non-nil, got %v (ok=%v)", v, ok)
	}

	// display 嵌套字段
	display, ok := m["display"].(map[string]any)
	if !ok {
		t.Fatalf("display expect object, got %T", m["display"])
	}
	if v := display["maxHeight"]; v != float64(300) {
		t.Fatalf("display.maxHeight unexpected: %v", v)
	}
	if v := display["coverImageUrl"]; v != "http://img.test/cover.png" {
		t.Fatalf("display.coverImageUrl unexpected: %v", v)
	}
	if v := display["coverImageUrlDark"]; v != "http://img.test/cover_dark.png" {
		t.Fatalf("display.coverImageUrlDark unexpected: %v", v)
	}
	if v := display["displayHeight"]; v != float64(200) {
		t.Fatalf("display.displayHeight unexpected: %v", v)
	}
	if v := display["displayRatio"]; v != float64(1.5) {
		t.Fatalf("display.displayRatio unexpected: %v", v)
	}

	mh, ok := display["measureHeight"].(map[string]any)
	if !ok {
		t.Fatalf("display.measureHeight expect object, got %T", display["measureHeight"])
	}
	if v := mh["expression"]; v != "min(200, screenHeight*0.5)" {
		t.Fatalf("measureHeight.expression unexpected: %v", v)
	}
	if v := mh["bundleVersionCode"]; v != "1000023" {
		t.Fatalf("measureHeight.bundleVersionCode unexpected: %v", v)
	}

	// displayBackup 嵌套字段
	db, ok := m["DisplayBackup"].(map[string]any)
	if !ok {
		t.Fatalf("DisplayBackup expect object, got %T", m["DisplayBackup"])
	}
	if v := db["displayHeight"]; v != float64(150) {
		t.Fatalf("DisplayBackup.displayHeight unexpected: %v", v)
	}
	if v := db["displayRatio"]; v != float64(1.2) {
		t.Fatalf("DisplayBackup.displayRatio unexpected: %v", v)
	}
	if v := db["coverImageUrl"]; v != "http://img.test/cover_backup.png" {
		t.Fatalf("DisplayBackup.coverImageUrl unexpected: %v", v)
	}
	if v := db["coverImageUrlDark"]; v != "http://img.test/cover_backup_dark.png" {
		t.Fatalf("DisplayBackup.coverImageUrlDark unexpected: %v", v)
	}
}
