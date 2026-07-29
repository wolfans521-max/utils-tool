package log

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/config"
	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"git.intra.weibo.com/search_fe/wbutil-go/internal/log/formatter"
	"git.intra.weibo.com/search_fe/wbutil-go/internal/log/owner"
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
	"github.com/bytedance/sonic"
	"github.com/goccy/go-yaml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志模块常量配置
const (
	// AsyncWriterBufferSize 异步写入缓冲区大小
	AsyncWriterBufferSize = 10000
	// AsyncWriterBatchSize 批量写入大小
	AsyncWriterBatchSize = 100
	// AsyncWriterFlushInterval 刷新间隔
	AsyncWriterFlushInterval = 100 * time.Millisecond
	// FileHandleCacheSize 文件句柄缓存大小
	FileHandleCacheSize = 100
	// FileHandleIdleTimeout 文件句柄空闲超时
	FileHandleIdleTimeout = 5 * time.Minute
	// FilePermission 日志文件权限 (0664允许组用户写入)
	FilePermission = 0o664
	// DirPermission 日志目录权限
	DirPermission = 0o775
)

// Level 日志级别
type Level string

const (
	LevelInfo  Level = "Info"
	LevelError Level = "Error"
	LevelDebug Level = "Debug"
)

func (lv Level) toZapLevel() zapcore.Level {
	switch lv {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// UIDFunc 提供 uid 获取逻辑（默认返回空字符串）。
type UIDFunc func() string

// SeqIDFunc 提供 seqId 获取逻辑（默认返回 "-"）。
type SeqIDFunc func() string
type Option func(*Logger)

func WithUIDFunc(fn UIDFunc) Option {
	return func(l *Logger) {
		if fn != nil {
			l.uidFn = fn
		}
	}
}

func WithSeqIDFunc(fn SeqIDFunc) Option {
	return func(l *Logger) {
		if fn != nil {
			l.seqIDFn = fn
		}
	}
}

// logEntry 日志条目
type logEntry struct {
	path string
	line string
}

// fileHandle 文件句柄缓存项
type fileHandle struct {
	file       *os.File
	path       string
	lastUsed   time.Time
	writeCount int64
	inode      uint64 // 文件 inode，用于检测文件是否被替换
}

// AsyncFileWriter 异步文件写入器
type AsyncFileWriter struct {
	// 写入队列
	queue     chan *logEntry
	batchSize int

	// 文件句柄缓存
	handles    map[string]*fileHandle
	handlesMu  sync.RWMutex
	lru        *list.List
	lruMap     map[string]*list.Element
	maxHandles int

	// 状态控制
	closed  atomic.Bool
	wg      sync.WaitGroup
	closeCh chan struct{}
	flushCh chan chan struct{}

	// 统计信息
	stats struct {
		totalWrites   atomic.Int64
		totalBytes    atomic.Int64
		droppedWrites atomic.Int64
		flushCount    atomic.Int64
		errorCount    atomic.Int64
	}
}

// AsyncWriterStats 异步写入器统计信息
type AsyncWriterStats struct {
	TotalWrites   int64
	TotalBytes    int64
	DroppedWrites int64
	FlushCount    int64
	ErrorCount    int64
	QueueSize     int
	HandleCount   int
}

// 全局异步写入器
var (
	globalAsyncWriter     *AsyncFileWriter
	globalAsyncWriterOnce sync.Once
)

// GetAsyncWriter 获取全局异步写入器
func GetAsyncWriter() *AsyncFileWriter {
	globalAsyncWriterOnce.Do(func() {
		globalAsyncWriter = NewAsyncFileWriter(AsyncWriterBufferSize, AsyncWriterBatchSize, FileHandleCacheSize)
	})
	return globalAsyncWriter
}

// NewAsyncFileWriter 创建异步文件写入器
func NewAsyncFileWriter(bufferSize, batchSize, maxHandles int) *AsyncFileWriter {
	if bufferSize <= 0 {
		bufferSize = AsyncWriterBufferSize
	}
	if batchSize <= 0 {
		batchSize = AsyncWriterBatchSize
	}
	if maxHandles <= 0 {
		maxHandles = FileHandleCacheSize
	}

	w := &AsyncFileWriter{
		queue:      make(chan *logEntry, bufferSize),
		batchSize:  batchSize,
		handles:    make(map[string]*fileHandle),
		lru:        list.New(),
		lruMap:     make(map[string]*list.Element),
		maxHandles: maxHandles,
		closeCh:    make(chan struct{}),
		flushCh:    make(chan chan struct{}, 10),
	}

	// 启动写入协程
	w.wg.Add(1)
	go w.writeLoop()

	// 启动清理协程
	w.wg.Add(1)
	go w.cleanupLoop()

	return w
}

// Write 异步写入日志
func (w *AsyncFileWriter) Write(path, line string) error {
	if w.closed.Load() {
		return errors.New("async writer is closed")
	}

	entry := &logEntry{
		path: path,
		line: line,
	}

	select {
	case w.queue <- entry:
		return nil
	default:
		// 队列满，丢弃日志并记录
		w.stats.droppedWrites.Add(1)
		return errors.New("log queue is full, entry dropped")
	}
}

// WriteSync 同步写入日志（用于关键日志）
func (w *AsyncFileWriter) WriteSync(path, line string) error {
	return w.writeToFile(path, line)
}

// Flush 刷新所有缓冲的日志
func (w *AsyncFileWriter) Flush() error {
	if w.closed.Load() {
		return errors.New("async writer is closed")
	}

	done := make(chan struct{})
	select {
	case w.flushCh <- done:
		<-done
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("flush timeout")
	}
}

// Close 关闭异步写入器
func (w *AsyncFileWriter) Close() error {
	if !w.closed.CompareAndSwap(false, true) {
		return nil
	}

	close(w.closeCh)
	w.wg.Wait()

	// 关闭所有文件句柄
	w.handlesMu.Lock()
	for _, h := range w.handles {
		if h.file != nil {
			h.file.Sync()
			h.file.Close()
		}
	}
	w.handles = nil
	w.handlesMu.Unlock()

	return nil
}

// Stats 获取统计信息
func (w *AsyncFileWriter) Stats() *AsyncWriterStats {
	w.handlesMu.RLock()
	handleCount := len(w.handles)
	w.handlesMu.RUnlock()

	return &AsyncWriterStats{
		TotalWrites:   w.stats.totalWrites.Load(),
		TotalBytes:    w.stats.totalBytes.Load(),
		DroppedWrites: w.stats.droppedWrites.Load(),
		FlushCount:    w.stats.flushCount.Load(),
		ErrorCount:    w.stats.errorCount.Load(),
		QueueSize:     len(w.queue),
		HandleCount:   handleCount,
	}
}

// writeLoop 写入循环
func (w *AsyncFileWriter) writeLoop() {
	defer w.wg.Done()

	batch := make([]*logEntry, 0, w.batchSize)
	ticker := time.NewTicker(AsyncWriterFlushInterval)
	defer ticker.Stop()

	for {
		select {
		case entry := <-w.queue:
			batch = append(batch, entry)
			// 批量写入
			if len(batch) >= w.batchSize {
				w.writeBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			// 定时刷新
			if len(batch) > 0 {
				w.writeBatch(batch)
				batch = batch[:0]
			}

		case done := <-w.flushCh:
			// 手动刷新
			if len(batch) > 0 {
				w.writeBatch(batch)
				batch = batch[:0]
			}
			// 刷新所有文件句柄
			w.syncAllHandles()
			w.stats.flushCount.Add(1)
			close(done)

		case <-w.closeCh:
			// 关闭前写入剩余日志
			// 先排空队列
		drainLoop:
			for {
				select {
				case entry := <-w.queue:
					batch = append(batch, entry)
				default:
					break drainLoop
				}
			}
			if len(batch) > 0 {
				w.writeBatch(batch)
			}
			return
		}
	}
}

// writeBatch 批量写入
func (w *AsyncFileWriter) writeBatch(batch []*logEntry) {
	// 按文件路径分组
	groups := make(map[string][]string)
	for _, entry := range batch {
		groups[entry.path] = append(groups[entry.path], entry.line)
	}

	// 批量写入每个文件
	for path, lines := range groups {
		content := strings.Join(lines, "")
		if err := w.writeToFile(path, content); err != nil {
			w.stats.errorCount.Add(1)
		} else {
			w.stats.totalWrites.Add(int64(len(lines)))
			w.stats.totalBytes.Add(int64(len(content)))
		}
	}
}

// writeToFile 写入文件
func (w *AsyncFileWriter) writeToFile(path, content string) error {
	// 最多重试2次
	var lastErr error
	for retry := 0; retry < 2; retry++ {
		handle, err := w.getOrCreateHandle(path)
		if err != nil {
			lastErr = err
			// 如果是目录不存在导致的错误，等待一下再重试
			if retry == 0 {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			return err
		}

		// 尝试获取文件锁，以兼容 PHP 的 flock
		// PHP 使用 LOCK_EX 写入日志，Go 需要等待 PHP 释放锁
		// 最多等待 100ms，避免长时间阻塞
		locked := false
		if err := FlockWithTimeout(handle.file, true, 100*time.Millisecond); err != nil {
			// 获取锁失败（超时），但不阻止写入
			// 因为可能 PHP 没有使用 flock，或者文件系统不支持锁
			// 继续尝试写入
		} else {
			locked = true
		}
		_, err = handle.file.WriteString(content)

		// 如果获取了锁，立即释放
		if locked {
			Funlock(handle.file)
		}

		if err != nil {
			// 文件写入失败，关闭句柄并重试
			w.closeHandle(path)
			lastErr = err
			if retry == 0 {
				// 第一次失败，尝试重新创建
				continue
			}
			return err
		}

		handle.lastUsed = time.Now()
		handle.writeCount++
		return nil
	}
	return lastErr
}

// getOrCreateHandle 获取或创建文件句柄
func (w *AsyncFileWriter) getOrCreateHandle(path string) (*fileHandle, error) {
	// 先尝试读锁获取
	w.handlesMu.RLock()
	if h, ok := w.handles[path]; ok {
		// 检查文件句柄是否仍然有效
		// 1. 检查文件句柄是否有效
		// 2. 检查 inode 是否变化（PHP 可能删除并重新创建文件）
		if h.file != nil {
			// 尝试获取文件状态，如果失败说明文件已被删除
			if _, err := h.file.Stat(); err == nil {
				// 检查 inode 是否变化
				currentInode := utils.GetPathInode(path)
				if currentInode != 0 && currentInode == h.inode {
					w.handlesMu.RUnlock()
					// 更新LRU
					w.updateLRU(path)
					return h, nil
				}
				// inode 变化，文件被替换，需要重新打开
				w.handlesMu.RUnlock()
				w.closeHandle(path)
			} else {
				// 文件已被删除，需要关闭并重新创建
				w.handlesMu.RUnlock()
				w.closeHandle(path)
			}
		} else {
			w.handlesMu.RUnlock()
			w.closeHandle(path)
		}
	} else {
		w.handlesMu.RUnlock()
	}

	// 需要创建，使用写锁
	w.handlesMu.Lock()
	defer w.handlesMu.Unlock()

	// 双重检查
	if h, ok := w.handles[path]; ok {
		// 再次检查文件有效性
		if h.file != nil {
			if _, err := h.file.Stat(); err == nil {
				return h, nil
			}
		}
		// 文件无效，删除缓存
		if h.file != nil {
			h.file.Close()
		}
		delete(w.handles, path)
		if elem, ok := w.lruMap[path]; ok {
			w.lru.Remove(elem)
			delete(w.lruMap, path)
		}
	}

	// 检查是否需要淘汰
	if len(w.handles) >= w.maxHandles {
		w.evictOldestLocked()
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, DirPermission); err != nil {
		return nil, fmt.Errorf("create log dir failed: %w", err)
	}
	// 尝试将目录所有者改为配置的用户/组
	owner.ChownToConfigOwner(dir)

	// 打开文件
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, FilePermission)
	if err != nil {
		return nil, fmt.Errorf("open log file failed: %w", err)
	}
	// 尝试修改文件权限为 0664，确保与 PHP 写入同一文件时权限正确
	// 如果文件已存在，OpenFile 不会修改权限，需要显式设置
	_ = os.Chmod(path, FilePermission)
	// 尝试将文件所有者改为配置的用户/组
	owner.ChownToConfigOwner(path)

	h := &fileHandle{
		file:     f,
		path:     path,
		lastUsed: time.Now(),
		inode:    utils.GetInode(f),
	}
	w.handles[path] = h

	// 添加到LRU
	elem := w.lru.PushFront(path)
	w.lruMap[path] = elem

	return h, nil
}

// updateLRU 更新LRU顺序
func (w *AsyncFileWriter) updateLRU(path string) {
	w.handlesMu.Lock()
	defer w.handlesMu.Unlock()

	if elem, ok := w.lruMap[path]; ok {
		w.lru.MoveToFront(elem)
	}
}

// evictOldestLocked 淘汰最旧的句柄（需要持有写锁）
func (w *AsyncFileWriter) evictOldestLocked() {
	if w.lru.Len() == 0 {
		return
	}

	// 从LRU尾部获取最旧的
	elem := w.lru.Back()
	if elem == nil {
		return
	}

	path := elem.Value.(string)
	w.lru.Remove(elem)
	delete(w.lruMap, path)

	if h, ok := w.handles[path]; ok {
		if h.file != nil {
			h.file.Sync()
			h.file.Close()
		}
		delete(w.handles, path)
	}
}

// closeHandle 关闭指定句柄
func (w *AsyncFileWriter) closeHandle(path string) {
	w.handlesMu.Lock()
	defer w.handlesMu.Unlock()

	if h, ok := w.handles[path]; ok {
		if h.file != nil {
			h.file.Close()
		}
		delete(w.handles, path)
	}

	if elem, ok := w.lruMap[path]; ok {
		w.lru.Remove(elem)
		delete(w.lruMap, path)
	}
}

// syncAllHandles 同步所有文件句柄
func (w *AsyncFileWriter) syncAllHandles() {
	w.handlesMu.RLock()
	defer w.handlesMu.RUnlock()

	for _, h := range w.handles {
		if h.file != nil {
			h.file.Sync()
		}
	}
}

// cleanupLoop 清理过期句柄
func (w *AsyncFileWriter) cleanupLoop() {
	defer w.wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.cleanupIdleHandles()
		case <-w.closeCh:
			return
		}
	}
}

// cleanupIdleHandles 清理空闲句柄
func (w *AsyncFileWriter) cleanupIdleHandles() {
	w.handlesMu.Lock()
	defer w.handlesMu.Unlock()

	now := time.Now()
	var toRemove []string

	for path, h := range w.handles {
		if now.Sub(h.lastUsed) > FileHandleIdleTimeout {
			toRemove = append(toRemove, path)
		}
	}

	for _, path := range toRemove {
		if h, ok := w.handles[path]; ok {
			if h.file != nil {
				h.file.Sync()
				h.file.Close()
			}
			delete(w.handles, path)
		}
		if elem, ok := w.lruMap[path]; ok {
			w.lru.Remove(elem)
			delete(w.lruMap, path)
		}
	}
}

// Logger 基于 zap 的日志封装。
type Logger struct {
	mu sync.Mutex

	z *zap.Logger

	uidFn   UIDFunc
	seqIDFn SeqIDFunc

	moduleRecords []moduleRecord
	asyncWriter   *AsyncFileWriter
}

// NewLogger 创建 logger。z 不能为空。
func NewLogger(z *zap.Logger, opts ...Option) *Logger {
	if z == nil {
		z = zap.NewNop()
	}
	l := &Logger{
		z:           z,
		uidFn:       func() string { return "" },
		seqIDFn:     func() string { return "-" },
		asyncWriter: GetAsyncWriter(),
	}
	for _, opt := range opts {
		if opt != nil {
			opt(l)
		}
	}
	return l
}

// Underlying 返回底层 zap logger。
func (l *Logger) Underlying() *zap.Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.z
}

