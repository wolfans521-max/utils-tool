package utils

import (
	"testing"
)

// 构建测试数据，模拟 JSON 反序列化后的 map[string]any 结构
// 注意：JSON 反序列化后数值类型统一为 float64
func buildTestMap() map[string]any {
	return map[string]any{
		"code": float64(1000),
		"data": map[string]any{
			"bigday_ad_id":     "bigday_123",
			"user_refresh_num": float64(23),
			"ad_vip_id":        "393",
			"ad_refresh_num":   float64(2),
			"is_has_ad":        true,
			"ad_valid":         false,
		},
	}
}

func TestMapValue_Float64(t *testing.T) {
	data := buildTestMap()

	// 直接取顶层数值字段（JSON 反序列化后为 float64）
	got := MapValue[float64](data, "code")
	if got != 1000 {
		t.Errorf("MapValue[float64](data, 'code') = %v, want 1000", got)
	}

	// 嵌套数值字段
	got = MapValue[float64](data, "data.user_refresh_num")
	if got != 23 {
		t.Errorf("MapValue[float64](data, 'data.user_refresh_num') = %v, want 23", got)
	}

	got = MapValue[float64](data, "data.ad_refresh_num")
	if got != 2 {
		t.Errorf("MapValue[float64](data, 'data.ad_refresh_num') = %v, want 2", got)
	}
}

func TestMapValue_String(t *testing.T) {
	data := buildTestMap()

	// 嵌套字符串字段
	got := MapValue[string](data, "data.bigday_ad_id")
	if got != "bigday_123" {
		t.Errorf("MapValue[string](data, 'data.bigday_ad_id') = %v, want bigday_123", got)
	}

	got = MapValue[string](data, "data.ad_vip_id")
	if got != "393" {
		t.Errorf("MapValue[string](data, 'data.ad_vip_id') = %v, want 393", got)
	}
}

func TestMapValue_Bool(t *testing.T) {
	data := buildTestMap()

	got := MapValue[bool](data, "data.is_has_ad")
	if got != true {
		t.Errorf("MapValue[bool](data, 'data.is_has_ad') = %v, want true", got)
	}

	got = MapValue[bool](data, "data.ad_valid")
	if got != false {
		t.Errorf("MapValue[bool](data, 'data.ad_valid') = %v, want false", got)
	}
}

func TestMapValue_Map(t *testing.T) {
	data := buildTestMap()

	// 取嵌套的 map[string]any
	got := MapValue[map[string]any](data, "data")
	if got == nil {
		t.Fatal("MapValue[map[string]any](data, 'data') = nil, want non-nil")
	}
	if got["bigday_ad_id"] != "bigday_123" {
		t.Errorf("data.bigday_ad_id = %v, want bigday_123", got["bigday_ad_id"])
	}
}

func TestMapValue_TypeMismatch(t *testing.T) {
	data := buildTestMap()

	// code 是 float64，用 string 取应返回零值
	got := MapValue[string](data, "code")
	if got != "" {
		t.Errorf("MapValue[string](data, 'code') = %q, want empty string (type mismatch)", got)
	}

	// data.bigday_ad_id 是 string，用 float64 取应返回零值
	gotFloat := MapValue[float64](data, "data.bigday_ad_id")
	if gotFloat != 0 {
		t.Errorf("MapValue[float64](data, 'data.bigday_ad_id') = %v, want 0 (type mismatch)", gotFloat)
	}

	// data.is_has_ad 是 bool，用 string 取应返回零值
	gotStr := MapValue[string](data, "data.is_has_ad")
	if gotStr != "" {
		t.Errorf("MapValue[string](data, 'data.is_has_ad') = %q, want empty string (type mismatch)", gotStr)
	}
}

func TestMapValue_NotFound(t *testing.T) {
	data := buildTestMap()

	// 不存在的 key，返回零值
	got := MapValue[string](data, "nonexistent")
	if got != "" {
		t.Errorf("MapValue[string](data, 'nonexistent') = %q, want empty string", got)
	}

	gotFloat := MapValue[float64](data, "data.nonexistent")
	if gotFloat != 0 {
		t.Errorf("MapValue[float64](data, 'data.nonexistent') = %v, want 0", gotFloat)
	}

	gotMap := MapValue[map[string]any](data, "nonexistent")
	if gotMap != nil {
		t.Errorf("MapValue[map[string]any](data, 'nonexistent') = %v, want nil", gotMap)
	}
}

func TestMapValue_DefaultValue(t *testing.T) {
	data := buildTestMap()

	// 不存在的 key，使用默认值
	got := MapValue(data, "nonexistent", "fallback")
	if got != "fallback" {
		t.Errorf("MapValue(data, 'nonexistent', 'fallback') = %q, want fallback", got)
	}

	gotFloat := MapValue(data, "data.nonexistent", -1.0)
	if gotFloat != -1.0 {
		t.Errorf("MapValue(data, 'data.nonexistent', -1) = %v, want -1", gotFloat)
	}

	// 类型不匹配时也使用默认值
	gotStr := MapValue(data, "code", "default_code")
	if gotStr != "default_code" {
		t.Errorf("MapValue(data, 'code', 'default_code') = %q, want default_code (type mismatch fallback)", gotStr)
	}
}

func TestMapValue_NilRoot(t *testing.T) {
	got := MapValue[string](nil, "code")
	if got != "" {
		t.Errorf("MapValue[string](nil, 'code') = %q, want empty string", got)
	}

	got = MapValue(nil, "code", "default")
	if got != "default" {
		t.Errorf("MapValue(nil, 'code', 'default') = %q, want default", got)
	}
}

func TestMapValue_EmptyKey(t *testing.T) {
	data := buildTestMap()

	got := MapValue[string](data, "")
	if got != "" {
		t.Errorf("MapValue[string](data, '') = %q, want empty string", got)
	}

	got = MapValue(data, "", "default")
	if got != "default" {
		t.Errorf("MapValue(data, '', 'default') = %q, want default", got)
	}
}

func TestMapValue_NilValueInPath(t *testing.T) {
	data := map[string]any{
		"key": nil,
	}

	got := MapValue[string](data, "key")
	if got != "" {
		t.Errorf("MapValue[string](data, 'key') = %q, want empty string (nil value)", got)
	}

	got = MapValue(data, "key", "fallback")
	if got != "fallback" {
		t.Errorf("MapValue(data, 'key', 'fallback') = %q, want fallback (nil value)", got)
	}
}
