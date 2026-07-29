package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard58 校验工厂方法是否正确设置 card_type
func TestNewCard58(t *testing.T) {
	c := NewCard58()
	if c == nil {
		t.Fatalf("NewCard58 returned nil")
	}
	if c.CardType != 58 {
		t.Fatalf("CardType expected 58, got %d", c.CardType)
	}
}

// TestC58JSONStructureType1 测试 type=1 样式
func TestC58JSONStructureType1(t *testing.T) {
	card := NewCard58()

	card.Type = 1
	card.ItemId = "ihlkhl-9797"
	card.Name = "查看更多热门微博"
	card.BackgroundColor = "#EEEEEE"
	card.BackgroundColorDark = "#1E1E1E"
	card.Scheme = "sinaweibo://cardlist?containerid=107803_1592515472"
	card.NameSize = 14
	card.NameColor = "#636363"
	card.NameColorDark = "#636363"
	card.NameBold = true
	card.Height = 80

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C58 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(58) {
		t.Fatalf("card_type expect 58, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["type"]; !ok || v != float64(1) {
		t.Fatalf("type expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["name"]; !ok || v != "查看更多热门微博" {
		t.Fatalf("name unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["background_color"]; !ok || v != "#EEEEEE" {
		t.Fatalf("background_color unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["name_size"]; !ok || v != float64(14) {
		t.Fatalf("name_size unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["name_bold"]; !ok || v != true {
		t.Fatalf("name_bold unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["height"]; !ok || v != float64(80) {
		t.Fatalf("height unexpected: %v (ok=%v)", v, ok)
	}
}

// TestC58JSONStructureType2 测试 type=2 样式（带富文本）
func TestC58JSONStructureType2(t *testing.T) {
	card := NewCard58()

	card.Type = 2
	card.Name = "热门推荐"
	card.BackgroundColor = "#FFFFFF"
	card.BackgroundColorDark = "#1E1E1E"
	card.NameSize = 14
	card.NameColor = "#636363"
	card.NameColorDark = "#636363"
	card.Height = 80
	card.RichText = "去推荐看看"
	card.RichTextColor = "#507DAF"
	card.RichTextDarkColor = "#7691B9"
	card.RichTextScheme = "231643_11_5009"
	card.RichArrow = true
	card.RichArrowImage = "https://h5.sinaimg.cn/upload/100/1473/2020/12/16/normalarrow.png"
	card.RichArrowDarkImage = "https://h5.sinaimg.cn/upload/100/1473/2020/12/16/dardarrow.png"
	card.DottedLine = true

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C58 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["type"]; !ok || v != float64(2) {
		t.Fatalf("type expect 2, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["rich_text"]; !ok || v != "去推荐看看" {
		t.Fatalf("rich_text unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["rich_text_color"]; !ok || v != "#507DAF" {
		t.Fatalf("rich_text_color unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["rich_arrow"]; !ok || v != true {
		t.Fatalf("rich_arrow unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["dotted_line"]; !ok || v != true {
		t.Fatalf("dotted_line unexpected: %v (ok=%v)", v, ok)
	}
}

// TestC58JSONStructureType3 测试 type=3 样式（带按钮）
func TestC58JSONStructureType3(t *testing.T) {
	card := NewCard58()

	card.Type = 3
	card.Name = "筛选"
	card.NameColor = "#636363"
	card.NameColorDark = "#636363"
	card.Height = 60
	card.DottedLine = true
	card.ButtonBackgroundColor = "#FFFFFF"
	card.ButtonDarkBackgroundColor = "#242424"
	card.ButtonSelectedBackgroundColor = "#F7F7F7"
	card.ButtonSelectedDarkBackgroundColor = "#1E1E1E"
	card.ButtonSelectedTextColor = "#FF8200"
	card.ButtonSelectedDarkTextColor = "#EA8011"
	card.ButtonScheme = "231643_11_5009"

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C58 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["type"]; !ok || v != float64(3) {
		t.Fatalf("type expect 3, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["button_background_color"]; !ok || v != "#FFFFFF" {
		t.Fatalf("button_background_color unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["button_selected_text_color"]; !ok || v != "#FF8200" {
		t.Fatalf("button_selected_text_color unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["button_scheme"]; !ok || v != "231643_11_5009" {
		t.Fatalf("button_scheme unexpected: %v (ok=%v)", v, ok)
	}
}

// TestC58JSONStructureType4 测试 type=4 样式（带用户图标）
func TestC58JSONStructureType4(t *testing.T) {
	card := NewCard58()

	card.Type = 4
	card.NameColor = "#636363"
	card.NameColorDark = "#636363"
	card.Height = 50
	card.ButtonBackgroundColor = "#FFFFFF"
	card.ButtonDarkBackgroundColor = "#242424"
	card.ButtonSelectedBackgroundColor = "#F7F7F7"
	card.ButtonSelectedDarkBackgroundColor = "#1E1E1E"
	card.ButtonScheme = "231643_11_5009"
	card.UserIconImage = "https://h5.sinaimg.cn/upload/100/1473/2020/12/16/normalarrow.png"
	card.UserIconImageDark = "https://h5.sinaimg.cn/upload/100/1473/2020/12/16/normalarrow.png"

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C58 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["type"]; !ok || v != float64(4) {
		t.Fatalf("type expect 4, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["user_icon_image"]; !ok {
		t.Fatalf("user_icon_image should exist, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["user_icon_image_dark"]; !ok {
		t.Fatalf("user_icon_image_dark should exist, got %v (ok=%v)", v, ok)
	}
}

// TestC58EmptyFields 测试空字段时 omitempty 是否生效
func TestC58EmptyFields(t *testing.T) {
	card := NewCard58()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C58 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(58) {
		t.Fatalf("card_type expect 58, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["name"]; ok {
		t.Fatalf("name should be omitted when empty")
	}
	if _, ok := m["rich_text"]; ok {
		t.Fatalf("rich_text should be omitted when empty")
	}
}

// TestC58PartialFields 测试部分字段填充
func TestC58PartialFields(t *testing.T) {
	card := NewCard58()
	card.Type = 0
	card.Name = "测试间隔"
	card.BackgroundColor = "#FFFFFF"

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C58 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["name"]; !ok || v != "测试间隔" {
		t.Fatalf("name expect '测试间隔', got %v (ok=%v)", v, ok)
	}
	if v, ok := m["background_color"]; !ok || v != "#FFFFFF" {
		t.Fatalf("background_color expect '#FFFFFF', got %v (ok=%v)", v, ok)
	}
}
