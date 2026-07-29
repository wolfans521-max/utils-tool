// Package cache 提供缓存核心实现，包含 Ristretto / BigCache / SharedCache(mmap) 三种后端
//
// 本文件实现基于 mmap 的 SharedCache（跨平台共享内存，支持进程间共享）
// 如需创建缓存实例，请使用 facades/cache 包的 NewCache() 工厂方法
package cache

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/edsrzf/mmap-go"
	"github.com/goccy/go-yaml"
	"go.uber.org/zap"
)

// 常量定义
const (
	// DefaultCacheSize 默认缓存大小 10MB
	DefaultCacheSize = 10 * 1024 * 1024
	// MaxCacheSize 最大缓存大小 1GB
	MaxCacheSize = 1024 * 1024 * 1024
	// MinCacheSize 最小缓存大小 1MB
	MinCacheSize = 1 * 1024 * 1024
	// MaxValueSize 单个value的最大大小 1MB（防止单个key过大导致序列化耗时飙升）
	MaxValueSize = 1 * 1024 * 1024
	// DefaultShardCount 默认分片数量，必须是2的幂次方
	DefaultShardCount = 32
	// HeaderSize 文件头大小（4字节长度 + 4字节CRC32）
	HeaderSize = 8
	// FilePermission 缓存文件权限（仅所有者可读写）
	FilePermission = 0600
	// DirPermission 缓存目录权限
	DirPermission = 0755
	// LocalCacheTTL 本地缓存默认TTL（秒）
	LocalCacheTTL = 1
	// CleanupInterval 过期清理间隔
	CleanupInterval = 30 * time.Second
	// MmapCleanupInterval mmap过期数据清理间隔
	MmapCleanupInterval = 5 * time.Minute
	// WriteBufferSize 写缓冲区大小
	WriteBufferSize = 1000
	// FlushInterval 刷盘间隔
	FlushInterval = 100 * time.Millisecond
	// MaxWriteRetries 最大写重试次数
	MaxWriteRetries = 3
	// ExpandThreshold 扩容阈值（使用率达到80%时扩容）
	ExpandThreshold = 0.8
	// ExpandFactor 扩容因子
	ExpandFactor = 2
	// GracefulShutdownTimeout 优雅关闭超时时间
	GracefulShutdownTimeout = 10 * time.Second
)

// ExpirePolicy 过期策略类型
type ExpirePolicy int

const (
	// ExpirePolicyAbsolute 绝对过期（到达指定时间点过期）
	ExpirePolicyAbsolute ExpirePolicy = iota
	// ExpirePolicyRelative 相对过期（从设置时刻开始计算TTL）
	ExpirePolicyRelative
	// ExpirePolicySliding 滑动过期（每次访问重置TTL）
	ExpirePolicySliding
)

// SharedCacheItem 共享缓存项
type SharedCacheItem struct {
	Value        any          `json:"value"`
	ExpireTime   int64        `json:"expire_time"`   // Unix 时间戳，0 表示永不过期
	ExpirePolicy ExpirePolicy `json:"expire_policy"` // 过期策略
	TTL          int64        `json:"ttl"`           // 原始TTL（秒），用于滑动过期
	CreatedAt    int64        `json:"created_at"`    // 创建时间
	AccessCount  int64        `json:"access_count"`  // 访问次数
}

// IsExpired 检查缓存项是否过期
func (item *SharedCacheItem) IsExpired() bool {
	if item.ExpireTime == 0 {
		return false
	}
	return time.Now().Unix() >= item.ExpireTime
}

// localCacheItem 本地缓存项（带本地过期时间）
type localCacheItem struct {
	item       *SharedCacheItem
	localExpAt int64 // 本地缓存过期时间（纳秒）
}

// isLocalExpired 检查本地缓存是否过期
func (l *localCacheItem) isLocalExpired() bool {
	return time.Now().UnixNano() >= l.localExpAt
}

// SharedCacheData 共享缓存数据结构
type SharedCacheData struct {
	Items   map[string]*SharedCacheItem `json:"items"`
	Version int64                       `json:"version"` // 数据版本号
}

// SharedConfig 共享缓存配置
type SharedConfig struct {
	// Name 缓存名称（用于生成文件名）
	Name string
	// Size 缓存文件大小（字节）
	Size int
	// Dir 缓存文件目录（可选，默认使用临时目录）
	Dir string
	// ShardCount 分片数量（可选，默认32）
	ShardCount int
	// LocalCacheTTLMs 本地缓存TTL（毫秒，可选，默认1000ms）
	LocalCacheTTLMs int
	// AsyncFlush 是否启用异步刷盘（可选，默认false）
	AsyncFlush bool
	// AutoExpand 是否启用自动扩容（可选，默认true）
	AutoExpand bool
	// EnableCRC 是否启用CRC校验（可选，默认true）
	EnableCRC bool
	// EnableLog 是否启用日志（可选，默认false）
	EnableLog bool
	// Logger 自定义日志器（可选）
	Logger *zap.Logger
}

// writeOperation 写操作
type writeOperation struct {
	key      string
	item     *SharedCacheItem
	isDel    bool
	doneChan chan error // 用于同步等待写入完成
}

// cacheShard 缓存分片，每个分片有独立的锁
type cacheShard struct {
	rwMutex    sync.RWMutex
	localCache sync.Map // key -> *localCacheItem
}

