// Package cache 提供基于 BigCache 的高性能缓存实现
// 适用于高并发场景，零 GC 压力，替代基于 mmap 的 SharedCache
package cache

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/bytedance/sonic"
)

// BigCacheConfig BigCache 配置
type BigCacheConfig struct {
	// Shards 分片数量，必须是2的幂，默认64
	Shards int
	// LifeWindow 全局最大TTL，条目最长存活时间，默认2小时
	LifeWindow time.Duration
	// CleanWindow 清理间隔，默认30秒
	CleanWindow time.Duration
	// MaxEntriesInWindow 窗口内最大条目数，影响初始分片大小，默认10000
	MaxEntriesInWindow int
	// MaxEntrySize 单条目最大字节数，默认1MB
	MaxEntrySize int
	// HardMaxCacheSize 缓存最大大小（MB），默认256MB
	HardMaxCacheSize int
	// Verbose 是否启用详细日志
	Verbose bool
}

// DefaultBigCacheConfig 返回默认 BigCache 配置
func DefaultBigCacheConfig() BigCacheConfig {
	return BigCacheConfig{
		Shards:             64,
		LifeWindow:         2 * time.Hour,
		CleanWindow:        30 * time.Second,
		MaxEntriesInWindow: 10000,
		MaxEntrySize:       1 * 1024 * 1024, // 1MB
		HardMaxCacheSize:   256,              // 256MB
		Verbose:            false,
	}
}


// BigCacheStats BigCache 统计信息
type BigCacheStats struct {
	ItemCount    int     // 缓存项数量
	HitCount     int64   // 命中次数
	MissCount    int64   // 未命中次数
	TotalGets    int64   // 总Get次数
	TotalSets    int64   // 总Set次数
	TotalDeletes int64   // 总Delete次数
	HitRate      float64 // 命中率
	Capacity     int     // 缓存容量（字节）
	UsedSize     int64   // 已使用字节数（序列化后的 entry 字节数累加）
}

// BigCacheAdapter 基于 bigcache 的缓存适配器
// 实现与 SharedCache 相同的方法集，支持透明切换
type BigCacheAdapter struct {
	client *bigcache.BigCache
	config BigCacheConfig
	closed atomic.Bool

	// 统计信息
	stats struct {
		totalGets    atomic.Int64
		totalSets    atomic.Int64
		totalDeletes atomic.Int64
		hitCount     atomic.Int64
		missCount    atomic.Int64
		usedBytes    atomic.Int64 // 已使用字节数（序列化后的 entry 字节数累加）
	}
}

var (
	bigCacheInstance *BigCacheAdapter
	bigCacheOnce     sync.Once
	bigCacheInitErr  error
)

// GetBigCacheInstance 获取 BigCache 单例实例
func GetBigCacheInstance(config BigCacheConfig) (*BigCacheAdapter, error) {
	bigCacheOnce.Do(func() {
		impl, err := NewBigCacheAdapter(config)
		if err != nil {
			bigCacheInitErr = err
			return
		}
		bigCacheInstance = impl
	})
	return bigCacheInstance, bigCacheInitErr
}

// NewBigCacheAdapter 创建新的 BigCache 适配器实例
func NewBigCacheAdapter(config BigCacheConfig) (*BigCacheAdapter, error) {
	// 设置默认值
	if config.Shards <= 0 {
		config.Shards = 64
	}
	if config.LifeWindow <= 0 {
		config.LifeWindow = 2 * time.Hour
	}
	if config.CleanWindow <= 0 {
		config.CleanWindow = 30 * time.Second
	}
	if config.MaxEntriesInWindow <= 0 {
		config.MaxEntriesInWindow = 10000
	}
	if config.MaxEntrySize <= 0 {
		config.MaxEntrySize = 1 * 1024 * 1024
	}
	if config.HardMaxCacheSize <= 0 {
		config.HardMaxCacheSize = 256
	}

	bigCacheConfig := bigcache.DefaultConfig(config.LifeWindow)
	bigCacheConfig.Shards = config.Shards
	bigCacheConfig.CleanWindow = config.CleanWindow
	bigCacheConfig.MaxEntriesInWindow = config.MaxEntriesInWindow
	bigCacheConfig.MaxEntrySize = config.MaxEntrySize
	bigCacheConfig.HardMaxCacheSize = config.HardMaxCacheSize
	bigCacheConfig.Verbose = config.Verbose

	client, err := bigcache.NewBigCache(bigCacheConfig)
	if err != nil {
		return nil, fmt.Errorf("创建 BigCache 失败: %w", err)
	}

	adapter := &BigCacheAdapter{
		client: client,
		config: config,
	}

	return adapter, nil
}

