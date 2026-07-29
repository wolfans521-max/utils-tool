package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard10 校验工厂方法是否正确设置 card_type
func TestNewCard10(t *testing.T) {
	c := NewCard10()
	if c == nil {
		t.Fatalf("NewCard10 returned nil")
	}
	if c.CardType != 10 {
		t.Fatalf("CardType expected 10, got %d", c.CardType)
	}
}

// TestC10JSONStructure 构造一个完整的 C10，校验关键字段的 JSON 结构
func TestC10JSONStructure(t *testing.T) {
	card := NewCard10()

	card.ItemId = "card10_test_item"
	card.Scheme = "sinaweibo://userinfo?uid=123456789"
	card.Title = "推荐用户"
	card.TitleExtraText = "查看更多"
	card.ShowTitleArrow = 1
	card.DisplayArrow = 1
	card.Desc1 = "知名科技博主"
	card.Desc2 = "粉丝 100万 · 微博 1000条"
	card.OpenUrl = "sinaweibo://userinfo?uid=123456789"

	card.User = &C10User{
		Id:              123456789,
		Idstr:           "123456789",
		ScreenName:      "科技达人",
		ProfileImageUrl: "http://example.com/avatar_small.jpg",
		AvatarLarge:     "http://example.com/avatar_large.jpg",
		AvatarHd:        "http://example.com/avatar_hd.jpg",
		Verified:        true,
		VerifiedType:    0,
		VerifiedReason:  "知名科技博主",
		Description:     "分享最新科技资讯",
	}

	card.Buttons = []*Button{
		{
			Type: "follow",
			Name: "关注",
			Pic:  "http://example.com/follow.png",
			Params: map[string]any{
				"uid": 123456789,
			},
			ActionLog: map[string]any{
				"act_code": 91,
			},
		},
	}

	card.Desc2Struct = []*C10Desc2Item{
		{
			Name:   "粉丝 100万",
			Scheme: "sinaweibo://fans?uid=123456789",
		},
		{
			Name:   "微博 1000条",
			Scheme: "sinaweibo://mblog?uid=123456789",
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C10 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(10) {
		t.Fatalf("card_type expect 10, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["title"]; !ok || v != "推荐用户" {
		t.Fatalf("title unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["title_extra_text"]; !ok || v != "查看更多" {
		t.Fatalf("title_extra_text unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["show_title_arrow"]; !ok || v != float64(1) {
		t.Fatalf("show_title_arrow unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc1"]; !ok || v != "知名科技博主" {
		t.Fatalf("desc1 unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc2"]; !ok || v != "粉丝 100万 · 微博 1000条" {
		t.Fatalf("desc2 unexpected: %v (ok=%v)", v, ok)
	}

	// user 字段
	user, ok := m["user"].(map[string]any)
	if !ok {
		t.Fatalf("user expect object, got %T", m["user"])
	}
	if v := user["id"]; v != float64(123456789) {
		t.Fatalf("user.id unexpected: %v", v)
	}
	if v := user["screen_name"]; v != "科技达人" {
		t.Fatalf("user.screen_name unexpected: %v", v)
	}
	if v := user["verified"]; v != true {
		t.Fatalf("user.verified unexpected: %v", v)
	}
	if v := user["verified_reason"]; v != "知名科技博主" {
		t.Fatalf("user.verified_reason unexpected: %v", v)
	}
	if v := user["description"]; v != "分享最新科技资讯" {
		t.Fatalf("user.description unexpected: %v", v)
	}

	// buttons 数组
	buttons, ok := m["buttons"].([]any)
	if !ok || len(buttons) != 1 {
		t.Fatalf("buttons expect array of 1, got %T len=%d", m["buttons"], len(buttons))
	}
	btn0, ok := buttons[0].(map[string]any)
	if !ok {
		t.Fatalf("buttons[0] expect object, got %T", buttons[0])
	}
	if v := btn0["type"]; v != "follow" {
		t.Fatalf("buttons[0].type unexpected: %v", v)
	}
	if v := btn0["name"]; v != "关注" {
		t.Fatalf("buttons[0].name unexpected: %v", v)
	}

	// desc2_struct 数组
	desc2Struct, ok := m["desc2_struct"].([]any)
	if !ok || len(desc2Struct) != 2 {
		t.Fatalf("desc2_struct expect array of 2, got %T len=%d", m["desc2_struct"], len(desc2Struct))
	}
	item0, ok := desc2Struct[0].(map[string]any)
	if !ok {
		t.Fatalf("desc2_struct[0] expect object, got %T", desc2Struct[0])
	}
	if v := item0["name"]; v != "粉丝 100万" {
		t.Fatalf("desc2_struct[0].name unexpected: %v", v)
	}
	if v := item0["scheme"]; v != "sinaweibo://fans?uid=123456789" {
		t.Fatalf("desc2_struct[0].scheme unexpected: %v", v)
	}
}

// TestC10EmptyFields 测试空字段时 omitempty 是否生效
func TestC10EmptyFields(t *testing.T) {
	card := NewCard10()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C10 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(10) {
		t.Fatalf("card_type expect 10, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["title"]; ok {
		t.Fatalf("title should be omitted when empty")
	}
	if _, ok := m["user"]; ok {
		t.Fatalf("user should be omitted when nil")
	}
	if _, ok := m["buttons"]; ok {
		t.Fatalf("buttons should be omitted when nil")
	}
	if _, ok := m["desc2_struct"]; ok {
		t.Fatalf("desc2_struct should be omitted when nil")
	}
}

// TestC10PartialFields 测试部分字段填充
func TestC10PartialFields(t *testing.T) {
	card := NewCard10()
	card.Title = "推荐关注"
	card.User = &C10User{
		Id:         987654321,
		ScreenName: "测试用户",
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C10 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["title"]; !ok || v != "推荐关注" {
		t.Fatalf("title expect '推荐关注', got %v (ok=%v)", v, ok)
	}

	user, ok := m["user"].(map[string]any)
	if !ok {
		t.Fatalf("user expect object, got %T", m["user"])
	}
	if v := user["id"]; v != float64(987654321) {
		t.Fatalf("user.id unexpected: %v", v)
	}
	if v := user["screen_name"]; v != "测试用户" {
		t.Fatalf("user.screen_name unexpected: %v", v)
	}

	// 验证不存在的字段
	if _, ok := m["buttons"]; ok {
		t.Fatalf("buttons should be omitted when nil")
	}
	if _, ok := m["desc2_struct"]; ok {
		t.Fatalf("desc2_struct should be omitted when nil")
	}
}

// TestC10UserStructure 测试用户结构完整性
func TestC10UserStructure(t *testing.T) {
	card := NewCard10()
	card.User = &C10User{
		Id:              111222333,
		Idstr:           "111222333",
		ScreenName:      "完整用户",
		ProfileImageUrl: "http://example.com/small.jpg",
		AvatarLarge:     "http://example.com/large.jpg",
		AvatarHd:        "http://example.com/hd.jpg",
		Verified:        true,
		VerifiedType:    1,
		VerifiedReason:  "企业认证",
		Description:     "这是用户简介",
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C10 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	user, ok := m["user"].(map[string]any)
	if !ok {
		t.Fatalf("user expect object")
	}

	// 验证所有用户字段
	if v := user["idstr"]; v != "111222333" {
		t.Fatalf("user.idstr unexpected: %v", v)
	}
	if v := user["profile_image_url"]; v != "http://example.com/small.jpg" {
		t.Fatalf("user.profile_image_url unexpected: %v", v)
	}
	if v := user["avatar_large"]; v != "http://example.com/large.jpg" {
		t.Fatalf("user.avatar_large unexpected: %v", v)
	}
	if v := user["avatar_hd"]; v != "http://example.com/hd.jpg" {
		t.Fatalf("user.avatar_hd unexpected: %v", v)
	}
	if v := user["verified_type"]; v != float64(1) {
		t.Fatalf("user.verified_type unexpected: %v", v)
	}
}

// TestC10Desc2Struct 测试来源文字配置
func TestC10Desc2Struct(t *testing.T) {
	card := NewCard10()
	card.Desc2Struct = []*C10Desc2Item{
		{Name: "链接1", Scheme: "scheme://link1"},
		{Name: "链接2", Scheme: "scheme://link2"},
		{Name: "链接3", Scheme: "scheme://link3"},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C10 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	desc2Struct, ok := m["desc2_struct"].([]any)
	if !ok || len(desc2Struct) != 3 {
		t.Fatalf("desc2_struct expect array of 3")
	}

	for i, item := range desc2Struct {
		itemMap, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("desc2_struct[%d] expect object", i)
		}
		expectedName := "链接" + string(rune('1'+i))
		if v := itemMap["name"]; v != expectedName {
			t.Fatalf("desc2_struct[%d].name expect %s, got %v", i, expectedName, v)
		}
	}
}