// Stats 缓存统计信息
type Stats struct {
	ItemCount          int     // 缓存项数量
	LocalCacheHits     int64   // 本地缓存命中次数
	MmapReads          int64   // mmap读取次数
	MmapWrites         int64   // mmap写入次数
	ShardCount         int     // 分片数量
	LocalCacheTTLMs    int64   // 本地缓存TTL（毫秒）
	HitRate            float64 // 命中率
	TotalGets          int64   // 总Get次数
	TotalSets          int64   // 总Set次数
	TotalDeletes       int64   // 总Delete次数
	WriteRetries       int64   // 写重试次数
	CRCErrors          int64   // CRC校验错误次数
	ExpandCount        int64   // 扩容次数
	CurrentSize        int64   // 当前文件大小
	UsedSize           int64   // 已使用大小
	DroppedWrites      int64   // 丢弃的写操作数
	AsyncFlushCount    int64   // 异步刷盘次数
	BatchWriteCount    int64   // 批量写入次数
	ExpiredItemsCount  int64   // 清理的过期项数量
	SyncFallbackWrites int64   // 降级同步写入次数
	ValueSizeExceeded  int64   // value大小超限次数
}

// SharedCache 基于内存映射文件的缓存管理器（跨平台，支持进程间共享）
// 优化设计V3：
// 1. 使用分片减少锁竞争
// 2. 使用RWMutex支持并发读
// 3. 添加本地缓存层减少mmap访问
// 4. 写缓冲池+批量刷盘
// 5. 异步刷盘模式（带数据安全保证）
// 6. CRC32数据校验
// 7. 自动扩容（修复内存泄漏）
// 8. 写重试机制
// 9. 高性能JSON序列化（sonic）
// 10. 日志埋点
// 11. mmap过期数据定期清理
// 12. 增强并发安全
type SharedCache struct {
	config        SharedConfig
	file          *os.File
	mmap          mmap.MMap
	filePath      string
	shards        []*cacheShard
	shardCount    uint32
	localCacheTTL time.Duration
	globalRWMutex sync.RWMutex // 用于保护mmap的全局读写锁
	closed        atomic.Bool
	cleanupDone   chan struct{}
	mmapCleanDone chan struct{}

	// 写缓冲相关
	writeBuffer chan *writeOperation
	flushDone   chan struct{}
	flushWg     sync.WaitGroup

	// 日志
	logger *zap.Logger

	// 统计信息
	stats struct {
		localCacheHits     atomic.Int64
		mmapReads          atomic.Int64
		mmapWrites         atomic.Int64
		totalGets          atomic.Int64
		totalSets          atomic.Int64
		totalDeletes       atomic.Int64
		writeRetries       atomic.Int64
		crcErrors          atomic.Int64
		expandCount        atomic.Int64
		droppedWrites      atomic.Int64
		asyncFlushCount    atomic.Int64
		batchWriteCount    atomic.Int64
		expiredItemsCount  atomic.Int64
		syncFallbackWrites atomic.Int64
		valueSizeExceeded  atomic.Int64
	}
}

var (
	sharedInstance *SharedCache
	sharedOnce     sync.Once
	sharedInitErr  error
)

// GetInstance 获取共享缓存单例实例
func GetInstance(name string) (*SharedCache, error) {
	sharedOnce.Do(func() {
		if name == "" {
			sharedInitErr = fmt.Errorf("缓存名称不能为空")
			return
		}

		// 读取配置文件
		size, dir, err := loadSharedCacheConfig(name)
		if err != nil {
			// 如果配置文件不存在，使用默认配置
			size = DefaultCacheSize
			dir = ""
		}

		// 创建共享缓存实例
		impl, err := NewSharedCache(SharedConfig{
			Name:       name,
			Size:       size,
			Dir:        dir,
			AutoExpand: true,
			EnableCRC:  true,
		})
		if err != nil {
			sharedInitErr = err
			return
		}

		sharedInstance = impl
	})
	return sharedInstance, sharedInitErr
}

// loadSharedCacheConfig 读取配置文件
func loadSharedCacheConfig(name string) (size int, dir string, err error) {
	// 读取 app.yaml 获取缓存配置
	appConfigPath := "config/app.yaml"
	appData, err := os.ReadFile(appConfigPath)
	if err != nil {
		return DefaultCacheSize, "", nil // 配置文件不存在时使用默认值
	}

	var appConfig map[string]any
	if err := yaml.Unmarshal(appData, &appConfig); err != nil {
		return DefaultCacheSize, "", nil // 解析失败时使用默认值
	}

	// 获取 cache 配置节
	cacheConfig, ok := appConfig["cache"].(map[string]any)
	if !ok {
		return DefaultCacheSize, "", nil
	}

	// 直接从 cache 配置节中获取指定名称的缓存配置
	memoryConfig, ok := cacheConfig[name].(map[string]any)
	if !ok {
		return DefaultCacheSize, "", nil
	}

	// 读取 max_memory / size
	size = DefaultCacheSize
	if maxMem, ok := memoryConfig["max_memory"]; ok {
		size = toIntShared(maxMem)
	}
	if s, ok := memoryConfig["size"]; ok {
		size = toIntShared(s)
	}

	// 确保 size 在 [MinCacheSize, MaxCacheSize] 范围内
	if size < MinCacheSize {
		size = MinCacheSize
	} else if size > MaxCacheSize {
		size = MaxCacheSize
	}

	// 读取 dir
	if d, ok := memoryConfig["dir"].(string); ok {
		dir = d
	}

	return size, dir, nil
}