// Set 设置缓存，支持可选的过期时间（秒）
func (a *BigCacheAdapter) Set(key string, value any, ttl ...int) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	a.stats.totalSets.Add(1)

	entry := &cacheEntry{
		Value: value,
	}

	if len(ttl) > 0 && ttl[0] > 0 {
		entry.ExpireTime = time.Now().Add(time.Duration(ttl[0]) * time.Second).Unix()
	}

	data, err := sonic.Marshal(entry)
	if err != nil {
		return fmt.Errorf("序列化缓存条目失败: %w", err)
	}

	if err := a.client.Set(key, data); err != nil {
		return fmt.Errorf("写入 BigCache 失败: %w", err)
	}

	// 累加已使用字节数（BigCache Set 是覆盖写，需要先减旧值再加新值）
	a.updateUsedBytesOnSet(key, len(data))

	return nil
}

// SetWithPolicy 使用指定过期策略设置缓存
// 注意：BigCache 适配器仅支持相对过期策略，其他策略按相对过期处理
func (a *BigCacheAdapter) SetWithPolicy(key string, value any, policy ExpirePolicy, ttl ...int) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	a.stats.totalSets.Add(1)

	entry := &cacheEntry{
		Value: value,
	}

	if len(ttl) > 0 && ttl[0] > 0 {
		switch policy {
		case ExpirePolicyAbsolute:
			// ttl 作为绝对时间戳
			entry.ExpireTime = int64(ttl[0])
		case ExpirePolicyRelative, ExpirePolicySliding:
			entry.ExpireTime = time.Now().Add(time.Duration(ttl[0]) * time.Second).Unix()
		}
	}

	data, err := sonic.Marshal(entry)
	if err != nil {
		return fmt.Errorf("序列化缓存条目失败: %w", err)
	}

	if err := a.client.Set(key, data); err != nil {
		return fmt.Errorf("写入 BigCache 失败: %w", err)
	}

	a.updateUsedBytesOnSet(key, len(data))

	return nil
}

// SetWithDuration 设置缓存，使用 time.Duration 作为过期时间
// 保留完整的 Duration 精度（包括亚秒级），避免秒级截断
// expiration <= 0 表示永不过期
func (a *BigCacheAdapter) SetWithDuration(key string, value any, expiration time.Duration) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	a.stats.totalSets.Add(1)

	entry := &cacheEntry{
		Value: value,
	}

	if expiration > 0 {
		entry.ExpireTime = time.Now().Add(expiration).Unix()
	}

	data, err := sonic.Marshal(entry)
	if err != nil {
		return fmt.Errorf("序列化缓存条目失败: %w", err)
	}

	if err := a.client.Set(key, data); err != nil {
		return fmt.Errorf("写入 BigCache 失败: %w", err)
	}

	a.updateUsedBytesOnSet(key, len(data))

	return nil
}

// SetBatch 批量设置缓存
func (a *BigCacheAdapter) SetBatch(items map[string]any, ttl ...int) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	for key, value := range items {
		if err := a.Set(key, value, ttl...); err != nil {
			return err
		}
	}
	return nil
}

// Get 获取缓存值
func (a *BigCacheAdapter) Get(key string) (any, bool) {
	if a.closed.Load() {
		return nil, false
	}

	a.stats.totalGets.Add(1)

	data, err := a.client.Get(key)
	if err != nil {
		a.stats.missCount.Add(1)
		return nil, false
	}

	var entry cacheEntry
	if err := sonic.Unmarshal(data, &entry); err != nil {
		a.stats.missCount.Add(1)
		return nil, false
	}
	// 检查 per-key TTL 过期
	if entry.isExpired() {
		a.stats.missCount.Add(1)
		// 同步删除过期条目，避免高并发下 goroutine 泄漏
		// BigCache Delete 是 O(1) 操作，开销极低
		a.stats.usedBytes.Add(-int64(len(data)))
		_ = a.client.Delete(key)
		return nil, false
	}

	a.stats.hitCount.Add(1)
	return entry.Value, true
}

// GetBatch 批量获取缓存值
func (a *BigCacheAdapter) GetBatch(keys []string) map[string]any {
	if a.closed.Load() {
		return nil
	}

	result := make(map[string]any, len(keys))
	for _, key := range keys {
		if val, ok := a.Get(key); ok {
			result[key] = val
		}
	}
	return result
}

