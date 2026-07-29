package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/internal/config"
	"git.intra.weibo.com/search_fe/wbutil-go/internal/db/hash"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// 连接池默认配置常量（针对高并发场景优化）
const (
	// DefaultPoolSize 默认连接池大小（支持10w QPS场景）
	DefaultPoolSize = 200
	// DefaultMinIdleConns 默认最小空闲连接数
	DefaultMinIdleConns = 50
	// DefaultMaxRetries 默认最大重试次数（设置为0，快速失败，避免超时堆积）
	DefaultMaxRetries = 0
	// DefaultDialTimeout 默认连接超时
	DefaultDialTimeout = 100 * time.Millisecond
	// DefaultReadTimeout 默认读超时
	DefaultReadTimeout = 100 * time.Millisecond
	// DefaultWriteTimeout 默认写超时
	DefaultWriteTimeout = 100 * time.Millisecond
	// DefaultIdleTimeout 默认空闲超时
	DefaultIdleTimeout = 5 * time.Minute
	// DefaultIdleCheckFrequency 默认空闲检查频率
	DefaultIdleCheckFrequency = 1 * time.Minute
	// DefaultPoolTimeout 默认获取连接超时（设置为30ms，与业务超时兼容，避免504堆积）
	DefaultPoolTimeout = 30 * time.Millisecond
	// DefaultMaxConnAge 默认最大连接年龄
	DefaultMaxConnAge = 30 * time.Minute
)

// NodeConfig Redis 节点配置
type NodeConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

// ClusterConfig Redis 集群配置
type ClusterConfig struct {
	Nodes   []NodeConfig `mapstructure:"nodes"`
	Timeout int          `mapstructure:"timeout"` // 毫秒
	Auth    string       `mapstructure:"auth"`
}

// PoolConfig 连接池配置（增强版）
type PoolConfig struct {
	// 基础配置
	MaxIdle     int `mapstructure:"max_idle"`
	MaxActive   int `mapstructure:"max_active"`
	IdleTimeout int `mapstructure:"idle_timeout"` // 毫秒

	// 高级配置（新增）
	MinIdleConns       int `mapstructure:"min_idle_conns"`       // 最小空闲连接数
	MaxRetries         int `mapstructure:"max_retries"`          // 最大重试次数
	DialTimeout        int `mapstructure:"dial_timeout"`         // 连接超时（毫秒）
	ReadTimeout        int `mapstructure:"read_timeout"`         // 读超时（毫秒）
	WriteTimeout       int `mapstructure:"write_timeout"`        // 写超时（毫秒）
	PoolTimeout        int `mapstructure:"pool_timeout"`         // 获取连接超时（毫秒）
	IdleCheckFrequency int `mapstructure:"idle_check_frequency"` // 空闲检查频率（毫秒）
	MaxConnAge         int `mapstructure:"max_conn_age"`         // 最大连接年龄（毫秒）
}

// Config Redis 总配置
type Config struct {
	DefaultTimeout int                      `mapstructure:"default_timeout"`
	DefaultDB      int                      `mapstructure:"default_db"`
	Pool           PoolConfig               `mapstructure:"pool"`
	Clusters       map[string]ClusterConfig `mapstructure:"clusters"`
}

// PoolStats 连接池统计信息
type PoolStats struct {
	Hits       uint32 // 连接池命中次数
	Misses     uint32 // 连接池未命中次数
	Timeouts   uint32 // 等待超时次数
	TotalConns uint32 // 总连接数
	IdleConns  uint32 // 空闲连接数
	StaleConns uint32 // 过期连接数
}

// Manager Redis 管理器
type Manager struct {
	config       *Config
	clients      map[string][]*redis.Client // busKey -> clients
	mu           sync.RWMutex
	configLoaded bool   // 配置是否已加载
	configPath   string // app.yaml 配置文件路径

	// 统计信息
	stats struct {
		totalRequests atomic.Int64
		totalErrors   atomic.Int64
		totalLatency  atomic.Int64 // 纳秒
	}
}

// ConnectArgs 连接参数
type ConnectArgs struct {
	BusKey  string // 业务集群 key
	HashKey string // hash 分片 key
	HashNo  int    // 指定节点编号（从1开始，0表示使用hash）
	Timeout int    // 超时时间（毫秒），0使用默认值
	DB      int    // 数据库编号，-1使用配置默认值
}