// toIntShared 将 any 转换为 int
func toIntShared(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case uint:
		return int(val)
	case uint64:
		return int(val)
	case float64:
		return int(val)
	default:
		return 0
	}
}

// NewSharedCache 创建新的共享缓存实例
func NewSharedCache(config SharedConfig) (*SharedCache, error) {
	// 设置默认值并验证配置
	if config.Size <= 0 {
		config.Size = DefaultCacheSize
	}
	if config.Size < MinCacheSize {
		config.Size = MinCacheSize
	}
	if config.Size > MaxCacheSize {
		config.Size = MaxCacheSize
	}
	if config.ShardCount <= 0 {
		config.ShardCount = DefaultShardCount
	}
	// 确保分片数是2的幂次方
	config.ShardCount = nextPowerOfTwo(config.ShardCount)

	localTTL := time.Duration(LocalCacheTTL) * time.Second
	if config.LocalCacheTTLMs > 0 {
		localTTL = time.Duration(config.LocalCacheTTLMs) * time.Millisecond
	}

	// 确定缓存文件路径
	var filePath string
	if config.Dir != "" {
		filePath = filepath.Join(config.Dir, config.Name+".cache")
	} else {
		// 使用临时目录
		filePath = filepath.Join(os.TempDir(), "wbutil_cache_"+config.Name+".cache")
	}

	// 初始化分片
	shards := make([]*cacheShard, config.ShardCount)
	for i := 0; i < config.ShardCount; i++ {
		shards[i] = &cacheShard{}
	}

	// 初始化日志
	var logger *zap.Logger
	if config.EnableLog {
		if config.Logger != nil {
			logger = config.Logger
		} else {
			logger, _ = zap.NewProduction()
		}
	}

	cache := &SharedCache{
		config:        config,
		filePath:      filePath,
		shards:        shards,
		shardCount:    uint32(config.ShardCount),
		localCacheTTL: localTTL,
		cleanupDone:   make(chan struct{}),
		mmapCleanDone: make(chan struct{}),
		writeBuffer:   make(chan *writeOperation, WriteBufferSize),
		flushDone:     make(chan struct{}),
		logger:        logger,
	}

	if err := cache.init(); err != nil {
		return nil, err
	}

	cache.logInfo("cache initialized", zap.String("path", filePath), zap.Int("size", config.Size))

	// 启动后台清理协程
	go cache.cleanupLoop()

	// 启动mmap过期数据清理协程
	go cache.mmapCleanupLoop()

	// 如果启用异步刷盘，启动写入协程
	if config.AsyncFlush {
		cache.flushWg.Add(1)
		go cache.asyncFlushLoop()
	}

	return cache, nil
}

// logInfo 记录信息日志
func (c *SharedCache) logInfo(msg string, fields ...zap.Field) {
	if c.logger != nil {
		c.logger.Info(msg, fields...)
	}
}

// logWarn 记录警告日志
func (c *SharedCache) logWarn(msg string, fields ...zap.Field) {
	if c.logger != nil {
		c.logger.Warn(msg, fields...)
	}
}

// logError 记录错误日志
func (c *SharedCache) logError(msg string, fields ...zap.Field) {
	if c.logger != nil {
		c.logger.Error(msg, fields...)
	}
}

// nextPowerOfTwo 返回大于等于n的最小2的幂次方
func nextPowerOfTwo(n int) int {
	if n <= 1 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}

// getShard 根据key获取对应的分片
func (c *SharedCache) getShard(key string) *cacheShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return c.shards[h.Sum32()&(c.shardCount-1)]
}

// init 初始化内存映射文件
func (c *SharedCache) init() error {
	// 确保目录存在
	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, DirPermission); err != nil {
		return fmt.Errorf("创建缓存目录失败: %w", err)
	}

	// 打开或创建文件（使用安全的文件权限）
	file, err := os.OpenFile(c.filePath, os.O_RDWR|os.O_CREATE, FilePermission)
	if err != nil {
		return fmt.Errorf("打开缓存文件失败: %w", err)
	}

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 如果文件大小不足，扩展文件
	if info.Size() < int64(c.config.Size) {
		if err := file.Truncate(int64(c.config.Size)); err != nil {
			file.Close()
			return fmt.Errorf("扩展缓存文件失败: %w", err)
		}
	}

	// 创建内存映射
	m, err := mmap.Map(file, mmap.RDWR, 0)
	if err != nil {
		file.Close()
		return fmt.Errorf("创建内存映射失败: %w", err)
	}

	c.file = file
	c.mmap = m

	return nil
}

// asyncFlushLoop 异步刷盘循环（修复数据丢失问题）
func (c *SharedCache) asyncFlushLoop() {
	defer c.flushWg.Done()

	ticker := time.NewTicker(FlushInterval)
	defer ticker.Stop()

	batch := make([]*writeOperation, 0, WriteBufferSize)
	pendingDoneChans := make([]chan error, 0, WriteBufferSize)

	flushAndNotify := func() {
		if len(batch) == 0 {
			return
		}

		err := c.flushBatchWithError(batch)

		// 通知所有等待的写操作
		for _, doneChan := range pendingDoneChans {
			if doneChan != nil {
				select {
				case doneChan <- err:
				default:
				}
			}
		}

		batch = batch[:0]
		pendingDoneChans = pendingDoneChans[:0]
	}

	for {
		select {
		case op := <-c.writeBuffer:
			batch = append(batch, op)
			if op.doneChan != nil {
				pendingDoneChans = append(pendingDoneChans, op.doneChan)
			}
			// 批量达到阈值，立即刷盘
			if len(batch) >= WriteBufferSize/2 {
				flushAndNotify()
			}

		case <-ticker.C:
			// 定时刷盘
			flushAndNotify()

		case <-c.flushDone:
			// 关闭前刷盘剩余数据 - 确保数据不丢失
			c.logInfo("async flush loop shutting down, draining buffer")

			// 排空缓冲区
		drainLoop:
			for {
				select {
				case op := <-c.writeBuffer:
					batch = append(batch, op)
					if op.doneChan != nil {
						pendingDoneChans = append(pendingDoneChans, op.doneChan)
					}
				default:
					break drainLoop
				}
			}

			// 最终刷盘
			flushAndNotify()
			c.logInfo("async flush loop shutdown complete", zap.Int("flushed", len(batch)))
			return
		}
	}
}

