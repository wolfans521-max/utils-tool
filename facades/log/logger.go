package log

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	internalConfig "git.intra.weibo.com/search_fe/wbutil-go/internal/config"
	internalLog "git.intra.weibo.com/search_fe/wbutil-go/internal/log"
	"git.intra.weibo.com/search_fe/wbutil-go/internal/log/formatter"
	"git.intra.weibo.com/search_fe/wbutil-go/internal/log/target"
)

type Config struct {
	File target.FileConfig `json:"file" yaml:"file"`
	Elk  target.ElkConfig  `json:"elk" yaml:"elk"`
}

func DefaultConfig() Config {
	return Config{
		File: target.FileConfig{
			Enable:    true,
			Path:      "/data1/apache2/phplogs",
			Prefix:    "search",
			Delimiter: "_",
			Postfix:   "2006-01-02-15",
			Extension: ".log",
			Tags:      nil,
			// 仅保留 Error / Info / Debug 三种级别
			Levels: []string{
				string(internalLog.LevelError),
				string(internalLog.LevelInfo),
				string(internalLog.LevelDebug),
			},
		},
		Elk: target.ElkConfig{
			Enable:    false,
			Prefix:    "search",
			Hosts:     nil,
			Tags:      nil,
			Levels:    []string{string(internalLog.LevelError)},
			WhiteList: nil,
		},
	}
}

var (
	globalMu sync.RWMutex
	global   *internalLog.Logger

	globalCfgMu     sync.RWMutex
	globalCfg       Config
	globalCfgInited bool
)

func Init(cfg Config, opts ...internalLog.Option) error {
	globalMu.Lock()
	defer globalMu.Unlock()

	z, err := target.BuildZapLogger(cfg.File, cfg.Elk)
	if err != nil {
		return err
	}

	global = internalLog.NewLogger(z, opts...)

	// 缓存全局配置，供 Mertics 等场景复用
	globalCfgMu.Lock()
	globalCfg = cfg
	globalCfgInited = true
	globalCfgMu.Unlock()

	return nil
}

// Get 返回全局 Logger
// 使用读写锁优化：大多数情况下只需要读锁，只有初始化时才需要写锁
func Get() (*internalLog.Logger, error) {
	// 先尝试读锁，快速路径
	globalMu.RLock()
	if global != nil {
		l := global
		globalMu.RUnlock()
		return l, nil
	}
	globalMu.RUnlock()

	// 需要初始化，使用写锁
	globalMu.Lock()
	defer globalMu.Unlock()

	// 双重检查，避免重复初始化
	if global != nil {
		return global, nil
	}

	if err := initFromComponentConfigLocked(); err != nil {
		return nil, err
	}
	if global == nil {
		return nil, fmt.Errorf("log: not initialized")
	}
	return global, nil
}

func initFromComponentConfigLocked() error {
	v, err := internalConfig.ComponentViper()
	if err != nil {
		return fmt.Errorf("log: auto init failed: %w", err)
	}

	cfg := DefaultConfig()
	if err := internalConfig.LoadSectionFromViper(v, "logger", &cfg); err != nil {
		return fmt.Errorf("log: auto init failed: %w", err)
	}
	cfg.File.Path = internalConfig.AbsPathFromConfigFile(v, cfg.File.Path)

	z, err := target.BuildZapLogger(cfg.File, cfg.Elk)
	if err != nil {
		return fmt.Errorf("log: auto init failed: %w", err)
	}

	global = internalLog.NewLogger(z)

	// 缓存全局配置，供 Mertics 等场景复用
	globalCfgMu.Lock()
	globalCfg = cfg
	globalCfgInited = true
	globalCfgMu.Unlock()

	return nil
}

// Log 向全局 logger 记一条日志。
func Log(message any, tag string, level internalLog.Level) {
	l, err := Get()
	if err != nil {
		return
	}
	l.Log(message, tag, level)
}

// WithCtx 向全局 logger 记一条日志，使用传入的 context 获取 uid/seqId。
func WithCtx(ctx *context.Context, message any, tag string, level internalLog.Level) {
	l, err := Get()
	if err != nil {
		return
	}
	l.LogWithCtx(ctx, message, tag, level)
}

