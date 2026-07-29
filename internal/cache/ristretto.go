// Package cache 提供基于 Ristretto 的高性能缓存实现
// 适用于读多写少、需要智能淘汰策略的场景，替代基于 mmap 的 SharedCache
//
// Ristretto 核心优势：
//  1. TinyLFU 淘汰策略 — 在内存满时自动淘汰冷数据，与 ngx.shared.DICT 的 LRU 行为对齐
//  2. MaxCost 硬上限 — 从根本上防止 OOM
//  3. 原生 per-key TTL — 通过 SetWithTTL 直接设置，Ristretto 内部管理过期
//  4. 零 GC 压力 — 只存储 []byte
//
// 性能优化要点（面向 C 端高并发场景）：
//   - Set 后不调用 Wait()，避免将异步写入变为同步阻塞（降级缓存场景可容忍微秒级延迟）
//   - 统一存储 []byte，Ristretto 内部不持有 any 指针，GC 零压力
//   - Get 时才反序列化，减少热路径 CPU 开销
//   - 维护 sync.Map keySet 跟踪已写入的 key，解决 Ristretto 不暴露 key 列表的问题
//
// 过期策略说明：
//   - Ristretto 通过 SetWithTTL 原生支持 per-key TTL，过期由 Ristretto 内部管理
//   - 不在应用层做过期检查，避免额外的 ExpireTime 字段和反序列化开销
//   - TTL 到期后 Ristretto 在内部淘汰时自动清理，Get 可能返回已过期但尚未被淘汰的条目
//   - 对于降级缓存场景，短暂返回过期数据是可接受的（降级数据本身就有时效性容忍）
//
// cost 策略说明：
//   - cost 设为 len(data)（序列化后的字节数），MaxCost 设为 1GB（1073741824）
//   - 这样 Ristretto 按实际内存使用量淘汰，总缓存内存不超过 MaxCost，防止 OOM
//   - TinyLFU 淘汰策略的准入决策基于 hit frequency（访问频率），与 cost 大小无关
//   - cost 只影响淘汰时需要释放多少空间，不影响准入判断
//   - 必须配合 IgnoreInternalCost=true 使用！否则 Ristretto 会在 cost 上加 itemSize（56字节），
//     导致实际 cost = len(data) + 56，使缓存被过早淘汰
package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/dgraph-io/ristretto/v2"
)

// RistrettoConfig Ristretto 配置
type RistrettoConfig struct {
	// MaxCacheItems 最大缓存条目数，建议设为 key 数量的 2-3 倍，默认 1500
	MaxCacheItems int64
	// MaxCost 最大缓存字节数，即缓存内存硬上限，默认 1GB（1073741824）
	// cost=len(data) 时，MaxCost 的语义为"最大字节数"，总缓存内存不超过此值
	MaxCost int64
	// BufferItems 写入缓冲区大小，默认 64
	BufferItems int64
	// SampleSize 淘汰策略内部采样数，默认 4
	SampleSize int64
	// Metrics 是否启用指标采集，默认 true
	Metrics bool
	// IgnoreInternalCost 是否忽略 Ristretto 内部存储开销，默认 true
	// 必须设为 true！原因：Ristretto v2 在 processItems() 中默认会在 cost 上加 itemSize（56字节）
	// 作为内部存储开销，导致实际 cost = len(data) + 56，使缓存被过早淘汰。
	// 设为 true 后，cost 保持为 len(data)，MaxCost 即为最大字节数，行为与预期一致。
	IgnoreInternalCost bool
}

// DefaultRistrettoConfig 返回默认 Ristretto 配置
func DefaultRistrettoConfig() RistrettoConfig {
	return RistrettoConfig{
		MaxCacheItems:      1500,
		MaxCost:            1024 * 1024 * 1024, // 1GB（cost=len(data) 时 MaxCost 为最大字节数）
		BufferItems:        64,
		SampleSize:         4,
		Metrics:            true,
		IgnoreInternalCost: true, // 必须为 true！忽略内部存储开销，否则实际 cost=len(data)+56
	}
}

