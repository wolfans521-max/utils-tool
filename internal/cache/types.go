// Package cache 提供缓存核心实现的公共类型和工具函数
package cache

import "time"

// cacheEntry 缓存条目，包含过期时间
// Ristretto 和 BigCache 适配器共用此结构体
// 序列化后存储为 []byte，缓存层内部不持有 any 指针，避免 GC 扫描压力
type cacheEntry struct {
	Value      any   `json:"value"`
	ExpireTime int64 `json:"expire_time"` // Unix 时间戳，0 表示永不过期
}

// isExpired 检查缓存条目是否过期
func (e *cacheEntry) isExpired() bool {
	if e.ExpireTime == 0 {
		return false
	}
	return time.Now().Unix() >= e.ExpireTime
}
