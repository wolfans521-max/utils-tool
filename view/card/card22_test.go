package card

import (
	"encoding/json"
	"testing"
)

// TestC22JSONStructure 用示例数据构造一个 card22，并校验 JSON 结构和字段名是否符合约定
// 同时通过工厂方法 NewCard22() 创建，以符合标准使用方式
func TestC22JSONStructure(t *testing.T) {
	card := NewCard22()

	// 基础展示字段
	card.Width = 100
	card.Height = 100

	card.FlowGap = 3
	card.AutoFlow = 0
	card.BottomPadding = 1
	card.TopPadding = 1
	card.LeftRightPadding = 1

	card.PicHWScale = 1
	card.PicUnenableClick = false
	card.IsShowCornerRadius = 1
	card.CardAdStyle = 1

	// 轮播图配置
	card.PicItems = []*C22PicItems{
		{
			Pic:    "https://h5.sinaimg.cn/upload/2016/12/06/398/message_birthday3x.png",
			Scheme: "sinaweibo://cardlist?containerid=1003003",
			ActionLog: map[string]any{
				"act_code": 1,
			},
		},
	}

	// 按钮配置
	card.CommonButton = &Button{}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal c22 failed: %v", err)
	}

	// 反序列化到 map 方便断言 JSON key
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(22) {
		t.Fatalf("card_type expect 22, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["width"]; !ok || v != float64(100) {
		t.Fatalf("width expect 100, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["height"]; !ok || v != float64(100) {
		t.Fatalf("height expect 100, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["flow_gap"]; !ok || v != float64(3) {
		t.Fatalf("flow_gap expect 3, got %v (ok=%v)", v, ok)
	}
	// AutoFlow 为 0 时是字段零值，对应 json:"auto_flow,omitempty"，因此不会出现在 JSON 中
	if _, ok := m["auto_flow"]; ok {
		t.Fatalf("auto_flow should be omitted when zero value")
	}
	if v, ok := m["left_right_padding"]; !ok || v != float64(1) {
		t.Fatalf("left_right_padding expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["top_padding"]; !ok || v != float64(1) {
		t.Fatalf("top_padding expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["bottom_padding"]; !ok || v != float64(1) {
		t.Fatalf("bottom_padding expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["pic_h_w_scale"]; !ok || v != float64(1) {
		t.Fatalf("pic_h_w_scale expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["pic_unenable_click"]; !ok || v != false {
		t.Fatalf("pic_unenable_click expect false, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["is_show_corner_radius"]; !ok || v != float64(1) {
		t.Fatalf("is_show_corner_radius expect 1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["card_ad_style"]; !ok || v != float64(1) {
		t.Fatalf("card_ad_style expect 1, got %v (ok=%v)", v, ok)
	}

	// 嵌套的轮播图字段
	picItems, ok := m["pic_items"].([]any)
	if !ok || len(picItems) == 0 {
		t.Fatalf("pic_items expect non-empty array, got %T len=%d", m["pic_items"], len(picItems))
	}

	first, ok := picItems[0].(map[string]any)
	if !ok {
		t.Fatalf("pic_items[0] expect object, got %T", picItems[0])
	}
	if v := first["pic"]; v != "https://h5.sinaimg.cn/upload/2016/12/06/398/message_birthday3x.png" {
		t.Fatalf("pic_items[0].pic unexpected: %v", v)
	}
	if v := first["scheme"]; v != "sinaweibo://cardlist?containerid=1003003" {
		t.Fatalf("pic_items[0].scheme unexpected: %v", v)
	}
	if v := first["actionlog"]; v == nil {
		t.Fatalf("pic_items[0].actionlog should not be nil")
	}
}