// RistrettoStats Ristretto 统计信息
type RistrettoStats struct {
	ItemCount    int     // 缓存项数量（来自 Ristretto Metrics: KeysAdded - KeysEvicted）
	HitCount     int64   // 命中次数
	MissCount    int64   // 未命中次数
	TotalGets    int64   // 总Get次数
	TotalSets    int64   // 总Set次数
	TotalDeletes int64   // 总Delete次数
	HitRate      float64 // 命中率
	MaxCost      int64   // 最大字节数（cost=len(data) 时即为最大内存字节数）
	UsedSize     int64   // 已使用字节数（来自 Ristretto Metrics: CostAdded - CostEvicted）
}

// ristrettoMetricsUsedSize 从 Ristretto Metrics 计算实际已使用字节数
// 公式：CostAdded - CostEvicted
// CostAdded: 所有成功 Set 的 cost 之和（包括 itemUpdate 的 cost 差值）
// CostEvicted: 所有被淘汰的 cost 之和（TinyLFU淘汰 + 显式Delete + TTL过期清理）
func ristrettoMetricsUsedSize(m *ristretto.Metrics) int64 {
	if m == nil {
		return 0
	}
	used := int64(m.CostAdded() - m.CostEvicted())
	if used < 0 {
		used = 0 // 防御性处理
	}
	return used
}

// ristrettoMetricsCount 从 Ristretto Metrics 计算当前缓存项数量
// 公式：KeysAdded - KeysEvicted
// KeysAdded: 所有成功新增的 key 数量（不含 itemUpdate）
// KeysEvicted: 所有被淘汰的 key 数量（TinyLFU淘汰 + 显式Delete + TTL过期清理）
// 注意：此值在 Delete/Del 后需要 Wait() 才能保证准确，因为 Del 是异步的
func ristrettoMetricsCount(m *ristretto.Metrics) int {
	if m == nil {
		return 0
	}
	count := int64(m.KeysAdded() - m.KeysEvicted())
	if count < 0 {
		count = 0 // 防御性处理
	}
	return int(count)
}

// RistrettoAdapter 基于 ristretto 的缓存适配器
// 实现与 SharedCache 相同的方法集，支持透明切换
type RistrettoAdapter struct {
	client *ristretto.Cache[string, []byte]
	config RistrettoConfig
	closed atomic.Bool

	// keySet 跟踪已写入的 key，弥补 Ristretto 不暴露 key 列表的限制
	// 仅用于 Keys() 方法返回 key 列表，不用于统计（统计由 Ristretto Metrics 负责）
	// 使用 sync.Map 保证并发安全，key 为 string，value 为 struct{}（仅标记存在性）
	keySet sync.Map

	// 统计信息（Ristretto 内置 Metrics 不包含的字段）
	// Count 和 UsedSize 由 Ristretto Metrics 提供（KeysAdded-KeysEvicted, CostAdded-CostEvicted），
	// 不再自维护 keyCount/usedBytes，因为 Ristretto 内部淘汰/TTL清理时无法感知，会导致统计偏大
	stats struct {
		totalGets    atomic.Int64
		totalSets    atomic.Int64
		totalDeletes atomic.Int64
		hitCount     atomic.Int64
		missCount    atomic.Int64
	}
}

var (
	ristrettoInstance *RistrettoAdapter
	ristrettoOnce     sync.Once
	ristrettoInitErr  error
)

// GetRistrettoInstance 获取 Ristretto 单例实例
func GetRistrettoInstance(config RistrettoConfig) (*RistrettoAdapter, error) {
	ristrettoOnce.Do(func() {
		impl, err := NewRistrettoAdapter(config)
		if err != nil {
			ristrettoInitErr = err
			return
		}
		ristrettoInstance = impl
	})
	return ristrettoInstance, ristrettoInitErr
}

