package redis

import (
	"errors"
	"fmt"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/log"
	internalRedis "git.intra.weibo.com/search_fe/wbutil-go/internal/db/redis"
	"github.com/go-redis/redis/v8"
)

// Manager Redis 管理器外观
type Manager struct {
	impl *internalRedis.Manager
}

// Connection Redis 连接外观
type Connection struct {
	impl    *internalRedis.Connection
	busKey  string       // 业务集群 key，用于日志记录
	request []RequestLog // 请求日志记录
}

// RequestLog 请求日志记录
type RequestLog struct {
	Command string // 命令名称
	Args    []any  // 命令参数
	Result  any    // 返回结果
	Ok      bool   // 是否成功
}

// ConnectArgs 连接参数
type ConnectArgs struct {
	BusKey  string // 业务集群 key
	HashKey string // hash 分片 key
	HashNo  int    // 指定节点编号（从1开始，0表示使用hash）
	Timeout int    // 超时时间（毫秒），0使用默认值
	DB      int    // 数据库编号，-1使用配置默认值
}

// GetManager 获取 Redis 管理器单例
func GetManager() *Manager {
	return &Manager{
		impl: internalRedis.GetManager(),
	}
}

// 示例:
//
//	conn, err := manager.GetConnect(ConnectArgs{
//	    BusKey: "fe_m_redis",
//	    HashKey: "user_123",
//	    Timeout: 500,
//	})
func (m *Manager) GetConnect(args ConnectArgs) (*Connection, error) {
	internalArgs := internalRedis.ConnectArgs{
		BusKey:  args.BusKey,
		HashKey: args.HashKey,
		HashNo:  args.HashNo,
		Timeout: args.Timeout,
		DB:      args.DB,
	}

	conn, err := m.impl.GetConnect(internalArgs)
	if err != nil {
		return nil, err
	}

	return &Connection{
		impl:    conn,
		busKey:  args.BusKey,
		request: make([]RequestLog, 0),
	}, nil
}

// Close 关闭所有连接
func (m *Manager) Close() error {
	return m.impl.Close()
}

// ===== Connection 方法 =====

// GetClient 获取原生 redis.Client，用于执行更复杂的操作
func (c *Connection) GetClient() *redis.Client {
	return c.impl.GetClient()
}