// flushBatchWithError 批量刷盘（返回错误）
func (c *SharedCache) flushBatchWithError(batch []*writeOperation) error {
	if len(batch) == 0 {
		return nil
	}

	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	data, err := c.readDataFromMmapLocked()
	if err != nil {
		c.logError("flush batch read error", zap.Error(err))
		return err
	}

	// 应用所有操作
	for _, op := range batch {
		if op.isDel {
			delete(data.Items, op.key)
		} else {
			data.Items[op.key] = op.item
		}
	}

	// 写入
	if err := c.writeDataToMmapWithRetry(data); err != nil {
		c.logError("flush batch write error", zap.Error(err))
		return err
	}

	c.stats.batchWriteCount.Add(1)
	c.stats.asyncFlushCount.Add(1)
	return nil
}

// flushBatch 批量刷盘（兼容旧接口）
func (c *SharedCache) flushBatch(batch []*writeOperation) {
	_ = c.flushBatchWithError(batch)
}

// cleanupLoop 后台清理过期的本地缓存
func (c *SharedCache) cleanupLoop() {
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupExpiredLocal()
		case <-c.cleanupDone:
			return
		}
	}
}

// mmapCleanupLoop 后台清理mmap中的过期数据
func (c *SharedCache) mmapCleanupLoop() {
	ticker := time.NewTicker(MmapCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanupExpiredMmap()
		case <-c.mmapCleanDone:
			return
		}
	}
}

// cleanupExpiredLocal 清理过期的本地缓存
func (c *SharedCache) cleanupExpiredLocal() {
	if c.closed.Load() {
		return
	}

	now := time.Now().UnixNano()
	for _, shard := range c.shards {
		shard.localCache.Range(func(key, value any) bool {
			if item, ok := value.(*localCacheItem); ok {
				if item.localExpAt <= now || item.item.IsExpired() {
					shard.localCache.Delete(key)
				}
			}
			return true
		})
	}
}

// cleanupExpiredMmap 清理mmap中的过期数据
func (c *SharedCache) cleanupExpiredMmap() {
	if c.closed.Load() {
		return
	}

	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	data, err := c.readDataFromMmapLocked()
	if err != nil {
		return
	}

	// 收集过期的key
	expiredKeys := make([]string, 0)
	for key, item := range data.Items {
		if item.IsExpired() {
			expiredKeys = append(expiredKeys, key)
		}
	}

	if len(expiredKeys) == 0 {
		return
	}

	// 删除过期项
	for _, key := range expiredKeys {
		delete(data.Items, key)
	}

	// 写回
	if err := c.writeDataToMmapWithRetry(data); err == nil {
		c.stats.expiredItemsCount.Add(int64(len(expiredKeys)))
		c.logInfo("cleaned expired items from mmap", zap.Int("count", len(expiredKeys)))
	}
}

// calculateCRC32 计算CRC32校验值
func calculateCRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

// readDataFromMmap 从内存映射读取数据（需要持有读锁）
func (c *SharedCache) readDataFromMmap() (*SharedCacheData, error) {
	c.globalRWMutex.RLock()
	defer c.globalRWMutex.RUnlock()
	return c.readDataFromMmapLocked()
}

// readDataFromMmapLocked 从内存映射读取数据（内部方法，需要已持有锁）
func (c *SharedCache) readDataFromMmapLocked() (*SharedCacheData, error) {
	c.stats.mmapReads.Add(1)

	if c.mmap == nil || len(c.mmap) < HeaderSize {
		return &SharedCacheData{Items: make(map[string]*SharedCacheItem)}, nil
	}

	// 读取数据长度（前4字节）
	length := int(binary.LittleEndian.Uint32(c.mmap[0:4]))

	if length == 0 || length > c.config.Size-HeaderSize {
		// 空数据或无效数据，返回空缓存
		return &SharedCacheData{Items: make(map[string]*SharedCacheItem)}, nil
	}

	// 读取存储的CRC32（4-8字节）
	storedCRC := binary.LittleEndian.Uint32(c.mmap[4:8])

	// 读取JSON数据
	dataBytes := make([]byte, length)
	copy(dataBytes, c.mmap[HeaderSize:HeaderSize+length])

	// 验证CRC32
	if c.config.EnableCRC {
		calculatedCRC := calculateCRC32(dataBytes)
		if calculatedCRC != storedCRC {
			c.stats.crcErrors.Add(1)
			c.logWarn("CRC32 mismatch", zap.Uint32("stored", storedCRC), zap.Uint32("calculated", calculatedCRC))
			// CRC校验失败，返回空缓存
			return &SharedCacheData{Items: make(map[string]*SharedCacheItem)}, nil
		}
	}

	// 使用sonic高性能反序列化
	var data SharedCacheData
	if err := sonic.Unmarshal(dataBytes, &data); err != nil {
		c.logWarn("unmarshal error", zap.Error(err))
		return &SharedCacheData{Items: make(map[string]*SharedCacheItem)}, nil
	}

	if data.Items == nil {
		data.Items = make(map[string]*SharedCacheItem)
	}

	return &data, nil
}

