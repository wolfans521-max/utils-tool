// Package cache 提供缓存外观接口
// 包含基于 Ristretto 的高性能实现（推荐）、基于 BigCache 的高性能实现和基于 mmap 的共享内存实现
package cache

import (
	"time"

	internalCache "git.intra.weibo.com/search_fe/wbutil-go/internal/cache"
)

// Ristretto 基于 Ristretto 的高性能缓存外观
// 适用于读多写少、需要智能淘汰策略的场景
// 核心优势：TinyLFU 淘汰策略 + MaxCost 硬上限 + 原生 per-key TTL
type Ristretto struct {
	impl *internalCache.RistrettoAdapter
}

// RistrettoConfig Ristretto 配置选项
type RistrettoConfig struct {
	MaxCacheItems      int64 // 最大缓存条目数，建议设为 key 数量的 2-3 倍，默认 1500
	MaxCost            int64 // 最大缓存字节数，cost=len(data) 时即为最大内存字节数，默认 1GB
	BufferItems        int64 // 写入缓冲区大小，默认 64
	SampleSize         int64 // 淘汰策略内部采样数，默认 4
	Metrics            bool  // 是否启用指标采集，默认 true
	IgnoreInternalCost bool  // 必须为 true！忽略内部存储开销，否则实际 cost=len(data)+56，缓存被过早淘汰
}

// NewRistretto 创建新的 Ristretto 实例
func NewRistretto(config RistrettoConfig) (*Ristretto, error) {
	impl, err := internalCache.NewRistrettoAdapter(internalCache.RistrettoConfig{
		MaxCacheItems:      config.MaxCacheItems,
		MaxCost:            config.MaxCost,
		BufferItems:        config.BufferItems,
		SampleSize:         config.SampleSize,
		Metrics:            config.Metrics,
		IgnoreInternalCost: config.IgnoreInternalCost,
	})
	if err != nil {
		return nil, err
	}
	return &Ristretto{impl: impl}, nil
}

// Set 设置缓存，支持可选的过期时间（秒）
// key: 缓存键
// value: 缓存值
// ttl: 过期时间（秒），0 表示永不过期
func (c *Ristretto) Set(key string, value any, ttl ...int) error {
	return c.impl.Set(key, value, ttl...)
}

// SetWithPolicy 使用指定过期策略设置缓存
// key: 缓存键
// value: 缓存值
// policy: 过期策略（ExpirePolicyAbsolute/ExpirePolicyRelative/ExpirePolicySliding）
// ttl: 过期时间（秒），对于绝对过期策略，ttl为Unix时间戳
func (c *Ristretto) SetWithPolicy(key string, value any, policy ExpirePolicy, ttl ...int) error {
	return c.impl.SetWithPolicy(key, value, policy, ttl...)
}

// SetBatch 批量设置缓存
// items: key-value映射
// ttl: 过期时间（秒），0 表示永不过期
func (c *Ristretto) SetBatch(items map[string]any, ttl ...int) error {
	return c.impl.SetBatch(items, ttl...)
}

// Get 获取缓存值
// key: 缓存键
// 返回值: 缓存值和是否存在的标志
func (c *Ristretto) Get(key string) (any, bool) {
	return c.impl.Get(key)
}

// GetBatch 批量获取缓存值
// keys: 缓存键列表
// 返回值: key-value映射（仅包含存在的键）
func (c *Ristretto) GetBatch(keys []string) map[string]any {
	return c.impl.GetBatch(keys)
}

// Delete 删除缓存
func (c *Ristretto) Delete(key string) error {
	return c.impl.Delete(key)
}

// DeleteBatch 批量删除缓存
// keys: 缓存键列表
func (c *Ristretto) DeleteBatch(keys []string) error {
	return c.impl.DeleteBatch(keys)
}

// Clear 清空所有缓存
func (c *Ristretto) Clear() error {
	return c.impl.Clear()
}

// Has 检查缓存键是否存在
func (c *Ristretto) Has(key string) bool {
	return c.impl.Has(key)
}

// Flush 手动刷盘（Ristretto 无需刷盘，保持接口兼容）
func (c *Ristretto) Flush() error {
	return c.impl.Flush()
}

// Close 关闭缓存
func (c *Ristretto) Close() error {
	return c.impl.Close()
}

// Count 返回缓存项数量
func (c *Ristretto) Count() int {
	return c.impl.Count()
}

// Keys 返回所有缓存键
// 注意：Ristretto 不提供遍历所有 key 的能力，返回 nil
func (c *Ristretto) Keys() []string {
	return c.impl.Keys()
}

// GetStats 获取缓存统计信息
func (c *Ristretto) GetStats() *internalCache.RistrettoStats {
	return c.impl.GetStats()
}

// SetWithDuration 设置缓存，使用time.Duration作为过期时间
// key: 缓存键
// value: 缓存值
// expiration: 过期时间，0 表示永不过期
func (c *Ristretto) SetWithDuration(key string, value any, expiration time.Duration) error {
	return c.impl.SetWithDuration(key, value, expiration)
}

// GetString 获取字符串类型缓存
func (c *Ristretto) GetString(key string) (string, bool) {
	return c.impl.GetString(key)
}

// GetInt 获取整数类型缓存
func (c *Ristretto) GetInt(key string) (int, bool) {
	return c.impl.GetInt(key)
}

// GetMap 获取map类型缓存
func (c *Ristretto) GetMap(key string) (map[string]any, bool) {
	return c.impl.GetMap(key)
}

// GetSlice 获取切片类型缓存
func (c *Ristretto) GetSlice(key string) ([]any, bool) {
	return c.impl.GetSlice(key)
}
