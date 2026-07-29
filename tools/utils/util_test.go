package utils

import (
	"testing"
)

// testMapper 实现 ToMapper 接口用于测试
type testMapper struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (t testMapper) ToMap() map[string]any {
	return map[string]any{
		"name": t.Name,
		"age":  t.Age,
	}
}

func TestStructToMap_ToMapperInterface(t *testing.T) {
	// 实现 ToMapper 接口的类型应优先使用 ToMap() 方法
	s := testMapper{Name: "alice", Age: 30}
	result := StructToMap(s)
	if result["name"] != "alice" {
		t.Errorf("name = %v, want alice", result["name"])
	}
	if result["age"] != 30 {
		t.Errorf("age = %v, want 30", result["age"])
	}
}

func TestStructToMap_JSONFallback(t *testing.T) {
	// 未实现 ToMapper 的 struct 使用 sonic JSON 往返兜底
	type TestStruct struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Hidden string `json:"hidden,omitempty"`
	}

	s := TestStruct{Name: "test", Age: 25, Hidden: ""}
	result := StructToMap(s)

	if result["name"] != "test" {
		t.Errorf("name = %v, want test", result["name"])
	}
	// JSON 往返后数值类型变为 float64
	if age, ok := result["age"].(float64); !ok || age != 25 {
		t.Errorf("age = %v, want 25 (float64)", result["age"])
	}
	// omitempty: 空字符串应被省略
	if _, ok := result["hidden"]; ok {
		t.Error("hidden field with omitempty and zero value should be excluded")
	}
}

func TestStructToMap_MapInput(t *testing.T) {
	// 已经是 map[string]any 的直接返回
	input := map[string]any{"key": "value"}
	result := StructToMap(input)
	if result["key"] != "value" {
		t.Errorf("key = %v, want value", result["key"])
	}
}

func TestStructToMap_NilInput(t *testing.T) {
	result := StructToMap(nil)
	if result != nil {
		t.Errorf("nil input should return nil, got %v", result)
	}
}

func TestStructToMap_PointerInput(t *testing.T) {
	type Simple struct {
		Name string `json:"name"`
	}

	s := &Simple{Name: "ptr_test"}
	result := StructToMap(s)
	if result["name"] != "ptr_test" {
		t.Errorf("name = %v, want ptr_test", result["name"])
	}
}

func TestStructToMap_BoolAndFloat(t *testing.T) {
	type Types struct {
		Flag  bool    `json:"flag"`
		Score float64 `json:"score"`
	}

	s := Types{Flag: true, Score: 3.14}
	result := StructToMap(s)
	if result["flag"] != true {
		t.Errorf("flag = %v, want true", result["flag"])
	}
	if result["score"] != 3.14 {
		t.Errorf("score = %v, want 3.14", result["score"])
	}
}

func TestStructToMap_NestedStruct(t *testing.T) {
	type Inner struct {
		Value string `json:"value"`
	}
	type Outer struct {
		Name  string `json:"name"`
		Inner Inner  `json:"inner"`
	}

	s := Outer{Name: "outer", Inner: Inner{Value: "inner_val"}}
	result := StructToMap(s)

	if result["name"] != "outer" {
		t.Errorf("name = %v, want outer", result["name"])
	}
	if inner, ok := result["inner"].(map[string]any); ok {
		if inner["value"] != "inner_val" {
			t.Errorf("inner.value = %v, want inner_val", inner["value"])
		}
	} else {
		t.Errorf("inner = %v, want map[string]any", result["inner"])
	}
}

func TestStructToMap_SliceField(t *testing.T) {
	type WithSlice struct {
		Items []string `json:"items"`
	}

	s := WithSlice{Items: []string{"a", "b", "c"}}
	result := StructToMap(s)

	if items, ok := result["items"].([]any); ok {
		if len(items) != 3 {
			t.Errorf("items length = %d, want 3", len(items))
		}
	} else {
		t.Errorf("items = %v, want []any", result["items"])
	}
}
