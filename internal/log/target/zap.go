package target

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/internal/log"
	"git.intra.weibo.com/search_fe/wbutil-go/internal/log/owner"
	"git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// BuildZapLogger 使用 File/ELK 配置构建 zap.Logger。
func BuildZapLogger(fileCfg FileConfig, elkCfg ElkConfig) (*zap.Logger, error) {
	var cores []zapcore.Core
	if fileCfg.Enable {
		c, err := newFileCore(fileCfg)
		if err != nil {
			return nil, err
		}
		cores = append(cores, c)
	}
	if elkCfg.Enable {
		c, err := newElkCore(elkCfg)
		if err != nil {
			return nil, err
		}
		cores = append(cores, c)
	}
	if len(cores) == 0 {
		return nil, errors.New("log: no zap cores enabled")
	}

	tee := zapcore.NewTee(cores...)
	// 不强制 AddCaller；调用方如需 caller 可在外层重新构建。
	return zap.New(tee), nil
}

// ========================= File Core =========================

type fileCore struct {
	cfg FileConfig

	tagsSet   map[string]struct{}
	levelsSet map[string]struct{}

	mu    sync.Mutex
	files map[string]*tagFile // tag -> file

	fields []zapcore.Field
}

type tagFile struct {
	fileName string
	f        *os.File
	inode    uint64 // 文件 inode，用于检测文件是否被替换
}

func newFileCore(cfg FileConfig) (*fileCore, error) {
	if cfg.Path == "" {
		return nil, errors.New("log: file.path is empty")
	}
	if cfg.Delimiter == "" {
		cfg.Delimiter = "_"
	}
	if cfg.Postfix == "" {
		cfg.Postfix = "20060102"
	}
	if cfg.Extension == "" {
		cfg.Extension = ".log"
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "app"
	}
	if err := os.MkdirAll(cfg.Path, 0o775); err != nil {
		return nil, err
	}
	// 尝试将目录所有者改为配置的用户/组
	owner.ChownToConfigOwner(cfg.Path)

	c := &fileCore{
		cfg:       cfg,
		tagsSet:   sliceToSet(cfg.Tags),
		levelsSet: sliceToSet(cfg.Levels),
		files:     map[string]*tagFile{},
	}
	return c, nil
}

func (c *fileCore) Enabled(lvl zapcore.Level) bool {
	// 最小级别判断交给 allow(levelStr)
	_ = lvl
	return true
}

func (c *fileCore) With(fs []zapcore.Field) zapcore.Core {
	nc := *c
	nc.fields = append(append([]zapcore.Field(nil), c.fields...), fs...)
	return &nc
}

func (c *fileCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	// 不在 Check 阶段做 tag/field 过滤：因为 tag 可能只在 fields 中（例如上层通过 zap.String("tag", ...) 传入）。
	// 统一在 Write 阶段解析 tag/level 后过滤，避免误判导致不落盘。
	return ce.AddCore(ent, c)
}

