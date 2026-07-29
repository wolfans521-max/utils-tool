package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard122 校验工厂方法是否正确设置 card_type
func TestNewCard122(t *testing.T) {
	c := NewCard122()
	if c == nil {
		t.Fatalf("NewCard122 returned nil")
	}
	if c.CardType != 122 {
		t.Fatalf("CardType expected 122, got %d", c.CardType)
	}
}

// TestC122JSONStructure 构造一个完整的 C122，校验关键字段的 JSON 结构
func TestC122JSONStructure(t *testing.T) {
	card := NewCard122()

	card.VoteObject = &C122VoteObject{
		Id:               "123333",
		VoteType:         0,
		Content:          "参与投票，快来选出你心目中的冠军队伍",
		Parted:           1,
		PartInfo:         "100 人参与",
		State:            1,
		UserId:           "2323242",
		UserNick:         "NBA",
		ExpireDate:       32232323,
		ChoiceMoreScheme: "sinaweibo://vote/more",
		ChoiceMoreActionLog: map[string]any{
			"act_code": 100,
		},
		VoteList: []*C122VoteItem{
			{
				Id:        "8877",
				Content:   "中国",
				Selected:  1,
				PartNum:   "30人",
				PartRatio: 0.3,
				Pic:       "https://example.com/avatar.png",
			},
			{
				Id:        "8878",
				Content:   "日本",
				Selected:  0,
				PartNum:   "20人",
				PartRatio: 0.21,
				Pic:       "https://example.com/avatar2.png",
			},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C122 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(122) {
		t.Fatalf("card_type expect 122, got %v (ok=%v)", v, ok)
	}

	// vote_object 字段
	voteObject, ok := m["vote_object"].(map[string]any)
	if !ok {
		t.Fatalf("vote_object expect object, got %T", m["vote_object"])
	}
	if v := voteObject["id"]; v != "123333" {
		t.Fatalf("vote_object.id unexpected: %v", v)
		// vote_type 为 0 时会被 omitempty 忽略，所以检查不存在或为 0
		if v, ok := voteObject["vote_type"]; ok && v != float64(0) {
			t.Fatalf("vote_object.vote_type unexpected: %v", v)
		}
	}
	if v := voteObject["content"]; v != "参与投票，快来选出你心目中的冠军队伍" {
		t.Fatalf("vote_object.content unexpected: %v", v)
	}
	if v := voteObject["parted"]; v != float64(1) {
		t.Fatalf("vote_object.parted unexpected: %v", v)
	}
	if v := voteObject["part_info"]; v != "100 人参与" {
		t.Fatalf("vote_object.part_info unexpected: %v", v)
	}
	if v := voteObject["state"]; v != float64(1) {
		t.Fatalf("vote_object.state unexpected: %v", v)
	}
	if v := voteObject["user_nick"]; v != "NBA" {
		t.Fatalf("vote_object.user_nick unexpected: %v", v)
	}
	if v := voteObject["expire_date"]; v != float64(32232323) {
		t.Fatalf("vote_object.expire_date unexpected: %v", v)
	}
	if v := voteObject["choice_more_scheme"]; v != "sinaweibo://vote/more" {
		t.Fatalf("vote_object.choice_more_scheme unexpected: %v", v)
	}

	// vote_list 数组
	voteList, ok := voteObject["vote_list"].([]any)
	if !ok || len(voteList) != 2 {
		t.Fatalf("vote_list expect array of 2, got %T len=%d", voteObject["vote_list"], len(voteList))
	}
	vote0, ok := voteList[0].(map[string]any)
	if !ok {
		t.Fatalf("vote_list[0] expect object, got %T", voteList[0])
	}
	if v := vote0["content"]; v != "中国" {
		t.Fatalf("vote_list[0].content unexpected: %v", v)
	}
	if v := vote0["selected"]; v != float64(1) {
		t.Fatalf("vote_list[0].selected unexpected: %v", v)
	}
	if v := vote0["part_num"]; v != "30人" {
		t.Fatalf("vote_list[0].part_num unexpected: %v", v)
	}
	if v := vote0["part_ratio"]; v != float64(0.3) {
		t.Fatalf("vote_list[0].part_ratio unexpected: %v", v)
	}

	vote1, ok := voteList[1].(map[string]any)
	if !ok {
		t.Fatalf("vote_list[1] expect object, got %T", voteList[1])
	}
	if v := vote1["content"]; v != "日本" {
		t.Fatalf("vote_list[1].content unexpected: %v", v)
	}
	// selected 为 0 时会被 omitempty 忽略
	if v, ok := vote1["selected"]; ok && v != float64(0) {
		t.Fatalf("vote_list[1].selected unexpected: %v", v)
	}
	if v := vote1["part_ratio"]; v != float64(0.21) {
		t.Fatalf("vote_list[1].part_ratio unexpected: %v", v)
	}
}

// TestC122VoteTypeImage 测试图片+文字投票类型
func TestC122VoteTypeImage(t *testing.T) {
	card := NewCard122()

	card.VoteObject = &C122VoteObject{
		Id:       "456789",
		VoteType: 1, // 图片+文字
		Content:  "选择你喜欢的明星",
		State:    1,
		VoteList: []*C122VoteItem{
			{
				Id:        "1001",
				Content:   "明星A",
				Selected:  0,
				PartNum:   "50人",
				PartRatio: 0.5,
				Pic:       "https://example.com/star_a.png",
			},
			{
				Id:        "1002",
				Content:   "明星B",
				Selected:  0,
				PartNum:   "50人",
				PartRatio: 0.5,
				Pic:       "https://example.com/star_b.png",
			},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C122 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	voteObject, ok := m["vote_object"].(map[string]any)
	if !ok {
		t.Fatalf("vote_object expect object, got %T", m["vote_object"])
	}
	if v := voteObject["vote_type"]; v != float64(1) {
		t.Fatalf("vote_object.vote_type expect 1, got %v", v)
	}

	voteList, ok := voteObject["vote_list"].([]any)
	if !ok || len(voteList) != 2 {
		t.Fatalf("vote_list expect array of 2")
	}
	vote0, ok := voteList[0].(map[string]any)
	if !ok {
		t.Fatalf("vote_list[0] expect object")
	}
	if v := vote0["pic"]; v != "https://example.com/star_a.png" {
		t.Fatalf("vote_list[0].pic unexpected: %v", v)
	}
}

// TestC122EmptyFields 测试空字段时 omitempty 是否生效
func TestC122EmptyFields(t *testing.T) {
	card := NewCard122()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C122 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(122) {
		t.Fatalf("card_type expect 122, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["vote_object"]; ok {
		t.Fatalf("vote_object should be omitted when nil")
	}
}

// TestC122PartialFields 测试部分字段填充
func TestC122PartialFields(t *testing.T) {
	card := NewCard122()
	card.VoteObject = &C122VoteObject{
		Id:      "test123",
		Content: "测试投票",
		State:   0, // 已结束
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C122 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	voteObject, ok := m["vote_object"].(map[string]any)
	if !ok {
		t.Fatalf("vote_object expect object, got %T", m["vote_object"])
	}
	if v := voteObject["id"]; v != "test123" {
		t.Fatalf("vote_object.id unexpected: %v", v)
	}
	if v := voteObject["content"]; v != "测试投票" {
		t.Fatalf("vote_object.content unexpected: %v", v)
	}

	// vote_list 应该被忽略
	if _, ok := voteObject["vote_list"]; ok {
		t.Fatalf("vote_list should be omitted when nil")
	}
}
