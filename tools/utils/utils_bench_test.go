package utils

import (
	"sync"
	"testing"
	"time"
)

// BenchmarkGenerateRequestId 测试 GenerateRequestId 在高并发下的性能
// 模拟 10w QPS 场景
func BenchmarkGenerateRequestId(b *testing.B) {
	b.Run("Single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GenerateRequestId()
		}
	})

	b.Run("Parallel-100", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				GenerateRequestId()
			}
		})
	})

	b.Run("Parallel-1000", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				GenerateRequestId()
			}
		})
	})

	b.Run("Parallel-10000", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				GenerateRequestId()
			}
		})
	})
}

// BenchmarkGenerateSrid 测试 GenerateSrid 在高并发下的性能
func BenchmarkGenerateSrid(b *testing.B) {
	b.Run("Single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GenerateSrid()
		}
	})

	b.Run("Parallel-100", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				GenerateSrid()
			}
		})
	})

	b.Run("Parallel-1000", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				GenerateSrid()
			}
		})
	})
}

// BenchmarkArrayRandom 测试 ArrayRandom 在高并发下的性能
func BenchmarkArrayRandom(b *testing.B) {
	testArray := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.Run("Single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ArrayRandom[int](testArray, 3)
		}
	})

	b.Run("Parallel-100", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ArrayRandom[int](testArray, 3)
			}
		})
	})

	b.Run("Parallel-1000", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ArrayRandom[int](testArray, 3)
			}
		})
	})
}

// BenchmarkShuffle 测试 Shuffle 在高并发下的性能
func BenchmarkShuffle(b *testing.B) {
	b.Run("Single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			Shuffle(&arr)
		}
	})

	b.Run("Parallel-100", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				Shuffle(&arr)
			}
		})
	})

	b.Run("Parallel-1000", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				Shuffle(&arr)
			}
		})
	})
}

// BenchmarkArrayColumn 测试 ArrayColumn 在高并发下的性能
func BenchmarkArrayColumn(b *testing.B) {
	testData := []map[string]any{
		{"id": 1, "name": "Alice", "age": 25},
		{"id": 2, "name": "Bob", "age": 30},
		{"id": 3, "name": "Charlie", "age": 35},
		{"id": 4, "name": "David", "age": 40},
		{"id": 5, "name": "Eve", "age": 45},
	}

	b.Run("Single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ArrayColumn(testData, "name")
		}
	})

	b.Run("Parallel-100", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ArrayColumn(testData, "name")
			}
		})
	})

	b.Run("Parallel-1000", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ArrayColumn(testData, "name")
			}
		})
	})
}

// TestGenerateRequestIdConcurrent 测试 GenerateRequestId 的并发唯一性
func TestGenerateRequestIdConcurrent(t *testing.T) {
	const numGoroutines = 1000
	const numIdsPerGoroutine = 1000

	idMap := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIdsPerGoroutine; j++ {
				id := GenerateRequestId()
				mu.Lock()
				if idMap[id] {
					t.Errorf("Duplicate ID generated: %s", id)
				}
				idMap[id] = true
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	t.Logf("Generated %d unique IDs", len(idMap))
}

// TestGenerateSridConcurrent 测试 GenerateSrid 的并发唯一性
func TestGenerateSridConcurrent(t *testing.T) {
	const numGoroutines = 1000
	const numIdsPerGoroutine = 1000

	idMap := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numIdsPerGoroutine; j++ {
				id := GenerateSrid()
				mu.Lock()
				if idMap[id] {
					t.Errorf("Duplicate SRID generated: %s", id)
				}
				idMap[id] = true
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	t.Logf("Generated %d unique SRIDs", len(idMap))
}

// TestArrayRandomConcurrent 测试 ArrayRandom 的并发安全性
func TestArrayRandomConcurrent(t *testing.T) {
	const numGoroutines = 1000
	const numCallsPerGoroutine = 1000

	testArray := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	start := time.Now()
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numCallsPerGoroutine; j++ {
				result := ArrayRandom[int](testArray, 3)
				if len(result) != 3 {
					t.Errorf("Expected 3 elements, got %d", len(result))
				}
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("Completed %d operations in %v", numGoroutines*numCallsPerGoroutine, elapsed)
}

// TestShuffleConcurrent 测试 Shuffle 的并发安全性
func TestShuffleConcurrent(t *testing.T) {
	const numGoroutines = 1000
	const numCallsPerGoroutine = 1000

	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	start := time.Now()
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numCallsPerGoroutine; j++ {
				arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				Shuffle(&arr)
				if len(arr) != 10 {
					t.Errorf("Expected 10 elements, got %d", len(arr))
				}
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("Completed %d operations in %v", numGoroutines*numCallsPerGoroutine, elapsed)
}

// BenchmarkArrayValue 测试 ArrayValue 在高并发下的性能
func BenchmarkArrayValue(b *testing.B) {
	testData := map[string]any{
		"user": map[string]any{
			"profile": map[string]any{
				"name": "Alice",
				"age":  25,
			},
		},
	}

	b.Run("Single", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ArrayValue(testData, "user.profile.name")
		}
	})

	b.Run("Parallel-100", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ArrayValue(testData, "user.profile.name")
			}
		})
	})

	b.Run("Parallel-1000", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ArrayValue(testData, "user.profile.name")
			}
		})
	})
}