func (c *fileCore) Write(ent zapcore.Entry, fs []zapcore.Field) error {
	// 以 fields 中的 tag 为准；loggerName 可能是 zap 的 base name 或层级 name（会影响分文件）。
	tag := fieldString(append(c.fields, fs...), "tag")
	if tag == "" {
		tag = ent.LoggerName
	}
	if tag == "" {
		tag = "default"
	}

	// 过滤与输出 level：优先使用业务透传的 levelName（例如 Page），否则退化为 zap level 映射。
	lvlStr := fieldString(append(c.fields, fs...), "levelName")
	if lvlStr == "" {
		lvlStr = levelString(ent.Level)
	}
	if !c.allow(tag, lvlStr) {
		return nil
	}

	seq := fieldString(append(c.fields, fs...), "seqId")
	if seq == "" {
		seq = "-"
	}
	msg := ent.Message

	line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", ent.Time.Format("2006-01-02 15:04:05"), seq, lvlStr, tag, msg)

	fileName := c.cfg.Prefix + c.cfg.Delimiter + tag + c.cfg.Delimiter + ent.Time.Format(c.cfg.Postfix) + c.cfg.Extension
	fp := filepath.Join(c.cfg.Path, fileName)

	c.mu.Lock()
	defer c.mu.Unlock()

	tf := c.files[tag]
	needReopen := tf == nil || tf.f == nil || tf.fileName != fp

	// 检查文件句柄是否仍然有效（文件是否被删除或替换）
	if !needReopen && tf.f != nil {
		if _, err := tf.f.Stat(); err != nil {
			// 文件已被删除，需要重新打开
			needReopen = true
		} else {
			// 检查 inode 是否变化（PHP 可能删除并重新创建文件）
				currentInode := utils.GetPathInode(fp)
			if currentInode != 0 && currentInode != tf.inode {
				fmt.Fprintf(os.Stderr, "[zap] File inode changed for %s: old=%d, new=%d, reopening\n", fp, tf.inode, currentInode)
				needReopen = true
			}
		}
	}

	if needReopen {
		if tf != nil && tf.f != nil {
			_ = tf.f.Close()
		}
		// 确保目录存在
		dir := filepath.Dir(fp)
		if err := os.MkdirAll(dir, 0o775); err != nil {
			return err
		}
		// 尝试将目录所有者改为配置的用户/组
		owner.ChownToConfigOwner(dir)

		// 使用 O_APPEND 模式打开文件，确保与 PHP 兼容
		// O_APPEND 是原子操作，确保多进程同时写入时不会互相覆盖
		f, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664)
		if err != nil {
			return err
		}
		// 尝试修改文件权限为 0664，确保与 PHP 写入同一文件时权限正确
		// 如果文件已存在，OpenFile 不会修改权限，需要显式设置
		_ = os.Chmod(fp, 0o664)
		// 尝试将文件所有者改为配置的用户/组
		owner.ChownToConfigOwner(fp)

		tf = &tagFile{fileName: fp, f: f, inode: utils.GetInode(f)}
		c.files[tag] = tf
	}

	// 尝试获取文件锁，以兼容 PHP 的 flock
	// PHP 使用 LOCK_EX 写入日志，Go 需要等待 PHP 释放锁
	locked := false
	if err := log.FlockWithTimeout(tf.f, true, 100*time.Millisecond); err != nil {
		// 获取锁失败（超时），但不阻止写入
		// 因为可能 PHP 没有使用 flock，或者文件系统不支持锁
		// 继续尝试写入
	} else {
		locked = true
	}

	// 使用 O_APPEND 模式打开的文件，写入操作是原子的（对于小于 PIPE_BUF 的数据）
	// 但这里仍然可能失败（如文件被 PHP 删除），所以需要错误处理
	_, err := tf.f.WriteString(line)

	// 如果获取了锁，立即释放
	if locked {
		log.Funlock(tf.f)
	}

	if err != nil {
		// 写入失败，可能是文件被删除或权限问题
		// 关闭并删除缓存，下次写入时会重新打开
		if tf.f != nil {
			_ = tf.f.Close()
		}
		delete(c.files, tag)
		return fmt.Errorf("write to log file failed: %w", err)
	}
	return nil
}

func (c *fileCore) Sync() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var first error
	for _, tf := range c.files {
		if tf != nil && tf.f != nil {
			if err := tf.f.Sync(); err != nil && first == nil {
				first = err
			}
		}
	}
	return first
}

func (c *fileCore) allow(tag string, level string) bool {
	if len(c.tagsSet) > 0 {
		if _, ok := c.tagsSet[tag]; !ok {
			return false
		}
	}
	if len(c.levelsSet) > 0 {
		if _, ok := c.levelsSet[level]; !ok {
			return false
		}
	}
	return true
}

// ========================= ELK Core =========================

type elkCore struct {
	cfg ElkConfig

	tagsSet   map[string]struct{}
	levelsSet map[string]struct{}
	whiteSet  map[int64]struct{}

	fields []zapcore.Field
}

func newElkCore(cfg ElkConfig) (*elkCore, error) {
	if len(cfg.Hosts) == 0 {
		return nil, errors.New("log: elk.hosts is empty")
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "app"
	}

	c := &elkCore{
		cfg:       cfg,
		tagsSet:   sliceToSet(cfg.Tags),
		levelsSet: sliceToSet(cfg.Levels),
	}
	if len(cfg.WhiteList) > 0 {
		c.whiteSet = make(map[int64]struct{}, len(cfg.WhiteList))
		for _, u := range cfg.WhiteList {
			c.whiteSet[u] = struct{}{}
		}
	}
	return c, nil
}

