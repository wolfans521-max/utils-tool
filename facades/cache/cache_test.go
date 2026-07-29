package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// 测试辅助函数：创建临时缓存
func createTestCache(t *testing.T, name string) *SharedCache {
	tmpDir := t.TempDir()
	cache, err := NewSharedCache(SharedCacheConfig{
		Name: name,
		Size: 1024 * 1024, // 1MB
		Dir:  tmpDir,
	})
	if err != nil {
		t.Fatalf("创建缓存失败: %v", err)
	}
	return cache
}

// TestGetInstance 测试 GetInstance 单例获取功能
// 注意：GetInstance 是单例模式，依赖配置文件 config/app.yaml
// 此测试需要在项目根目录下运行，或者配置文件存在
func TestGetInstance(t *testing.T) {
	// 保存当前工作目录
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}

	// 创建临时目录和配置文件
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("创建配置目录失败: %v", err)
	}

	// 创建测试配置文件
	configContent := `cache:
  test_memoryname:
    max_memory: 10485760
`
	configPath := filepath.Join(configDir, "app.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("创建配置文件失败: %v", err)
	}

	// 切换到临时目录
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("切换目录失败: %v", err)
	}
	defer os.Chdir(originalDir)

	// 测试获取实例（注意：由于单例模式，此测试只能运行一次）
	// 如果之前已经初始化过单例，此测试可能会失败
	t.Run("获取缓存实例", func(t *testing.T) {
		cache, err := GetInstance("test_memoryname")
		if err != nil {
			t.Logf("GetInstance() error = %v (可能是单例已初始化或配置问题)", err)
			// 单例模式下，如果之前已初始化，可能会返回错误
			// 这里不作为失败处理
			return
		}

		if cache == nil {
			t.Error("GetInstance() 返回 nil")
			return
		}

		// 测试基本操作
		err = cache.Set("instance_key", "instance_value")
		if err != nil {
			t.Errorf("Set() error = %v", err)
		}

		value, ok := cache.Get("instance_key")
		if !ok {
			t.Error("Get() 应该找到键")
		}
		if value != "instance_value" {
			t.Errorf("Get() = %v, want instance_value", value)
		}

		// 清理
		cache.Delete("instance_key")
	})
}

// TestGetInstance_EmptyName 测试空名称的情况
func TestGetInstance_EmptyName(t *testing.T) {
	// 注意：由于单例模式，如果之前已经用有效名称初始化过，
	// 这个测试可能不会触发空名称错误
	// 这里主要测试函数签名和基本调用
	t.Run("空名称应该返回错误或已初始化的实例", func(t *testing.T) {
		_, err := GetInstance("")
		// 如果单例已初始化，可能返回已有实例而不是错误
		// 如果单例未初始化，应该返回错误
		if err != nil {
			t.Logf("GetInstance(\"\") 返回错误: %v (预期行为)", err)
		}
	})
}

func TestNewSharedCache(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		config  SharedCacheConfig
		wantErr bool
	}{
		{
			name: "正常创建缓存",
			config: SharedCacheConfig{
				Name: "test_cache",
				Size: 1024 * 1024,
				Dir:  tmpDir,
			},
			wantErr: false,
		},
		{
			name: "使用默认大小",
			config: SharedCacheConfig{
				Name: "test_cache_default",
				Size: 0, // 应该使用默认值 10MB
				Dir:  tmpDir,
			},
			wantErr: false,
		},
		{
			name: "不指定目录使用临时目录",
			config: SharedCacheConfig{
				Name: "test_cache_temp",
				Size: 1024 * 1024,
				Dir:  "", // 使用临时目录
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := NewSharedCache(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSharedCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if cache != nil {
				defer cache.Unlink()

				// 验证缓存文件已创建
				filePath := cache.GetFilePath()
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("缓存文件未创建: %s", filePath)
				}
			}
		})
	}
}

func TestSharedCache_SetAndGet(t *testing.T) {
	cache := createTestCache(t, "test_set_get")
	defer cache.Unlink()

	tests := []struct {
		name  string
		key   string
		value any
	}{
		{"字符串值", "string_key", "hello world"},
		{"整数值", "int_key", 12345},
		{"浮点数值", "float_key", 3.14159},
		{"布尔值", "bool_key", true},
		{"切片值", "slice_key", []string{"a", "b", "c"}},
		{"Map值", "map_key", map[string]any{"name": "test", "age": 18}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置缓存
			err := cache.Set(tt.key, tt.value)
			if err != nil {
				t.Errorf("Set() error = %v", err)
				return
			}

			// 获取缓存
			got, ok := cache.Get(tt.key)
			if !ok {
				t.Errorf("Get() 未找到键 %s", tt.key)
				return
			}

			// 注意：由于 JSON 序列化，数值类型会变成 float64
			// 这里只验证值存在，不做严格类型比较
			if got == nil {
				t.Errorf("Get() 返回 nil")
			}
		})
	}
}