// ReplaceUnderlying 替换底层 zap logger。
func (l *Logger) ReplaceUnderlying(z *zap.Logger) {
	if z == nil {
		z = zap.NewNop()
	}
	l.mu.Lock()
	l.z = z
	l.mu.Unlock()
}

func (l *Logger) Log(message any, tag string, level Level) {
	if tag == "" {
		tag = "default"
	}
	if level == "" {
		level = LevelInfo
	}

	msg := normalizeMessage(message)
	uid := safeCallUID(l.uidFn)
	seq := safeCallSeqID(l.seqIDFn)

	fields := []zap.Field{
		zap.String("uid", uid),
		zap.String("seqId", seq),
		zap.String("tag", tag),
		zap.String("levelName", string(level)),
	}

	l.mu.Lock()
	z := l.z
	l.mu.Unlock()

	if ce := z.Check(level.toZapLevel(), msg); ce != nil {
		ce.Write(fields...)
	}
}

// LogWithCtx 使用传入的请求上下文 ctx 写 uid/seqId 字段。
func (l *Logger) LogWithCtx(ctx *context.Context, message any, tag string, level Level) {
	if tag == "" {
		tag = "default"
	}
	if level == "" {
		level = LevelInfo
	}

	msg := normalizeMessage(message)
	uid := ""
	seq := "-"
	if ctx != nil {
		uid = ctx.GetUserId()
		seq = ctx.GetSeqId()
	}

	fields := []zap.Field{
		zap.String("uid", uid),
		zap.String("seqId", seq),
		zap.String("tag", tag),
		zap.String("levelName", string(level)),
	}

	l.mu.Lock()
	z := l.z
	l.mu.Unlock()

	if ce := z.Check(level.toZapLevel(), msg); ce != nil {
		ce.Write(fields...)
	}
}

