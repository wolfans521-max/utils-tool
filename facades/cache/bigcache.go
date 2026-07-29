// Package cache 提供缓存外观接口
// 包含基于 BigCache 的高性能实现和基于 mmap 的共享内存实现
package cache

import (
	"time"

	internalCache "git.intra.weibo.com/search_fe/wbutil-go/internal/cache"
)

// BigCache 基于 BigCache 的高性能缓存外观
// 适用于高并发场景，零 GC 压力
type BigCache struct {
	impl *internalCache.BigCacheAdapter
}

// BigCacheConfig BigCache 配置选项
type BigCacheConfig struct {
	Shards             int           // 分片数量，必须是2的幂，默认64
	LifeWindow         time.Duration // 全局最大TTL，默认2小时
	CleanWindow        time.Duration // 清理间隔，默认30秒
	MaxEntriesInWindow int           // 窗口内最大条目数，默认10000
	MaxEntrySize       int           // 单条目最大字节数，默认1MB
	HardMaxCacheSize   int           // 缓存最大大小（MB），默认256MB
	Verbose            bool          // 是否启用详细日志
}

// NewBigCache 创建新的 BigCache 实例
func NewBigCache(config BigCacheConfig) (*BigCache, error) {
	impl, err := internalCache.NewBigCacheAdapter(internalCache.BigCacheConfig{
		Shards:             config.Shards,
		LifeWindow:         config.LifeWindow,
		CleanWindow:        config.CleanWindow,
		MaxEntriesInWindow: config.MaxEntriesInWindow,
		MaxEntrySize:       config.MaxEntrySize,
		HardMaxCacheSize:   config.HardMaxCacheSize,
		Verbose:            config.Verbose,
	})
	if err != nil {
		return nil, err
	}
	return &BigCache{impl: impl}, nil
}

// Set 设置缓存，支持可选的过期时间（秒）
// key: 缓存键
// value: 缓存值
// ttl: 过期时间（秒），0 表示永不过期
func (c *BigCache) Set(key string, value any, ttl ...int) error {
	return c.impl.Set(key, value, ttl...)
}

// SetWithPolicy 使用指定过期策略设置缓存
// key: 缓存键
// value: 缓存值
// policy: 过期策略（ExpirePolicyAbsolute/ExpirePolicyRelative/ExpirePolicySliding）
// ttl: 过期时间（秒），对于绝对过期策略，ttl为Unix时间戳
func (c *BigCache) SetWithPolicy(key string, value any, policy ExpirePolicy, ttl ...int) error {
	return c.impl.SetWithPolicy(key, value, policy, ttl...)
}

// SetBatch 批量设置缓存
// items: key-value映射
// ttl: 过期时间（秒），0 表示永不过期
func (c *BigCache) SetBatch(items map[string]any, ttl ...int) error {
	return c.impl.SetBatch(items, ttl...)
}

// Get 获取缓存值
// key: 缓存键
// 返回值: 缓存值和是否存在的标志
func (c *BigCache) Get(key string) (any, bool) {
	return c.impl.Get(key)
}

// GetBatch 批量获取缓存值
// keys: 缓存键列表
// 返回值: key-value映射（仅包含存在的键）
func (c *BigCache) GetBatch(keys []string) map[string]any {
	return c.impl.GetBatch(keys)
}

// Delete 删除缓存
func (c *BigCache) Delete(key string) error {
	return c.impl.Delete(key)
}

// DeleteBatch 批量删除缓存
// keys: 缓存键列表
func (c *BigCache) DeleteBatch(keys []string) error {
	return c.impl.DeleteBatch(keys)
}

// Clear 清空所有缓存
func (c *BigCache) Clear() error {
	return c.impl.Clear()
}

// Has 检查缓存键是否存在
func (c *BigCache) Has(key string) bool {
	return c.impl.Has(key)
}

// Flush 手动刷盘（BigCache 无需刷盘，保持接口兼容）
func (c *BigCache) Flush() error {
	return c.impl.Flush()
}

// Close 关闭缓存
func (c *BigCache) Close() error {
	return c.impl.Close()
}

// Count 返回缓存项数量
func (c *BigCache) Count() int {
	return c.impl.Count()
}

// Keys 返回所有缓存键
func (c *BigCache) Keys() []string {
	return c.impl.Keys()
}

// GetStats 获取缓存统计信息
func (c *BigCache) GetStats() *internalCache.BigCacheStats {
	return c.impl.GetStats()
}

// SetWithDuration 设置缓存，使用time.Duration作为过期时间
// key: 缓存键
// value: 缓存值
// expiration: 过期时间，0 表示永不过期
func (c *BigCache) SetWithDuration(key string, value any, expiration time.Duration) error {
	return c.impl.SetWithDuration(key, value, expiration)
}

// GetString 获取字符串类型缓存
func (c *BigCache) GetString(key string) (string, bool) {
	return c.impl.GetString(key)
}

// GetInt 获取整数类型缓存
func (c *BigCache) GetInt(key string) (int, bool) {
	return c.impl.GetInt(key)
}

// GetMap 获取map类型缓存
func (c *BigCache) GetMap(key string) (map[string]any, bool) {
	return c.impl.GetMap(key)
}

// GetSlice 获取切片类型缓存
func (c *BigCache) GetSlice(key string) ([]any, bool) {
	return c.impl.GetSlice(key)
}