// Connection Redis 连接封装
type Connection struct {
	client    *redis.Client
	busKey    string
	nodeIndex int
	startTime time.Time
	config    *NodeConfig
	timeout   time.Duration
	db        int
	manager   *Manager // 引用管理器用于统计
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager 获取 Redis 管理器单例
func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			clients:      make(map[string][]*redis.Client),
			configLoaded: false,
			configPath:   "config/app.yaml", // 默认配置路径
		}
	})
	return instance
}

// SetConfigPath 设置 app.yaml 配置文件路径（可选，默认为 "config/app.yaml"）
func (m *Manager) SetConfigPath(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.configPath = path
	m.configLoaded = false // 重置加载状态
}

// loadConfig 从 app.yaml 加载 Redis 配置（内部方法）
// 使用 config.ComponentViper() 懒加载方式读取配置
func (m *Manager) loadConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果已加载，直接返回
	if m.configLoaded {
		return nil
	}

	// 1. 使用 ComponentViper 获取 app.yaml 配置
	v, err := config.ComponentViper()
	if err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	// 2. 获取 redis.config_path
	redisConfigPath := v.GetString("redis.config_path")
	if redisConfigPath == "" {
		return fmt.Errorf("app.yaml 中未配置 redis.config_path")
	}

	// 3. 将相对路径转为相对配置文件目录的绝对路径
	redisConfigPath = config.AbsPathFromConfigFile(v, redisConfigPath)

	// 4. 读取 redis 配置文件
	redisViper := viper.New()
	redisViper.SetConfigFile(redisConfigPath)
	redisViper.SetConfigType("yaml")

	if err := redisViper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取 Redis 配置文件失败 %s: %w", redisConfigPath, err)
	}

	// 5. 解析 redis 配置
	var cfg Config
	if err := redisViper.UnmarshalKey("redis", &cfg); err != nil {
		return fmt.Errorf("解析 Redis 配置失败: %w", err)
	}

	// 6. 设置默认值
	m.applyDefaultPoolConfig(&cfg.Pool)

	m.config = &cfg
	m.configLoaded = true
	return nil
}

// applyDefaultPoolConfig 应用默认连接池配置
func (m *Manager) applyDefaultPoolConfig(pool *PoolConfig) {
	if pool.MaxActive <= 0 {
		pool.MaxActive = DefaultPoolSize
	}
	if pool.MinIdleConns <= 0 {
		pool.MinIdleConns = DefaultMinIdleConns
	}
	// MaxRetries 默认为 0，表示不重试，快速失败
	// 注意：只有 < 0 时才使用默认值，允许设置为 0
	if pool.MaxRetries < 0 {
		pool.MaxRetries = DefaultMaxRetries
	}
	if pool.DialTimeout <= 0 {
		pool.DialTimeout = int(DefaultDialTimeout.Milliseconds())
	}
	if pool.ReadTimeout <= 0 {
		pool.ReadTimeout = int(DefaultReadTimeout.Milliseconds())
	}
	if pool.WriteTimeout <= 0 {
		pool.WriteTimeout = int(DefaultWriteTimeout.Milliseconds())
	}
	if pool.IdleTimeout <= 0 {
		pool.IdleTimeout = int(DefaultIdleTimeout.Milliseconds())
	}
	if pool.PoolTimeout <= 0 {
		pool.PoolTimeout = int(DefaultPoolTimeout.Milliseconds())
	}
	if pool.IdleCheckFrequency <= 0 {
		pool.IdleCheckFrequency = int(DefaultIdleCheckFrequency.Milliseconds())
	}
	if pool.MaxConnAge <= 0 {
		pool.MaxConnAge = int(DefaultMaxConnAge.Milliseconds())
	}
}