// Flush 日志最终写入文件（包含 fsync）
func (l *Logger) Flush(ctx *context.Context) error {
	l.flushModulesLocked()

	// 刷新异步写入器 + zap Sync
	if l.asyncWriter != nil {
		l.asyncWriter.Flush()
	}

	l.mu.Lock()
	z := l.z
	l.mu.Unlock()

	if z == nil {
		return errors.New("log: underlying zap logger is nil")
	}
	return z.Sync()
}

// FlushModules 只将 Push 缓冲的 moduleRecords 写出到 asyncWriter，
// 不触发 fsync——磁盘同步由后台 ticker（100ms）负责。
func (l *Logger) FlushModules() {
	l.flushModulesLocked()
}

// flushModulesLocked 内部实现：将 moduleRecords 写出
func (l *Logger) flushModulesLocked() {
	l.mu.Lock()
	records := l.moduleRecords
	l.moduleRecords = nil
	l.mu.Unlock()

	if len(records) > 0 {
		for _, rec := range records {
			l.WriteModule(rec.Module, rec.Message, rec.Ctx)
		}
	}
}

func safeCallUID(fn UIDFunc) (uid string) {
	defer func() {
		if recover() != nil {
			uid = ""
		}
	}()
	if fn == nil {
		return ""
	}
	return fn()
}

