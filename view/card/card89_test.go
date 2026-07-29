package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard89 校验工厂方法是否正确设置 card_type
func TestNewCard89(t *testing.T) {
	c := NewCard89()
	if c == nil {
		t.Fatalf("NewCard89 returned nil")
	}
	if c.CardType != 89 {
		t.Fatalf("CardType expected 89, got %d", c.CardType)
	}
}

// TestC89JSONStructure 构造一个完整的 C89，校验关键字段的 JSON 结构
func TestC89JSONStructure(t *testing.T) {
	card := NewCard89()

	card.ItemId = "2311590001_"
	card.Scheme = "sinaweibo://detail/?mblogid=FikhxjeIr"
	card.HideBottom = 0
	card.HotDiscussion = "https://h5.sinaimg.cn/upload/1019/115/2018/12/06/reyibiaoqian.png"
	card.NeedFeedback = 1
	card.HeaderTag = &C89HeaderTag{
		Icon:  "https://h5.sinaimg.cn/upload/1010/306/2018/11/13/469BA32C3868A0CDD39BFE81E1FA5F10.png",
		Title: "根据你关注的博主推荐",
	}
	card.Mblog = &C89Mblog{
		User: &C89User{
			Id:             6037687121,
			Idstr:          "6037687121",
			ScreenName:     "Bigger研究所",
			AvatarLarge:    "https://tvax1.sinaimg.cn/crop.90.84.709.709.180/006ABwelly1fffxpcyei3j30p00p0n2k.jpg",
			VerifiedReason: "微博原创视频博主",
		},
		PageInfo: &C89PageInfo{
			Type:       "11",
			PageId:     "230444cd033b50e72cae18908b2a1ece80a79a",
			ObjectType: "video",
			ObjectId:   "1034:cd033b50e72cae18908b2a1ece80a79a",
			Content1:   "篮球教学论坛的秒拍视频",
			Content2:   "美国Ardsley高中球员后场抢断后命中超远绝杀！",
			ActStatus:  1,
			PagePic:    "https://wx3.sinaimg.cn/large/8c89acbfly1fp13ohjzl2j20dc0dc0th.jpg",
			PageTitle:  "篮球教学论坛的秒拍视频",
			PageUrl:    "sinaweibo://infopage?containerid=230444cd033b50e72cae18908b2a1ece80a79a",
			Oid:        "2357832895",
			AuthorId:   "2357832895",
			MediaInfo: &C89MediaInfo{
				Name:              "篮球教学论坛的秒拍视频",
				Autoplay:          0,
				Duration:          44,
				StreamUrl:         "http://f.us.sinaimg.cn/003qLXMDlx07iDv8MGC401040200a0uh0k010.mp4",
				StreamUrlHd:       "http://f.us.sinaimg.cn/003wKDCglx07iDv8KQwg01040200f33M0k010.mp4",
				Mp4SdUrl:          "http://f.us.sinaimg.cn/003qLXMDlx07iDv8MGC401040200a0uh0k010.mp4",
				Mp4HdUrl:          "http://f.us.sinaimg.cn/003wKDCglx07iDv8KQwg01040200f33M0k010.mp4",
				PrefetchType:      1,
				PrefetchSize:      193019,
				ActStatus:         1,
				AuthorMid:         "4213959676281407",
				AuthorName:        "篮球教学论坛",
				HasRecommendVideo: 1,
				OnlineUsers:       "5万次观看",
				OnlineUsersNumber: 53755,
				Ttl:               3600,
				StorageType:       "unistore",
				PlayCompletionActions: []*C89PlayCompletionAction{
					{
						Type:         "1",
						Icon:         "http://img.t.sinajs.cn/t6/style/images/face/feed_c_r.png",
						Text:         "重播",
						BtnCode:      1000,
						ShowPosition: 1,
						ActionLog: map[string]any{
							"act_code": 1221,
						},
					},
				},
				Titles: []*C89Title{
					{Title: "美国高中生绝杀", Default: true},
				},
				VideoTags: []*C89VideoTag{
					{
						Name:   "盛唐幻夜",
						Scheme: "sinaweibo://videotimeline?kid=7052",
						ActionLog: map[string]any{
							"act_code": 3089,
						},
					},
				},
				ScrubberPreview: &C89ScrubberPreview{
					Previews: []*C89Preview{
						{
							Width:    80,
							Height:   45,
							Row:      10,
							Col:      10,
							Interval: 1,
							Label:    "ld",
							Urls:     []string{"http://wx4.sinaimg.cn/large/8c89acbfly1fp13pnytmvj20m80cimy7.jpg"},
						},
					},
				},
			},
			PicInfo: &C89PicInfo{
				PicBig: &C89Pic{
					Height: "480",
					Width:  "480",
					Url:    "https://wx3.sinaimg.cn/large/8c89acbfly1fp13ohjzl2j20dc0dc0th.jpg",
				},
			},
			ActionLog: map[string]any{
				"act_code": 799,
			},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C89 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(89) {
		t.Fatalf("card_type expect 89, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["itemid"]; !ok || v != "2311590001_" {
		t.Fatalf("itemid unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["scheme"]; !ok || v != "sinaweibo://detail/?mblogid=FikhxjeIr" {
		t.Fatalf("scheme unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["hot_discussion"]; !ok {
		t.Fatalf("hot_discussion should exist, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["need_feedback"]; !ok || v != float64(1) {
		t.Fatalf("need_feedback unexpected: %v (ok=%v)", v, ok)
	}

	// header_tag 字段
	headerTag, ok := m["header_tag"].(map[string]any)
	if !ok {
		t.Fatalf("header_tag expect object, got %T", m["header_tag"])
	}
	if v := headerTag["title"]; v != "根据你关注的博主推荐" {
		t.Fatalf("header_tag.title unexpected: %v", v)
	}

	// mblog 字段
	mblog, ok := m["mblog"].(map[string]any)
	if !ok {
		t.Fatalf("mblog expect object, got %T", m["mblog"])
	}

	// mblog.user 字段
	user, ok := mblog["user"].(map[string]any)
	if !ok {
		t.Fatalf("mblog.user expect object, got %T", mblog["user"])
	}
	if v := user["id"]; v != float64(6037687121) {
		t.Fatalf("mblog.user.id unexpected: %v", v)
	}
	if v := user["screen_name"]; v != "Bigger研究所" {
		t.Fatalf("mblog.user.screen_name unexpected: %v", v)
	}

	// mblog.page_info 字段
	pageInfo, ok := mblog["page_info"].(map[string]any)
	if !ok {
		t.Fatalf("mblog.page_info expect object, got %T", mblog["page_info"])
	}
	if v := pageInfo["type"]; v != "11" {
		t.Fatalf("mblog.page_info.type unexpected: %v", v)
	}
	if v := pageInfo["object_type"]; v != "video" {
		t.Fatalf("mblog.page_info.object_type unexpected: %v", v)
	}

	// mblog.page_info.media_info 字段
	mediaInfo, ok := pageInfo["media_info"].(map[string]any)
	if !ok {
		t.Fatalf("mblog.page_info.media_info expect object, got %T", pageInfo["media_info"])
	}
	if v := mediaInfo["duration"]; v != float64(44) {
		t.Fatalf("media_info.duration unexpected: %v", v)
	}
	if v := mediaInfo["online_users"]; v != "5万次观看" {
		t.Fatalf("media_info.online_users unexpected: %v", v)
	}

	// play_completion_actions 数组
	actions, ok := mediaInfo["play_completion_actions"].([]any)
	if !ok || len(actions) == 0 {
		t.Fatalf("play_completion_actions expect non-empty array")
	}

	// titles 数组
	titles, ok := mediaInfo["titles"].([]any)
	if !ok || len(titles) == 0 {
		t.Fatalf("titles expect non-empty array")
	}

	// video_tags 数组
	videoTags, ok := mediaInfo["video_tags"].([]any)
	if !ok || len(videoTags) == 0 {
		t.Fatalf("video_tags expect non-empty array")
	}
}

// TestC89EmptyFields 测试空字段时 omitempty 是否生效
func TestC89EmptyFields(t *testing.T) {
	card := NewCard89()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C89 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(89) {
		t.Fatalf("card_type expect 89, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["mblog"]; ok {
		t.Fatalf("mblog should be omitted when nil")
	}
	if _, ok := m["header_tag"]; ok {
		t.Fatalf("header_tag should be omitted when nil")
	}
}

// TestC89PartialFields 测试部分字段填充
func TestC89PartialFields(t *testing.T) {
	card := NewCard89()
	card.HideBottom = 1
	card.HeaderTag = &C89HeaderTag{
		Title: "推荐视频",
	}
	card.Mblog = &C89Mblog{
		User: &C89User{
			ScreenName: "测试用户",
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C89 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["hide_bottom"]; !ok || v != float64(1) {
		t.Fatalf("hide_bottom expect 1, got %v (ok=%v)", v, ok)
	}

	headerTag, ok := m["header_tag"].(map[string]any)
	if !ok {
		t.Fatalf("header_tag expect object, got %T", m["header_tag"])
	}
	if v := headerTag["title"]; v != "推荐视频" {
		t.Fatalf("header_tag.title unexpected: %v", v)
	}

	mblog, ok := m["mblog"].(map[string]any)
	if !ok {
		t.Fatalf("mblog expect object, got %T", m["mblog"])
	}
	user, ok := mblog["user"].(map[string]any)
	if !ok {
		t.Fatalf("mblog.user expect object, got %T", mblog["user"])
	}
	if v := user["screen_name"]; v != "测试用户" {
		t.Fatalf("mblog.user.screen_name unexpected: %v", v)
	}
}