// NewRistrettoAdapter 创建新的 Ristretto 适配器实例
func NewRistrettoAdapter(config RistrettoConfig) (*RistrettoAdapter, error) {
	// 设置默认值
	if config.MaxCacheItems <= 0 {
		config.MaxCacheItems = 1500
		if config.MaxCost <= 0 {
			config.MaxCost = 1024 * 1024 * 1024 // 默认最大 1GB（cost=len(data) 时 MaxCost 为最大字节数）
		}
	}
	if config.BufferItems <= 0 {
		config.BufferItems = 64
	}
	if config.SampleSize <= 0 {
		config.SampleSize = 4
	}
	// IgnoreInternalCost 默认必须为 true
	// 原因：Ristretto v2 在 processItems() 中会自动给 cost 加上 itemSize（56字节），
	// 如果 IgnoreInternalCost=false，实际 cost = len(data) + 56，
	// 导致缓存被过早淘汰（例如 MaxCost=1GB 时，实际可用空间被 56 字节/条的额外开销蚕食）。
	// 设为 true 后，cost 保持为 len(data)，MaxCost 即为最大字节数，行为与预期一致。
	config.IgnoreInternalCost = true

	ristrettoConfig := &ristretto.Config[string, []byte]{
		NumCounters:        config.MaxCacheItems * 10, // TinyLFU 计数器数量，推荐 10x MaxCacheItems
		MaxCost:            config.MaxCost,
		BufferItems:        config.BufferItems,
		Metrics:            config.Metrics,
		IgnoreInternalCost: config.IgnoreInternalCost,
	}

	client, err := ristretto.NewCache[string, []byte](ristrettoConfig)
	if err != nil {
		return nil, fmt.Errorf("创建 Ristretto 失败: %w", err)
	}

	adapter := &RistrettoAdapter{
		client: client,
		config: config,
	}

	return adapter, nil
}

// Set 设置缓存，支持可选的过期时间（秒）
func (a *RistrettoAdapter) Set(key string, value any, ttl ...int) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	a.stats.totalSets.Add(1)

	data, err := sonic.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化缓存条目失败: %w", err)
	}

	// cost 设为 len(data)：让 Ristretto 按实际内存使用量淘汰，MaxCost 为最大字节数
	// 详见包注释中的 cost 策略说明
	cost := int64(len(data))

	var admitted bool
	if len(ttl) > 0 && ttl[0] > 0 {
		admitted = a.client.SetWithTTL(key, data, cost, time.Duration(ttl[0])*time.Second)
	} else {
		admitted = a.client.Set(key, data, cost)
	}

	if !admitted {
		// 条目被 TinyLFU 拒绝准入，不计入 keySet
		// 这不是错误，是正常的淘汰行为（缓存满时新条目替换旧条目）
		return nil
	}

	// 不调用 Wait()！
	// Ristretto 的 Set 是异步的（写入内部 buffer 后立即返回）
	// 在 C 端高并发场景下，调用 Wait() 会将异步写入变为同步阻塞，严重拖慢写路径
	// 降级缓存场景：写入后下次请求才读取，Ristretto 内部 buffer 消费是微秒级，
	// 下次请求到来时数据早已可见，无需 Wait()
	// 如果确实需要立即可见（如写后立即读），调用方应使用 SetAndWait 方法

	// 跟踪 key（仅用于 Keys() 方法，统计由 Ristretto Metrics 负责）
	a.keySet.Store(key, struct{}{})

	return nil
}

// SetAndWait 设置缓存并等待写入完成（确保立即可见）
// 仅在写后需要立即读取的场景使用，普通场景使用 Set 即可
func (a *RistrettoAdapter) SetAndWait(key string, value any, ttl ...int) error {
	if err := a.Set(key, value, ttl...); err != nil {
		return err
	}
	a.client.Wait()
	return nil
}