// writeDataToMmap 将数据写入内存映射（需要持有写锁）
func (c *SharedCache) writeDataToMmap(data *SharedCacheData) error {
	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()
	return c.writeDataToMmapWithRetry(data)
}

// writeDataToMmapWithRetry 带重试的写入（内部方法，需要已持有写锁）
func (c *SharedCache) writeDataToMmapWithRetry(data *SharedCacheData) error {
	var lastErr error

	for retry := 0; retry < MaxWriteRetries; retry++ {
		err := c.writeDataToMmapLocked(data)
		if err == nil {
			return nil
		}

		lastErr = err
		c.stats.writeRetries.Add(1)

		// 如果是空间不足，尝试扩容
		if c.config.AutoExpand && isSpaceError(err) {
			if expandErr := c.expandLocked(); expandErr != nil {
				c.logError("expand failed", zap.Error(expandErr))
				continue
			}
			// 扩容成功，重试写入
			continue
		}

		// 其他错误，短暂等待后重试
		time.Sleep(time.Duration(retry+1) * 10 * time.Millisecond)
	}

	c.logError("write failed after retries", zap.Error(lastErr), zap.Int("retries", MaxWriteRetries))
	return fmt.Errorf("写入失败（已重试%d次）: %w", MaxWriteRetries, lastErr)
}

// isSpaceError 检查是否是空间不足错误
func isSpaceError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "数据超过缓存大小限制"
}

// writeDataToMmapLocked 将数据写入内存映射（内部方法，需要已持有写锁）
func (c *SharedCache) writeDataToMmapLocked(data *SharedCacheData) error {
	if c.mmap == nil {
		return fmt.Errorf("mmap未初始化")
	}

	c.stats.mmapWrites.Add(1)

	// 更新版本号
	data.Version++

	// 使用sonic高性能序列化
	jsonData, err := sonic.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	length := len(jsonData)
	if length > c.config.Size-HeaderSize {
		return fmt.Errorf("数据超过缓存大小限制")
	}

	// 计算CRC32
	crcValue := calculateCRC32(jsonData)

	// 写入数据长度（前4字节）
	binary.LittleEndian.PutUint32(c.mmap[0:4], uint32(length))

	// 写入CRC32（4-8字节）
	binary.LittleEndian.PutUint32(c.mmap[4:8], crcValue)

	// 写入JSON数据
	copy(c.mmap[HeaderSize:], jsonData)

	// 刷新到文件
	if err := c.mmap.Flush(); err != nil {
		return fmt.Errorf("刷新内存映射失败: %w", err)
	}

	return nil
}

// expandLocked 扩容（内部方法，需要已持有写锁）- 修复内存泄漏
func (c *SharedCache) expandLocked() error {
	newSize := c.config.Size * ExpandFactor
	if newSize > MaxCacheSize {
		newSize = MaxCacheSize
	}
	if newSize == c.config.Size {
		return fmt.Errorf("已达到最大缓存大小限制")
	}

	c.logInfo("expanding cache", zap.Int("oldSize", c.config.Size), zap.Int("newSize", newSize))

	// 保存旧的mmap引用
	oldMmap := c.mmap

	// 取消当前映射 - 必须先Unmap再重新Map，否则会内存泄漏
	if oldMmap != nil {
		if err := oldMmap.Unmap(); err != nil {
			return fmt.Errorf("取消内存映射失败: %w", err)
		}
		c.mmap = nil // 立即置空，防止重复释放
	}

	// 强制GC回收旧映射内存
	runtime.GC()

	// 扩展文件
	if err := c.file.Truncate(int64(newSize)); err != nil {
		// 扩展失败，尝试恢复旧映射
		if m, mapErr := mmap.Map(c.file, mmap.RDWR, 0); mapErr == nil {
			c.mmap = m
		}
		return fmt.Errorf("扩展文件失败: %w", err)
	}

	// 重新映射
	m, err := mmap.Map(c.file, mmap.RDWR, 0)
	if err != nil {
		return fmt.Errorf("重新映射失败: %w", err)
	}

	c.mmap = m
	c.config.Size = newSize
	c.stats.expandCount.Add(1)

	c.logInfo("cache expanded successfully", zap.Int("newSize", newSize))
	return nil
}

// Set 设置缓存，支持可选的过期时间（秒）
func (c *SharedCache) Set(key string, value any, ttl ...int) error {
	return c.SetWithPolicy(key, value, ExpirePolicyRelative, ttl...)
}

