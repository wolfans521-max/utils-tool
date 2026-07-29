package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard3 校验工厂方法是否正确设置 card_type
func TestNewCard3(t *testing.T) {
	c := NewCard3()
	if c == nil {
		t.Fatalf("NewCard3 returned nil")
	}
	if c.CardType != 3 {
		t.Fatalf("CardType expected 3, got %d", c.CardType)
	}
}

// TestC3JSONStructure 构造一个完整的 C3，校验关键字段的 JSON 结构
func TestC3JSONStructure(t *testing.T) {
	card := NewCard3()

	card.ItemId = "card3_test_item"
	card.Scheme = "sinaweibo://topic?name=测试话题"
	card.Title = "热门话题"
	card.TitleExtraText = "查看更多"
	card.ShowTitleArrow = 1
	card.ShowAvatar = 1
	card.RoundedCorner = 1
	card.DisplayArrow = 1
	card.MaxItemCount = 4
	card.PicSum = "+99"
	card.ShowLayer = 1
	card.FlagPic = "http://example.com/flag.png"
	card.OpenUrl = "sinaweibo://topic?name=测试"

	card.Pics = []*C3Pic{
		{
			Pic:   "http://example.com/pic1.jpg",
			Desc1: "图片标题1",
			Desc2: "图片描述1",
			ActionLog: map[string]any{
				"act_code": 100,
			},
		},
		{
			Pic:   "http://example.com/pic2.jpg",
			Desc1: "图片标题2",
			Desc2: "图片描述2",
		},
	}

	card.Users = []*C3User{
		{
			Id:    123456789,
			Idstr: "123456789",
		},
		{
			Id:    987654321,
			Idstr: "987654321",
		},
	}

	card.Elements = []*C3Element{
		{
			Uid:        "123456789",
			Desc1:      "用户名1",
			Desc2:      "用户描述1",
			Scheme:     "sinaweibo://userinfo?uid=123456789",
			ShowBorder: true,
			Highlight: &C3Highlight{
				DescFont: 14,
				DescEm:   [][]int{{0, 3}},
			},
			ActionLog: map[string]any{
				"act_code": 200,
			},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C3 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(3) {
		t.Fatalf("card_type expect 3, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title"]; !ok || v != "热门话题" {
		t.Fatalf("title unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_extra_text"]; !ok || v != "查看更多" {
		t.Fatalf("title_extra_text unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["show_title_arrow"]; !ok || v != float64(1) {
		t.Fatalf("show_title_arrow unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["max_item_count"]; !ok || v != float64(4) {
		t.Fatalf("max_item_count unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["pic_sum"]; !ok || v != "+99" {
		t.Fatalf("pic_sum unexpected: %v (ok=%v)", v, ok)
	}

	// pics 数组
	pics, ok := m["pics"].([]any)
	if !ok || len(pics) != 2 {
		t.Fatalf("pics expect array of 2, got %T len=%d", m["pics"], len(pics))
	}
	pic0, ok := pics[0].(map[string]any)
	if !ok {
		t.Fatalf("pics[0] expect object, got %T", pics[0])
	}
	if v := pic0["pic"]; v != "http://example.com/pic1.jpg" {
		t.Fatalf("pics[0].pic unexpected: %v", v)
	}
	if v := pic0["desc1"]; v != "图片标题1" {
		t.Fatalf("pics[0].desc1 unexpected: %v", v)
	}

	// users 数组
	users, ok := m["users"].([]any)
	if !ok || len(users) != 2 {
		t.Fatalf("users expect array of 2, got %T len=%d", m["users"], len(users))
	}
	user0, ok := users[0].(map[string]any)
	if !ok {
		t.Fatalf("users[0] expect object, got %T", users[0])
	}
	if v := user0["id"]; v != float64(123456789) {
		t.Fatalf("users[0].id unexpected: %v", v)
	}

	// elements 数组
	elements, ok := m["elements"].([]any)
	if !ok || len(elements) != 1 {
		t.Fatalf("elements expect array of 1, got %T len=%d", m["elements"], len(elements))
	}
	elem0, ok := elements[0].(map[string]any)
	if !ok {
		t.Fatalf("elements[0] expect object, got %T", elements[0])
	}
	if v := elem0["uid"]; v != "123456789" {
		t.Fatalf("elements[0].uid unexpected: %v", v)
	}
	if v := elem0["desc1"]; v != "用户名1" {
		t.Fatalf("elements[0].desc1 unexpected: %v", v)
	}
	if v := elem0["show_border"]; v != true {
		t.Fatalf("elements[0].show_border unexpected: %v", v)
	}

	// highlight 嵌套结构
	highlight, ok := elem0["highlight"].(map[string]any)
	if !ok {
		t.Fatalf("elements[0].highlight expect object, got %T", elem0["highlight"])
	}
	if v := highlight["desc_font"]; v != float64(14) {
		t.Fatalf("highlight.desc_font unexpected: %v", v)
	}
}

// TestC3EmptyFields 测试空字段时 omitempty 是否生效
func TestC3EmptyFields(t *testing.T) {
	card := NewCard3()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C3 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(3) {
		t.Fatalf("card_type expect 3, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["title"]; ok {
		t.Fatalf("title should be omitted when empty")
	}
	if _, ok := m["pics"]; ok {
		t.Fatalf("pics should be omitted when nil")
	}
	if _, ok := m["users"]; ok {
		t.Fatalf("users should be omitted when nil")
	}
	if _, ok := m["elements"]; ok {
		t.Fatalf("elements should be omitted when nil")
	}
}

// TestC3PartialFields 测试部分字段填充
func TestC3PartialFields(t *testing.T) {
	card := NewCard3()
	card.Title = "测试标题"
	card.Pics = []*C3Pic{
		{
			Pic:   "http://example.com/test.jpg",
			Desc1: "测试图片",
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C3 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["title"]; !ok || v != "测试标题" {
		t.Fatalf("title expect '测试标题', got %v (ok=%v)", v, ok)
	}

	pics, ok := m["pics"].([]any)
	if !ok || len(pics) != 1 {
		t.Fatalf("pics expect array of 1")
	}

	// 验证不存在的字段
	if _, ok := m["users"]; ok {
		t.Fatalf("users should be omitted when nil")
	}
}
