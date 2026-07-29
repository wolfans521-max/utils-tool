package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard57 校验工厂方法是否正确设置 card_type
func TestNewCard57(t *testing.T) {
	c := NewCard57()
	if c == nil {
		t.Fatalf("NewCard57 returned nil")
	}
	if c.CardType != 57 {
		t.Fatalf("CardType expected 57, got %d", c.CardType)
	}
}

// TestC57JSONStructure 构造一个完整的 C57，校验关键字段的 JSON 结构
func TestC57JSONStructure(t *testing.T) {
	card := NewCard57()

	card.Date = "Fri, 24 Jul 2015 11:35:48 +0800"
	card.DisplayArrow = 0
	card.OscillatePrice = "44.014"
	card.OscillateRate = "1.07%"
	card.Price = "4167.94"
	card.PriceColorType = 1
	card.StateDesc = "正常"
	card.Scheme = "sinaweibo://cardlist?containerid:230771_-_STOCKINDEX"
	card.StockInfo = &C57StockInfo{
		Scheme: "sinaweibo://cardlist?containerid:230771_-_HANGQING_SECOND_PAGE",
		StockPriceInfos: []*C57StockPriceInfo{
			{Price: "4.89", Desc: "今开"},
			{Price: "4.90", Desc: "昨收"},
			{Price: "4.32", Desc: "最低"},
			{Price: "5.36", Desc: "最高"},
			{Price: "5.32", Desc: "最牛"},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C57 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(57) {
		t.Fatalf("card_type expect 57, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["date"]; !ok || v != "Fri, 24 Jul 2015 11:35:48 +0800" {
		t.Fatalf("date unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["oscillate_price"]; !ok || v != "44.014" {
		t.Fatalf("oscillate_price unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["oscillate_rate"]; !ok || v != "1.07%" {
		t.Fatalf("oscillate_rate unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["price"]; !ok || v != "4167.94" {
		t.Fatalf("price unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["price_color_type"]; !ok || v != float64(1) {
		t.Fatalf("price_color_type unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["state_desc"]; !ok || v != "正常" {
		t.Fatalf("state_desc unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["scheme"]; !ok || v != "sinaweibo://cardlist?containerid:230771_-_STOCKINDEX" {
		t.Fatalf("scheme unexpected: %v (ok=%v)", v, ok)
	}

	// stock_info 字段
	stockInfo, ok := m["stock_info"].(map[string]any)
	if !ok {
		t.Fatalf("stock_info expect object, got %T", m["stock_info"])
	}
	if v := stockInfo["scheme"]; v != "sinaweibo://cardlist?containerid:230771_-_HANGQING_SECOND_PAGE" {
		t.Fatalf("stock_info.scheme unexpected: %v", v)
	}

	// stock_price_infos 数组
	priceInfos, ok := stockInfo["stock_price_infos"].([]any)
	if !ok || len(priceInfos) != 5 {
		t.Fatalf("stock_price_infos expect array of 5, got %T len=%d", stockInfo["stock_price_infos"], len(priceInfos))
	}

	// 验证第一个价格信息
	p0, ok := priceInfos[0].(map[string]any)
	if !ok {
		t.Fatalf("stock_price_infos[0] expect object, got %T", priceInfos[0])
	}
	if v := p0["price"]; v != "4.89" {
		t.Fatalf("stock_price_infos[0].price unexpected: %v", v)
	}
	if v := p0["desc"]; v != "今开" {
		t.Fatalf("stock_price_infos[0].desc unexpected: %v", v)
	}

	// 验证最后一个价格信息
	p4, ok := priceInfos[4].(map[string]any)
	if !ok {
		t.Fatalf("stock_price_infos[4] expect object, got %T", priceInfos[4])
	}
	if v := p4["price"]; v != "5.32" {
		t.Fatalf("stock_price_infos[4].price unexpected: %v", v)
	}
	if v := p4["desc"]; v != "最牛" {
		t.Fatalf("stock_price_infos[4].desc unexpected: %v", v)
	}
}

// TestC57EmptyFields 测试空字段时 omitempty 是否生效
func TestC57EmptyFields(t *testing.T) {
	card := NewCard57()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C57 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(57) {
		t.Fatalf("card_type expect 57, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["date"]; ok {
		t.Fatalf("date should be omitted when empty")
	}
	if _, ok := m["price"]; ok {
		t.Fatalf("price should be omitted when empty")
	}
	if _, ok := m["stock_info"]; ok {
		t.Fatalf("stock_info should be omitted when nil")
	}
}

// TestC57PartialFields 测试部分字段填充
func TestC57PartialFields(t *testing.T) {
	card := NewCard57()
	card.Price = "100.00"
	card.StateDesc = "交易中"
	card.StockInfo = &C57StockInfo{
		Scheme: "sinaweibo://stock",
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C57 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["price"]; !ok || v != "100.00" {
		t.Fatalf("price expect '100.00', got %v (ok=%v)", v, ok)
	}
	if v, ok := m["state_desc"]; !ok || v != "交易中" {
		t.Fatalf("state_desc expect '交易中', got %v (ok=%v)", v, ok)
	}

	stockInfo, ok := m["stock_info"].(map[string]any)
	if !ok {
		t.Fatalf("stock_info expect object, got %T", m["stock_info"])
	}
	if v := stockInfo["scheme"]; v != "sinaweibo://stock" {
		t.Fatalf("stock_info.scheme unexpected: %v", v)
	}

	// stock_price_infos 应该被忽略
	if _, ok := stockInfo["stock_price_infos"]; ok {
		t.Fatalf("stock_price_infos should be omitted when nil")
	}
}
