// Package cache 提供缓存外观接口
// 本文件实现基于 mmap 的 SharedCache 外观（跨平台共享内存，支持进程间共享）
// 工厂入口请参考 factory.go，其他缓存实现请参考 ristretto.go / bigcache.go
package cache

import (
	"time"

	internalCache "git.intra.weibo.com/search_fe/wbutil-go/internal/cache"
)

// ExpirePolicy 过期策略类型
type ExpirePolicy = internalCache.ExpirePolicy

// 过期策略常量
const (
	// ExpirePolicyAbsolute 绝对过期（到达指定时间点过期）
	ExpirePolicyAbsolute = internalCache.ExpirePolicyAbsolute
	// ExpirePolicyRelative 相对过期（从设置时刻开始计算TTL）
	ExpirePolicyRelative = internalCache.ExpirePolicyRelative
	// ExpirePolicySliding 滑动过期（每次访问重置TTL）
	ExpirePolicySliding = internalCache.ExpirePolicySliding
)

// Stats 缓存统计信息
type Stats = internalCache.Stats

// SharedCache 共享缓存外观接口（跨平台，支持进程间共享）
type SharedCache struct {
	impl *internalCache.SharedCache
}

// SharedCacheConfig 共享缓存配置选项
type SharedCacheConfig struct {
	Name            string // 缓存名称（用于生成文件名）
	Size            int    // 缓存文件大小（字节）
	Dir             string // 缓存文件目录（可选，默认使用临时目录）
	ShardCount      int    // 分片数量（可选，默认32）
	LocalCacheTTLMs int    // 本地缓存TTL（毫秒，可选，默认1000ms）
	AsyncFlush      bool   // 是否启用异步刷盘（可选，默认false）
	AutoExpand      bool   // 是否启用自动扩容（可选，默认true）
	EnableCRC       bool   // 是否启用CRC校验（可选，默认true）
}

// GetInstance GetSharedInstance 获取共享缓存单例实例
// 使用示例
// c,_ := cache.GetInstance("test_memoryname")
// c.Set("key", "value")
func GetInstance(name string) (*SharedCache, error) {
	impl, err := internalCache.GetInstance(name)
	if err != nil {
		return nil, err
	}
	return &SharedCache{impl: impl}, nil
}

// NewSharedCache 创建新的共享缓存实例
func NewSharedCache(config SharedCacheConfig) (*SharedCache, error) {
	impl, err := internalCache.NewSharedCache(internalCache.SharedConfig{
		Name:            config.Name,
		Size:            config.Size,
		Dir:             config.Dir,
		ShardCount:      config.ShardCount,
		LocalCacheTTLMs: config.LocalCacheTTLMs,
		AsyncFlush:      config.AsyncFlush,
		AutoExpand:      config.AutoExpand,
		EnableCRC:       config.EnableCRC,
	})
	if err != nil {
		return nil, err
	}
	return &SharedCache{impl: impl}, nil
}

// Set 设置缓存，支持可选的过期时间（秒）
// key: 缓存键
// value: 缓存值
// ttl: 过期时间（秒），0 表示永不过期
func (c *SharedCache) Set(key string, value any, ttl ...int) error {
	return c.impl.Set(key, value, ttl...)
}

// SetWithPolicy 使用指定过期策略设置缓存
// key: 缓存键
// value: 缓存值
// policy: 过期策略（ExpirePolicyAbsolute/ExpirePolicyRelative/ExpirePolicySliding）
// ttl: 过期时间（秒），对于绝对过期策略，ttl为Unix时间戳
func (c *SharedCache) SetWithPolicy(key string, value any, policy ExpirePolicy, ttl ...int) error {
	return c.impl.SetWithPolicy(key, value, policy, ttl...)
}

// SetBatch 批量设置缓存
// items: key-value映射
// ttl: 过期时间（秒），0 表示永不过期
func (c *SharedCache) SetBatch(items map[string]any, ttl ...int) error {
	return c.impl.SetBatch(items, ttl...)
}

// Get 获取缓存值
// key: 缓存键
// 返回值: 缓存值和是否存在的标志
func (c *SharedCache) Get(key string) (any, bool) {
	return c.impl.Get(key)
}

// GetBatch 批量获取缓存值
// keys: 缓存键列表
// 返回值: key-value映射（仅包含存在的键）
func (c *SharedCache) GetBatch(keys []string) map[string]any {
	return c.impl.GetBatch(keys)
}

// Delete 删除缓存
func (c *SharedCache) Delete(key string) error {
	return c.impl.Delete(key)
}

// DeleteBatch 批量删除缓存
// keys: 缓存键列表
func (c *SharedCache) DeleteBatch(keys []string) error {
	return c.impl.DeleteBatch(keys)
}

// Clear 清空所有缓存
func (c *SharedCache) Clear() error {
	return c.impl.Clear()
}

// Has 检查缓存键是否存在
func (c *SharedCache) Has(key string) bool {
	return c.impl.Has(key)
}

// Flush 手动刷盘（用于异步模式）
func (c *SharedCache) Flush() error {
	return c.impl.Flush()
}

// Close 关闭缓存
func (c *SharedCache) Close() error {
	return c.impl.Close()
}

// Unlink 删除缓存文件
func (c *SharedCache) Unlink() error {
	return c.impl.Unlink()
}

// Count 返回缓存项数量
func (c *SharedCache) Count() int {
	return c.impl.Count()
}

// Keys 返回所有缓存键
func (c *SharedCache) Keys() []string {
	return c.impl.Keys()
}

// GetFilePath 返回缓存文件路径
func (c *SharedCache) GetFilePath() string {
	return c.impl.GetFilePath()
}

// GetStats 获取缓存统计信息
func (c *SharedCache) GetStats() *Stats {
	return c.impl.GetStats()
}

// Expand 手动扩容
// newSize: 新的缓存大小（字节）
func (c *SharedCache) Expand(newSize int) error {
	return c.impl.Expand(newSize)
}

// SetWithDuration 设置缓存，使用time.Duration作为过期时间
// key: 缓存键
// value: 缓存值
// expiration: 过期时间，0 表示永不过期
func (c *SharedCache) SetWithDuration(key string, value any, expiration time.Duration) error {
	if expiration <= 0 {
		return c.impl.Set(key, value)
	}
	return c.impl.Set(key, value, int(expiration.Seconds()))
}

// GetString 获取字符串类型缓存
func (c *SharedCache) GetString(key string) (string, bool) {
	if val, ok := c.impl.Get(key); ok {
		if str, ok := val.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetInt 获取整数类型缓存
func (c *SharedCache) GetInt(key string) (int, bool) {
	if val, ok := c.impl.Get(key); ok {
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

// GetMap 获取map类型缓存
func (c *SharedCache) GetMap(key string) (map[string]any, bool) {
	if val, ok := c.impl.Get(key); ok {
		if mapVal, ok := val.(map[string]any); ok {
			return mapVal, true
		}
	}
	return nil, false
}

// GetSlice 获取切片类型缓存
func (c *SharedCache) GetSlice(key string) ([]any, bool) {
	if val, ok := c.impl.Get(key); ok {
		if s, ok := val.([]any); ok {
			return s, true
		}
	}
	return nil, false
}