// GetConnect 获取 Redis 连接（自动加载配置）
func (m *Manager) GetConnect(args ConnectArgs) (*Connection, error) {
	startTime := time.Now()
	m.stats.totalRequests.Add(1)

	// 懒加载：首次调用时自动加载配置
	if !m.configLoaded {
		if err := m.loadConfig(); err != nil {
			m.stats.totalErrors.Add(1)
			return nil, fmt.Errorf("自动加载配置失败: %w", err)
		}
	}

	if m.config == nil {
		m.stats.totalErrors.Add(1)
		return nil, fmt.Errorf("配置加载失败")
	}

	if args.BusKey == "" {
		m.stats.totalErrors.Add(1)
		return nil, fmt.Errorf("busKey 不能为空")
	}

	cluster, ok := m.config.Clusters[args.BusKey]
	if !ok {
		m.stats.totalErrors.Add(1)
		return nil, fmt.Errorf("集群 %s 不存在", args.BusKey)
	}

	if len(cluster.Nodes) == 0 {
		m.stats.totalErrors.Add(1)
		return nil, fmt.Errorf("集群 %s 没有配置节点", args.BusKey)
	}

	// 确定节点索引
	nodeIndex := 0
	if args.HashNo > 0 {
		// 指定节点
		nodeIndex = args.HashNo - 1
		if nodeIndex >= len(cluster.Nodes) {
			m.stats.totalErrors.Add(1)
			return nil, fmt.Errorf("节点索引 %d 超出范围", args.HashNo)
		}
	} else if len(cluster.Nodes) > 1 {
		// 多节点需要 hash
		if args.HashKey == "" {
			m.stats.totalErrors.Add(1)
			return nil, fmt.Errorf("多节点集群需要提供 hashKey")
		}
		nodeIndex = m.selectNode(cluster.Nodes, args.HashKey)
	}

	node := cluster.Nodes[nodeIndex]

	// 确定超时时间
	timeout := time.Duration(m.config.DefaultTimeout) * time.Millisecond
	if args.Timeout > 0 {
		timeout = time.Duration(args.Timeout) * time.Millisecond
	} else if cluster.Timeout > 0 {
		timeout = time.Duration(cluster.Timeout) * time.Millisecond
	}

	// 确定数据库
	db := m.config.DefaultDB
	if node.DB > 0 {
		db = node.DB
	}
	if args.DB > 0 {
		db = args.DB
	}

	// 创建或获取客户端
	client, err := m.getOrCreateClient(args.BusKey, nodeIndex, &node, cluster.Auth, db, timeout)
	if err != nil {
		m.stats.totalErrors.Add(1)
		return nil, err
	}

	// 记录延迟
	m.stats.totalLatency.Add(time.Since(startTime).Nanoseconds())

	return &Connection{
		client:    client,
		busKey:    args.BusKey,
		nodeIndex: nodeIndex,
		startTime: time.Now(),
		config:    &node,
		timeout:   timeout,
		db:        db,
		manager:   m,
	}, nil
}

// selectNode 使用 CRC64 选择节点，与 Lua/C 实现一致
func (m *Manager) selectNode(nodes []NodeConfig, hashKey string) int {
	return hash.CRC64StringMod(hashKey, len(nodes))
}

// getOrCreateClient 获取或创建客户端（优化版）
func (m *Manager) getOrCreateClient(busKey string, nodeIndex int, node *NodeConfig, auth string, db int, timeout time.Duration) (*redis.Client, error) {
	// 先用读锁检查
	m.mu.RLock()
	if clients, ok := m.clients[busKey]; ok && nodeIndex < len(clients) && clients[nodeIndex] != nil {
		client := clients[nodeIndex]
		m.mu.RUnlock()
		return client, nil
	}
	m.mu.RUnlock()

	// 需要创建，使用写锁
	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查
	if clients, ok := m.clients[busKey]; ok && nodeIndex < len(clients) && clients[nodeIndex] != nil {
		return clients[nodeIndex], nil
	}

	// 创建新客户端（使用优化的连接池配置）
	addr := fmt.Sprintf("%s:%d", node.Host, node.Port)
	pool := m.config.Pool

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: auth,
		DB:       db,

		// 超时配置
		DialTimeout:  time.Duration(pool.DialTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(pool.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(pool.WriteTimeout) * time.Millisecond,

		// 连接池配置（针对高并发优化）
		PoolSize:           pool.MaxActive,
		MinIdleConns:       pool.MinIdleConns,
		MaxConnAge:         time.Duration(pool.MaxConnAge) * time.Millisecond,
		PoolTimeout:        time.Duration(pool.PoolTimeout) * time.Millisecond,
		IdleTimeout:        time.Duration(pool.IdleTimeout) * time.Millisecond,
		IdleCheckFrequency: time.Duration(pool.IdleCheckFrequency) * time.Millisecond,

		// 重试配置
		MaxRetries:      pool.MaxRetries,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,

		// 连接钩子（可用于监控）
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			// 可以在这里添加连接建立时的监控逻辑
			return nil
		},
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("连接 Redis 失败 %s: %w", addr, err)
	}

	// 保存客户端
	if _, ok := m.clients[busKey]; !ok {
		m.clients[busKey] = make([]*redis.Client, len(m.config.Clusters[busKey].Nodes))
	}
	m.clients[busKey][nodeIndex] = client

	return client, nil
}

