package api

import "testing"

func TestRequestBuilderAddAndLookup(t *testing.T) {
	rb := NewRequestBuilder()
	rb.Add("band", &BatchRequestItem{URLKey: "band_api"}).
		Add("window", &BatchRequestItem{URLKey: "window_api"}).
		Add("more_channel", &BatchRequestItem{URLKey: "more_channel_api"})

	if rb.Len() != 3 {
		t.Fatalf("Len = %d, want 3", rb.Len())
	}
	if !rb.Has("band") || !rb.Has("window") || !rb.Has("more_channel") {
		t.Errorf("Has miss, want all three registered: %+v", rb.idx)
	}
	if rb.Has("not_exist") {
		t.Error("Has(not_exist) = true, want false")
	}

	// IndexOf 反映注册顺序，与 Build 切片一致。
	if rb.IndexOf("band") != 0 {
		t.Errorf("IndexOf(band) = %d, want 0", rb.IndexOf("band"))
	}
	if rb.IndexOf("window") != 1 {
		t.Errorf("IndexOf(window) = %d, want 1", rb.IndexOf("window"))
	}
	if rb.IndexOf("more_channel") != 2 {
		t.Errorf("IndexOf(more_channel) = %d, want 2", rb.IndexOf("more_channel"))
	}
	if rb.IndexOf("not_exist") != -1 {
		t.Errorf("IndexOf(not_exist) = %d, want -1", rb.IndexOf("not_exist"))
	}
}

func TestRequestBuilderBuildOrder(t *testing.T) {
	rb := NewRequestBuilder()
	rb.Add("band", &BatchRequestItem{URLKey: "band_api"})
	rb.Add("window", &BatchRequestItem{URLKey: "window_api"})

	built := rb.Build()
	if len(built) != 2 {
		t.Fatalf("Build len = %d, want 2", len(built))
	}
	// Build 顺序与 Add 一致，IndexOf 下标可直接索引 results。
	if built[rb.IndexOf("band")].URLKey != "band_api" {
		t.Errorf("built[band] = %q, want band_api", built[rb.IndexOf("band")].URLKey)
	}
	if built[rb.IndexOf("window")].URLKey != "window_api" {
		t.Errorf("built[window] = %q, want window_api", built[rb.IndexOf("window")].URLKey)
	}
}

func TestRequestBuilderNilItemSkipped(t *testing.T) {
	// 对应 dao.OtherChannelFeedOption 返回 nil 的场景：nil 请求项必须被跳过。
	rb := NewRequestBuilder()
	rb.Add("band", &BatchRequestItem{URLKey: "band_api"})
	rb.Add("other_channel_feed", nil) // 应跳过，不登记 key
	rb.Add("window", &BatchRequestItem{URLKey: "window_api"})

	if rb.Len() != 2 {
		t.Fatalf("Len = %d, want 2 (nil item must be skipped)", rb.Len())
	}
	if rb.Has("other_channel_feed") {
		t.Error("Has(other_channel_feed) = true, want false for nil item")
	}
	// window 仍紧随 band，未被 nil 项顶出一个空位。
	if rb.IndexOf("window") != 1 {
		t.Errorf("IndexOf(window) = %d, want 1", rb.IndexOf("window"))
	}
}

func TestRequestBuilderDuplicateKeyOverwrites(t *testing.T) {
	rb := NewRequestBuilder()
	rb.Add("band", &BatchRequestItem{URLKey: "first"})
	rb.Add("band", &BatchRequestItem{URLKey: "second"}) // 同名覆盖

	if rb.Len() != 1 {
		t.Fatalf("Len = %d, want 1 (duplicate key must not append)", rb.Len())
	}
	if rb.Build()[0].URLKey != "second" {
		t.Errorf("duplicate Add should overwrite, got %q", rb.Build()[0].URLKey)
	}
}

func TestRequestBuilderEmpty(t *testing.T) {
	rb := NewRequestBuilder()
	if rb.Len() != 0 {
		t.Fatalf("Len = %d, want 0", rb.Len())
	}
	if rb.Has("anything") {
		t.Error("Has on empty builder = true, want false")
	}
	if rb.IndexOf("anything") != -1 {
		t.Error("IndexOf on empty builder != -1")
	}
	if built := rb.Build(); len(built) != 0 {
		t.Errorf("Build len = %d, want 0", len(built))
	}
}

// TestRequestBuilderResultsAccess 模拟 MRequest 返回后按名字取结果，
// 验证 IndexOf 与 results 切片的对齐关系（即取代 result[bandPos] 的写法）。
func TestRequestBuilderResultsAccess(t *testing.T) {
	rb := NewRequestBuilder()
	rb.Add("band", &BatchRequestItem{URLKey: "band_api"})
	rb.Add("window", &BatchRequestItem{URLKey: "window_api"})

	// 模拟 MRequest 按 Build 顺序返回结果。
	results := []map[string]any{
		{"url": "band_api"},
		{"url": "window_api"},
	}

	band := results[rb.IndexOf("band")]
	if band["url"] != "band_api" {
		t.Errorf("results[band] = %v, want band_api", band["url"])
	}
	window := results[rb.IndexOf("window")]
	if window["url"] != "window_api" {
		t.Errorf("results[window] = %v, want window_api", window["url"])
	}
}