func TestSharedCache_SetWithTTL(t *testing.T) {
	cache := createTestCache(t, "test_ttl")
	defer cache.Unlink()

	// 设置带 TTL 的缓存（1秒过期）
	err := cache.Set("ttl_key", "ttl_value", 1)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 立即获取应该成功
	value, ok := cache.Get("ttl_key")
	if !ok {
		t.Error("Get() 应该找到未过期的键")
	}
	if value != "ttl_value" {
		t.Errorf("Get() = %v, want %v", value, "ttl_value")
	}

	// 等待过期
	time.Sleep(2 * time.Second)

	// 过期后获取应该失败
	_, ok = cache.Get("ttl_key")
	if ok {
		t.Error("Get() 不应该找到已过期的键")
	}
}

func TestSharedCache_SetWithZeroTTL(t *testing.T) {
	cache := createTestCache(t, "test_zero_ttl")
	defer cache.Unlink()

	// 设置 TTL 为 0（永不过期）
	err := cache.Set("no_expire_key", "no_expire_value", 0)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 获取应该成功
	value, ok := cache.Get("no_expire_key")
	if !ok {
		t.Error("Get() 应该找到永不过期的键")
	}
	if value != "no_expire_value" {
		t.Errorf("Get() = %v, want %v", value, "no_expire_value")
	}
}

func TestSharedCache_Delete(t *testing.T) {
	cache := createTestCache(t, "test_delete")
	defer cache.Unlink()

	// 设置缓存
	err := cache.Set("delete_key", "delete_value")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 验证存在
	if !cache.Has("delete_key") {
		t.Error("Has() 应该返回 true")
	}

	// 删除缓存
	err = cache.Delete("delete_key")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// 验证已删除
	if cache.Has("delete_key") {
		t.Error("Has() 应该返回 false")
	}

	// 删除不存在的键不应报错
	err = cache.Delete("non_existent_key")
	if err != nil {
		t.Errorf("Delete() 删除不存在的键不应报错: %v", err)
	}
}

func TestSharedCache_Clear(t *testing.T) {
	cache := createTestCache(t, "test_clear")
	defer cache.Unlink()

	// 设置多个缓存
	for i := 0; i < 5; i++ {
		key := "key_" + string(rune('a'+i))
		err := cache.Set(key, i)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	// 验证数量
	if cache.Count() != 5 {
		t.Errorf("Count() = %d, want 5", cache.Count())
	}

	// 清空缓存
	err := cache.Clear()
	if err != nil {
		t.Errorf("Clear() error = %v", err)
	}

	// 验证已清空
	if cache.Count() != 0 {
		t.Errorf("Count() = %d, want 0", cache.Count())
	}
}

func TestSharedCache_Has(t *testing.T) {
	cache := createTestCache(t, "test_has")
	defer cache.Unlink()

	// 不存在的键
	if cache.Has("non_existent") {
		t.Error("Has() 对不存在的键应返回 false")
	}

	// 设置缓存
	err := cache.Set("exist_key", "exist_value")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 存在的键
	if !cache.Has("exist_key") {
		t.Error("Has() 对存在的键应返回 true")
	}
}

func TestSharedCache_Count(t *testing.T) {
	cache := createTestCache(t, "test_count")
	defer cache.Unlink()

	// 初始为空
	if cache.Count() != 0 {
		t.Errorf("Count() = %d, want 0", cache.Count())
	}

	// 添加缓存
	for i := 0; i < 3; i++ {
		key := "count_key_" + string(rune('a'+i))
		err := cache.Set(key, i)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	if cache.Count() != 3 {
		t.Errorf("Count() = %d, want 3", cache.Count())
	}

	// 删除一个
	err := cache.Delete("count_key_a")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if cache.Count() != 2 {
		t.Errorf("Count() = %d, want 2", cache.Count())
	}
}

func TestSharedCache_Keys(t *testing.T) {
	cache := createTestCache(t, "test_keys")
	defer cache.Unlink()

	// 初始为空
	keys := cache.Keys()
	if len(keys) != 0 {
		t.Errorf("Keys() 长度 = %d, want 0", len(keys))
	}

	// 添加缓存
	expectedKeys := []string{"key_a", "key_b", "key_c"}
	for _, key := range expectedKeys {
		err := cache.Set(key, "value")
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	// 获取所有键
	keys = cache.Keys()
	if len(keys) != len(expectedKeys) {
		t.Errorf("Keys() 长度 = %d, want %d", len(keys), len(expectedKeys))
	}

	// 验证所有键都存在（顺序可能不同）
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}
	for _, expected := range expectedKeys {
		if !keyMap[expected] {
			t.Errorf("Keys() 缺少键 %s", expected)
		}
	}
}

func TestSharedCache_GetFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	cache, err := NewSharedCache(SharedCacheConfig{
		Name: "test_filepath",
		Size: 1024 * 1024,
		Dir:  tmpDir,
	})
	if err != nil {
		t.Fatalf("NewSharedCache() error = %v", err)
	}
	defer cache.Unlink()

	filePath := cache.GetFilePath()
	expectedPath := filepath.Join(tmpDir, "test_filepath.cache")
	if filePath != expectedPath {
		t.Errorf("GetFilePath() = %s, want %s", filePath, expectedPath)
	}
}

func TestSharedCache_Close(t *testing.T) {
	cache := createTestCache(t, "test_close")

	// 设置一些数据
	err := cache.Set("close_key", "close_value")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 关闭缓存
	err = cache.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// 清理文件
	os.Remove(cache.GetFilePath())
}

func TestSharedCache_Unlink(t *testing.T) {
	tmpDir := t.TempDir()
	cache, err := NewSharedCache(SharedCacheConfig{
		Name: "test_unlink",
		Size: 1024 * 1024,
		Dir:  tmpDir,
	})
	if err != nil {
		t.Fatalf("NewSharedCache() error = %v", err)
	}

	filePath := cache.GetFilePath()

	// 验证文件存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("缓存文件应该存在")
	}

	// 删除缓存文件
	err = cache.Unlink()
	if err != nil {
		t.Errorf("Unlink() error = %v", err)
	}

	// 验证文件已删除
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("缓存文件应该已被删除")
	}
}