// GetPoolStats 获取指定集群的连接池统计信息
func (m *Manager) GetPoolStats(busKey string, nodeIndex int) (*PoolStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clients, ok := m.clients[busKey]
	if !ok || nodeIndex >= len(clients) || clients[nodeIndex] == nil {
		return nil, fmt.Errorf("客户端不存在")
	}

	stats := clients[nodeIndex].PoolStats()
	return &PoolStats{
		Hits:       stats.Hits,
		Misses:     stats.Misses,
		Timeouts:   stats.Timeouts,
		TotalConns: stats.TotalConns,
		IdleConns:  stats.IdleConns,
		StaleConns: stats.StaleConns,
	}, nil
}

// GetAllPoolStats 获取所有连接池统计信息
func (m *Manager) GetAllPoolStats() map[string][]PoolStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string][]PoolStats)
	for busKey, clients := range m.clients {
		stats := make([]PoolStats, len(clients))
		for i, client := range clients {
			if client != nil {
				s := client.PoolStats()
				stats[i] = PoolStats{
					Hits:       s.Hits,
					Misses:     s.Misses,
					Timeouts:   s.Timeouts,
					TotalConns: s.TotalConns,
					IdleConns:  s.IdleConns,
					StaleConns: s.StaleConns,
				}
			}
		}
		result[busKey] = stats
	}
	return result
}

// GetManagerStats 获取管理器统计信息
func (m *Manager) GetManagerStats() map[string]int64 {
	return map[string]int64{
		"total_requests": m.stats.totalRequests.Load(),
		"total_errors":   m.stats.totalErrors.Load(),
		"avg_latency_ns": func() int64 {
			total := m.stats.totalRequests.Load()
			if total == 0 {
				return 0
			}
			return m.stats.totalLatency.Load() / total
		}(),
	}
}

// Close 关闭所有连接
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for busKey, clients := range m.clients {
		for i, client := range clients {
			if client != nil {
				if err := client.Close(); err != nil {
					lastErr = fmt.Errorf("关闭 %s[%d] 失败: %w", busKey, i, err)
				}
			}
		}
	}
	m.clients = make(map[string][]*redis.Client)
	return lastErr
}

// ===== Connection 方法 =====

// GetClient 获取原生 redis.Client
func (c *Connection) GetClient() *redis.Client {
	return c.client
}

// Exec 执行 Redis 命令
func (c *Connection) Exec(ctx context.Context, cmd string, args ...any) (any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}

	// 使用 Do 方法执行任意命令
	cmdArgs := make([]any, len(args)+1)
	cmdArgs[0] = cmd
	copy(cmdArgs[1:], args)

	result := c.client.Do(ctx, cmdArgs...)
	return result.Val(), result.Err()
}

// Get 获取字符串值
// Get 获取字符串值
// 注意：如果 key 不存在，返回 ("", redis.Nil)，这是正常业务情况
func (c *Connection) Get(ctx context.Context, key string) (string, error) {
	// 始终使用连接创建时设置的超时时间，确保超时控制一致
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.Get(ctx, key).Result()
}
// Set 设置字符串值
func (c *Connection) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.Set(ctx, key, value, expiration).Err()
}

// SetNX 设置字符串值（仅当key不存在时）
func (c *Connection) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.SetNX(ctx, key, value, expiration).Result()
}

// SetEX 设置字符串值并指定过期时间
func (c *Connection) SetEX(ctx context.Context, key string, value any, expiration time.Duration) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.SetEX(ctx, key, value, expiration).Err()
}

// MGet 批量获取
func (c *Connection) MGet(ctx context.Context, keys ...string) ([]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.MGet(ctx, keys...).Result()
}

