package cache

import (
	"fmt"
	"testing"
	"time"
)

func newTestBigCacheAdapter(t *testing.T) *BigCacheAdapter {
	t.Helper()
	config := BigCacheConfig{
		Shards:             64,
		LifeWindow:         2 * time.Hour,
		CleanWindow:        30 * time.Second,
		MaxEntriesInWindow: 10000,
		MaxEntrySize:       1024 * 1024,
		HardMaxCacheSize:   256,
	}
	adapter, err := NewBigCacheAdapter(config)
	if err != nil {
		t.Fatalf("创建 BigCacheAdapter 失败: %v", err)
	}
	t.Cleanup(func() {
		adapter.Close()
	})
	return adapter
}

func TestBigCacheAdapter_SetAndGet(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	// 测试基本 Set/Get
	err := adapter.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Set 失败: %v", err)
	}

	val, ok := adapter.Get("key1")
	if !ok {
		t.Fatal("Get 未命中")
	}
	if val != "value1" {
		t.Fatalf("期望 value1, 实际 %v", val)
	}
}

func TestBigCacheAdapter_SetWithTTL(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	// 设置 1 秒过期
	err := adapter.Set("key_ttl", "value_ttl", 1)
	if err != nil {
		t.Fatalf("Set 失败: %v", err)
	}

	// 立即获取应该命中
	val, ok := adapter.Get("key_ttl")
	if !ok {
		t.Fatal("Get 未命中")
	}
	if val != "value_ttl" {
		t.Fatalf("期望 value_ttl, 实际 %v", val)
	}

	// 等待过期
	time.Sleep(1100 * time.Millisecond)

	// 过期后应该未命中
	_, ok = adapter.Get("key_ttl")
	if ok {
		t.Fatal("过期后应该未命中")
	}
}

func TestBigCacheAdapter_SetWithDuration(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	err := adapter.SetWithDuration("key_dur", "value_dur", 5*time.Second)
	if err != nil {
		t.Fatalf("SetWithDuration 失败: %v", err)
	}

	val, ok := adapter.Get("key_dur")
	if !ok {
		t.Fatal("Get 未命中")
	}
	if val != "value_dur" {
		t.Fatalf("期望 value_dur, 实际 %v", val)
	}

	// 0 duration 表示永不过期
	err = adapter.SetWithDuration("key_no_exp", "value_no_exp", 0)
	if err != nil {
		t.Fatalf("SetWithDuration 失败: %v", err)
	}
	val, ok = adapter.Get("key_no_exp")
	if !ok || val != "value_no_exp" {
		t.Fatal("永不过期的 key 应该命中")
	}
}

func TestBigCacheAdapter_GetMiss(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	_, ok := adapter.Get("nonexistent")
	if ok {
		t.Fatal("不存在的 key 不应该命中")
	}
}

func TestBigCacheAdapter_Delete(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	adapter.Set("key_del", "value_del")
	err := adapter.Delete("key_del")
	if err != nil {
		t.Fatalf("Delete 失败: %v", err)
	}

	_, ok := adapter.Get("key_del")
	if ok {
		t.Fatal("删除后不应该命中")
	}
}

func TestBigCacheAdapter_DeleteBatch(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	adapter.Set("key1", "val1")
	adapter.Set("key2", "val2")
	adapter.Set("key3", "val3")

	err := adapter.DeleteBatch([]string{"key1", "key3"})
	if err != nil {
		t.Fatalf("DeleteBatch 失败: %v", err)
	}

	_, ok1 := adapter.Get("key1")
	_, ok2 := adapter.Get("key2")
	_, ok3 := adapter.Get("key3")
	if ok1 {
		t.Fatal("key1 应该已删除")
	}
	if !ok2 {
		t.Fatal("key2 不应该被删除")
	}
	if ok3 {
		t.Fatal("key3 应该已删除")
	}
}

func TestBigCacheAdapter_Clear(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	adapter.Set("key1", "val1")
	adapter.Set("key2", "val2")

	err := adapter.Clear()
	if err != nil {
		t.Fatalf("Clear 失败: %v", err)
	}

	if adapter.Count() != 0 {
		t.Fatalf("清空后 Count 应该为 0, 实际 %d", adapter.Count())
	}
}

func TestBigCacheAdapter_Has(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	if adapter.Has("key1") {
		t.Fatal("不存在的 key 不应该 Has")
	}

	adapter.Set("key1", "val1")
	if !adapter.Has("key1") {
		t.Fatal("存在的 key 应该 Has")
	}
}

func TestBigCacheAdapter_Count(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	if adapter.Count() != 0 {
		t.Fatalf("初始 Count 应该为 0, 实际 %d", adapter.Count())
	}

	adapter.Set("key1", "val1")
	adapter.Set("key2", "val2")
	if adapter.Count() != 2 {
		t.Fatalf("Count 应该为 2, 实际 %d", adapter.Count())
	}
}