// SetWithPolicy 使用指定过期策略设置缓存
// 注意：Ristretto 适配器仅支持相对过期策略，其他策略按相对过期处理
func (a *RistrettoAdapter) SetWithPolicy(key string, value any, policy ExpirePolicy, ttl ...int) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	a.stats.totalSets.Add(1)

	data, err := sonic.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化缓存条目失败: %w", err)
	}

	// cost 设为 len(data)：让 Ristretto 按实际内存使用量淘汰（同 Set 方法说明）
	cost := int64(len(data))

	var ttlDuration time.Duration
	if len(ttl) > 0 && ttl[0] > 0 {
		switch policy {
		case ExpirePolicyAbsolute:
			// ttl 作为绝对时间戳，计算剩余时间作为 Ristretto TTL
			remaining := time.Until(time.Unix(int64(ttl[0]), 0))
			if remaining > 0 {
				ttlDuration = remaining
			}
		case ExpirePolicyRelative, ExpirePolicySliding:
			ttlDuration = time.Duration(ttl[0]) * time.Second
		}
	}

	var admitted bool
	if ttlDuration > 0 {
		admitted = a.client.SetWithTTL(key, data, cost, ttlDuration)
	} else {
		admitted = a.client.Set(key, data, cost)
	}

	if !admitted {
		return nil
	}

	// 不调用 Wait()，同 Set()

	// 跟踪 key（仅用于 Keys() 方法，统计由 Ristretto Metrics 负责）
	a.keySet.Store(key, struct{}{})

	return nil
}

// SetWithDuration 设置缓存，使用 time.Duration 作为过期时间
// 保留纳秒精度，直接使用 Ristretto 的 SetWithTTL
func (a *RistrettoAdapter) SetWithDuration(key string, value any, expiration time.Duration) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	a.stats.totalSets.Add(1)

	data, err := sonic.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化缓存条目失败: %w", err)
	}

	// cost 设为 len(data)：让 Ristretto 按实际内存使用量淘汰（同 Set 方法说明）
	cost := int64(len(data))

	var admitted bool
	if expiration > 0 {
		admitted = a.client.SetWithTTL(key, data, cost, expiration)
	} else {
		admitted = a.client.Set(key, data, cost)
	}

	if !admitted {
		return nil
	}

	// 不调用 Wait()

	// 跟踪 key（仅用于 Keys() 方法，统计由 Ristretto Metrics 负责）
	a.keySet.Store(key, struct{}{})

	return nil
}

// SetBatch 批量设置缓存
func (a *RistrettoAdapter) SetBatch(items map[string]any, ttl ...int) error {
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
// Ristretto 内部存储 []byte，Get 时反序列化，避免缓存层持有 any 指针导致 GC 压力
// 注意：Ristretto v2 的 Get 会主动检查 TTL 过期（store.go:178），过期条目返回 false
func (a *RistrettoAdapter) Get(key string) (any, bool) {
	if a.closed.Load() {
		return nil, false
	}

	a.stats.totalGets.Add(1)

	data, ok := a.client.Get(key)
	if !ok {
		a.stats.missCount.Add(1)
		// keySet 同步：如果 keySet 中存在但 Ristretto 内部已淘汰/过期，清理 keySet
		// 这解决了 Ristretto TinyLFU 静默淘汰/TTL过期导致 keySet 与实际存储不一致的问题
		a.keySet.Delete(key)
		return nil, false
	}

	var value any
	if err := sonic.Unmarshal(data, &value); err != nil {
		a.stats.missCount.Add(1)
		return nil, false
	}

	a.stats.hitCount.Add(1)
	return value, true
}

// GetBatch 批量获取缓存值
func (a *RistrettoAdapter) GetBatch(keys []string) map[string]any {
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
// 调用 client.Del() 后必须 Wait()，因为 Del 是异步的（写入 setBuf → processItems 消费），
// 只有 Wait() 后 Ristretto Metrics（KeysEvicted/CostEvicted）才会更新，Count()/UsedSize 才准确
func (a *RistrettoAdapter) Delete(key string) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	a.stats.totalDeletes.Add(1)

	a.client.Del(key)
	// Wait() 确保 processItems() 处理完 itemDelete，Metrics 已更新
	a.client.Wait()

	// 从 keySet 中移除
	a.keySet.Delete(key)

	return nil
}

// DeleteBatch 批量删除缓存
func (a *RistrettoAdapter) DeleteBatch(keys []string) error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	for _, key := range keys {
		a.stats.totalDeletes.Add(1)
		a.client.Del(key)
		// 从 keySet 中移除
		a.keySet.Delete(key)
	}
	// 批量删除后统一 Wait()，减少同步开销
	a.client.Wait()

	return nil
}