func safeCallSeqID(fn SeqIDFunc) (seq string) {
	defer func() {
		if recover() != nil {
			seq = "-"
		}
	}()
	if fn == nil {
		return "-"
	}
	seq = fn()
	if seq == "" {
		return "-"
	}
	return seq
}

func normalizeMessage(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case error:
		return x.Error()
	case fmt.Stringer:
		return x.String()
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(x)
	}

	rv := reflect.ValueOf(v)
	if rv.IsValid() {
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.Bool:
			return fmt.Sprint(v)
		case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct, reflect.Ptr:
			b, err := sonic.Marshal(v)
			if err != nil {
				return fmt.Sprint(v)
			}
			return string(b)
		default:
			return fmt.Sprint(v)
		}
	}

	b, err := sonic.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}

// joinErrors 保持与旧实现一致
func joinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	parts := make([]string, 0, len(errs))
	for _, e := range errs {
		if e != nil {
			parts = append(parts, e.Error())
		}
	}
	if len(parts) == 0 {
		return nil
	}
	return errors.New("log: " + strings.Join(parts, "; "))
}

// moduleDef 描述 config/log.yaml 中单个模块配置。
type moduleDef struct {
	Name        string `yaml:"name"`
	FilePath    string `yaml:"file_path"`
	Formatter   string `yaml:"formatter"`
	NoTimestamp *bool  `yaml:"no_timestamp"`
}

