// Package cache 提供缓存外观接口和工厂方法
// 支持三种缓存后端：Ristretto、BigCache、mmap（SharedCache）
// 通过 NewCache() 工厂方法根据配置类型自动创建对应的缓存实例
package cache

import (
	"fmt"
	"time"
)

// CacheType 缓存类型
type CacheType string

const (
	// TypeRistretto Ristretto 缓存类型（推荐）
	// TinyLFU 淘汰策略 + MaxCost 硬限制防 OOM + 原生 per-key TTL
	TypeRistretto CacheType = "ristretto"
	// TypeBigCache BigCache 缓存类型
	// 高性能进程内缓存，零 GC 压力，适用于高并发场景
	TypeBigCache CacheType = "bigcache"
	// TypeMmap mmap 共享缓存类型
	// 基于内存映射文件，支持进程间共享和持久化
	TypeMmap CacheType = "mmap"
)

// Cache 统一缓存接口
// 所有缓存后端（Ristretto / BigCache / SharedCache）均实现此接口
// 使用方通过此接口操作缓存，无需关心底层实现
type Cache interface {
	// Set 设置缓存，支持可选的过期时间（秒）
	// key: 缓存键
	// value: 缓存值
	// ttl: 过期时间（秒），0 表示永不过期
	Set(key string, value any, ttl ...int) error

	// SetWithDuration 设置缓存，使用 time.Duration 作为过期时间
	// key: 缓存键
	// value: 缓存值
	// expiration: 过期时间，0 表示永不过期
	SetWithDuration(key string, value any, expiration time.Duration) error

	// SetWithPolicy 使用指定过期策略设置缓存
	// key: 缓存键
	// value: 缓存值
	// policy: 过期策略（ExpirePolicyAbsolute/ExpirePolicyRelative/ExpirePolicySliding）
	// ttl: 过期时间（秒），对于绝对过期策略，ttl为Unix时间戳
	SetWithPolicy(key string, value any, policy ExpirePolicy, ttl ...int) error

	// SetBatch 批量设置缓存
	// items: key-value映射
	// ttl: 过期时间（秒），0 表示永不过期
	SetBatch(items map[string]any, ttl ...int) error

	// Get 获取缓存值
	// key: 缓存键
	// 返回值: 缓存值和是否存在的标志
	Get(key string) (any, bool)

	// GetBatch 批量获取缓存值
	// keys: 缓存键列表
	// 返回值: key-value映射（仅包含存在的键）
	GetBatch(keys []string) map[string]any

	// Delete 删除缓存
	Delete(key string) error

	// DeleteBatch 批量删除缓存
	DeleteBatch(keys []string) error

	// Clear 清空所有缓存
	Clear() error

	// Has 检查缓存键是否存在
	Has(key string) bool

	// Flush 手动刷盘
	Flush() error

	// Close 关闭缓存
	Close() error

	// Count 返回缓存项数量
	Count() int

	// Keys 返回所有缓存键
	Keys() []string

	// GetString 获取字符串类型缓存
	GetString(key string) (string, bool)

	// GetInt 获取整数类型缓存
	GetInt(key string) (int, bool)

	// GetMap 获取map类型缓存
	GetMap(key string) (map[string]any, bool)

	// GetSlice 获取切片类型缓存
	GetSlice(key string) ([]any, bool)
}

// CacheConfig 缓存配置（工厂方法使用）
type CacheConfig struct {
	// Type 缓存类型: ristretto | bigcache | mmap
	Type CacheType

	// Ristretto 配置（当 Type=ristretto 时生效）
	Ristretto RistrettoConfig

	// BigCache 配置（当 Type=bigcache 时生效）
	BigCache BigCacheConfig

	// SharedCache / mmap 配置（当 Type=mmap 时生效）
	SharedCache SharedCacheConfig
}

// NewCache 工厂方法：根据配置类型创建对应的缓存实例
// 返回 Cache 接口，使用方无需关心底层实现
func NewCache(config CacheConfig) (Cache, error) {
	switch config.Type {
	case TypeRistretto:
		return NewRistretto(config.Ristretto)
	case TypeBigCache:
		return NewBigCache(config.BigCache)
	case TypeMmap:
		return NewSharedCache(config.SharedCache)
	default:
		return nil, fmt.Errorf("不支持的缓存类型: %s", config.Type)
	}
}

// 确保 Ristretto / BigCache / SharedCache 均实现 Cache 接口
var (
	_ Cache = (*Ristretto)(nil)
	_ Cache = (*BigCache)(nil)
	_ Cache = (*SharedCache)(nil)
)