// Append 追加日志（请求级缓冲）：
// 默认语义为请求 Info 日志：tag=req，level=Info。
func Append(ctx *context.Context, message string) {
	appendWith(ctx, message, "req", internalLog.LevelInfo)
}

// Error 记录错误追加日志：tag=err，level=Error。
func Error(ctx *context.Context, message string) {
	appendWith(ctx, message, "err", internalLog.LevelError)
}

// Debug 记录调试追加日志：tag=debug，level=Debug。
func Debug(ctx *context.Context, message string) {
	appendWith(ctx, message, "debug", internalLog.LevelDebug)
}

// appendWith 用于少量需要自定义 tag/level 的场景（仅在本包内部使用）。
// 注意：此函数是并发安全的，使用 ctx.AppendMu 保护 AppendBuf 的访问。
// 使用 []string 收集消息，Flush 时 Join，避免 O(n²) 字符串拼接开销。
func appendWith(ctx *context.Context, message, tag string, level internalLog.Level) {
	if ctx == nil || message == "" {
		return
	}
	if tag == "" {
		// 默认视为请求日志
		tag = "req"
	}
	if level == "" {
		level = internalLog.LevelInfo
	}

	// 加锁保护 AppendBuf 的并发访问
	ctx.AppendMu.Lock()
	defer ctx.AppendMu.Unlock()

	if ctx.AppendBuf == nil {
		ctx.AppendBuf = make(map[string]map[string][]string)
	}
	byLevel := ctx.AppendBuf[tag]
	if byLevel == nil {
		byLevel = make(map[string][]string)
		ctx.AppendBuf[tag] = byLevel
	}

	lvKey := string(level)
	byLevel[lvKey] = append(byLevel[lvKey], message)
}

// Write 通过 config/log.yaml 的模块配置写入文件日志
func Write(module string, message any, ctx *context.Context) {
	l, err := Get()
	if err != nil {
		return
	}
	l.WriteModule(module, message, ctx)
}

// Action 通过行为码缓冲到 action_code 模块（push，Flush 时落盘）。
// actCode 行为码；oid 对象 id；ext 为 ext 字段内容（可以是 string 或 map[string]any）；ctx 为请求上下文。
func Action(actCode int, oid string, ext any, ctx *context.Context) {
	l, err := Get()
	if err != nil {
		return
	}

	// 将 actCode、oid、ext 临时设置到 context 中，供 formatter.Action 使用
	ctx.SetRequestParam("act_code", strconv.Itoa(actCode))
	if oid != "" {
		ctx.SetRequestParam("oid", oid)
	}
	if ext != nil {
		ctx.Set("ext", ext)
	}

	// 在调用时就格式化好日志内容，避免 Flush 时参数被覆盖
	line := formatter.Action(nil, ctx)

	// 将格式化后的日志内容传给 Push
	l.Push("action_code", line, ctx)
}

// Flush 日志写入：
// 1. 将当前请求 ctx 的 AppendBuf 聚合写出；
// 2. 调用底层 logger.Flush 执行 zap.Sync 以及模块日志 flush。
// 注意：此函数是并发安全的，使用 ctx.AppendMu 保护 AppendBuf 的访问。
func Flush(ctx *context.Context) error {
	l, err := Get()
	if err != nil {
		return err
	}

	if ctx != nil {
		// 加锁保护 AppendBuf 的并发访问
		ctx.AppendMu.Lock()
		appendBuf := ctx.AppendBuf
		ctx.AppendBuf = nil
		ctx.AppendMu.Unlock()

		// 在锁外处理日志写入，避免长时间持有锁
		flushAppendBuf(l, ctx, appendBuf)
	}

	return l.Flush(ctx)
}

