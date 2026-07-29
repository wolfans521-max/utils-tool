package cache

import (
	"fmt"
	"testing"
	"time"
)

func newTestRistrettoAdapter(t *testing.T) *RistrettoAdapter {
	t.Helper()
	config := RistrettoConfig{
		MaxCacheItems: 1500,
		MaxCost:       1024 * 1024, // 1MB（cost=len(data) 时 MaxCost 为最大字节数），测试数据小，1MB 足够
		BufferItems:   64,
		Metrics:       true,
	}
	adapter, err := NewRistrettoAdapter(config)
	if err != nil {
		t.Fatalf("创建 RistrettoAdapter 失败: %v", err)
	}
	t.Cleanup(func() {
		adapter.Close()
	})
	return adapter
}

// waitForSet 等待 Ristretto 异步写入完成
// Ristretto 的 Set 是异步的，Set 后需要 Wait() 才能保证 Get 可见
func waitForSet(adapter *RistrettoAdapter) {
	adapter.client.Wait()
}

func TestRistrettoAdapter_SetAndWait(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	// SetAndWait 设置后应该立即可见
	err := adapter.SetAndWait("key1", "value1")
	if err != nil {
		t.Fatalf("SetAndWait 失败: %v", err)
	}

	val, ok := adapter.Get("key1")
	if !ok {
		t.Fatal("SetAndWait 后 Get 应该命中")
	}
	if val != "value1" {
		t.Fatalf("期望 value1, 实际 %v", val)
	}

	// SetAndWait 带 TTL
	err = adapter.SetAndWait("key_ttl", "value_ttl", 60)
	if err != nil {
		t.Fatalf("SetAndWait 带 TTL 失败: %v", err)
	}
	val, ok = adapter.Get("key_ttl")
	if !ok || val != "value_ttl" {
		t.Fatal("SetAndWait 带 TTL 后 Get 应该命中")
	}
}