func TestBigCacheAdapter_Keys(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	adapter.Set("key1", "val1")
	adapter.Set("key2", "val2")

	keys := adapter.Keys()
	if len(keys) != 2 {
		t.Fatalf("Keys 长度应该为 2, 实际 %d", len(keys))
	}

	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}
	if !keyMap["key1"] || !keyMap["key2"] {
		t.Fatal("Keys 应该包含 key1 和 key2")
	}
}

func TestBigCacheAdapter_GetString(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	adapter.Set("str_key", "hello")
	val, ok := adapter.GetString("str_key")
	if !ok || val != "hello" {
		t.Fatalf("期望 hello, 实际 %s, ok=%v", val, ok)
	}

	// 非 string 类型
	adapter.Set("int_key", 123)
	_, ok = adapter.GetString("int_key")
	if ok {
		t.Fatal("非 string 类型不应该 GetString 成功")
	}
}

func TestBigCacheAdapter_GetInt(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	adapter.Set("int_key", 42)
	val, ok := adapter.GetInt("int_key")
	if !ok || val != 42 {
		t.Fatalf("期望 42, 实际 %d, ok=%v", val, ok)
	}

	// float64 类型（JSON 反序列化数字默认为 float64）
	adapter.Set("float_key", float64(99))
	val, ok = adapter.GetInt("float_key")
	if !ok || val != 99 {
		t.Fatalf("期望 99, 实际 %d, ok=%v", val, ok)
	}
}

func TestBigCacheAdapter_GetMap(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	testMap := map[string]any{"a": 1, "b": "hello"}
	adapter.Set("map_key", testMap)

	val, ok := adapter.GetMap("map_key")
	if !ok {
		t.Fatal("GetMap 未命中")
	}
	if val["a"] == nil || val["b"] == nil {
		t.Fatalf("GetMap 结果不完整: %v", val)
	}
}

func TestBigCacheAdapter_GetSlice(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	testSlice := []any{1, "hello", true}
	adapter.Set("slice_key", testSlice)

	val, ok := adapter.GetSlice("slice_key")
	if !ok {
		t.Fatal("GetSlice 未命中")
	}
	if len(val) != 3 {
		t.Fatalf("GetSlice 长度应该为 3, 实际 %d", len(val))
	}
}

func TestBigCacheAdapter_SetBatch(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	items := map[string]any{
		"batch1": "val1",
		"batch2": "val2",
	}
	err := adapter.SetBatch(items, 60)
	if err != nil {
		t.Fatalf("SetBatch 失败: %v", err)
	}

	val1, ok1 := adapter.Get("batch1")
	val2, ok2 := adapter.Get("batch2")
	if !ok1 || val1 != "val1" {
		t.Fatal("batch1 未命中或值不正确")
	}
	if !ok2 || val2 != "val2" {
		t.Fatal("batch2 未命中或值不正确")
	}
}

func TestBigCacheAdapter_GetBatch(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	adapter.Set("batch1", "val1")
	adapter.Set("batch2", "val2")

	result := adapter.GetBatch([]string{"batch1", "batch2", "batch3"})
	if len(result) != 2 {
		t.Fatalf("GetBatch 应该返回 2 个结果, 实际 %d", len(result))
	}
	if result["batch1"] != "val1" {
		t.Fatal("batch1 值不正确")
	}
	if result["batch2"] != "val2" {
		t.Fatal("batch2 值不正确")
	}
}

func TestBigCacheAdapter_SetWithPolicy(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	// 相对过期
	err := adapter.SetWithPolicy("relative_key", "val", ExpirePolicyRelative, 2)
	if err != nil {
		t.Fatalf("SetWithPolicy 失败: %v", err)
	}
	val, ok := adapter.Get("relative_key")
	if !ok || val != "val" {
		t.Fatal("相对过期策略 Set/Get 失败")
	}

	// 滑动过期
	err = adapter.SetWithPolicy("sliding_key", "val2", ExpirePolicySliding, 2)
	if err != nil {
		t.Fatalf("SetWithPolicy 失败: %v", err)
	}
	val, ok = adapter.Get("sliding_key")
	if !ok || val != "val2" {
		t.Fatal("滑动过期策略 Set/Get 失败")
	}
}

func TestBigCacheAdapter_Close(t *testing.T) {
	config := BigCacheConfig{
		Shards:             64,
		LifeWindow:         2 * time.Hour,
		CleanWindow:        30 * time.Second,
		MaxEntriesInWindow: 10000,
		MaxEntrySize:       1024 * 1024,
		HardMaxCacheSize:   256,
	}
	adapter, err := NewBigCacheAdapter(config)
	if err != nil {
		t.Fatalf("创建 BigCacheAdapter 失败: %v", err)
	}

	adapter.Set("key1", "val1")

	err = adapter.Close()
	if err != nil {
		t.Fatalf("Close 失败: %v", err)
	}

	// 关闭后操作应该失败
	err = adapter.Set("key2", "val2")
	if err == nil {
		t.Fatal("关闭后 Set 应该失败")
	}

	_, ok := adapter.Get("key1")
	if ok {
		t.Fatal("关闭后 Get 应该失败")
	}
}

