package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard61 校验工厂方法是否正确设置 card_type
func TestNewCard61(t *testing.T) {
	c := NewCard61()
	if c == nil {
		t.Fatalf("NewCard61 returned nil")
	}
	if c.CardType != 61 {
		t.Fatalf("CardType expected 61, got %d", c.CardType)
	}
}

// TestC61JSONStructure 构造一个完整的 C61，校验关键字段的 JSON 结构
func TestC61JSONStructure(t *testing.T) {
	card := NewCard61()

	card.CardCommTitle = ""
	card.CardBackgroudcolor = "#FFFFFF"
	card.IsShowBorder = "1"
	card.CardId = "1028031201_1_45.4_1499104401_2280198017"
	card.ItemId = "1028031201_1_45.4_1499104401_2280198017"
	card.Desc1 = "解放军报、中国军网法人微博"
	card.Scheme = "sinaweibo://userinfo?uid=2280198017"
	card.ActionLog = map[string]any{
		"act_code": 1338,
		"ext":      "growth:1|terminal:|reason:45.4",
	}
	card.DeleteAction = &C61DeleteAction{
		CanDelete: 1,
		ActionLog: map[string]any{
			"act_code": 1338,
			"ext":      "growth:1|terminal:|reason:45.4|act:uninterested",
		},
	}
	card.Buttons = []*Button{
		{
			SkipFormat: 1,
			Type:       "follow",
			Name:       "关注",
			SubType:    0,
			Pic:        "http://u1.sinaimg.cn/upload/2013/07/01/userinfo_relationship_indicator_discuss.png",
			Params: map[string]any{
				"uid":           2280198017,
				"api_type":      "pagecardlist",
				"disable_group": 1,
			},
			ActionLog: map[string]any{
				"act_code": 1338,
			},
		},
	}
	card.User = &C61User{
		Id:              2280198017,
		ScreenName:      "军报记者",
		Description:     "解放军报、中国军网法人微博。 瞭望军事风云，关注国家安全。",
		Gender:          "m",
		ProfileImageUrl: "http://tva2.sinaimg.cn/crop.6.4.370.370.50/87e90f81jw8evofygvmgjj20aj0ajdh7.jpg",
		AvatarLarge:     "http://tva2.sinaimg.cn/crop.6.4.370.370.180/87e90f81jw8evofygvmgjj20aj0ajdh7.jpg",
		AvatarHd:        "http://tva2.sinaimg.cn/crop.6.4.370.370.1024/87e90f81jw8evofygvmgjj20aj0ajdh7.jpg",
		Verified:        true,
		VerifiedType:    3,
		FollowersCount:  15534777,
		FriendsCount:    800,
		VerifiedTypeExt: 0,
		VerifiedReason:  "解放军报、中国军网法人微博",
		Level:           2,
		Ptype:           0,
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C61 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(61) {
		t.Fatalf("card_type expect 61, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["card_backgroudcolor"]; !ok || v != "#FFFFFF" {
		t.Fatalf("card_backgroudcolor unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["is_show_border"]; !ok || v != "1" {
		t.Fatalf("is_show_border unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["desc1"]; !ok || v != "解放军报、中国军网法人微博" {
		t.Fatalf("desc1 unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["scheme"]; !ok || v != "sinaweibo://userinfo?uid=2280198017" {
		t.Fatalf("scheme unexpected: %v (ok=%v)", v, ok)
	}

	// delete_action 字段
	deleteAction, ok := m["delete_action"].(map[string]any)
	if !ok {
		t.Fatalf("delete_action expect object, got %T", m["delete_action"])
	}
	if v := deleteAction["can_delete"]; v != float64(1) {
		t.Fatalf("delete_action.can_delete unexpected: %v", v)
	}
	if v := deleteAction["actionlog"]; v == nil {
		t.Fatalf("delete_action.actionlog should not be nil")
	}

	// buttons 数组
	buttons, ok := m["buttons"].([]any)
	if !ok || len(buttons) == 0 {
		t.Fatalf("buttons expect non-empty array, got %T len=%d", m["buttons"], len(buttons))
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
	if v := btn0["params"]; v == nil {
		t.Fatalf("buttons[0].params should not be nil")
	}

	// user 字段
	user, ok := m["user"].(map[string]any)
	if !ok {
		t.Fatalf("user expect object, got %T", m["user"])
	}
	if v := user["id"]; v != float64(2280198017) {
		t.Fatalf("user.id unexpected: %v", v)
	}
	if v := user["screen_name"]; v != "军报记者" {
		t.Fatalf("user.screen_name unexpected: %v", v)
	}
	if v := user["verified"]; v != true {
		t.Fatalf("user.verified unexpected: %v", v)
	}
	if v := user["verified_type"]; v != float64(3) {
		t.Fatalf("user.verified_type unexpected: %v", v)
	}
	if v := user["followers_count"]; v != float64(15534777) {
		t.Fatalf("user.followers_count unexpected: %v", v)
	}
}

// TestC61EmptyFields 测试空字段时 omitempty 是否生效
func TestC61EmptyFields(t *testing.T) {
	card := NewCard61()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C61 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(61) {
		t.Fatalf("card_type expect 61, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["desc1"]; ok {
		t.Fatalf("desc1 should be omitted when empty")
	}
	if _, ok := m["user"]; ok {
		t.Fatalf("user should be omitted when nil")
	}
	if _, ok := m["buttons"]; ok {
		t.Fatalf("buttons should be omitted when nil")
	}
}

// TestC61PartialFields 测试部分字段填充
func TestC61PartialFields(t *testing.T) {
	card := NewCard61()
	card.Desc1 = "测试描述"
	card.User = &C61User{
		Id:         12345,
		ScreenName: "测试用户",
		Verified:   true,
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C61 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
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
	if v := user["verified"]; v != true {
		t.Fatalf("user.verified unexpected: %v", v)
	}

	// buttons 应该被忽略
	if _, ok := m["buttons"]; ok {
		t.Fatalf("buttons should be omitted when nil")
	}
}
