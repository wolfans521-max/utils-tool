package utils

import (
    "runtime"
    "sync"
    "testing"
)

// 压力测试：高并发生成 ID，验证长度、首位非 0、无重复
func TestGenerateRequestId_Stress(t *testing.T) {
    const goroutines = 200
    const perG = 20000 // 每个协程 2 万，总计 400 万，跨度多毫秒

    ids := make([]string, 0, goroutines*perG)
    idsMu := sync.Mutex{}
    wg := sync.WaitGroup{}

    // 提前触发 instanceCode 推导
    _ = GenerateRequestId()

    for i := 0; i < goroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            local := make([]string, 0, perG)
            for j := 0; j < perG; j++ {
                local = append(local, GenerateRequestId())
            }
            idsMu.Lock()
            ids = append(ids, local...)
            idsMu.Unlock()
        }()
    }

    wg.Wait()

    // 基本检查：长度与首位
    for idx, id := range ids {
        if len(id) != 16 {
            t.Fatalf("id length mismatch at %d: %s", idx, id)
        }
        if id[0] == '0' {
            t.Fatalf("id leading zero at %d: %s", idx, id)
        }
    }

    // 重复检测
    seen := make(map[string]struct{}, len(ids))
    for idx, id := range ids {
        if _, ok := seen[id]; ok {
            t.Fatalf("duplicate id at %d: %s", idx, id)
        }
        seen[id] = struct{}{}
    }
}

// Benchmark: 并发基准，观察吞吐
func BenchmarkGenerateRequestId_Parallel(b *testing.B) {
    runtime.GOMAXPROCS(runtime.NumCPU())
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _ = GenerateRequestId()
        }
    })
}