// Clear 清空所有缓存
func (a *RistrettoAdapter) Clear() error {
	if a.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}
	a.client.Clear()
	// client.Clear() 内部会清空 Metrics，无需手动重置

	// 清空 keySet
	a.keySet.Range(func(key, _ any) bool {
		a.keySet.Delete(key)
		return true
	})

	return nil
}

// Has 检查缓存键是否存在
func (a *RistrettoAdapter) Has(key string) bool {
	_, ok := a.Get(key)
	return ok
}

// Flush 手动刷盘（Ristretto 无需刷盘，保持接口兼容）
func (a *RistrettoAdapter) Flush() error {
	return nil
}

// Close 关闭缓存
func (a *RistrettoAdapter) Close() error {
	if !a.closed.CompareAndSwap(false, true) {
		return nil
	}
	a.client.Close()
	return nil
}

// Count 返回缓存项数量
// 使用 Ristretto 内置 Metrics（KeysAdded - KeysEvicted），确保与实际存储一致
// 注意：Delete/Del 后需要 Wait() 才能保证准确，因为 Del 是异步的
func (a *RistrettoAdapter) Count() int {
	if a.closed.Load() {
		return 0
	}
	return ristrettoMetricsCount(a.client.Metrics)
}

// Keys 返回所有缓存键
// 通过维护的 keySet 返回，弥补 Ristretto 不暴露 key 列表的限制
func (a *RistrettoAdapter) Keys() []string {
	if a.closed.Load() {
		return nil
	}

	keys := make([]string, 0)
	a.keySet.Range(func(key, _ any) bool {
		if k, ok := key.(string); ok {
			keys = append(keys, k)
		}
		return true
	})
	return keys
}

// GetStats 获取缓存统计信息
// UsedSize 和 ItemCount 使用 Ristretto 内置 Metrics，确保与实际存储一致
// 自维护的 usedBytes/keyCount 已移除，因为 Ristretto 内部淘汰/TTL清理时无法感知，会导致统计偏大
func (a *RistrettoAdapter) GetStats() *RistrettoStats {
	totalGets := a.stats.totalGets.Load()
	hitCount := a.stats.hitCount.Load()

	var hitRate float64
	if totalGets > 0 {
		hitRate = float64(hitCount) / float64(totalGets)
	}

	return &RistrettoStats{
		ItemCount:    a.Count(),
		HitCount:     hitCount,
		MissCount:    a.stats.missCount.Load(),
		TotalGets:    totalGets,
		TotalSets:    a.stats.totalSets.Load(),
		TotalDeletes: a.stats.totalDeletes.Load(),
		HitRate:      hitRate,
		MaxCost:      a.config.MaxCost,
		UsedSize:     ristrettoMetricsUsedSize(a.client.Metrics),
	}
}

// GetString 获取字符串类型缓存
func (a *RistrettoAdapter) GetString(key string) (string, bool) {
	if val, ok := a.Get(key); ok {
		if str, ok := val.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetInt 获取整数类型缓存
// 注意：JSON round-trip 后数字类型会变成 float64，这里做安全转换
func (a *RistrettoAdapter) GetInt(key string) (int, bool) {
	if val, ok := a.Get(key); ok {
		switch v := val.(type) {
		case int:
			return v, true
		case int64:
			return int(v), true
		case float64:
			return int(v), true
		case json.Number:
			n, err := v.Int64()
			if err == nil {
				return int(n), true
			}
		}
	}
	return 0, false
}

// GetMap 获取 map 类型缓存
func (a *RistrettoAdapter) GetMap(key string) (map[string]any, bool) {
	if val, ok := a.Get(key); ok {
		if mapVal, ok := val.(map[string]any); ok {
			return mapVal, true
		}
	}
	return nil, false
}

// GetSlice 获取切片类型缓存
func (a *RistrettoAdapter) GetSlice(key string) ([]any, bool) {
	if val, ok := a.Get(key); ok {
		if s, ok := val.([]any); ok {
			return s, true
		}
	}
	return nil, false
}