// moduleConfig 表示从 config/log.yaml 解析出的整体配置。
type moduleConfig struct {
	RootPath string
	Modules  map[string]*moduleDef
}

type moduleRecord struct {
	Module  string
	Message any
	Ctx     *context.Context
}

var (
	moduleCfgOnce sync.Once
	moduleCfg     *moduleConfig
	moduleCfgErr  error

	configDirOnce sync.Once
	configDir     string
	configDirErr  error
)

// Push 按 config/log.yaml 模块配置缓冲一条文件日志
func (l *Logger) Push(module string, message any, ctx *context.Context) {
	if module == "" {
		return
	}

	l.mu.Lock()
	l.moduleRecords = append(l.moduleRecords, moduleRecord{
		Module:  module,
		Message: message,
		Ctx:     ctx,
	})
	l.mu.Unlock()
}

// WriteModule 按 config/log.yaml 中的模块配置写入文件日志（异步）
func (l *Logger) WriteModule(module string, message any, ctx *context.Context) {
	if module == "" {
		return
	}

	cfg, err := loadModuleConfig()
	if err != nil || cfg == nil {
		return
	}
	mod, ok := cfg.Modules[module]
	if !ok || mod == nil || mod.FilePath == "" {
		return
	}

	// 设置 json_log 标记到 ctx
	if ctx != nil {
		ctx.SetRequestParam("json_log", "true")
		if mod.NoTimestamp != nil {
			if *mod.NoTimestamp {
				ctx.SetRequestParam("no_timestamp", "true")
			} else {
				ctx.SetRequestParam("no_timestamp", "false")
			}
		}
	}

	var line string
	if s, ok := message.(string); ok && s != "" {
		// 确保字符串末尾有换行符
		if !strings.HasSuffix(s, "\n") {
			line = s + "\n"
		} else {
			line = s
		}
	} else {
		switch strings.ToLower(mod.Formatter) {
		case "action":
			line = formatter.Action(message, ctx)
		default:
			line = formatter.Common(message, ctx)
		}
	}
	if line == "" {
		return
	}

	// 使用异步写入
	path := buildLogPath(cfg.RootPath, mod.FilePath)
	if path != "" && l.asyncWriter != nil {
		_ = l.asyncWriter.Write(path, line)
	}
}

