package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard53 校验工厂方法是否正确设置 card_type
func TestNewCard53(t *testing.T) {
	c := NewCard53()
	if c == nil {
		t.Fatalf("NewCard53 returned nil")
	}
	if c.CardType != 53 {
		t.Fatalf("CardType expected 53, got %d", c.CardType)
	}
}

// TestC53JSONStructureReceived 测试已领取红包样式
func TestC53JSONStructureReceived(t *testing.T) {
	card := NewCard53()

	card.ItemId = "1087030002_417_417_495621"
	card.Scheme = "sinaweibo://userinfo?uid=1252397723"
	card.BgUrl = "http://example.com/hongbaofei_card_background.png"
	card.Desc1Color = "#ffffff"
	card.Desc1Size = 12
	card.Desc1ColorSkin = "#ffffff"
	card.Desc1 = "手气不错！已领取"
	card.Price = &C53Price{
		TextColor:     "#fed65a",
		TextSize:      12,
		TextColorSkin: "#fed65a",
		Text:          "100元",
		Number: &C53PriceNumber{
			Num:           "100",
			TextColor:     "#fed65a",
			TextColorSkin: "#fed65a",
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C53 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(53) {
		t.Fatalf("card_type expect 53, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["bg_url"]; !ok {
		t.Fatalf("bg_url should exist, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc1"]; !ok || v != "手气不错！已领取" {
		t.Fatalf("desc1 unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc1_color"]; !ok || v != "#ffffff" {
		t.Fatalf("desc1_color unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc1_size"]; !ok || v != float64(12) {
		t.Fatalf("desc1_size unexpected: %v (ok=%v)", v, ok)
	}

	// price 字段
	price, ok := m["price"].(map[string]any)
	if !ok {
		t.Fatalf("price expect object, got %T", m["price"])
	}
	if v := price["text"]; v != "100元" {
		t.Fatalf("price.text unexpected: %v", v)
	}
	if v := price["text_color"]; v != "#fed65a" {
		t.Fatalf("price.text_color unexpected: %v", v)
	}

	// price.number 字段
	number, ok := price["number"].(map[string]any)
	if !ok {
		t.Fatalf("price.number expect object, got %T", price["number"])
	}
	if v := number["num"]; v != "100" {
		t.Fatalf("price.number.num unexpected: %v", v)
	}
}

// TestC53JSONStructureUnreceived 测试未领取红包样式
func TestC53JSONStructureUnreceived(t *testing.T) {
	card := NewCard53()

	card.ItemId = "1087030002_417_417_495621"
	card.Scheme = "sinaweibo://userinfo?uid=1252397723"
	card.BgUrl = "http://example.com/hongbaofei_card_background.png"
	card.TxtTitle = "微博红包"
	card.Desc1 = "恭喜发财  大吉大利"
	card.TitleColor = "#fed65a"
	card.TitleSize = 16
	card.TitleColorSkin = "#fed65a"
	card.Desc1Color = "#ffffff"
	card.Desc1Size = 12
	card.Desc1ColorSkin = "#ffffff"
	card.ButtonText = "点击领取"
	card.Buttons = []*C53Button{
		{
			Type: "red_envelope",
			Name: "点击领取",
			Pic:  "http://example.com/hongbaofei_button_money.png",
			Params: &C53ButtonParams{
				Scheme:       "http://www.baidu.com",
				TxtBg:        "http://example.com/hongbaofei_card_button.png",
				TxtColor:     "#ff000000",
				TxtColorSkin: "#ff849684",
			},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C53 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(53) {
		t.Fatalf("card_type expect 53, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["txt_title"]; !ok || v != "微博红包" {
		t.Fatalf("txt_title unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_color"]; !ok || v != "#fed65a" {
		t.Fatalf("title_color unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_size"]; !ok || v != float64(16) {
		t.Fatalf("title_size unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["button_text"]; !ok || v != "点击领取" {
		t.Fatalf("button_text unexpected: %v (ok=%v)", v, ok)
	}

	// buttons 数组
	buttons, ok := m["buttons"].([]any)
	if !ok || len(buttons) == 0 {
		t.Fatalf("buttons expect non-empty array")
	}
	btn0, ok := buttons[0].(map[string]any)
	if !ok {
		t.Fatalf("buttons[0] expect object, got %T", buttons[0])
	}
	if v := btn0["type"]; v != "red_envelope" {
		t.Fatalf("buttons[0].type unexpected: %v", v)
	}
	if v := btn0["name"]; v != "点击领取" {
		t.Fatalf("buttons[0].name unexpected: %v", v)
	}

	// buttons[0].params 字段
	params, ok := btn0["params"].(map[string]any)
	if !ok {
		t.Fatalf("buttons[0].params expect object, got %T", btn0["params"])
	}
	if v := params["scheme"]; v != "http://www.baidu.com" {
		t.Fatalf("buttons[0].params.scheme unexpected: %v", v)
	}
	if v := params["txt_color"]; v != "#ff000000" {
		t.Fatalf("buttons[0].params.txt_color unexpected: %v", v)
	}
}

// TestC53EmptyFields 测试空字段时 omitempty 是否生效
func TestC53EmptyFields(t *testing.T) {
	card := NewCard53()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C53 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(53) {
		t.Fatalf("card_type expect 53, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["price"]; ok {
		t.Fatalf("price should be omitted when nil")
	}
	if _, ok := m["buttons"]; ok {
		t.Fatalf("buttons should be omitted when nil")
	}
	if _, ok := m["txt_title"]; ok {
		t.Fatalf("txt_title should be omitted when empty")
	}
}

// TestC53PartialFields 测试部分字段填充
func TestC53PartialFields(t *testing.T) {
	card := NewCard53()
	card.TxtTitle = "测试红包"
	card.Desc1 = "测试描述"
	card.User = &C53User{
		Id:         12345,
		ScreenName: "测试用户",
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C53 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["txt_title"]; !ok || v != "测试红包" {
		t.Fatalf("txt_title expect '测试红包', got %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc1"]; !ok || v != "测试描述" {
		t.Fatalf("desc1 expect '测试描述', got %v (ok=%v)", v, ok)
	}

	user, ok := m["user"].(map[string]any)
	if !ok {
		t.Fatalf("user expect object, got %T", m["user"])
	}
	if v := user["id"]; v != float64(12345) {
		t.Fatalf("user.id unexpected: %v", v)
	}
	if v := user["screen_name"]; v != "测试用户" {
		t.Fatalf("user.screen_name unexpected: %v", v)
	}
}
