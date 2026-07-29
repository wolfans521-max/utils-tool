package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/log"
	internalConfig "git.intra.weibo.com/search_fe/wbutil-go/internal/config"
	vintage "git.intra.weibo.com/wbsearch/go-vintage"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/sync/singleflight"
)

const vintageAddr = "http://config.api.weibo.com"

// ========================= 配置中心 =========================

// InitFromViper 从使用方 viper 里初始化 discovery.group。
//
// 注意：业务方一般不需要调用（只用 ConfigGet 即可）；仅在你希望由业务方显式初始化时使用。
// 这里使用通用配置加载工具 [`internal/config.LoadSectionFromViper()`](internal/config/config.go:24)
// 读取 discovery 段，避免在各组件里重复实现 sub/unmarshal 逻辑。
func InitFromViper(v *viper.Viper) error {
	if v == nil {
		return errors.New("discovery: viper is nil")
	}

	type cfg struct {
		Group string `mapstructure:"group"`
	}
	var c cfg
	if err := internalConfig.LoadSectionFromViper(v, "discovery", &c); err != nil {
		return err
	}

	grp := strings.TrimSpace(c.Group)
	if grp == "" {
		return errors.New("discovery: config discovery.group is empty")
	}

	SetGroup(grp)
	return nil
}

// SetGroup 设置分组（仅供本仓库内部/测试使用；推荐使用 InitFromViper）。
func SetGroup(groupName string) {
	groupMu.Lock()
	defer groupMu.Unlock()
	groupID = groupName
}

var (
	groupMu sync.RWMutex
	groupID string

	groupInitOnce sync.Once
	groupInitErr  error

	configOnce sync.Once
	cfgCli     *vintage.ConfigServiceClient
	cfgErr     error

	errGroupNotSet = errors.New("discovery: config group not set; ensure config/app.yaml contains discovery.group")
)

type configCacheEntry struct {
	raw string
	err error
}

const configRefreshInterval = 5 * time.Second

var (
	// configCachePtr 指向当前配置缓存（copy-on-write：刷新时构建新 map，原子替换）。
	configCachePtr    atomic.Pointer[map[string]*configCacheEntry]
	configRefreshOnce sync.Once
)

// refreshConfigOnce 执行一次全量刷新（GetKeyListByGroup + LookupKey），填充/更新 configCache。
// 供 init 启动阶段和后续心跳周期复用。
func refreshConfigOnce() {
	w := httptest.NewRecorder()
	gin.SetMode(gin.ReleaseMode)
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request, _ = http.NewRequest("GET", "/internal/discovery/refresh_config", nil)
	ctx := context.New(ginCtx)
	grp, err := getGroup()
	if err != nil {
		log.Error(ctx, "refresh config err:"+err.Error())
		return
	}

	c, err := getConfigClient()
	if err != nil {
		log.Error(ctx, "refresh config err:"+err.Error())
		return
	}

	effectiveGroup := effectiveGroupName(grp)

	keys, err := c.GetKeyListByGroup(effectiveGroup)
	if err != nil {
		log.Error(ctx, "refresh config err:"+err.Error())
		return
	}

	// 构建 new cache（copy-on-write：失败时保留旧值避免配置闪断，旧 map 整体替换自然删旧 key）
	old := configCachePtr.Load()
	newCache := make(map[string]*configCacheEntry, len(keys))
	for _, k := range keys {
		raw, err := c.LookupKey(effectiveGroup, k)
		if err != nil {
			// 失败时继承旧值（若有），仅记录 err
			if old != nil {
				if oldEntry, ok := (*old)[k]; ok && oldEntry != nil && oldEntry.raw != "" {
					newCache[k] = &configCacheEntry{raw: oldEntry.raw, err: err}
					continue
				}
			}
			newCache[k] = &configCacheEntry{err: err}
		} else {
			newCache[k] = &configCacheEntry{raw: raw}
		}
	}
	configCachePtr.Store(&newCache)
	log.Mertics("[Discovery]", "refresh config cache")
}