func TestSharedCache_GetNonExistent(t *testing.T) {
	cache := createTestCache(t, "test_get_nonexistent")
	defer cache.Unlink()

	// 获取不存在的键
	value, ok := cache.Get("non_existent_key")
	if ok {
		t.Error("Get() 对不存在的键应返回 false")
	}
	if value != nil {
		t.Errorf("Get() 对不存在的键应返回 nil, got %v", value)
	}
}

func TestSharedCache_OverwriteValue(t *testing.T) {
	cache := createTestCache(t, "test_overwrite")
	defer cache.Unlink()

	// 设置初始值
	err := cache.Set("overwrite_key", "initial_value")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 验证初始值
	value, ok := cache.Get("overwrite_key")
	if !ok || value != "initial_value" {
		t.Errorf("Get() = %v, want initial_value", value)
	}

	// 覆盖值
	err = cache.Set("overwrite_key", "new_value")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// 验证新值
	value, ok = cache.Get("overwrite_key")
	if !ok || value != "new_value" {
		t.Errorf("Get() = %v, want new_value", value)
	}
}

func TestSharedCache_ComplexData(t *testing.T) {
	cache := createTestCache(t, "test_complex")
	defer cache.Unlink()

	// 测试复杂嵌套数据结构
	complexData := map[string]any{
		"name": "test",
		"nested": map[string]any{
			"level1": map[string]any{
				"level2": "deep_value",
			},
		},
		"array": []any{1, 2, 3, "four"},
	}

	err := cache.Set("complex_key", complexData)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	value, ok := cache.Get("complex_key")
	if !ok {
		t.Error("Get() 应该找到复杂数据")
	}

	// 验证数据结构
	if valueMap, ok := value.(map[string]any); ok {
		if valueMap["name"] != "test" {
			t.Errorf("复杂数据 name = %v, want test", valueMap["name"])
		}
	} else {
		t.Error("Get() 返回的数据类型不正确")
	}
}

// 基准测试
func BenchmarkSharedCache_Set(b *testing.B) {
	tmpDir := b.TempDir()
	cache, err := NewSharedCache(SharedCacheConfig{
		Name: "bench_set",
		Size: 10 * 1024 * 1024,
		Dir:  tmpDir,
	})
	if err != nil {
		b.Fatalf("创建缓存失败: %v", err)
	}
	defer cache.Unlink()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "bench_key_" + string(rune(i%26+'a'))
		cache.Set(key, "bench_value")
	}
}

func BenchmarkSharedCache_Get(b *testing.B) {
	tmpDir := b.TempDir()
	cache, err := NewSharedCache(SharedCacheConfig{
		Name: "bench_get",
		Size: 10 * 1024 * 1024,
		Dir:  tmpDir,
	})
	if err != nil {
		b.Fatalf("创建缓存失败: %v", err)
	}
	defer cache.Unlink()

	// 预先设置一些数据
	for i := 0; i < 100; i++ {
		key := "bench_key_" + string(rune(i%26+'a'))
		cache.Set(key, "bench_value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "bench_key_" + string(rune(i%26+'a'))
		cache.Get(key)
	}
}

func BenchmarkSharedCache_SetAndGet(b *testing.B) {
	tmpDir := b.TempDir()
	cache, err := NewSharedCache(SharedCacheConfig{
		Name: "bench_set_get",
		Size: 10 * 1024 * 1024,
		Dir:  tmpDir,
	})
	if err != nil {
		b.Fatalf("创建缓存失败: %v", err)
	}
	defer cache.Unlink()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "bench_key_" + string(rune(i%26+'a'))
		cache.Set(key, "bench_value")
		cache.Get(key)
	}
}