// WriteModuleSync 同步写入模块日志（用于关键日志）
func (l *Logger) WriteModuleSync(module string, message any, ctx *context.Context) {
	if module == "" {
		return
	}

	cfg, err := loadModuleConfig()
	if err != nil || cfg == nil {
		return
	}
	mod, ok := cfg.Modules[module]
	if !ok || mod == nil || mod.FilePath == "" {
		return
	}

	if ctx != nil {
		ctx.SetRequestParam("json_log", "true")
		if mod.NoTimestamp != nil {
			if *mod.NoTimestamp {
				ctx.SetRequestParam("no_timestamp", "true")
			} else {
				ctx.SetRequestParam("no_timestamp", "false")
			}
		}
	}

	var line string
	if s, ok := message.(string); ok && s != "" {
		// 确保字符串末尾有换行符
		if !strings.HasSuffix(s, "\n") {
			line = s + "\n"
		} else {
			line = s
		}
	} else {
		switch strings.ToLower(mod.Formatter) {
		case "action":
			line = formatter.Action(message, ctx)
		default:
			line = formatter.Common(message, ctx)
		}
	}
	if line == "" {
		return
	}

	path := buildLogPath(cfg.RootPath, mod.FilePath)
	if path != "" && l.asyncWriter != nil {
		l.asyncWriter.WriteSync(path, line)
	}
}