// 包初始化时启动配置心跳刷新协程：
//
//	1）先同步执行一次 refreshConfigOnce，使服务启动后尽量能直接命中缓存；
//	2）再启动 goroutine 每隔 configRefreshInterval 周期性刷新。
func init() {
	configRefreshOnce.Do(func() {
		// 异步刷新：不阻塞启动（启动期间 ConfigGet 返回 defaultValue 兜底），
		// 首次刷新完成后缓存生效，后续按 configRefreshInterval 周期刷新。
		go func() {
			refreshConfigOnce()

			ticker := time.NewTicker(configRefreshInterval)
			defer ticker.Stop()

			for range ticker.C {
				refreshConfigOnce()
			}
		}()
	})
}

func getGroup() (string, error) {
	groupMu.RLock()
	id := groupID
	groupMu.RUnlock()
	if id == "" {
		// 业务方只调用 ConfigGet/NamingGet：这里懒加载读取 app.yaml / config/app.yaml。
		if err := initGroupFromComponentConfigIfNeeded(); err != nil {
			return "", err
		}
		groupMu.RLock()
		id = groupID
		groupMu.RUnlock()
	}
	if id == "" {
		return "", errGroupNotSet
	}
	return id, nil
}

func initGroupFromComponentConfigIfNeeded() error {
	groupInitOnce.Do(func() {
		v, err := internalConfig.ComponentViper()
		if err != nil {
			groupInitErr = err
			return
		}

		type cfg struct {
			Group string `mapstructure:"group"`
		}
		var c cfg
		if err := internalConfig.LoadSectionFromViper(v, "discovery", &c); err != nil {
			groupInitErr = err
			return
		}

		grp := strings.TrimSpace(c.Group)
		if grp == "" {
			groupInitErr = errors.New("discovery: config discovery.group is empty")
			return
		}

		SetGroup(grp)
	})

	return groupInitErr
}

// getConfigClient 懒加载
func getConfigClient() (*vintage.ConfigServiceClient, error) {
	configOnce.Do(func() {
		if cfgCli != nil || cfgErr != nil {
			return
		}
		cfgCli = vintage.NewConfigServiceClient(vintageAddr, 5*time.Second, 5*time.Second)
	})
	if cfgErr != nil {
		return nil, cfgErr
	}
	if cfgCli == nil {
		return nil, errors.New("discovery: ConfigServiceClient is nil")
	}
	return cfgCli, nil
}

// effectiveGroupName 将配置中心 group 名转换为实际请求使用的 group。
// 若未以 "online-" 前缀开头，补全为 "online-{group}-default"。
func effectiveGroupName(grp string) string {
	if strings.HasPrefix(grp, "online-") {
		return grp
	}
	return fmt.Sprintf("online-%s-default", grp)
}

// lookupRaw 根据当前group和key获取原始配置字符串（当前不暴露 ctx 语义，仅内部使用）。
func lookupRaw(key string) (string, error) {
	grp, err := getGroup()
	if err != nil {
		return "", err
	}

	c, err := getConfigClient()
	if err != nil {
		return "", err
	}

	effectiveGroup := effectiveGroupName(grp)

	return c.LookupKey(effectiveGroup, key)
}

// ConfigGet 配置获取入口
//   - serialization=false：直接返回原始字符串 raw；
//
// ConfigGet 配置获取入口
//   - serialization=false：直接返回原始字符串 raw；
//   - defaultValue：获取或解析失败时的兜底值；
//   - 返回值为 (any, error)：
//   - serialization=false 时返回 string(raw)
//   - serialization=true 时返回 map[string]any / []any（或 defaultValue），其中数字尽量保持为整型，避免科学计数。
//   - 内部通过心跳协程按 group 维度定时刷新所有 key 的 raw 值，调用方始终从本地缓存读取，避免过期抖动。
func ConfigGet(key string, serialization bool, defaultValue any) (any, error) {
	m := configCachePtr.Load()
	if m != nil {
		if entry, ok := (*m)[key]; ok && entry != nil {
			if entry.err != nil {
				return defaultValue, entry.err
			}
			raw := entry.raw
			if raw == "" {
				return defaultValue, nil
			}
			if !serialization {
				return raw, nil
			}
			return decodeConfigJSON(raw, key, defaultValue)
		}
	}

	return defaultValue, nil
}