// Exec 执行任意 Redis 命令
// ctx: 上下文，nil 则使用默认超时
// cmd: 命令名称，如 "GET", "SET"
// args: 命令参数
// 自动记录日志，与 Lua 版本 _exec 行为一致
func (c *Connection) Exec(ctx *context.Context, cmd string, args ...any) (any, error) {
	res, err := c.impl.Exec(ctx.GinContext(), cmd, args...)
	c.addRequest(cmd, args, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// Get 获取字符串值
// 自动记录日志
func (c *Connection) Get(ctx *context.Context, key string) (string, error) {
	res, err := c.impl.Get(ctx.GinContext(), key)
	c.addRequest("GET", []any{key}, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// Set 设置字符串值
// expiration: 过期时间，0 表示永不过期
// 自动记录日志
func (c *Connection) Set(ctx *context.Context, key string, value any, expiration time.Duration) error {
	err := c.impl.Set(ctx.GinContext(), key, value, expiration)
	c.addRequest("SET", []any{key, value, expiration}, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// Del 删除键
// 自动记录日志
func (c *Connection) Del(ctx *context.Context, keys ...string) error {
	args := make([]any, len(keys))
	for i, k := range keys {
		args[i] = k
	}
	err := c.impl.Del(ctx.GinContext(), keys...)
	c.addRequest("DEL", args, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// Exists 检查键是否存在
// 返回存在的键数量
// 自动记录日志
func (c *Connection) Exists(ctx *context.Context, keys ...string) (int64, error) {
	args := make([]any, len(keys))
	for i, k := range keys {
		args[i] = k
	}
	res, err := c.impl.Exists(ctx.GinContext(), keys...)
	c.addRequest("EXISTS", args, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// Expire 设置过期时间
// 自动记录日志
func (c *Connection) Expire(ctx *context.Context, key string, expiration time.Duration) error {
	err := c.impl.Expire(ctx.GinContext(), key, expiration)
	c.addRequest("EXPIRE", []any{key, expiration}, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// HGet 获取 hash 字段值
// 自动记录日志
func (c *Connection) HGet(ctx *context.Context, key, field string) (string, error) {
	res, err := c.impl.HGet(ctx.GinContext(), key, field)
	c.addRequest("HGET", []any{key, field}, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// HSet 设置 hash 字段值
// values: 字段和值交替出现，如 "field1", "value1", "field2", "value2"
// 自动记录日志
func (c *Connection) HSet(ctx *context.Context, key string, values ...any) error {
	args := append([]any{key}, values...)
	err := c.impl.HSet(ctx.GinContext(), key, values...)
	c.addRequest("HSET", args, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// HGetAll 获取 hash 所有字段
// 自动记录日志
func (c *Connection) HGetAll(ctx *context.Context, key string) (map[string]string, error) {
	res, err := c.impl.HGetAll(ctx.GinContext(), key)
	c.addRequest("HGETALL", []any{key}, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// HDel 删除 hash 字段
// 自动记录日志
func (c *Connection) HDel(ctx *context.Context, key string, fields ...string) error {
	args := make([]any, len(fields)+1)
	args[0] = key
	for i, f := range fields {
		args[i+1] = f
	}
	err := c.impl.HDel(ctx.GinContext(), key, fields...)
	c.addRequest("HDEL", args, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// 成员
// 自动记录日志
func (c *Connection) ZAdd(ctx *context.Context, key string, members ...*redis.Z) error {
	args := make([]any, len(members)+1)
	args[0] = key
	for i, m := range members {
		args[i+1] = m
	}
	err := c.impl.ZAdd(ctx.GinContext(), key, members...)
	c.addRequest("ZADD", args, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// ZRange 获取有序集合范围
// 自动记录日志
func (c *Connection) ZRange(ctx *context.Context, key string, start, stop int64) ([]string, error) {
	res, err := c.impl.ZRange(ctx.GinContext(), key, start, stop)
	c.addRequest("ZRANGE", []any{key, start, stop}, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// ZRangeWithScores 获取有序集合范围（带分数）
// 自动记录日志
func (c *Connection) ZRangeWithScores(ctx *context.Context, key string, start, stop int64) ([]redis.Z, error) {
	res, err := c.impl.ZRangeWithScores(ctx.GinContext(), key, start, stop)
	c.addRequest("ZRANGEWITHSCORES", []any{key, start, stop}, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// ZRem 删除有序集合成员
// 自动记录日志
func (c *Connection) ZRem(ctx *context.Context, key string, members ...any) error {
	args := append([]any{key}, members...)
	err := c.impl.ZRem(ctx.GinContext(), key, members...)
	c.addRequest("ZREM", args, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// LPush 从列表左侧推入
// 自动记录日志
func (c *Connection) LPush(ctx *context.Context, key string, values ...any) error {
	args := append([]any{key}, values...)
	err := c.impl.LPush(ctx.GinContext(), key, values...)
	c.addRequest("LPUSH", args, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// RPush 从列表右侧推入
// 自动记录日志
func (c *Connection) RPush(ctx *context.Context, key string, values ...any) error {
	args := append([]any{key}, values...)
	err := c.impl.RPush(ctx.GinContext(), key, values...)
	c.addRequest("RPUSH", args, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// LRange 获取列表范围
// 自动记录日志
func (c *Connection) LRange(ctx *context.Context, key string, start, stop int64) ([]string, error) {
	res, err := c.impl.LRange(ctx.GinContext(), key, start, stop)
	c.addRequest("LRANGE", []any{key, start, stop}, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// SAdd 添加集合成员
// 自动记录日志
func (c *Connection) SAdd(ctx *context.Context, key string, members ...any) error {
	args := append([]any{key}, members...)
	err := c.impl.SAdd(ctx.GinContext(), key, members...)
	c.addRequest("SADD", args, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// SMembers 获取集合所有成员
// 自动记录日志
func (c *Connection) SMembers(ctx *context.Context, key string) ([]string, error) {
	res, err := c.impl.SMembers(ctx.GinContext(), key)
	c.addRequest("SMEMBERS", []any{key}, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// SRem 删除集合成员
// 自动记录日志
func (c *Connection) SRem(ctx *context.Context, key string, members ...any) error {
	args := append([]any{key}, members...)
	err := c.impl.SRem(ctx.GinContext(), key, members...)
	c.addRequest("SREM", args, nil, err == nil)
	c.logResult(ctx, err)
	return err
}

// Incr 自增
// 自动记录日志
func (c *Connection) Incr(ctx *context.Context, key string) (int64, error) {
	res, err := c.impl.Incr(ctx.GinContext(), key)
	c.addRequest("INCR", []any{key}, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// Decr 自减
// 自动记录日志
func (c *Connection) Decr(ctx *context.Context, key string) (int64, error) {
	res, err := c.impl.Decr(ctx.GinContext(), key)
	c.addRequest("DECR", []any{key}, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// IncrBy 增加指定值
// 自动记录日志
func (c *Connection) IncrBy(ctx *context.Context, key string, value int64) (int64, error) {
	res, err := c.impl.IncrBy(ctx.GinContext(), key, value)
	c.addRequest("INCRBY", []any{key, value}, res, err == nil)
	c.logResult(ctx, err)
	return res, err
}

// GetCost 获取连接耗时（毫秒）
func (c *Connection) GetCost() int64 {
	return c.impl.GetCost()
}

// GetNodeInfo 获取节点信息
func (c *Connection) GetNodeInfo() string {
	return c.impl.GetNodeInfo()
}

// GetBusKey 获取业务集群 key
func (c *Connection) GetBusKey() string {
	return c.busKey
}

// GetRequest 获取请求日志记录
func (c *Connection) GetRequest() []RequestLog {
	return c.request
}

// addRequest 添加请求日志记录（内部方法）
func (c *Connection) addRequest(command string, args []any, result any, ok bool) {
	c.request = append(c.request, RequestLog{
		Command: command,
		Args:    args,
		Result:  result,
		Ok:      ok,
	})
}

// logResult 统一记录操作结果日志（内部方法）
// 根据 err 是否为 nil 自动调用 LogSuccess 或 LogError
// 注意：redis.Nil（key 不存在）是正常业务情况，不记录为错误
func (c *Connection) logResult(ctx *context.Context, err error) {
	if err != nil && !errors.Is(err, redis.Nil) {
		c.LogError(ctx, 504, err)
	} else {
		c.LogSuccess(ctx)
	}
}

// ===== 日志方法 =====

// Log 记录 Redis 操作日志
// 日志格式与 Lua 版本保持一致：network:redis|{busKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
// statusCode: HTTP 状态码，200 表示成功，502 表示连接失败，504 表示执行失败
// errMsg: 错误信息，成功时为 "success"
func (c *Connection) Log(ctx *context.Context, statusCode int, errMsg string) {
	// 默认值处理，与 Lua 版本保持一致
	if statusCode == 0 {
		statusCode = 200
	}
	if errMsg == "" {
		errMsg = "success"
	}
	// 如果有错误信息且不是 success，状态码设为 502
	if errMsg != "success" && statusCode == 200 {
		statusCode = 502
	}

	// 计算耗时（秒，保留3位小数）
	cost := float64(c.impl.GetCost()) / 1000.0

	// 格式化日志消息，与 Lua 版本格式一致
	// network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
	msg := fmt.Sprintf("network:%s|%s|%.3f|%s|%d|%d|%d|0|msg:%s",
		"redis",
		c.busKey,
		cost,
		c.impl.GetPort(),
		0,
		1,
		statusCode,
		errMsg,
	)

	// 追加到请求日志
	if ctx != nil {
		log.Append(ctx, msg)
	}
}

// LogSuccess 记录成功的 Redis 操作日志
func (c *Connection) LogSuccess(ctx *context.Context) {
	c.Log(ctx, 200, "success")
}

// LogError 记录失败的 Redis 操作日志
// statusCode: 504 表示执行失败，502 表示连接失败
func (c *Connection) LogError(ctx *context.Context, statusCode int, err error) {
	errMsg := "timeout"
	if err != nil {
		// 记录具体错误信息，便于排查问题
		errMsg = err.Error()
		// 如果错误信息为空，则使用默认的 timeout
		if errMsg == "" {
			errMsg = "timeout"
		}
	}
	if statusCode == 0 {
		statusCode = 504
	}
	c.Log(ctx, statusCode, errMsg)
}

// ===== 辅助方法 =====

// Format 格式化有序集合数据为结构化数组
// data: angeWithScores 返回的数据
// fields: 字段名称，默认为 ["data", "score"]
//
// 示例:
//
//	scores, _ := conn.ctx, "key", 0, -1)
//	formatted := FormatZSet(scores, "member", "score")
//	// 返回: [{"member": "a", "score": 1.0}, {"member": "b", "score": 2.0}]
func FormatZSet(data []redis.Z, fields ...string) []map[string]any {
	if len(data) == 0 {
		return []map[string]any{}
	}

	firstField, secondField := "data", "score"
	if len(fields) >= 2 {
		firstField, secondField = fields[0], fields[1]
	}

	result := make([]map[string]any, len(data))
	for i, item := range data {
		result[i] = map[string]any{
			firstField:  item.Member,
			secondField: item.Score,
		}
	}

	return result
}

// FormatHash 格式化 hash 数据（已经是 map 格式，直接返回）
func FormatHash(data map[string]string) map[string]string {
	return data
}

// ===== 包级别便捷函数 =====

// GetConnect 包级别便捷函数，直接获取 Redis 连接
// 自动使用单例 Manager，配置会在首次调用时自动加载
// 连接失败时会自动记录日志（如果提供了 ctx）
//
// 示例:
//
//	conn, err := redis.GetConnect(ctx, redis.ConnectArgs{
//	    BusKey: "fe_m_redis",
//	})
func GetConnect(ctx *context.Context, args ConnectArgs) (*Connection, error) {
	internalArgs := internalRedis.ConnectArgs{
		BusKey:  args.BusKey,
		HashKey: args.HashKey,
		HashNo:  args.HashNo,
		Timeout: args.Timeout,
		DB:      args.DB,
	}

	conn, err := internalRedis.GetConnect(internalArgs)
	if err != nil {
		// 连接失败时自动记录日志（使用完整参数获取正确的节点信息）
		LogConnectErrorWithArgs(ctx, args, err)
		return nil, err
	}

	return &Connection{
		impl:    conn,
		busKey:  args.BusKey,
		request: make([]RequestLog, 0),
	}, nil
}

// LogConnectError 记录 Redis 连接失败的日志
// 日志格式与 Lua 版本保持一致：network:redis|{busKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
func LogConnectError(ctx *context.Context, busKey string, hashKey string, err error) {
	LogConnectErrorWithArgs(ctx, ConnectArgs{BusKey: busKey, HashKey: hashKey}, err)
}

// LogConnectErrorWithArgs 记录 Redis 连接失败的日志（使用完整的连接参数）
// 日志格式与 Lua 版本保持一致：network:redis|{busKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
func LogConnectErrorWithArgs(ctx *context.Context, args ConnectArgs, err error) {
	errMsg := "timeout"
	if err != nil {
		// 暂时全部默认timeout
		// errMsg = err.Error()
	}

	// 获取节点端口号（即使连接失败也能获取）
	port := internalRedis.GetNodePort(internalRedis.ConnectArgs{
		BusKey:  args.BusKey,
		HashKey: args.HashKey,
		HashNo:  args.HashNo,
	})
	if port == "" {
		port = args.HashKey // 如果无法获取节点信息，使用 hashKey 作为备选
	}

	// 格式化日志消息，与 Lua 版本格式一致
	// network:{protocol}|{urlKey}|{cost}|{nodeLog}|{len}|{tryTimes}|{statusCode}|{web_degrade}|msg:{err}
	msg := fmt.Sprintf("network:%s|%s|%.3f|%s|%d|%d|%d|0|msg:%s",
		"redis",
		args.BusKey,
		0.0, // 连接失败时耗时为 0
		port,
		0,
		1,
		502, // 502 表示连接失败
		errMsg,
	)

	// 追加到请求日志
	if ctx != nil {
		log.Append(ctx, msg)
	}
}

// SetConfigPath 包级别便捷函数，设置配置文件路径
// 需要在第一次调用 GetConnect 前设置
//
// 示例:
//
//	redis.SetConfigPath("/custom/path/app.yaml")
func SetConfigPath(path string) {
	internalRedis.SetConfigPath(path)
}

// Close 包级别便捷函数，关闭所有 Redis 连接
func Close() error {
	return internalRedis.Close()
}
