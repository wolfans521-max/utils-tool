package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard116 校验工厂方法是否正确设置 card_type
func TestNewCard116(t *testing.T) {
	c := NewCard116()
	if c == nil {
		t.Fatalf("NewCard116 returned nil")
	}
	if c.CardType != 116 {
		t.Fatalf("CardType expected 116, got %d", c.CardType)
	}
}

// TestC116JSONStructure 构造一个完整的 C116，校验关键字段的 JSON 结构
func TestC116JSONStructure(t *testing.T) {
	card := NewCard116()

	card.PageSize = 6
	card.Events = []*C116Event{
		{
			Id:     "312edsd8s9d8s9d",
			Title:  "小梅哭了",
			Scheme: "sinaweibo://detail?id=1",
			Color:  "#636363",
			ActionLog: map[string]any{
				"act_code": 100,
			},
		},
		{
			Id:     "312edsd8s9d8s9e",
			Title:  "龙猫是个大胖子",
			Scheme: "sinaweibo://detail?id=2",
			Color:  "#636363",
		},
		{
			Id:     "312edsd8s9d8s9f",
			Title:  "猫猫车跑的飞快",
			Scheme: "sinaweibo://detail?id=3",
			Color:  "#636363",
		},
		{
			Id:     "312edsd8s9d8s9g",
			Title:  "橡果子的故事",
			Scheme: "sinaweibo://detail?id=4",
			Color:  "#636363",
		},
	}
	card.ActionButton = &C116ActionButton{
		Text:   "换一换",
		Scheme: "sinaweibo://refresh",
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C116 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(116) {
		t.Fatalf("card_type expect 116, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["page_size"]; !ok || v != float64(6) {
		t.Fatalf("page_size expect 6, got %v (ok=%v)", v, ok)
	}

	// events 数组
	events, ok := m["events"].([]any)
	if !ok || len(events) != 4 {
		t.Fatalf("events expect array of 4, got %T len=%d", m["events"], len(events))
	}
	event0, ok := events[0].(map[string]any)
	if !ok {
		t.Fatalf("events[0] expect object, got %T", events[0])
	}
	if v := event0["id"]; v != "312edsd8s9d8s9d" {
		t.Fatalf("events[0].id unexpected: %v", v)
	}
	if v := event0["title"]; v != "小梅哭了" {
		t.Fatalf("events[0].title unexpected: %v", v)
	}
	if v := event0["color"]; v != "#636363" {
		t.Fatalf("events[0].color unexpected: %v", v)
	}
	if v := event0["actionlog"]; v == nil {
		t.Fatalf("events[0].actionlog should not be nil")
	}

	// action_button 字段
	actionButton, ok := m["action_button"].(map[string]any)
	if !ok {
		t.Fatalf("action_button expect object, got %T", m["action_button"])
	}
	if v := actionButton["text"]; v != "换一换" {
		t.Fatalf("action_button.text unexpected: %v", v)
	}
	if v := actionButton["scheme"]; v != "sinaweibo://refresh" {
		t.Fatalf("action_button.scheme unexpected: %v", v)
	}
}

// TestC116EmptyFields 测试空字段时 omitempty 是否生效
func TestC116EmptyFields(t *testing.T) {
	card := NewCard116()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C116 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(116) {
		t.Fatalf("card_type expect 116, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["events"]; ok {
		t.Fatalf("events should be omitted when nil")
	}
	if _, ok := m["action_button"]; ok {
		t.Fatalf("action_button should be omitted when nil")
	}
}

// TestC116WithoutActionButton 测试不下发换一换按钮
func TestC116WithoutActionButton(t *testing.T) {
	card := NewCard116()
	card.PageSize = 4
	card.Events = []*C116Event{
		{
			Id:    "test1",
			Title: "测试事件1",
			Color: "#333333",
		},
		{
			Id:    "test2",
			Title: "测试事件2",
			Color: "#333333",
		},
	}
	// 不设置 ActionButton

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C116 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["page_size"]; !ok || v != float64(4) {
		t.Fatalf("page_size expect 4, got %v (ok=%v)", v, ok)
	}

	events, ok := m["events"].([]any)
	if !ok || len(events) != 2 {
		t.Fatalf("events expect array of 2")
	}

	// action_button 应该被忽略
	if _, ok := m["action_button"]; ok {
		t.Fatalf("action_button should be omitted when nil")
	}
}

// TestC116PartialFields 测试部分字段填充
func TestC116PartialFields(t *testing.T) {
	card := NewCard116()
	card.PageSize = 3
	card.Events = []*C116Event{
		{
			Title: "测试标题",
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C116 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["page_size"]; !ok || v != float64(3) {
		t.Fatalf("page_size expect 3, got %v (ok=%v)", v, ok)
	}

	events, ok := m["events"].([]any)
	if !ok || len(events) != 1 {
		t.Fatalf("events expect array of 1")
	}
	event0, ok := events[0].(map[string]any)
	if !ok {
		t.Fatalf("events[0] expect object, got %T", events[0])
	}
	if v := event0["title"]; v != "测试标题" {
		t.Fatalf("events[0].title unexpected: %v", v)
	}
}