// decodeConfigJSON 使用 sonic 进行 JSON 反序列化，并统一经过 normalizeJSONNumbers 处理。
func decodeConfigJSON(raw, key string, defaultValue any) (any, error) {
	var v any
	if err := sonic.UnmarshalString(raw, &v); err != nil {
		return defaultValue, fmt.Errorf("discovery: json unmarshal failed for key %q: %w", key, err)
	}
	return normalizeJSONNumbers(v), nil
}

// normalizeJSONNumbers 遍历 JSON 结构，将数字尽量转为整型，避免打印时出现科学计数。
// 适配两种来源：
//   - encoding/json：会产出 json.Number；
//   - sonic：直接产出 float64。
func normalizeJSONNumbers(v any) any {
	switch x := v.(type) {
	case json.Number:
		// 优先按整型解析，失败再退化为 float64，保持数值语义正确
		if i, err := x.Int64(); err == nil {
			return i
		}
		if f, err := strconv.ParseFloat(string(x), 64); err == nil {
			return f
		}
		return string(x)
	case float64:
		// 对 sonic 解出的 float64，若恰好为整数且在 int64 范围内，则转为 int64，避免科学计数。
		if x == float64(int64(x)) {
			return int64(x)
		}
		return x
	case map[string]any:
		for k, v2 := range x {
			x[k] = normalizeJSONNumbers(v2)
		}
		return x
	case []any:
		for i, v2 := range x {
			x[i] = normalizeJSONNumbers(v2)
		}
		return x
	default:
		return v
	}
}

// DumpCache 返回当前缓存的配置快照（包含最近一次心跳/拉取或冷启动写入的 raw 值）。
// 为避免外部修改内部缓存，这里返回的是 map 的拷贝。
// 注意：即便某些 key 最近一次拉取失败（entry.err != nil），也会出现在返回结果中，其 value 可能为空字符串。
func DumpCache() map[string]any {
	out := make(map[string]any)
	if m := configCachePtr.Load(); m != nil {
		for key, entry := range *m {
			if entry != nil {
				out[key] = entry.raw
			}
		}
	}
	return out
}

// ========================= 注册中心 =========================

var (
	namingOnce sync.Once
	namingCli  *vintage.NamingServiceClient
	namingErr  error

	// 缓存 ExtInfo 解析结果，key=原始 extInfo 字符串，value=*extCacheEntry（带 TTL）。
	extCache sync.Map // string -> *extCacheEntry
)

// extCacheEntry ExtInfo 缓存项（带 TTL，防止节点轮换导致无界增长）
type extCacheEntry struct {
	value    map[string]any
	expireAt int64 // UnixNano
}

const extCacheTTL = 60 * time.Second

// getExtFromCache 解析并缓存 ExtInfo，避免重复 json.Unmarshal。
// 带 60s TTL + 周期清理，防止节点轮换导致无界增长。
// 为避免业务侧修改缓存对象，这里会对 map 做一次浅拷贝后返回。
func getExtFromCache(extInfo string) map[string]any {
	if extInfo == "" {
		return nil
	}

	now := time.Now().UnixNano()

	if v, ok := extCache.Load(extInfo); ok {
		if entry, ok2 := v.(*extCacheEntry); ok2 && entry != nil {
			if now < entry.expireAt {
				src := entry.value
				if src == nil {
					return nil
				}
				// 浅拷贝一份，防止调用方修改缓存内部结构
				dst := make(map[string]any, len(src))
				for k, val := range src {
					dst[k] = val
				}
				return dst
			}
			// 过期，删除
			extCache.Delete(extInfo)
		}
	}

	var tmp map[string]any
	if err := sonic.Unmarshal([]byte(extInfo), &tmp); err != nil || tmp == nil {
		extCache.Store(extInfo, &extCacheEntry{value: nil, expireAt: now + int64(extCacheTTL)})
		return nil
	}

	extCache.Store(extInfo, &extCacheEntry{value: tmp, expireAt: now + int64(extCacheTTL)})

	dst := make(map[string]any, len(tmp))
	for k, val := range tmp {
		dst[k] = val
	}
	return dst
}