func TestBigCacheAdapter_GetStats(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	adapter.Set("key1", "val1")
	adapter.Get("key1") // hit
	adapter.Get("key1") // hit
	adapter.Get("nonexistent") // miss

	stats := adapter.GetStats()
	if stats.TotalSets != 1 {
		t.Fatalf("TotalSets 应该为 1, 实际 %d", stats.TotalSets)
	}
	if stats.TotalGets != 3 {
		t.Fatalf("TotalGets 应该为 3, 实际 %d", stats.TotalGets)
	}
	if stats.HitCount != 2 {
		t.Fatalf("HitCount 应该为 2, 实际 %d", stats.HitCount)
	}
	if stats.MissCount != 1 {
		t.Fatalf("MissCount 应该为 1, 实际 %d", stats.MissCount)
	}
	if stats.ItemCount != 1 {
		t.Fatalf("ItemCount 应该为 1, 实际 %d", stats.ItemCount)
	}
}

func TestBigCacheAdapter_ComplexValue(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	// 测试复杂嵌套结构
	complexValue := map[string]any{
		"string_field": "hello",
		"int_field":    42,
		"nested": map[string]any{
			"inner": "world",
		},
		"array_field": []any{1, 2, 3},
	}

	err := adapter.Set("complex_key", complexValue)
	if err != nil {
		t.Fatalf("Set 复杂值失败: %v", err)
	}

	val, ok := adapter.Get("complex_key")
	if !ok {
		t.Fatal("Get 复杂值未命中")
	}

	result, ok := val.(map[string]any)
	if !ok {
		t.Fatal("类型断言失败")
	}
	if result["string_field"] != "hello" {
		t.Fatal("string_field 值不正确")
	}
}

func TestBigCacheAdapter_Concurrent(t *testing.T) {
	adapter := newTestBigCacheAdapter(t)

	const goroutines = 100
	const opsPerGoroutine = 100

	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				adapter.Set(key, j)
				adapter.Get(key)
				if j%10 == 0 {
					adapter.Delete(key)
				}
			}
			done <- true
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}
}

// Benchmark 对比

func BenchmarkBigCacheAdapter_Set(b *testing.B) {
	config := BigCacheConfig{
		Shards:             64,
		LifeWindow:         2 * time.Hour,
		CleanWindow:        30 * time.Second,
		MaxEntriesInWindow: 10000,
		MaxEntrySize:       1024 * 1024,
		HardMaxCacheSize:   256,
	}
	adapter, _ := NewBigCacheAdapter(config)
	defer adapter.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i%10000)
		adapter.Set(key, "bench_value")
	}
}

func BenchmarkBigCacheAdapter_Get(b *testing.B) {
	config := BigCacheConfig{
		Shards:             64,
		LifeWindow:         2 * time.Hour,
		CleanWindow:        30 * time.Second,
		MaxEntriesInWindow: 10000,
		MaxEntrySize:       1024 * 1024,
		HardMaxCacheSize:   256,
	}
	adapter, _ := NewBigCacheAdapter(config)
	defer adapter.Close()

	// 预填充
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		adapter.Set(key, "bench_value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i%10000)
		adapter.Get(key)
	}
}

func BenchmarkBigCacheAdapter_SetParallel(b *testing.B) {
	config := BigCacheConfig{
		Shards:             64,
		LifeWindow:         2 * time.Hour,
		CleanWindow:        30 * time.Second,
		MaxEntriesInWindow: 10000,
		MaxEntrySize:       1024 * 1024,
		HardMaxCacheSize:   256,
	}
	adapter, _ := NewBigCacheAdapter(config)
	defer adapter.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench_key_%d", i%10000)
			adapter.Set(key, "bench_value")
			i++
		}
	})
}

func BenchmarkBigCacheAdapter_GetParallel(b *testing.B) {
	config := BigCacheConfig{
		Shards:             64,
		LifeWindow:         2 * time.Hour,
		CleanWindow:        30 * time.Second,
		MaxEntriesInWindow: 10000,
		MaxEntrySize:       1024 * 1024,
		HardMaxCacheSize:   256,
	}
	adapter, _ := NewBigCacheAdapter(config)
	defer adapter.Close()

	// 预填充
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		adapter.Set(key, "bench_value")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench_key_%d", i%10000)
			adapter.Get(key)
			i++
		}
	})
}