// MSet 批量设置
func (c *Connection) MSet(ctx context.Context, values ...any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.MSet(ctx, values...).Err()
}

// Del 删除键
func (c *Connection) Del(ctx context.Context, keys ...string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (c *Connection) Exists(ctx context.Context, keys ...string) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (c *Connection) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (c *Connection) TTL(ctx context.Context, key string) (time.Duration, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.TTL(ctx, key).Result()
}

// HGet 获取 hash 字段值
// 注意：如果 field 不存在，返回 ("", redis.Nil)，这是正常业务情况
func (c *Connection) HGet(ctx context.Context, key, field string) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.HGet(ctx, key, field).Result()
}

// HSet 设置 hash 字段值
func (c *Connection) HSet(ctx context.Context, key string, values ...any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.HSet(ctx, key, values...).Err()
}

// HGetAll 获取 hash 所有字段
func (c *Connection) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.HGetAll(ctx, key).Result()
}

// HMGet 批量获取 hash 字段
func (c *Connection) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.HMGet(ctx, key, fields...).Result()
}

// HMSet 批量设置 hash 字段
func (c *Connection) HMSet(ctx context.Context, key string, values ...any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.HMSet(ctx, key, values...).Err()
}

// HDel 删除 hash 字段
func (c *Connection) HDel(ctx context.Context, key string, fields ...string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.HDel(ctx, key, fields...).Err()
}

// HExists 检查 hash 字段是否存在
func (c *Connection) HExists(ctx context.Context, key, field string) (bool, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.HExists(ctx, key, field).Result()
}

// HIncrBy hash 字段自增
func (c *Connection) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.HIncrBy(ctx, key, field, incr).Result()
}

// ZAdd 添加有序集合成员
func (c *Connection) ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合范围
func (c *Connection) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取有序集合范围（带分数）
func (c *Connection) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZRevRange 逆序获取有序集合范围
func (c *Connection) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.ZRevRange(ctx, key, start, stop).Result()
}

// ZRevRangeWithScores 逆序获取有序集合范围（带分数）
func (c *Connection) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

// ZRem 删除有序集合成员
func (c *Connection) ZRem(ctx context.Context, key string, members ...any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.ZRem(ctx, key, members...).Err()
}

// ZScore 获取有序集合成员分数
func (c *Connection) ZScore(ctx context.Context, key, member string) (float64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.ZScore(ctx, key, member).Result()
}

// ZCard 获取有序集合成员数量
func (c *Connection) ZCard(ctx context.Context, key string) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.ZCard(ctx, key).Result()
}

// LPush 从列表左侧推入
func (c *Connection) LPush(ctx context.Context, key string, values ...any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.LPush(ctx, key, values...).Err()
}

// RPush 从列表右侧推入
func (c *Connection) RPush(ctx context.Context, key string, values ...any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.RPush(ctx, key, values...).Err()
}

// LPop 从列表左侧弹出
func (c *Connection) LPop(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.LPop(ctx, key).Result()
}

// RPop 从列表右侧弹出
func (c *Connection) RPop(ctx context.Context, key string) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.RPop(ctx, key).Result()
}

// LRange 获取列表范围
func (c *Connection) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.LRange(ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func (c *Connection) LLen(ctx context.Context, key string) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.LLen(ctx, key).Result()
}

// SAdd 添加集合成员
func (c *Connection) SAdd(ctx context.Context, key string, members ...any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (c *Connection) SMembers(ctx context.Context, key string) ([]string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.SMembers(ctx, key).Result()
}

// SRem 删除集合成员
func (c *Connection) SRem(ctx context.Context, key string, members ...any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.SRem(ctx, key, members...).Err()
}

// SIsMember 检查是否是集合成员
func (c *Connection) SIsMember(ctx context.Context, key string, member any) (bool, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.SIsMember(ctx, key, member).Result()
}

// SCard 获取集合成员数量
func (c *Connection) SCard(ctx context.Context, key string) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.SCard(ctx, key).Result()
}

// Incr 自增
func (c *Connection) Incr(ctx context.Context, key string) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.Incr(ctx, key).Result()
}

// Decr 自减
func (c *Connection) Decr(ctx context.Context, key string) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.Decr(ctx, key).Result()
}

// IncrBy 增加指定值
func (c *Connection) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.IncrBy(ctx, key, value).Result()
}

