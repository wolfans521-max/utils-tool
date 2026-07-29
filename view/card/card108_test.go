package card

import (
	"encoding/json"
	"testing"
)

// TestNewCard108 校验工厂方法是否正确设置 card_type
func TestNewCard108(t *testing.T) {
	c := NewCard108()
	if c == nil {
		t.Fatalf("NewCard108 returned nil")
	}
	if c.CardType != 108 {
		t.Fatalf("CardType expected 108, got %d", c.CardType)
	}
}

// TestC108JSONStructure 构造一个完整的 C108，校验关键字段的 JSON 结构
func TestC108JSONStructure(t *testing.T) {
	card := NewCard108()

	card.ItemId = "4311320938387675"
	card.Scheme = "sinaweibo://detail?mblogid=H4Kvzzc15"
	card.NeedFeedback = 1
	card.FeedBackExt = "mid:4331767399540293|pid:4311320938387675"
	card.Playlist = &C108Playlist{
		Id:           4311320938387675,
		Name:         "人民日报",
		VideoCount:   1147,
		IsSubscribed: 1,
		RecomCode:    "1108",
		FeedBackExt:  "mid:4331767399540293|pid:4311320938387675",
		Index:        1,
		Author: &C108Author{
			Id:              2803301701,
			ScreenName:      "人民日报",
			Name:            "人民日报",
			ProfileImageUrl: "http://tva1.sinaimg.cn/crop.0.3.1018.1018.50/a716fd45gw1ev7q2k8japj20sg0sg779.jpg",
			CoverImage:      "http://wx4.sinaimg.cn/crop.0.0.920.300/a716fd45ly1fpjoldh9kaj20pk08cal8.jpg",
			Following:       true,
			Verified:        true,
			VerifiedType:    3,
			AvatarLarge:     "http://tva1.sinaimg.cn/crop.0.3.1018.1018.180/a716fd45gw1ev7q2k8japj20sg0sg779.jpg",
			VerifiedReason:  "《人民日报》法人微博",
			FollowersCount:  85230322,
		},
		Statuses: []*C108Status{
			{
				CreatedAt:      "Wed Jan 23 23:28:11 +0800 2019",
				Id:             4331767399540293,
				Idstr:          "4331767399540293",
				Mid:            "4331767399540293",
				Text:           "【我不要什么玩具，只要你回来就好了】快过年了...",
				TextLength:     268,
				RepostsCount:   863,
				CommentsCount:  435,
				AttitudesCount: 3725,
				MblogId:        "HdkpJE1RH",
				Scheme:         "sinaweibo://detail?mblogid=HdkpJE1RH",
				ObjExt:         "107万次观看",
				User: &C108StatusUser{
					Id:              2803301701,
					Idstr:           "2803301701",
					ScreenName:      "人民日报",
					Name:            "人民日报",
					Verified:        true,
					VerifiedType:    3,
					VerifiedReason:  "《人民日报》法人微博",
					FollowersCount:  85230322,
				},
				PageInfo: &C108PageInfo{
					Type:       "5",
					PageId:     "230442410c63fc6085cb35438f0846a60d4b3b",
					ObjectType: "video",
					Oid:        2803301701,
					PageTitle:  "人民日报的秒拍视频",
					PagePic:    "http://imgaliyuncdn.miaopai.com/images/test.jpg",
					PageUrl:    "sinaweibo://infopage?containerid=230442410c63fc6085cb35438f0846a60d4b3b",
					ObjectId:   "2017607:410c63fc6085cb35438f0846a60d4b3b",
					AuthorId:   2803301701,
					MediaInfo: &C108MediaInfo{
						Name:              "人民日报的秒拍视频",
						Duration:          277,
						StreamUrl:         "http://gslb.miaopai.com/stream/test.mp4",
						PrefetchType:      1,
						PrefetchSize:      524288,
						AuthorName:        "人民日报",
						PlaylistId:        4311320938387675,
						IsPlaylist:        1,
						HasRecommendVideo: 1,
						OnlineUsers:       "107万次观看",
						OnlineUsersNumber: 1077243,
						Titles: []*C108Title{
							{Default: true, Title: "留守儿童想家齐"},
						},
						PlayCompletionActions: []*C108PlayCompletionAction{
							{
								Type:         "1",
								Icon:         "http://img.t.sinajs.cn/t6/style/images/face/feed_c_r.png",
								Text:         "重播",
								BtnCode:      1000,
								ShowPosition: 1,
							},
						},
						ExtraInfo: &C108ExtraInfo{
							Sceneid:          "videotab_playlist_mblog",
							PlaylistDesc:     "共1147个视频 · 《人民日报》法人微博",
							PlaylistItemDesc: "昨天 23:28·107万次观看",
						},
					},
				},
				MultiAttitude: []*C108Attitude{
					{Type: 1, Count: 3654},
					{Type: 2, Count: 45},
				},
			},
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C108 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层字段
	if v, ok := m["card_type"]; !ok || v != float64(108) {
		t.Fatalf("card_type expect 108, got %v (ok=%v)", v, ok)
	}
	if v, ok := m["need_feedback"]; !ok || v != float64(1) {
		t.Fatalf("need_feedback unexpected: %v (ok=%v)", v, ok)
	}
	if v, ok := m["feed_back_ext"]; !ok {
		t.Fatalf("feed_back_ext should exist, got %v (ok=%v)", v, ok)
	}

	// playlist 字段
	playlist, ok := m["playlist"].(map[string]any)
	if !ok {
		t.Fatalf("playlist expect object, got %T", m["playlist"])
	}
	if v := playlist["id"]; v != float64(4311320938387675) {
		t.Fatalf("playlist.id unexpected: %v", v)
	}
	if v := playlist["name"]; v != "人民日报" {
		t.Fatalf("playlist.name unexpected: %v", v)
	}
	if v := playlist["video_count"]; v != float64(1147) {
		t.Fatalf("playlist.video_count unexpected: %v", v)
	}
	if v := playlist["is_subscribed"]; v != float64(1) {
		t.Fatalf("playlist.is_subscribed unexpected: %v", v)
	}

	// playlist.author 字段
	author, ok := playlist["author"].(map[string]any)
	if !ok {
		t.Fatalf("playlist.author expect object, got %T", playlist["author"])
	}
	if v := author["screen_name"]; v != "人民日报" {
		t.Fatalf("playlist.author.screen_name unexpected: %v", v)
	}
	if v := author["verified"]; v != true {
		t.Fatalf("playlist.author.verified unexpected: %v", v)
	}
	if v := author["followers_count"]; v != float64(85230322) {
		t.Fatalf("playlist.author.followers_count unexpected: %v", v)
	}

	// playlist.statuses 数组
	statuses, ok := playlist["statuses"].([]any)
	if !ok || len(statuses) == 0 {
		t.Fatalf("playlist.statuses expect non-empty array")
	}
	status0, ok := statuses[0].(map[string]any)
	if !ok {
		t.Fatalf("statuses[0] expect object, got %T", statuses[0])
	}
	if v := status0["mid"]; v != "4331767399540293" {
		t.Fatalf("statuses[0].mid unexpected: %v", v)
	}
	if v := status0["obj_ext"]; v != "107万次观看" {
		t.Fatalf("statuses[0].obj_ext unexpected: %v", v)
	}

	// statuses[0].page_info 字段
	pageInfo, ok := status0["page_info"].(map[string]any)
	if !ok {
		t.Fatalf("statuses[0].page_info expect object, got %T", status0["page_info"])
	}
	if v := pageInfo["object_type"]; v != "video" {
		t.Fatalf("page_info.object_type unexpected: %v", v)
	}

	// page_info.media_info 字段
	mediaInfo, ok := pageInfo["media_info"].(map[string]any)
	if !ok {
		t.Fatalf("page_info.media_info expect object, got %T", pageInfo["media_info"])
	}
	if v := mediaInfo["duration"]; v != float64(277) {
		t.Fatalf("media_info.duration unexpected: %v", v)
	}
	if v := mediaInfo["online_users"]; v != "107万次观看" {
		t.Fatalf("media_info.online_users unexpected: %v", v)
	}
}

// TestC108EmptyFields 测试空字段时 omitempty 是否生效
func TestC108EmptyFields(t *testing.T) {
	card := NewCard108()

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C108 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 只应该有 card_type 字段
	if v, ok := m["card_type"]; !ok || v != float64(108) {
		t.Fatalf("card_type expect 108, got %v (ok=%v)", v, ok)
	}

	// 其他字段应该被 omitempty 忽略
	if _, ok := m["playlist"]; ok {
		t.Fatalf("playlist should be omitted when nil")
	}
}

// TestC108PartialFields 测试部分字段填充
func TestC108PartialFields(t *testing.T) {
	card := NewCard108()
	card.NeedFeedback = 1
	card.Playlist = &C108Playlist{
		Id:         12345,
		Name:       "测试专辑",
		VideoCount: 10,
		Author: &C108Author{
			ScreenName: "测试作者",
		},
	}

	b, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal C108 failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 验证存在的字段
	if v, ok := m["need_feedback"]; !ok || v != float64(1) {
		t.Fatalf("need_feedback expect 1, got %v (ok=%v)", v, ok)
	}

	playlist, ok := m["playlist"].(map[string]any)
	if !ok {
		t.Fatalf("playlist expect object, got %T", m["playlist"])
	}
	if v := playlist["name"]; v != "测试专辑" {
		t.Fatalf("playlist.name unexpected: %v", v)
	}
	if v := playlist["video_count"]; v != float64(10) {
		t.Fatalf("playlist.video_count unexpected: %v", v)
	}

	author, ok := playlist["author"].(map[string]any)
	if !ok {
		t.Fatalf("playlist.author expect object, got %T", playlist["author"])
	}
	if v := author["screen_name"]; v != "测试作者" {
		t.Fatalf("playlist.author.screen_name unexpected: %v", v)
	}
}
