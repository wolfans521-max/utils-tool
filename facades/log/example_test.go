package log

import (
	"net/http/httptest"

	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ctxfacade "git.intra.weibo.com/search_fe/wbutil-go/facades/context"

	internalLog "git.intra.weibo.com/search_fe/wbutil-go/internal/log"
	"github.com/gin-gonic/gin"
)

// Example 日志使用示例
func Example() {
	// 构造一个用于测试的 gin.Context，并确保 Request 不为 nil
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = httptest.NewRequest("GET", "http://example.com/test", nil)

	ctx := ctxfacade.New(ginCtx)

	// 增加日志（多次追加到同一个 tag/level 下）
	Append(ctx, "logmessage")
	Append(ctx, "logmessagev2")
	Append(ctx, "logmessagev3")
	// 写入日志文件（传入请求上下文，以便 Flush 时使用 ctx 中的 uid/seqid 等信息）
	_ = Flush(ctx)

	// Output:
}

// Test_Mertics 模拟真实使用场景：业务侧只调用 Mertics，不做额外初始化或文件检查。
// 是否落盘由实际 app.yaml/logger 配置和运行环境决定。
func Test_Mertics(t *testing.T) {
	module := "test_metrics_module"
	longMsg := strings.Repeat("x", 10050) // 超过 10000，应被截断

	Mertics(module, longMsg)
}

// BenchmarkAppendAndFlush 压力测试：大量 Append，最后 Flush 一次
func BenchmarkAppendAndFlush(b *testing.B) {
	// 构造一个用于测试的 gin.Context，并确保 Request 不为 nil
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = httptest.NewRequest("GET", "http://example.com/test", nil)

	ctx := ctxfacade.New(ginCtx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 使用与线上一致的 tag "req"，避免被文件 target 的 tags 过滤掉
		Append(ctx, "benchmark log message"+toString(i))
	}
	b.StopTimer()

	// 压力测试完成后统一 Flush，将缓冲日志落盘
	_ = Flush(ctx)
}

// TestAppendFlush_NoDuplicate 判断 Append 聚合 + Flush 写入后是否产生重复数据
func TestAppendFlush_NoDuplicate(t *testing.T) {
	// 使用临时目录作为日志输出目录，避免污染真实日志
	dir := t.TempDir()

	cfg := DefaultConfig()
	cfg.File.Path = dir
	cfg.Elk.Enable = false
	cfg.File.Tags = nil
	cfg.File.Levels = []string{string(internalLog.LevelInfo)}

	if err := Init(cfg); err != nil {
		t.Fatalf("Init logger failed: %v", err)
	}

	// 构造带请求信息的上下文，用于 Flush 时写入 uid/seqid
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = httptest.NewRequest("GET", "http://example.com/test", nil)
	ctx := ctxfacade.New(ginCtx)

	const n = 1000
	tag := "req"

	for i := 0; i < n; i++ {
		msg := fmt.Sprintf("|msg-%04d|", i)
		Append(ctx, msg)
	}

	if err := Flush(ctx); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// 查找当前 tag 对应的日志文件
	entries, err := os.ReadDir(cfg.File.Path)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	var logPath string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.Contains(name, "_"+tag+"_") && strings.HasSuffix(name, cfg.File.Extension) {
			logPath = filepath.Join(cfg.File.Path, name)
			break
		}
	}
	if logPath == "" {
		t.Fatalf("no log file found for tag %q", tag)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	content := string(data)

	// 每条消息应当只出现一次
	for i := 0; i < n; i++ {
		msg := fmt.Sprintf("|msg-%04d|", i)
		cnt := strings.Count(content, msg)
		if cnt != 1 {
			t.Fatalf("message %q appears %d times, expect 1", msg, cnt)
		}
	}
}