// FlushBuffers 将当前请求 ctx 的 AppendBuf 聚合写出到 zap/asyncWriter，
// 但不触发 fsync——磁盘同步由后台 ticker（100ms）负责。
// 用于替代请求路径上的 Flush，避免每次请求都触发 fsync。
// 注意：此函数是并发安全的，使用 ctx.AppendMu 保护 AppendBuf 的访问。
func FlushBuffers(ctx *context.Context) error {
	l, err := Get()
	if err != nil {
		return err
	}

	if ctx != nil {
		// 加锁保护 AppendBuf 的并发访问
		ctx.AppendMu.Lock()
		appendBuf := ctx.AppendBuf
		ctx.AppendBuf = nil
		ctx.AppendMu.Unlock()

		// 在锁外处理日志写入，避免长时间持有锁
		flushAppendBuf(l, ctx, appendBuf)
	}

	// 只刷新 moduleRecords（Push 缓冲的日志），不做 fsync
	if l != nil {
		l.FlushModules()
	}

	return nil
}

// flushAppendBuf 将 AppendBuf 中的 []string 消息 Join 后写入日志。
// 将 O(n²) 的字符串拼接开销降为 O(n)。
func flushAppendBuf(l *internalLog.Logger, ctx *context.Context, appendBuf map[string]map[string][]string) {
	if appendBuf == nil {
		return
	}
	for tag, levels := range appendBuf {
		for lvStr, msgs := range levels {
			if len(msgs) == 0 {
				continue
			}
			// Join 所有消息为一条日志
			var msg string
			if len(msgs) == 1 {
				msg = msgs[0]
			} else {
				msg = strings.Join(msgs, " ")
			}
			l.LogWithCtx(ctx, msg, tag, internalLog.Level(lvStr))
		}
	}
}

func ensureGlobalCfgReady() (Config, bool) {
	globalCfgMu.RLock()
	inited := globalCfgInited
	cfg := globalCfg
	globalCfgMu.RUnlock()
	if inited && cfg.File.Path != "" {
		return cfg, true
	}

	// 尝试触发懒加载初始化（内部会调用 initFromComponentConfigLocked）
	if _, err := Get(); err != nil {
		return Config{}, false
	}

	globalCfgMu.RLock()
	cfg = globalCfg
	inited = globalCfgInited && cfg.File.Path != ""
	globalCfgMu.RUnlock()
	return cfg, inited
}

// Mertics 监控类日志（异步写入）：
//   - 行格式：YYYY-MM-DD hh:mm:ss\tmodule\tmessage；
//   - 文件名：prefix + delimiter + "metrics" + delimiter + time(postfix) + extension，
//     其中各字段复用 app.yaml.logger.file 的配置；
//   - 使用异步写入器，避免高并发场景下的性能瓶颈。
func Mertics(module, message string) {
	if module == "" || message == "" {
		return
	}

	// 限制字段长度，避免单条日志过大
	const moduleMaxLen = 30
	const messageMaxLen = 10000
	if len(module) > moduleMaxLen {
		module = module[:moduleMaxLen-3] + "..."
	}
	if len(message) > messageMaxLen {
		message = message[:messageMaxLen-3] + "..."
	}

	cfg, ok := ensureGlobalCfgReady()
	if !ok || !cfg.File.Enable || cfg.File.Path == "" {
		return
	}

	now := time.Now()
	timeStr := now.Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s\t%s\t%s\n", timeStr, module, message)

	delim := cfg.File.Delimiter
	if delim == "" {
		delim = "_"
	}
	postfix := cfg.File.Postfix
	if postfix == "" {
		postfix = "2006-01-02-15"
	}
	ext := cfg.File.Extension
	if ext == "" {
		ext = ".log"
	}

	tag := "metrics"
	fileTimePart := now.Format(postfix)
	fileName := fmt.Sprintf("%s%s%s%s%s%s", cfg.File.Prefix, delim, tag, delim, fileTimePart, ext)
	fullPath := filepath.Join(cfg.File.Path, fileName)

	// 使用异步写入器，避免高并发场景下的性能瓶颈
	asyncWriter := internalLog.GetAsyncWriter()
	if asyncWriter != nil {
		_ = asyncWriter.Write(fullPath, line)
	}
}