// DecrBy 减少指定值
func (c *Connection) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.DecrBy(ctx, key, value).Result()
}

// Pipeline 创建管道
func (c *Connection) Pipeline() redis.Pipeliner {
	return c.client.Pipeline()
}

// TxPipeline 创建事务管道
func (c *Connection) TxPipeline() redis.Pipeliner {
	return c.client.TxPipeline()
}

// Eval 执行 Lua 脚本
func (c *Connection) Eval(ctx context.Context, script string, keys []string, args ...any) (any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.Eval(ctx, script, keys, args...).Result()
}

// EvalSha 执行已缓存的 Lua 脚本
func (c *Connection) EvalSha(ctx context.Context, sha1 string, keys []string, args ...any) (any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
	}
	return c.client.EvalSha(ctx, sha1, keys, args...).Result()
}

// GetCost 获取连接耗时（毫秒）
func (c *Connection) GetCost() int64 {
	return time.Since(c.startTime).Milliseconds()
}

// GetNodeInfo 获取节点信息
func (c *Connection) GetNodeInfo() string {
	return fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
}

func (c *Connection) GetPort() string {
	return strconv.Itoa(c.config.Port)
}

// GetPoolStats 获取当前连接的连接池统计
func (c *Connection) GetPoolStats() *PoolStats {
	stats := c.client.PoolStats()
	return &PoolStats{
		Hits:       stats.Hits,
		Misses:     stats.Misses,
		Timeouts:   stats.Timeouts,
		TotalConns: stats.TotalConns,
		IdleConns:  stats.IdleConns,
		StaleConns: stats.StaleConns,
	}
}

// ===== 包级别便捷函数 =====

// GetNodeInfo 获取节点信息（即使连接失败也能获取）
func (m *Manager) GetNodeInfo(args ConnectArgs) string {
	// 懒加载：首次调用时自动加载配置
	if !m.configLoaded {
		if err := m.loadConfig(); err != nil {
			return ""
		}
	}

	if m.config == nil || args.BusKey == "" {
		return ""
	}

	cluster, ok := m.config.Clusters[args.BusKey]
	if !ok || len(cluster.Nodes) == 0 {
		return ""
	}

	// 确定节点索引
	nodeIndex := 0
	if args.HashNo > 0 {
		nodeIndex = args.HashNo - 1
		if nodeIndex >= len(cluster.Nodes) {
			return ""
		}
	} else if len(cluster.Nodes) > 1 {
		if args.HashKey == "" {
			return ""
		}
		nodeIndex = m.selectNode(cluster.Nodes, args.HashKey)
	}

	node := cluster.Nodes[nodeIndex]
	return fmt.Sprintf("%s:%d", node.Host, node.Port)
}

// GetConnect 包级别便捷函数，直接获取 Redis 连接
func GetConnect(args ConnectArgs) (*Connection, error) {
	return GetManager().GetConnect(args)
}

// GetNodeInfo 包级别便捷函数，获取节点信息
func GetNodeInfo(args ConnectArgs) string {
	return GetManager().GetNodeInfo(args)
}

// GetNodePort 包级别便捷函数，获取节点端口号
func GetNodePort(args ConnectArgs) string {
	nodeInfo := GetManager().GetNodeInfo(args)
	if nodeInfo == "" {
		return ""
	}
	// 从 host:port 格式中提取端口号
	if idx := strings.LastIndex(nodeInfo, ":"); idx != -1 {
		return nodeInfo[idx+1:]
	}
	return nodeInfo
}

// SetConfigPath 包级别便捷函数，设置配置文件路径
func SetConfigPath(path string) {
	GetManager().SetConfigPath(path)
}

// Close 包级别便捷函数，关闭所有 Redis 连接
func Close() error {
	return GetManager().Close()
}

// GetPoolStats 包级别便捷函数，获取连接池统计
func GetPoolStats(busKey string, nodeIndex int) (*PoolStats, error) {
	return GetManager().GetPoolStats(busKey, nodeIndex)
}

// GetAllPoolStats 包级别便捷函数，获取所有连接池统计
func GetAllPoolStats() map[string][]PoolStats {
	return GetManager().GetAllPoolStats()
}

// GetManagerStats 包级别便捷函数，获取管理器统计
func GetManagerStats() map[string]int64 {
	return GetManager().GetManagerStats()
}