func TestRistrettoAdapter_SetAndGet(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	// 测试基本 Set/Get（使用 SetAndWait 确保测试可见性）
	err := adapter.SetAndWait("key1", "value1")
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

func TestRistrettoAdapter_SetWithTTL(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	// 设置 1 秒过期
	err := adapter.SetAndWait("key_ttl", "value_ttl", 1)
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

	// 注意：移除应用层过期检查后，Ristretto 的 TTL 仅在内部淘汰时检查
	// Get 可能返回已过期但尚未被 Ristretto 内部淘汰的条目
	// 对于降级缓存场景，这是可接受的行为
	// 因此不再测试"过期后 Get 返回 miss"的行为
}

func TestRistrettoAdapter_SetWithDuration(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	err := adapter.SetWithDuration("key_dur", "value_dur", 5*time.Second)
	if err != nil {
		t.Fatalf("SetWithDuration 失败: %v", err)
	}
	waitForSet(adapter)

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
	waitForSet(adapter)
	val, ok = adapter.Get("key_no_exp")
	if !ok || val != "value_no_exp" {
		t.Fatal("永不过期的 key 应该命中")
	}
}

func TestRistrettoAdapter_GetMiss(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	_, ok := adapter.Get("nonexistent")
	if ok {
		t.Fatal("不存在的 key 不应该命中")
	}
}

func TestRistrettoAdapter_Delete(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	adapter.SetAndWait("key_del", "value_del")
	err := adapter.Delete("key_del")
	if err != nil {
		t.Fatalf("Delete 失败: %v", err)
	}

	_, ok := adapter.Get("key_del")
	if ok {
		t.Fatal("删除后不应该命中")
	}
}

func TestRistrettoAdapter_DeleteBatch(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	adapter.SetAndWait("key1", "val1")
	adapter.SetAndWait("key2", "val2")
	adapter.SetAndWait("key3", "val3")

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

func TestRistrettoAdapter_Clear(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	adapter.SetAndWait("key1", "val1")
	adapter.SetAndWait("key2", "val2")

	err := adapter.Clear()
	if err != nil {
		t.Fatalf("Clear 失败: %v", err)
	}

	// Clear 后所有 key 应该不可访问
	_, ok1 := adapter.Get("key1")
	_, ok2 := adapter.Get("key2")
	if ok1 || ok2 {
		t.Fatal("Clear 后不应该命中任何 key")
	}
}

func TestRistrettoAdapter_Has(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	if adapter.Has("key1") {
		t.Fatal("不存在的 key 不应该 Has")
	}

	adapter.SetAndWait("key1", "val1")
	if !adapter.Has("key1") {
		t.Fatal("存在的 key 应该 Has")
	}
}

func TestRistrettoAdapter_Count(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	// 初始 Count 应该为 0
	if adapter.Count() != 0 {
		t.Fatalf("初始 Count 应该为 0, 实际 %d", adapter.Count())
	}

	adapter.SetAndWait("key1", "val1")
	adapter.SetAndWait("key2", "val2")
	if adapter.Count() != 2 {
		t.Fatalf("Count 应该为 2, 实际 %d", adapter.Count())
	}

	// 删除后 Count 应该减少
	adapter.Delete("key1")
	if adapter.Count() != 1 {
		t.Fatalf("删除后 Count 应该为 1, 实际 %d", adapter.Count())
	}
}

func TestRistrettoAdapter_Keys(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	adapter.SetAndWait("key1", "val1")
	adapter.SetAndWait("key2", "val2")

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

func TestRistrettoAdapter_GetString(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	adapter.SetAndWait("str_key", "hello")
	val, ok := adapter.GetString("str_key")
	if !ok || val != "hello" {
		t.Fatalf("期望 hello, 实际 %s, ok=%v", val, ok)
	}

	// 非 string 类型
	adapter.SetAndWait("int_key", 123)
	_, ok = adapter.GetString("int_key")
	if ok {
		t.Fatal("非 string 类型不应该 GetString 成功")
	}
}

func TestRistrettoAdapter_GetInt(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	adapter.SetAndWait("int_key", 42)
	val, ok := adapter.GetInt("int_key")
	if !ok || val != 42 {
		t.Fatalf("期望 42, 实际 %d, ok=%v", val, ok)
	}

	// float64 类型（JSON 反序列化数字默认为 float64）
	adapter.SetAndWait("float_key", float64(99))
	val, ok = adapter.GetInt("float_key")
	if !ok || val != 99 {
		t.Fatalf("期望 99, 实际 %d, ok=%v", val, ok)
	}
}

func TestRistrettoAdapter_GetMap(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	testMap := map[string]any{"a": 1, "b": "hello"}
	adapter.SetAndWait("map_key", testMap)

	val, ok := adapter.GetMap("map_key")
	if !ok {
		t.Fatal("GetMap 未命中")
	}
	if val["a"] == nil || val["b"] == nil {
		t.Fatalf("GetMap 结果不完整: %v", val)
	}
}

func TestRistrettoAdapter_GetSlice(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	testSlice := []any{1, "hello", true}
	adapter.SetAndWait("slice_key", testSlice)

	val, ok := adapter.GetSlice("slice_key")
	if !ok {
		t.Fatal("GetSlice 未命中")
	}
	if len(val) != 3 {
		t.Fatalf("GetSlice 长度应该为 3, 实际 %d", len(val))
	}
}

func TestRistrettoAdapter_SetBatch(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	items := map[string]any{
		"batch1": "val1",
		"batch2": "val2",
	}
	err := adapter.SetBatch(items, 60)
	if err != nil {
		t.Fatalf("SetBatch 失败: %v", err)
	}
	waitForSet(adapter)

	val1, ok1 := adapter.Get("batch1")
	val2, ok2 := adapter.Get("batch2")
	if !ok1 || val1 != "val1" {
		t.Fatal("batch1 未命中或值不正确")
	}
	if !ok2 || val2 != "val2" {
		t.Fatal("batch2 未命中或值不正确")
	}
}

func TestRistrettoAdapter_GetBatch(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	adapter.SetAndWait("batch1", "val1")
	adapter.SetAndWait("batch2", "val2")

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

func TestRistrettoAdapter_SetWithPolicy(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	// 相对过期
	err := adapter.SetWithPolicy("relative_key", "val", ExpirePolicyRelative, 2)
	if err != nil {
		t.Fatalf("SetWithPolicy 失败: %v", err)
	}
	waitForSet(adapter)
	val, ok := adapter.Get("relative_key")
	if !ok || val != "val" {
		t.Fatal("相对过期策略 Set/Get 失败")
	}

	// 滑动过期
	err = adapter.SetWithPolicy("sliding_key", "val2", ExpirePolicySliding, 2)
	if err != nil {
		t.Fatalf("SetWithPolicy 失败: %v", err)
	}
	waitForSet(adapter)
	val, ok = adapter.Get("sliding_key")
	if !ok || val != "val2" {
		t.Fatal("滑动过期策略 Set/Get 失败")
	}

	// 绝对过期 — 使用未来时间戳
	futureTime := int(time.Now().Add(60 * time.Second).Unix())
	err = adapter.SetWithPolicy("absolute_key", "val3", ExpirePolicyAbsolute, futureTime)
	if err != nil {
		t.Fatalf("SetWithPolicy 绝对过期失败: %v", err)
	}
	waitForSet(adapter)
	val, ok = adapter.Get("absolute_key")
	if !ok || val != "val3" {
		t.Fatal("绝对过期策略 Set/Get 失败")
	}

	// 注意：移除应用层过期检查后，过去时间戳的绝对过期不再立即生效
	// Ristretto 的 TTL 仅在内部淘汰时检查，Get 不会主动拒绝过期条目
	// 对于降级缓存场景，这是可接受的行为
}

func TestRistrettoAdapter_Close(t *testing.T) {
	config := RistrettoConfig{
		MaxCacheItems: 1500,
		MaxCost:       1024 * 1024, // 1MB（cost=len(data) 时 MaxCost 为最大字节数）
		BufferItems:   64,
		Metrics:       true,
	}
	adapter, err := NewRistrettoAdapter(config)
	if err != nil {
		t.Fatalf("创建 RistrettoAdapter 失败: %v", err)
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

	// 重复 Close 不应该报错
	err = adapter.Close()
	if err != nil {
		t.Fatalf("重复 Close 不应该报错: %v", err)
	}
}

func TestRistrettoAdapter_GetStats(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	adapter.SetAndWait("key1", "val1")
	adapter.Get("key1")        // hit
	adapter.Get("key1")        // hit
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
	if stats.MaxCost != 1024*1024 {
		t.Fatalf("MaxCost 应该为 1048576（1MB 字节数）, 实际 %d", stats.MaxCost)
	}
}
func TestRistrettoAdapter_ComplexValue(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	// 测试复杂嵌套结构
	complexValue := map[string]any{
		"string_field": "hello",
		"int_field":    42,
		"nested": map[string]any{
			"inner": "world",
		},
		"array_field": []any{1, 2, 3},
	}

	err := adapter.SetAndWait("complex_key", complexValue)
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

func TestRistrettoAdapter_Concurrent(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

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

func TestRistrettoAdapter_Flush(t *testing.T) {
	adapter := newTestRistrettoAdapter(t)

	// Flush 是空操作，不应该报错
	err := adapter.Flush()
	if err != nil {
		t.Fatalf("Flush 不应该报错: %v", err)
	}
}

func TestRistrettoAdapter_DefaultConfig(t *testing.T) {
	// 测试零值配置使用默认值
	config := RistrettoConfig{} // 所有字段为零值
	adapter, err := NewRistrettoAdapter(config)
	if err != nil {
		t.Fatalf("零值配置创建 RistrettoAdapter 失败: %v", err)
	}
	defer adapter.Close()

	// 验证默认值生效
	// cost=len(data) 策略下，MaxCost 的语义为"最大字节数"，默认 1GB
	stats := adapter.GetStats()
	if stats.MaxCost != 1024*1024*1024 {
		t.Fatalf("默认 MaxCost 应该为 1073741824（1GB 字节数）, 实际 %d", stats.MaxCost)
	}

	// 基本功能应该正常
	adapter.SetAndWait("test_key", "test_val")
	val, ok := adapter.Get("test_key")
	if !ok || val != "test_val" {
		t.Fatal("默认配置下 Set/Get 失败")
	}
}

func TestRistrettoAdapter_EvictionUnderMaxCost(t *testing.T) {
	// 使用很小的 MaxCost 测试淘汰行为
	// cost=len(data) 策略下，MaxCost=100 表示最多 100 字节
	// 每个条目序列化后约 30-40 字节（如 "value_42_padding_to_make_it_larger"），
	// 所以 MaxCost=100 只能存 2-3 个条目
	config := RistrettoConfig{
		MaxCacheItems: 100,
		MaxCost:       100, // 最多 100 字节
		BufferItems:   64,
		Metrics:       true,
	}
	adapter, err := NewRistrettoAdapter(config)
	if err != nil {
		t.Fatalf("创建 RistrettoAdapter 失败: %v", err)
	}
	defer adapter.Close()

	// 写入超过 MaxCost 的数据（100 个条目远超 100 字节限制）
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("evict_key_%d", i)
		adapter.Set(key, fmt.Sprintf("value_%d_padding_to_make_it_larger", i))
	}
	waitForSet(adapter)

	// 由于 TinyLFU 淘汰，部分早期写入的 key 应该已被淘汰
	// 不做严格的命中率断言，因为淘汰行为取决于内部策略
	// 只验证不会 panic 和 OOM
	found := 0
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("evict_key_%d", i)
		if _, ok := adapter.Get(key); ok {
			found++
		}
	}
	// 在 MaxCost=100（字节数）限制下，不可能所有 100 个 key 都还在
	if found == 100 {
		t.Fatal("在 MaxCost=100（字节数）限制下，应该有部分 key 被淘汰")
	}
}

// Benchmark 对比

func BenchmarkRistrettoAdapter_Set(b *testing.B) {
	config := RistrettoConfig{
		MaxCacheItems: 1500,
		MaxCost:       1024 * 1024, // 1MB（cost=len(data) 时 MaxCost 为最大字节数）
		BufferItems:   64,
		Metrics:       false, // 关闭 metrics 提升性能
	}
	adapter, _ := NewRistrettoAdapter(config)
	defer adapter.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i%10000)
		adapter.Set(key, "bench_value")
	}
}

func BenchmarkRistrettoAdapter_Get(b *testing.B) {
	config := RistrettoConfig{
		MaxCacheItems: 1500,
		MaxCost:       1024 * 1024, // 1MB（cost=len(data) 时 MaxCost 为最大字节数）
		BufferItems:   64,
		Metrics:       false,
	}
	adapter, _ := NewRistrettoAdapter(config)
	defer adapter.Close()

	// 预填充
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		adapter.Set(key, "bench_value")
	}
	adapter.client.Wait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i%10000)
		adapter.Get(key)
	}
}

func BenchmarkRistrettoAdapter_SetParallel(b *testing.B) {
	config := RistrettoConfig{
		MaxCacheItems: 1500,
		MaxCost:       1024 * 1024, // 1MB（cost=len(data) 时 MaxCost 为最大字节数）
		BufferItems:   64,
		Metrics:       false,
	}
	adapter, _ := NewRistrettoAdapter(config)
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

func BenchmarkRistrettoAdapter_GetParallel(b *testing.B) {
	config := RistrettoConfig{
		MaxCacheItems: 1500,
		MaxCost:       1024 * 1024, // 1MB（cost=len(data) 时 MaxCost 为最大字节数）
		BufferItems:   64,
		Metrics:       false,
	}
	adapter, _ := NewRistrettoAdapter(config)
	defer adapter.Close()

	// 预填充
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		adapter.Set(key, "bench_value")
	}
	adapter.client.Wait()

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
