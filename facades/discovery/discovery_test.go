package discovery

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// Test_ConfigGet：验证 facades/discovery -> internal/discovery 的链路可用（真实 vintage 数据）。
func Test_ConfigGet(t *testing.T) {
	time.Sleep(10 * time.Second)
	// serialization=true：返回值应为已 json_decode 的结构（map[string]any）
	def := map[string]any{}
	v, err := ConfigGet("hot_feed_black_button_gray", true, def)
	if err != nil {
		t.Fatalf("ConfigGet error: %v", err)
	}

	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("ConfigGet returned non-map value: %#v", v)
	}

	fmt.Println(m)
	t.Logf("ConfigGet parsed: %#v", m)

	v, err = ConfigGet("finder_fore_grey", false, 0)
	if err != nil {
		t.Fatalf("ConfigGet error: %v", err)
	}

	t.Logf("ConfigGet finder_fore_grey: %#v", v)

	v, err = ConfigGet("should_throttle_switch", false, 0)

	if err != nil {
		t.Fatalf("ConfigGet error: %v", err)
	}

	t.Logf("ConfigGet should_throttle_switch: %#v", v)
}

// Test_NamingGet：验证 facades/discovery -> internal/discovery 的链路可用（真实 vintage 数据）。
func Test_NamingGet(t *testing.T) {
	service := "top-search-plus"
	cluster := "top.search.plus"

	nodes, err := NamingGet(service, cluster)
	if err != nil {
		t.Fatalf("NamingGet error: %v", err)
	}

	for i, n := range nodes {
		t.Logf("node[%d]: host=%s port=%d domain=%s weight=%d ext=%v", i, n.Host, n.Port, n.Domain, n.Weight, n.Ext)
	}
}

// Benchmark_ConfigGet：对 ConfigGet 做基准测试，关注平均耗时与分配。
func Benchmark_ConfigGet(b *testing.B) {
	b.ReportAllocs()

	// 使用真实的 key，defaultValue 仅在失败时回退。
	def := map[string]any{}

	for i := 0; i < b.N; i++ {
		if _, err := ConfigGet("hot_feed_black_button_gray", true, def); err != nil {
			b.Fatalf("ConfigGet error: %v", err)
		}
	}
}

// Benchmark_NamingGet：对 NamingGet 做基准测试，关注平均耗时与分配。
func Benchmark_NamingGet(b *testing.B) {
	b.ReportAllocs()

	service := "top-search-plus"
	cluster := "top.search.plus"

	for i := 0; i < b.N; i++ {
		if _, err := NamingGet(service, cluster); err != nil {
			b.Fatalf("NamingGet error: %v", err)
		}
	}
}

// Test_DumpCache：验证 internal/discovery.DumpCache 可用，输出当前缓存快照。
func Test_DumpCache(t *testing.T) {
	m := DumpCache()
	fmt.Println(m)
	t.Logf("DumpCache snapshot: %#v", m)
}

func Test_ConfigGet_CacheHit(t *testing.T) {
	key := "hot_feed_black_button_gray" // 用你已经确认存在的 key

	def := map[string]any{}
	// 第一次调用：预期走远程（vintage），并写入缓存
	v1, err := ConfigGet(key, true, def)
	if err != nil {
		t.Fatalf("first ConfigGet error: %v", err)
	}

	v3, err := ConfigGet("finder_fore_grey", false, 0)
	if err != nil {
		t.Fatalf("ConfigGet finder_fore_grey error: %v", err)
	}

	fmt.Println("v3 =", v3)

	Test_DumpCache(t)

	// 第二次调用：预期命中缓存
	v2, err := ConfigGet(key, true, def)
	if err != nil {
		t.Fatalf("second ConfigGet error: %v", err)
	}

	if !reflect.DeepEqual(v1, v2) {
		t.Fatalf("ConfigGet cache hit mismatch, first=%#v, second=%#v", v1, v2)
	}

	t.Logf("ConfigGet cache hit value: %#v", v1)
}