// SetWithPolicy 使用指定过期策略设置缓存
func (c *SharedCache) SetWithPolicy(key string, value any, policy ExpirePolicy, ttl ...int) error {
	if c.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	// 检查value大小限制
	if err := c.checkValueSize(value); err != nil {
		c.stats.valueSizeExceeded.Add(1)
		return err
	}

	c.stats.totalSets.Add(1)

	var expireTime int64
	var ttlValue int64
	if len(ttl) > 0 && ttl[0] > 0 {
		ttlValue = int64(ttl[0])
		switch policy {
		case ExpirePolicyAbsolute:
			expireTime = int64(ttl[0]) // ttl作为绝对时间戳
		case ExpirePolicyRelative, ExpirePolicySliding:
			expireTime = time.Now().Add(time.Duration(ttl[0]) * time.Second).Unix()
		}
	}

	item := &SharedCacheItem{
		Value:        value,
		ExpireTime:   expireTime,
		ExpirePolicy: policy,
		TTL:          ttlValue,
		CreatedAt:    time.Now().Unix(),
		AccessCount:  0,
	}

	// 更新本地缓存
	shard := c.getShard(key)
	localItem := &localCacheItem{
		item:       item,
		localExpAt: time.Now().Add(c.localCacheTTL).UnixNano(),
	}
	shard.localCache.Store(key, localItem)

	// 异步刷盘模式
	if c.config.AsyncFlush {
		op := &writeOperation{key: key, item: item, isDel: false}
		select {
		case c.writeBuffer <- op:
			return nil
		default:
			// 缓冲区满时降级为同步写入，避免丢数据
			c.stats.droppedWrites.Add(1)
			c.stats.syncFallbackWrites.Add(1)
			c.logWarn("write buffer full, falling back to sync write", zap.String("key", key))
			// 降级为同步写入
			c.globalRWMutex.Lock()
			defer c.globalRWMutex.Unlock()
			data, err := c.readDataFromMmapLocked()
			if err != nil {
				return err
			}
			data.Items[key] = item
			return c.writeDataToMmapWithRetry(data)
		}
	}

	// 同步写入mmap
	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	data, err := c.readDataFromMmapLocked()
	if err != nil {
		return err
	}

	data.Items[key] = item
	return c.writeDataToMmapWithRetry(data)
}

// checkValueSize 检查value大小是否超限
func (c *SharedCache) checkValueSize(value any) error {
	// 序列化value检查大小
	data, err := sonic.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化value失败: %w", err)
	}
	if len(data) > MaxValueSize {
		return fmt.Errorf("value大小(%d bytes)超过限制(%d bytes)", len(data), MaxValueSize)
	}
	return nil
}

// SetSync 同步设置缓存（即使在异步模式下也等待写入完成）
func (c *SharedCache) SetSync(key string, value any, ttl ...int) error {
	if c.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	c.stats.totalSets.Add(1)

	var expireTime int64
	var ttlValue int64
	if len(ttl) > 0 && ttl[0] > 0 {
		ttlValue = int64(ttl[0])
		expireTime = time.Now().Add(time.Duration(ttl[0]) * time.Second).Unix()
	}

	item := &SharedCacheItem{
		Value:        value,
		ExpireTime:   expireTime,
		ExpirePolicy: ExpirePolicyRelative,
		TTL:          ttlValue,
		CreatedAt:    time.Now().Unix(),
		AccessCount:  0,
	}

	// 更新本地缓存
	shard := c.getShard(key)
	localItem := &localCacheItem{
		item:       item,
		localExpAt: time.Now().Add(c.localCacheTTL).UnixNano(),
	}
	shard.localCache.Store(key, localItem)

	// 异步刷盘模式下等待写入完成
	if c.config.AsyncFlush {
		doneChan := make(chan error, 1)
		op := &writeOperation{key: key, item: item, isDel: false, doneChan: doneChan}

		select {
		case c.writeBuffer <- op:
			// 等待写入完成
			select {
			case err := <-doneChan:
				return err
			case <-time.After(5 * time.Second):
				return fmt.Errorf("写入超时")
			}
		default:
			c.stats.droppedWrites.Add(1)
			return fmt.Errorf("写缓冲区已满")
		}
	}

	// 同步写入mmap
	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	data, err := c.readDataFromMmapLocked()
	if err != nil {
		return err
	}

	data.Items[key] = item
	return c.writeDataToMmapWithRetry(data)
}

// SetBatch 批量设置缓存
func (c *SharedCache) SetBatch(items map[string]any, ttl ...int) error {
	if c.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	var expireTime int64
	if len(ttl) > 0 && ttl[0] > 0 {
		expireTime = time.Now().Add(time.Duration(ttl[0]) * time.Second).Unix()
	}

	now := time.Now()
	cacheItems := make(map[string]*SharedCacheItem, len(items))

	for key, value := range items {
		item := &SharedCacheItem{
			Value:        value,
			ExpireTime:   expireTime,
			ExpirePolicy: ExpirePolicyRelative,
			CreatedAt:    now.Unix(),
			AccessCount:  0,
		}
		cacheItems[key] = item

		// 更新本地缓存
		shard := c.getShard(key)
		localItem := &localCacheItem{
			item:       item,
			localExpAt: now.Add(c.localCacheTTL).UnixNano(),
		}
		shard.localCache.Store(key, localItem)

		c.stats.totalSets.Add(1)
	}

	// 批量写入mmap
	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	data, err := c.readDataFromMmapLocked()
	if err != nil {
		return err
	}

	for key, item := range cacheItems {
		data.Items[key] = item
	}

	return c.writeDataToMmapWithRetry(data)
}