// getNamingClient 懒加载创建
func getNamingClient() (*vintage.NamingServiceClient, error) {
	namingOnce.Do(func() {
		if namingCli != nil || namingErr != nil {
			return
		}
		// timeout=5s, heartInterval=5s，避免 go-vintage 内部 watcher 使用 0 interval 导致 panic
		namingCli = vintage.NewNamingServiceClient(vintageAddr, 5*time.Second, 5*time.Second)
	})
	if namingErr != nil {
		return nil, namingErr
	}
	if namingCli == nil {
		return nil, errors.New("discovery: NamingServiceClient is nil")
	}
	return namingCli, nil
}

type NamingNode struct {
	Host   string         // ip
	Port   int            // 端口
	Domain string         // 原始 host（ip:port）
	Weight int            // extInfo.weight 或 1
	Ext    map[string]any // 解析后的 extInfo（只读，不要在业务侧修改）
}

// NamingGet 获取某个service/cluster下的节点列表。
// 🔥 优化：使用本地缓存 + 超时控制，避免高并发下服务发现调用阻塞导致 goroutine 泄漏
func NamingGet(service, cluster string) ([]NamingNode, error) {
	// 使用高性能版本（针对 10w+ QPS 优化）
	return NamingGetHighPerformance(service, cluster)
}

// namingGetDirect 直接调用服务发现（不经过缓存）
func namingGetDirect(service, cluster string) ([]NamingNode, error) {
	cli, err := getNamingClient()
	if err != nil {
		return nil, err
	}

	sn, err := cli.Lookup(service, fmt.Sprintf("%s/service", cluster))
	if err != nil {
		return nil, err
	}
	if sn == nil {
		return nil, nil
	}

	working := sn.Nodes.Working
	res := make([]NamingNode, 0, len(working))

	for _, n := range working {
		ext := getExtFromCache(n.ExtInfo)

		weight := 1
		if v, ok := ext["weight"]; ok {
			switch vv := v.(type) {
			case float64:
				if vv > 0 {
					weight = int(vv)
				}
			case int:
				if vv > 0 {
					weight = vv
				}
			}
		}

		domain := n.Host
		host := domain
		port := 0
		if parts := strings.Split(domain, ":"); len(parts) == 2 {
			host = parts[0]
			if p, e := strconv.Atoi(parts[1]); e == nil {
				port = p
			}
		}

		res = append(res, NamingNode{
			Host:   host,
			Port:   port,
			Domain: domain,
			Weight: weight,
			Ext:    ext,
		})
	}

	return res, nil
}

// ========================= 高性能服务发现缓存 (10w+ QPS 优化) =========================

// 缓存项类型
type cacheItem struct {
	nodes      []NamingNode
	expireAt   int64 // UnixNano
	err        error
	lastAccess int64 // UnixNano，用于LRU
}

// 缓存配置
const (
	defaultCacheTTL      = 5 * time.Second  // 默认缓存有效期
	cacheCleanupInterval = 10 * time.Second // 清理过期缓存间隔
	maxCacheSize         = 1000             // 最大缓存条目数
)

var (
	// 本地缓存: key = service/cluster
	namingCache      sync.Map // map[string]*cacheItem
	cacheHitCount    int64
	cacheMissCount   int64
	cacheHitBytes    int64
	cacheMissBytes   int64
	cacheCleanupOnce sync.Once

	// singleflight 合并并发 miss，避免缓存过期瞬间的惊群效应
	namingSF singleflight.Group
)

// initCacheCleanup 启动缓存清理协程
func initCacheCleanup() {
	cacheCleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(cacheCleanupInterval)
			defer ticker.Stop()
			for range ticker.C {
				cleanupExpiredCache()
			}
		}()
	})
}

// cleanupExpiredCache 清理过期缓存
func cleanupExpiredCache() {
	now := time.Now().UnixNano()
	var deleted int
	namingCache.Range(func(key, value interface{}) bool {
		if item, ok := value.(*cacheItem); ok {
			if now > item.expireAt {
				namingCache.Delete(key)
				deleted++
			}
		}
		return true
	})
	if deleted > 0 {
		log.Mertics("[Discovery]", fmt.Sprintf("cleaned up %d expired cache entries", deleted))
	}

	// 同时清理过期的 extCache 条目
	var extDeleted int
	extCache.Range(func(key, value any) bool {
		if entry, ok := value.(*extCacheEntry); ok && now > entry.expireAt {
			extCache.Delete(key)
			extDeleted++
		}
		return true
	})
	if extDeleted > 0 {
		log.Mertics("[Discovery]", fmt.Sprintf("cleaned up %d expired extCache entries", extDeleted))
	}
}

