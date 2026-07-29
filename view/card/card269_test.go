package card

import (
	"encoding/json"
	"testing"
)

// TestNew269 校验工厂方法是否正确设置 card_type
func TestNew269(t *testing.T) {
	c := NewCard269()
	if c == nil {
		t.Fatalf("New269 returned nil")
	}
	if c.CardType != 269 {
		t.Fatalf("CardType expected 269, got %d", c.CardType)
	}
}

// TestC269JSONStructure 构造一个 C269，校验嵌套 display / displayBackup 以及 wbox 字段的 JSON 结构
func TestC269JSONStructure(t *testing.T) {
	card := &C269{}
	card.BaseItem = BaseItem{Base: Base{CardType: 269}}

	card.Display.MaxHeight = 400
	card.Display.CoverImageUrl = "http://img.test/cover269.png"
	card.Display.CoverImageUrlDark = "http://img.test/cover269_dark.png"
	card.Display.DisplayHeight = 250
	card.Display.DisplayRatio = 1.8
	card.Display.MeasureHeight = &MeasureHeight{
		Expression:        "min(250, screenHeight*0.6)",
		BundleVersionCode: "2000033",
	}

	card.DisplayBackup.DisplayHeight = 180
	card.DisplayBackup.DisplayRatio = 1.4
	card.DisplayBackup.CoverImageUrl = "http://img.test/cover269_backup.png"
	card.DisplayBackup.CoverImageUrlDark = "http://img.test/cover269_backup_dark.png"

	card.Scheme = "sinaweibo://cardlist?containerid=1006001"
	card.WboxScheme = "sinaweibo://wbox?app=bar"
	card.WboxMinRuntimeVersion = 11
	card.WboxMinSDKVersion = 21
	card.WboxParam = map[string]any{
		"foo": "bar269",
		"num": 2,
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C269 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(269) {
		t.Fatalf("card_type expect 269, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["scheme"]; !ok || v != "sinaweibo://cardlist?containerid=1006001" {
		t.Fatalf("scheme expect sinaweibo://cardlist?containerid=1006001, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["wboxScheme"]; !ok || v != "sinaweibo://wbox?app=bar" {
		t.Fatalf("wboxScheme expect sinaweibo://wbox?app=bar, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["wboxMinRuntimeVersion"]; !ok || v != float64(11) {
		t.Fatalf("wboxMinRuntimeVersion expect 11, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["wboxMinSDKVersion"]; !ok || v != float64(21) {
		t.Fatalf("wboxMinSDKVersion expect 21, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["wboxParam"]; !ok || v == nil {
		t.Fatalf("wboxParam expect non-nil, got %v (ok=%v)", v, ok)
	}

	// display 嵌套字段
	display, ok := m["display"].(map[string]any)
	if !ok {
		t.Fatalf("display expect object, got %T", m["display"])
	}
	if v := display["maxHeight"]; v != float64(400) {
		t.Fatalf("display.maxHeight unexpected: %v", v)
	}
	if v := display["coverImageUrl"]; v != "http://img.test/cover269.png" {
		t.Fatalf("display.coverImageUrl unexpected: %v", v)
	}
	if v := display["coverImageUrlDark"]; v != "http://img.test/cover269_dark.png" {
		t.Fatalf("display.coverImageUrlDark unexpected: %v", v)
	}
	if v := display["displayHeight"]; v != float64(250) {
		t.Fatalf("display.displayHeight unexpected: %v", v)
	}
	if v := display["displayRatio"]; v != float64(1.8) {
		t.Fatalf("display.displayRatio unexpected: %v", v)
	}

	mh, ok := display["measureHeight"].(map[string]any)
	if !ok {
		t.Fatalf("display.measureHeight expect object, got %T", display["measureHeight"])
	}
	if v := mh["expression"]; v != "min(250, screenHeight*0.6)" {
		t.Fatalf("measureHeight.expression unexpected: %v", v)
	}
	if v := mh["bundleVersionCode"]; v != "2000033" {
		t.Fatalf("measureHeight.bundleVersionCode unexpected: %v", v)
	}

	// displayBackup 嵌套字段
	db, ok := m["DisplayBackup"].(map[string]any)
	if !ok {
		t.Fatalf("DisplayBackup expect object, got %T", m["DisplayBackup"])
	}
	if v := db["displayHeight"]; v != float64(180) {
		t.Fatalf("DisplayBackup.displayHeight unexpected: %v", v)
	}
	if v := db["displayRatio"]; v != float64(1.4) {
		t.Fatalf("DisplayBackup.displayRatio unexpected: %v", v)
	}
	if v := db["coverImageUrl"]; v != "http://img.test/cover269_backup.png" {
		t.Fatalf("DisplayBackup.coverImageUrl unexpected: %v", v)
	}
	if v := db["coverImageUrlDark"]; v != "http://img.test/cover269_backup_dark.png" {
		t.Fatalf("DisplayBackup.coverImageUrlDark unexpected: %v", v)
	}
}