// Get 获取缓存值
func (c *SharedCache) Get(key string) (any, bool) {
	if c.closed.Load() {
		return nil, false
	}

	c.stats.totalGets.Add(1)
	shard := c.getShard(key)

	// 先查本地缓存（无锁）
	if cached, ok := shard.localCache.Load(key); ok {
		if localItem, ok := cached.(*localCacheItem); ok {
			if !localItem.isLocalExpired() && !localItem.item.IsExpired() {
				c.stats.localCacheHits.Add(1)
				// 使用原子操作更新访问计数，避免竞态
				atomic.AddInt64(&localItem.item.AccessCount, 1)

				// 滑动过期策略：更新过期时间
				if localItem.item.ExpirePolicy == ExpirePolicySliding && localItem.item.TTL > 0 {
					atomic.StoreInt64(&localItem.item.ExpireTime, time.Now().Add(time.Duration(localItem.item.TTL)*time.Second).Unix())
				}

				return localItem.item.Value, true
			}
			// 本地缓存过期，删除
			shard.localCache.Delete(key)
		}
	}

	// 本地缓存未命中，从mmap读取
	data, err := c.readDataFromMmap()
	if err != nil {
		return nil, false
	}

	item, ok := data.Items[key]
	if !ok {
		return nil, false
	}

	// 检查是否过期
	if item.IsExpired() {
		// 异步删除过期项，避免阻塞读操作
		go c.deleteExpiredKey(key)
		return nil, false
	}

	item.AccessCount++

	// 滑动过期策略：更新过期时间
	if item.ExpirePolicy == ExpirePolicySliding && item.TTL > 0 {
		item.ExpireTime = time.Now().Add(time.Duration(item.TTL) * time.Second).Unix()
		// 异步更新mmap中的过期时间
		go c.updateItemExpireTime(key, item.ExpireTime)
	}

	// 更新本地缓存
	localItem := &localCacheItem{
		item:       item,
		localExpAt: time.Now().Add(c.localCacheTTL).UnixNano(),
	}
	shard.localCache.Store(key, localItem)

	return item.Value, true
}

// GetBatch 批量获取缓存值
func (c *SharedCache) GetBatch(keys []string) map[string]any {
	if c.closed.Load() {
		return nil
	}

	result := make(map[string]any, len(keys))

	// 先从本地缓存获取
	var missedKeys []string
	for _, key := range keys {
		c.stats.totalGets.Add(1)
		shard := c.getShard(key)

		if cached, ok := shard.localCache.Load(key); ok {
			if localItem, ok := cached.(*localCacheItem); ok {
				if !localItem.isLocalExpired() && !localItem.item.IsExpired() {
					c.stats.localCacheHits.Add(1)
					result[key] = localItem.item.Value
					continue
				}
				shard.localCache.Delete(key)
			}
		}
		missedKeys = append(missedKeys, key)
	}

	// 从mmap获取未命中的
	if len(missedKeys) > 0 {
		data, err := c.readDataFromMmap()
		if err == nil {
			for _, key := range missedKeys {
				if item, ok := data.Items[key]; ok && !item.IsExpired() {
					result[key] = item.Value

					// 更新本地缓存
					shard := c.getShard(key)
					localItem := &localCacheItem{
						item:       item,
						localExpAt: time.Now().Add(c.localCacheTTL).UnixNano(),
					}
					shard.localCache.Store(key, localItem)
				}
			}
		}
	}

	return result
}

// updateItemExpireTime 更新缓存项的过期时间（仅当新时间 > 旧时间时更新，避免回退）
func (c *SharedCache) updateItemExpireTime(key string, expireTime int64) {
	if c.closed.Load() {
		return
	}

	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	data, err := c.readDataFromMmapLocked()
	if err != nil {
		return
	}

	if item, ok := data.Items[key]; ok {
		// 仅当新时间 > 旧时间时更新，避免并发场景下覆盖更晚的时间戳
		if expireTime > item.ExpireTime {
			item.ExpireTime = expireTime
			_ = c.writeDataToMmapWithRetry(data)
		}
	}
}

// deleteExpiredKey 异步删除过期的key
func (c *SharedCache) deleteExpiredKey(key string) {
	if c.closed.Load() {
		return
	}

	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	data, err := c.readDataFromMmapLocked()
	if err != nil {
		return
	}

	if item, ok := data.Items[key]; ok && item.IsExpired() {
		delete(data.Items, key)
		_ = c.writeDataToMmapWithRetry(data)
	}
}

// Delete 删除缓存
func (c *SharedCache) Delete(key string) error {
	if c.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	c.stats.totalDeletes.Add(1)

	// 删除本地缓存
	shard := c.getShard(key)
	shard.localCache.Delete(key)

	// 异步刷盘模式
	if c.config.AsyncFlush {
		select {
		case c.writeBuffer <- &writeOperation{key: key, isDel: true}:
			return nil
		default:
			c.stats.droppedWrites.Add(1)
			return fmt.Errorf("写缓冲区已满")
		}
	}

	// 删除mmap中的数据
	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	data, err := c.readDataFromMmapLocked()
	if err != nil {
		return err
	}

	delete(data.Items, key)
	return c.writeDataToMmapWithRetry(data)
}

// DeleteBatch 批量删除缓存
func (c *SharedCache) DeleteBatch(keys []string) error {
	if c.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	// 删除本地缓存
	for _, key := range keys {
		c.stats.totalDeletes.Add(1)
		shard := c.getShard(key)
		shard.localCache.Delete(key)
	}

	// 删除mmap中的数据
	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	data, err := c.readDataFromMmapLocked()
	if err != nil {
		return err
	}

	for _, key := range keys {
		delete(data.Items, key)
	}

	return c.writeDataToMmapWithRetry(data)
}

// Clear 清空所有缓存
func (c *SharedCache) Clear() error {
	if c.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	// 清空所有本地缓存
	for _, shard := range c.shards {
		shard.localCache.Range(func(key, _ any) bool {
			shard.localCache.Delete(key)
			return true
		})
	}

	// 清空mmap
	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	data := &SharedCacheData{Items: make(map[string]*SharedCacheItem)}
	return c.writeDataToMmapWithRetry(data)
}