// getCacheKey 生成缓存key
func getCacheKey(service, cluster string) string {
	return service + "/" + cluster
}

// getFromCache 从缓存获取
func getFromCache(key string) (*cacheItem, bool) {
	if v, ok := namingCache.Load(key); ok {
		if item, ok2 := v.(*cacheItem); ok2 {
			now := time.Now().UnixNano()
			if now < item.expireAt {
				// 更新最后访问时间
				atomic.StoreInt64(&item.lastAccess, now)
				return item, true
			}
		}
	}
	return nil, false
}

// setCache 设置缓存
func setCache(key string, nodes []NamingNode, err error, ttl time.Duration) {
	now := time.Now().UnixNano()
	item := &cacheItem{
		nodes:      nodes,
		expireAt:   now + int64(ttl),
		err:        err,
		lastAccess: now,
	}
	namingCache.Store(key, item)
}

// NamingGetHighPerformance 高性能版本（本地缓存 + singleflight 合并并发 miss）
func NamingGetHighPerformance(service, cluster string) ([]NamingNode, error) {
	// 初始化缓存清理
	initCacheCleanup()

	key := getCacheKey(service, cluster)

	// 1. 先查本地缓存
	if item, hit := getFromCache(key); hit {
		atomic.AddInt64(&cacheHitCount, 1)
		atomic.AddInt64(&cacheHitBytes, int64(len(item.nodes))*64) // 估算
		return item.nodes, item.err
	}

	atomic.AddInt64(&cacheMissCount, 1)

	// 2. 缓存未命中，使用 singleflight 合并同一 key 的并发请求
	//    避免缓存过期瞬间 N 个请求同时穿透到配置中心
	v, err, _ := namingSF.Do(key, func() (any, error) {
		nodes, fetchErr := namingGetWithTimeout(service, cluster, 100*time.Millisecond)
		// 写入缓存（即使失败也缓存，避免雪崩）
		setCache(key, nodes, fetchErr, defaultCacheTTL)
		return nodes, fetchErr
	})

	if err != nil {
		atomic.AddInt64(&cacheMissBytes, 1)
		return nil, err
	}

	nodes := v.([]NamingNode)
	atomic.AddInt64(&cacheMissBytes, int64(len(nodes))*64)

	return nodes, nil
}

// namingGetWithTimeout 带超时的服务发现调用
func namingGetWithTimeout(service, cluster string, timeout time.Duration) ([]NamingNode, error) {
	type result struct {
		nodes []NamingNode
		err   error
	}

	ch := make(chan result, 1)

	go func() {
		nodes, err := namingGetDirect(service, cluster)
		ch <- result{nodes: nodes, err: err}
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case res := <-ch:
		return res.nodes, res.err
	case <-timer.C:
		return nil, fmt.Errorf("naming get timeout: service=%s, cluster=%s", service, cluster)
	}
}

// GetCacheStats 获取缓存统计信息
func GetCacheStats() map[string]int64 {
	return map[string]int64{
		"hit":       atomic.LoadInt64(&cacheHitCount),
		"miss":      atomic.LoadInt64(&cacheMissCount),
		"hitBytes":  atomic.LoadInt64(&cacheHitBytes),
		"missBytes": atomic.LoadInt64(&cacheMissBytes),
	}
}

// ResetCacheStats 重置缓存统计
func ResetCacheStats() {
	atomic.StoreInt64(&cacheHitCount, 0)
	atomic.StoreInt64(&cacheMissCount, 0)
	atomic.StoreInt64(&cacheHitBytes, 0)
	atomic.StoreInt64(&cacheMissBytes, 0)
}

// ClearNamingCache 清空服务发现缓存
func ClearNamingCache() {
	namingCache.Range(func(key, _ any) bool {
		namingCache.Delete(key)
		return true
	})
}