func (c *elkCore) Enabled(lvl zapcore.Level) bool {
	_ = lvl
	return true
}

func (c *elkCore) With(fs []zapcore.Field) zapcore.Core {
	nc := *c
	nc.fields = append(append([]zapcore.Field(nil), c.fields...), fs...)
	return &nc
}

func (c *elkCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	// 同 fileCore：不在 Check 阶段过滤，避免 tag/uid 仅在 fields 中时被误过滤。
	return ce.AddCore(ent, c)
}

func (c *elkCore) Write(ent zapcore.Entry, fs []zapcore.Field) error {
	if !c.allow(ent, fs) {
		return nil
	}

	// 同 fileCore：优先 fields.tag。
	tag := fieldString(append(c.fields, fs...), "tag")
	if tag == "" {
		tag = ent.LoggerName
	}
	if tag == "" {
		tag = "default"
	}

	uid := fieldString(append(c.fields, fs...), "uid")
	seq := fieldString(append(c.fields, fs...), "seqId")
	if seq == "" {
		seq = "-"
	}

	source := c.cfg.Prefix + "-" + tag
	msg := ent.Message
	if len(msg) > 60000 {
		msg = msg[:60000] + "..."
	}

	payload := fmt.Sprintf("[text][%s][%s][%s][%s][%s]%s", ent.Time.Format("2006-01-02 15:04:05"), seq, source, tag, uid, msg)

	host := c.pickHost()
	addr := fmt.Sprintf("%s:%d", host.Host, host.Port)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	if host.TimeoutMS > 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(host.TimeoutMS) * time.Millisecond))
	}
	_, err = conn.Write([]byte(payload))
	return err
}

func (c *elkCore) Sync() error { return nil }

func (c *elkCore) pickHost() ElkHost {
	if len(c.cfg.Hosts) == 1 {
		return c.cfg.Hosts[0]
	}
	idx := rand.IntN(len(c.cfg.Hosts))
	return c.cfg.Hosts[idx]
}

func (c *elkCore) allow(ent zapcore.Entry, fs []zapcore.Field) bool {
	// allow 阶段尽量用真实 tag：优先 fields.tag，其次 loggerName。
	tag := fieldString(append(c.fields, fs...), "tag")
	if tag == "" {
		tag = ent.LoggerName
	}
	if tag == "" {
		tag = "default"
	}
	lvlStr := fieldString(append(c.fields, fs...), "levelName")
	if lvlStr == "" {
		lvlStr = levelString(ent.Level)
	}

	if len(c.tagsSet) > 0 {
		if _, ok := c.tagsSet[tag]; !ok {
			return false
		}
	}
	if len(c.levelsSet) > 0 {
		if _, ok := c.levelsSet[lvlStr]; !ok {
			return false
		}
	}

	// 白名单过滤：为空则关闭
	if len(c.whiteSet) > 0 {
		uid := fieldString(append(c.fields, fs...), "uid")
		n := parseInt64(uid)
		if n == 0 {
			return false
		}
		if _, ok := c.whiteSet[n]; !ok {
			return false
		}
	}

	return true
}

// ========================= helpers =========================

func levelString(lvl zapcore.Level) string {
	switch lvl {
	case zapcore.DebugLevel:
		return "Debug"
	case zapcore.InfoLevel:
		return "Info"
	case zapcore.WarnLevel:
		return "Warning"
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return "Error"
	default:
		return "Info"
	}
}

func fieldString(fields []zapcore.Field, key string) string {
	// zapcore.Field 是 union 类型；这里用反射读取其 Key/Type/String 等并不稳。
	// 简化：只支持 zap.String(...) 产生的字段（Type=StringType）。
	for i := len(fields) - 1; i >= 0; i-- {
		f := fields[i]
		if f.Key != key {
			continue
		}
		// 只取 String 字段
		if f.Type == zapcore.StringType {
			return f.String
		}
	}
	return ""
}

func parseInt64(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var n int64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int64(c-'0')
	}
	return n
}