// buildLogPath 构建日志文件路径
func buildLogPath(root, pattern string) string {
	if root == "" || pattern == "" {
		return ""
	}
	now := time.Now()
	ymd := now.Format("2006-01-02")
	hour := now.Format("15")
	rel := strings.ReplaceAll(pattern, "{Ymd}", ymd)
	rel = strings.ReplaceAll(rel, "{hour}", hour)
	return filepath.Join(root, rel)
}

// GetAsyncWriterStats 获取异步写入器统计信息
func (l *Logger) GetAsyncWriterStats() *AsyncWriterStats {
	if l.asyncWriter != nil {
		return l.asyncWriter.Stats()
	}
	return nil
}

func loadModuleConfig() (*moduleConfig, error) {
	moduleCfgOnce.Do(func() {
		data := config.LogYaml
		if len(data) == 0 {
			moduleCfgErr = fmt.Errorf("log: embedded log.yaml is empty")
			return
		}

		// 使用 map[string]any 来解析，因为 log.yaml 包含不同类型的值
		// 如 user/group 是数组格式，而模块配置是 map 格式
		raw := make(map[string]any)
		if err := yaml.Unmarshal(data, &raw); err != nil {
			moduleCfgErr = err
			return
		}

		cfg := &moduleConfig{
			RootPath: "",
			Modules:  make(map[string]*moduleDef),
		}

		configDir, err := findConfigDirUpward()
		if err == nil {
			appPath := filepath.Join(configDir, "app.yaml")
			if appData, readErr := os.ReadFile(appPath); readErr == nil {
				appRaw := make(map[string]any)
				if yaml.Unmarshal(appData, &appRaw) == nil {
					if logger, ok := appRaw["logger"].(map[string]any); ok {
						if file, ok := logger["file"].(map[string]any); ok {
							if p, ok := file["path"].(string); ok && p != "" {
								if !filepath.IsAbs(p) {
									p = filepath.Join(configDir, p)
								}
								cfg.RootPath = p
							}
						}
					}
				}
			}
		}

		// 遍历解析后的配置，跳过非模块配置（如 user、group）
		for name, val := range raw {
			// 跳过 user 和 group 配置
			if name == "user" || name == "group" {
				continue
			}

			// 只处理 map 类型的值（模块配置）
			m, ok := val.(map[string]any)
			if !ok {
				continue
			}

			mod := &moduleDef{}
			if v, ok := m["name"].(string); ok {
				mod.Name = v
			}
			if v, ok := m["file_path"].(string); ok && v != "" {
				mod.FilePath = v
			}
			if v, ok := m["formatter"].(string); ok && v != "" {
				mod.Formatter = v
			} else {
				mod.Formatter = "common"
			}
			if v, ok := m["no_timestamp"].(bool); ok {
				mod.NoTimestamp = &v
			}
			if mod.FilePath != "" {
				cfg.Modules[name] = mod
			}
		}

		moduleCfg = cfg
	})
	return moduleCfg, moduleCfgErr
}

// findConfigDirUpward 在当前工作目录向上查找 config 目录路径
func findConfigDirUpward() (string, error) {
	configDirOnce.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			configDirErr = fmt.Errorf("log: getwd failed: %w", err)
			return
		}
		for {
			candidate := filepath.Join(dir, "config")
			if st, e := os.Stat(candidate); e == nil && st.IsDir() {
				configDir = candidate
				return
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
		configDirErr = fmt.Errorf("log: config directory not found")
	})
	return configDir, configDirErr
}

// 兼容旧接口：writeModuleLine 改为使用异步写入
func writeModuleLine(root, pattern, line string) {
	path := buildLogPath(root, pattern)
	if path == "" || line == "" {
		return
	}
	GetAsyncWriter().Write(path, line)
}