// Delete 删除缓存
func (a *BigCacheAdapter) Delete(key string) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	a.stats.totalDeletes.Add(1)

	// 先读取旧值大小用于 usedBytes 扣减
	if oldData, err := a.client.Get(key); err == nil {
		a.stats.usedBytes.Add(-int64(len(oldData)))
	}

	return a.client.Delete(key)
}

// DeleteBatch 批量删除缓存
func (a *BigCacheAdapter) DeleteBatch(keys []string) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	for _, key := range keys {
		a.stats.totalDeletes.Add(1)
		if err := a.client.Delete(key); err != nil {
			// 忽略 key 不存在的错误
			if err != bigcache.ErrEntryNotFound {
				return err
			}
		}
	}
	return nil
}

// Clear 清空所有缓存
func (a *BigCacheAdapter) Clear() error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}
	err := a.client.Reset()
	if err == nil {
		a.stats.usedBytes.Store(0)
	}
	return err
}

// Has 检查缓存键是否存在
func (a *BigCacheAdapter) Has(key string) bool {
	_, ok := a.Get(key)
	return ok
}

// Flush 手动刷盘（BigCache 无需刷盘，保持接口兼容）
func (a *BigCacheAdapter) Flush() error {
	return nil
}

// Close 关闭缓存
func (a *BigCacheAdapter) Close() error {
	if !a.closed.CompareAndSwap(false, true) {
		return nil
	}
	return a.client.Close()
}

// Count 返回缓存项数量
func (a *BigCacheAdapter) Count() int {
	if a.closed.Load() {
		return 0
	}
	return a.client.Len()
}

// Keys 返回所有缓存键
func (a *BigCacheAdapter) Keys() []string {
	if a.closed.Load() {
		return nil
	}

	keys := make([]string, 0)
	iterator := a.client.Iterator()
	for iterator.SetNext() {
		current, err := iterator.Value()
		if err != nil {
			break
		}
		keys = append(keys, current.Key())
	}
	return keys
}

// GetStats 获取缓存统计信息
func (a *BigCacheAdapter) GetStats() *BigCacheStats {
	totalGets := a.stats.totalGets.Load()
	hitCount := a.stats.hitCount.Load()

	var hitRate float64
	if totalGets > 0 {
		hitRate = float64(hitCount) / float64(totalGets)
	}

	return &BigCacheStats{
		ItemCount:    a.Count(),
		HitCount:     hitCount,
		MissCount:    a.stats.missCount.Load(),
		TotalGets:    totalGets,
		TotalSets:    a.stats.totalSets.Load(),
		TotalDeletes: a.stats.totalDeletes.Load(),
		HitRate:      hitRate,
		Capacity:     a.config.HardMaxCacheSize * 1024 * 1024,
		UsedSize:     a.stats.usedBytes.Load(),
	}
}

// updateUsedBytesOnSet BigCache Set 是覆盖写，需要先减旧值再加新值
func (a *BigCacheAdapter) updateUsedBytesOnSet(key string, newDataLen int) {
	// 先尝试读取旧值大小
	if oldData, err := a.client.Get(key); err == nil {
		a.stats.usedBytes.Add(int64(newDataLen - len(oldData)))
	} else {
		// key 不存在，直接加新值
		a.stats.usedBytes.Add(int64(newDataLen))
	}
}

// GetString 获取字符串类型缓存
func (a *BigCacheAdapter) GetString(key string) (string, bool) {
	if val, ok := a.Get(key); ok {
		if str, ok := val.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetInt 获取整数类型缓存
func (a *BigCacheAdapter) GetInt(key string) (int, bool) {
	if val, ok := a.Get(key); ok {
		switch v := val.(type) {
		case int:
			return v, true
		case int64:
			return int(v), true
		case float64:
			return int(v), true
		}
	}
	return 0, false
}

// GetMap 获取 map 类型缓存
func (a *BigCacheAdapter) GetMap(key string) (map[string]any, bool) {
	if val, ok := a.Get(key); ok {
		if mapVal, ok := val.(map[string]any); ok {
			return mapVal, true
		}
	}
	return nil, false
}

// GetSlice 获取切片类型缓存
func (a *BigCacheAdapter) GetSlice(key string) ([]any, bool) {
	if val, ok := a.Get(key); ok {
		if s, ok := val.([]any); ok {
			return s, true
		}
	}
	return nil, false
}