// MerticsSync 同步写入监控日志（用于关键监控信息，确保立即写入）
func MerticsSync(module, message string) {
	if module == "" || message == "" {
		return
	}

	// 限制字段长度，避免单条日志过大
	const moduleMaxLen = 30
	const messageMaxLen = 10000
	if len(module) > moduleMaxLen {
		module = module[:moduleMaxLen-3] + "..."
	}
	if len(message) > messageMaxLen {
		message = message[:messageMaxLen-3] + "..."
	}

	cfg, ok := ensureGlobalCfgReady()
	if !ok || !cfg.File.Enable || cfg.File.Path == "" {
		return
	}

	now := time.Now()
	timeStr := now.Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s\t%s\t%s\n", timeStr, module, message)

	delim := cfg.File.Delimiter
	if delim == "" {
		delim = "_"
	}
	postfix := cfg.File.Postfix
	if postfix == "" {
		postfix = "2006-01-02-15"
	}
	ext := cfg.File.Extension
	if ext == "" {
		ext = ".log"
	}

	tag := "metrics"
	fileTimePart := now.Format(postfix)
	fileName := fmt.Sprintf("%s%s%s%s%s%s", cfg.File.Prefix, delim, tag, delim, fileTimePart, ext)
	fullPath := filepath.Join(cfg.File.Path, fileName)

	// 使用同步写入，确保关键日志立即落盘
	asyncWriter := internalLog.GetAsyncWriter()
	if asyncWriter != nil {
		_ = asyncWriter.WriteSync(fullPath, line)
	}
}

// WriteSync 同步写入模块日志（用于关键日志，确保立即写入）
func WriteSync(module string, message any, ctx *context.Context) {
	l, err := Get()
	if err != nil {
		return
	}
	l.WriteModuleSync(module, message, ctx)
}

// Stats 日志统计信息
type Stats struct {
	TotalWrites   int64  `json:"total_writes"`
	TotalBytes    int64  `json:"total_bytes"`
	DroppedWrites int64  `json:"dropped_writes"`
	FlushCount    int64  `json:"flush_count"`
	ErrorCount    int64  `json:"error_count"`
	QueueSize     int    `json:"queue_size"`
	HandleCount   int    `json:"handle_count"`
	Status        string `json:"status"`
}

// GetStats 获取日志系统统计信息
func GetStats() *Stats {
	asyncWriter := internalLog.GetAsyncWriter()
	if asyncWriter == nil {
		return &Stats{Status: "not_initialized"}
	}

	internalStats := asyncWriter.Stats()
	if internalStats == nil {
		return &Stats{Status: "stats_unavailable"}
	}

	return &Stats{
		TotalWrites:   internalStats.TotalWrites,
		TotalBytes:    internalStats.TotalBytes,
		DroppedWrites: internalStats.DroppedWrites,
		FlushCount:    internalStats.FlushCount,
		ErrorCount:    internalStats.ErrorCount,
		QueueSize:     internalStats.QueueSize,
		HandleCount:   internalStats.HandleCount,
		Status:        "running",
	}
}

// Close 优雅关闭日志系统
// 应在应用退出前调用，确保所有缓冲的日志都被写入
func Close() error {
	asyncWriter := internalLog.GetAsyncWriter()
	if asyncWriter != nil {
		return asyncWriter.Close()
	}
	return nil
}

// FlushAll 刷新所有缓冲的日志
// 用于在关键时刻确保日志已写入磁盘
func FlushAll() error {
	asyncWriter := internalLog.GetAsyncWriter()
	if asyncWriter != nil {
		return asyncWriter.Flush()
	}
	return nil
}

// FormatStats 格式化统计信息为字符串（用于日志输出）
func FormatStats() string {
	stats := GetStats()
	if stats == nil {
		return "log_stats: unavailable"
	}

	var sb strings.Builder
	sb.WriteString("log_stats: ")
	sb.WriteString(fmt.Sprintf("status=%s ", stats.Status))
	sb.WriteString(fmt.Sprintf("total_writes=%d ", stats.TotalWrites))
	sb.WriteString(fmt.Sprintf("total_bytes=%d ", stats.TotalBytes))
	sb.WriteString(fmt.Sprintf("dropped_writes=%d ", stats.DroppedWrites))
	sb.WriteString(fmt.Sprintf("flush_count=%d ", stats.FlushCount))
	sb.WriteString(fmt.Sprintf("error_count=%d ", stats.ErrorCount))
	sb.WriteString(fmt.Sprintf("queue_size=%d ", stats.QueueSize))
	sb.WriteString(fmt.Sprintf("handle_count=%d", stats.HandleCount))
	return sb.String()
}