// Has 检查缓存键是否存在
func (c *SharedCache) Has(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// Flush 手动刷盘（用于异步模式）
func (c *SharedCache) Flush() error {
	if c.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	if !c.config.AsyncFlush {
		return nil // 同步模式无需手动刷盘
	}

	// 等待写缓冲区清空
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("刷盘超时")
		case <-ticker.C:
			if len(c.writeBuffer) == 0 {
				return nil
			}
		}
	}
}

// Close 关闭缓存（优雅关闭，确保数据不丢失）
func (c *SharedCache) Close() error {
	if !c.closed.CompareAndSwap(false, true) {
		return nil // 已经关闭
	}

	c.logInfo("closing cache")

	// 停止清理协程
	close(c.cleanupDone)
	close(c.mmapCleanDone)

	// 停止异步刷盘协程并等待完成
	if c.config.AsyncFlush {
		close(c.flushDone)

		// 等待异步刷盘完成，带超时
		done := make(chan struct{})
		go func() {
			c.flushWg.Wait()
			close(done)
		}()

		select {
		case <-done:
			c.logInfo("async flush completed")
		case <-time.After(GracefulShutdownTimeout):
			c.logWarn("async flush timeout, some data may be lost")
		}
	}

	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	var errs []error

	if c.mmap != nil {
		// 最终刷盘
		if err := c.mmap.Flush(); err != nil {
			errs = append(errs, fmt.Errorf("最终刷盘失败: %w", err))
		}
		if err := c.mmap.Unmap(); err != nil {
			errs = append(errs, fmt.Errorf("取消内存映射失败: %w", err))
		}
		c.mmap = nil
	}

	if c.file != nil {
		if err := c.file.Sync(); err != nil {
			errs = append(errs, fmt.Errorf("同步文件失败: %w", err))
		}
		if err := c.file.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭文件失败: %w", err))
		}
		c.file = nil
	}

	if len(errs) > 0 {
		c.logError("close errors", zap.Errors("errors", errs))
		return errs[0]
	}

	c.logInfo("cache closed successfully")
	return nil
}

// Unlink 删除缓存文件
func (c *SharedCache) Unlink() error {
	if err := c.Close(); err != nil {
		return err
	}
	return os.Remove(c.filePath)
}

// Count 返回缓存项数量
func (c *SharedCache) Count() int {
	if c.closed.Load() {
		return 0
	}

	data, err := c.readDataFromMmap()
	if err != nil {
		return 0
	}
	return len(data.Items)
}

// Keys 返回所有缓存键
func (c *SharedCache) Keys() []string {
	if c.closed.Load() {
		return nil
	}

	data, err := c.readDataFromMmap()
	if err != nil {
		return nil
	}

	keys := make([]string, 0, len(data.Items))
	for k := range data.Items {
		keys = append(keys, k)
	}
	return keys
}

// GetFilePath 返回缓存文件路径
func (c *SharedCache) GetFilePath() string {
	return c.filePath
}

// GetStats 获取缓存统计信息
func (c *SharedCache) GetStats() *Stats {
	totalGets := c.stats.totalGets.Load()
	localHits := c.stats.localCacheHits.Load()

	var hitRate float64
	if totalGets > 0 {
		hitRate = float64(localHits) / float64(totalGets)
	}

	// 计算已使用大小
	var usedSize int64
	data, err := c.readDataFromMmap()
	if err == nil {
		jsonData, _ := sonic.Marshal(data)
		usedSize = int64(len(jsonData) + HeaderSize)
	}

	return &Stats{
		ItemCount:          c.Count(),
		LocalCacheHits:     localHits,
		MmapReads:          c.stats.mmapReads.Load(),
		MmapWrites:         c.stats.mmapWrites.Load(),
		ShardCount:         int(c.shardCount),
		LocalCacheTTLMs:    c.localCacheTTL.Milliseconds(),
		HitRate:            hitRate,
		TotalGets:          totalGets,
		TotalSets:          c.stats.totalSets.Load(),
		TotalDeletes:       c.stats.totalDeletes.Load(),
		WriteRetries:       c.stats.writeRetries.Load(),
		CRCErrors:          c.stats.crcErrors.Load(),
		ExpandCount:        c.stats.expandCount.Load(),
		CurrentSize:        int64(c.config.Size),
		UsedSize:           usedSize,
		DroppedWrites:      c.stats.droppedWrites.Load(),
		AsyncFlushCount:    c.stats.asyncFlushCount.Load(),
		BatchWriteCount:    c.stats.batchWriteCount.Load(),
		ExpiredItemsCount:  c.stats.expiredItemsCount.Load(),
		SyncFallbackWrites: c.stats.syncFallbackWrites.Load(),
		ValueSizeExceeded:  c.stats.valueSizeExceeded.Load(),
	}
}

// Expand 手动扩容
func (c *SharedCache) Expand(newSize int) error {
	if c.closed.Load() {
		return fmt.Errorf("缓存已关闭")
	}

	if newSize <= c.config.Size {
		return fmt.Errorf("新大小必须大于当前大小")
	}
	if newSize > MaxCacheSize {
		return fmt.Errorf("超过最大缓存大小限制")
	}

	c.globalRWMutex.Lock()
	defer c.globalRWMutex.Unlock()

	return c.expandLocked()
}
