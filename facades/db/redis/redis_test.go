package redis

import (
	"testing"

	"github.com/go-redis/redis/v8"
)

// TestFormatZSet 测试有序集合格式化函数
func TestFormatZSet(t *testing.T) {
	tests := []struct {
		name     string
		data     []redis.Z
		fields   []string
		expected []map[string]any
	}{
		{
			name:     "空数据",
			data:     []redis.Z{},
			fields:   nil,
			expected: []map[string]any{},
		},
		{
			name: "默认字段名",
			data: []redis.Z{
				{Score: 1.0, Member: "a"},
				{Score: 2.0, Member: "b"},
			},
			fields: nil,
			expected: []map[string]any{
				{"data": "a", "score": 1.0},
				{"data": "b", "score": 2.0},
			},
		},
		{
			name: "自定义字段名",
			data: []redis.Z{
				{Score: 100.5, Member: "member1"},
				{Score: 200.5, Member: "member2"},
			},
			fields: []string{"member", "value"},
			expected: []map[string]any{
				{"member": "member1", "value": 100.5},
				{"member": "member2", "value": 200.5},
			},
		},
		{
			name: "单个元素",
			data: []redis.Z{
				{Score: 99.9, Member: "single"},
			},
			fields: []string{"key", "rank"},
			expected: []map[string]any{
				{"key": "single", "rank": 99.9},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []map[string]any
			if len(tt.fields) >= 2 {
				result = FormatZSet(tt.data, tt.fields[0], tt.fields[1])
			} else {
				result = FormatZSet(tt.data)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("FormatZSet() 返回长度 = %d, 期望 %d", len(result), len(tt.expected))
				return
			}

			for i, item := range result {
				for key, expectedVal := range tt.expected[i] {
					if item[key] != expectedVal {
						t.Errorf("FormatZSet()[%d][%s] = %v, 期望 %v", i, key, item[key], expectedVal)
					}
				}
			}
		})
	}
}

// TestFormatHash 测试 hash 格式化函数
func TestFormatHash(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]string
		expected map[string]string
	}{
		{
			name:     "空 map",
			data:     map[string]string{},
			expected: map[string]string{},
		},
		{
			name: "单个字段",
			data: map[string]string{
				"field1": "value1",
			},
			expected: map[string]string{
				"field1": "value1",
			},
		},
		{
			name: "多个字段",
			data: map[string]string{
				"name":  "test",
				"age":   "25",
				"email": "test@example.com",
			},
			expected: map[string]string{
				"name":  "test",
				"age":   "25",
				"email": "test@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatHash(tt.data)

			if len(result) != len(tt.expected) {
				t.Errorf("FormatHash() 返回长度 = %d, 期望 %d", len(result), len(tt.expected))
				return
			}

			for key, expectedVal := range tt.expected {
				if result[key] != expectedVal {
					t.Errorf("FormatHash()[%s] = %v, 期望 %v", key, result[key], expectedVal)
				}
			}
		})
	}
}

// TestRequestLog 测试请求日志结构
func TestRequestLog(t *testing.T) {
	tests := []struct {
		name    string
		log     RequestLog
		wantCmd string
		wantOk  bool
	}{
		{
			name: "GET 命令成功",
			log: RequestLog{
				Command: "GET",
				Args:    []any{"test_key"},
				Result:  "test_value",
				Ok:      true,
			},
			wantCmd: "GET",
			wantOk:  true,
		},
		{
			name: "SET 命令成功",
			log: RequestLog{
				Command: "SET",
				Args:    []any{"key", "value", 0},
				Result:  nil,
				Ok:      true,
			},
			wantCmd: "SET",
			wantOk:  true,
		},
		{
			name: "命令失败",
			log: RequestLog{
				Command: "GET",
				Args:    []any{"nonexistent"},
				Result:  nil,
				Ok:      false,
			},
			wantCmd: "GET",
			wantOk:  false,
		},
		{
			name: "HSET 命令",
			log: RequestLog{
				Command: "HSET",
				Args:    []any{"hash_key", "field1", "value1"},
				Result:  nil,
				Ok:      true,
			},
			wantCmd: "HSET",
			wantOk:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.log.Command != tt.wantCmd {
				t.Errorf("RequestLog.Command = %v, 期望 %v", tt.log.Command, tt.wantCmd)
			}
			if tt.log.Ok != tt.wantOk {
				t.Errorf("RequestLog.Ok = %v, 期望 %v", tt.log.Ok, tt.wantOk)
			}
		})
	}
}

// TestConnectArgs 测试连接参数结构
func TestConnectArgs(t *testing.T) {
	tests := []struct {
		name string
		args ConnectArgs
	}{
		{
			name: "基本参数",
			args: ConnectArgs{
				BusKey: "fe_m_redis",
			},
		},
		{
			name: "带 hash key",
			args: ConnectArgs{
				BusKey:  "user_read",
				HashKey: "user_123",
			},
		},
		{
			name: "指定节点",
			args: ConnectArgs{
				BusKey: "counter_redis",
				HashNo: 2,
			},
		},
		{
			name: "完整参数",
			args: ConnectArgs{
				BusKey:  "fe_m_redis",
				HashKey: "test_hash",
				HashNo:  0,
				Timeout: 500,
				DB:      1,
			},
		},
		{
			name: "使用默认 DB",
			args: ConnectArgs{
				BusKey:  "fe_s_redis",
				HashKey: "key",
				DB:      -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.BusKey == "" {
				t.Error("BusKey 不应为空")
			}
		})
	}
}

// TestConnectionMethods 测试 Connection 结构的辅助方法
func TestConnectionMethods(t *testing.T) {
	conn := &Connection{
		impl:    nil,
		busKey:  "test_bus",
		request: make([]RequestLog, 0),
	}

	t.Run("GetBusKey", func(t *testing.T) {
		if conn.GetBusKey() != "test_bus" {
			t.Errorf("GetBusKey() = %v, 期望 %v", conn.GetBusKey(), "test_bus")
		}
	})

	t.Run("GetRequest_empty", func(t *testing.T) {
		requests := conn.GetRequest()
		if len(requests) != 0 {
			t.Errorf("GetRequest() 长度 = %d, 期望 0", len(requests))
		}
	})

	t.Run("addRequest", func(t *testing.T) {
		conn.addRequest("GET", []any{"key1"}, "value1", true)
		conn.addRequest("SET", []any{"key2", "value2"}, nil, true)
		conn.addRequest("DEL", []any{"key3"}, nil, false)

		requests := conn.GetRequest()
		if len(requests) != 3 {
			t.Errorf("GetRequest() 长度 = %d, 期望 3", len(requests))
		}

		if requests[0].Command != "GET" {
			t.Errorf("requests[0].Command = %v, 期望 GET", requests[0].Command)
		}
		if requests[0].Result != "value1" {
			t.Errorf("requests[0].Result = %v, 期望 value1", requests[0].Result)
		}
		if !requests[0].Ok {
			t.Error("requests[0].Ok 应为 true")
		}

		if requests[2].Ok {
			t.Error("requests[2].Ok 应为 false")
		}
	})
}

// TestManagerSingleton 测试 Manager 单例模式
func TestManagerSingleton(t *testing.T) {
	m1 := GetManager()
	m2 := GetManager()

	// facades 层的 Manager 是包装器，每次调用会创建新的包装器
	// 但内部的 internalRedis.Manager 是单例
	// 所以这里只验证 GetManager 返回非 nil
	if m1 == nil {
		t.Error("GetManager() 不应返回 nil")
	}

	if m2 == nil {
		t.Error("GetManager() 不应返回 nil")
	}

	// 验证内部实现是同一个单例（通过 impl 字段）
	// 由于 impl 是私有字段，我们通过功能测试来验证
	// 这里只验证 Manager 结构正确初始化
	if m1.impl == nil {
		t.Error("Manager.impl 不应为 nil")
	}

	// 验证两个 Manager 包装器共享同一个内部实现
	if m1.impl != m2.impl {
		t.Error("Manager 应共享同一个内部实现")
	}
}

// TestFormatZSetEdgeCases 测试 FormatZSet 边界情况
func TestFormatZSetEdgeCases(t *testing.T) {
	t.Run("only_one_field", func(t *testing.T) {
		data := []redis.Z{
			{Score: 1.0, Member: "test"},
		}
		result := FormatZSet(data, "onlyOne")
		if result[0]["data"] != "test" {
			t.Errorf("期望使用默认字段名 data")
		}
	})

	t.Run("nil_data", func(t *testing.T) {
		result := FormatZSet(nil)
		if len(result) != 0 {
			t.Errorf("nil 数据应返回空切片，实际长度 = %d", len(result))
		}
	})

	t.Run("special_chars", func(t *testing.T) {
		data := []redis.Z{
			{Score: 1.0, Member: "中文测试"},
			{Score: 2.0, Member: "special!@#$%"},
			{Score: 3.0, Member: ""},
		}
		result := FormatZSet(data)
		if len(result) != 3 {
			t.Errorf("期望长度 3，实际 %d", len(result))
		}
		if result[0]["data"] != "中文测试" {
			t.Errorf("中文成员处理错误")
		}
	})

	t.Run("negative_score", func(t *testing.T) {
		data := []redis.Z{
			{Score: -100.5, Member: "negative"},
		}
		result := FormatZSet(data)
		if result[0]["score"] != -100.5 {
			t.Errorf("负数分数处理错误，期望 -100.5，实际 %v", result[0]["score"])
		}
	})
}

// TestFormatHashEdgeCases 测试 FormatHash 边界情况
func TestFormatHashEdgeCases(t *testing.T) {
	t.Run("nil_map", func(t *testing.T) {
		result := FormatHash(nil)
		if result != nil {
			t.Error("nil 输入应返回 nil")
		}
	})

	t.Run("empty_values", func(t *testing.T) {
		data := map[string]string{
			"empty": "",
			"space": " ",
		}
		result := FormatHash(data)
		if result["empty"] != "" {
			t.Error("空值字段处理错误")
		}
		if result["space"] != " " {
			t.Error("空格值字段处理错误")
		}
	})

	t.Run("special_keys", func(t *testing.T) {
		data := map[string]string{
			"中文键":       "中文值",
			"key:with:": "colon",
		}
		result := FormatHash(data)
		if result["中文键"] != "中文值" {
			t.Error("中文键值处理错误")
		}
	})
}

// TestRequestLogArgs 测试 RequestLog 的 Args 字段
func TestRequestLogArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []any
		expected int
	}{
		{
			name:     "空参数",
			args:     []any{},
			expected: 0,
		},
		{
			name:     "单个参数",
			args:     []any{"key"},
			expected: 1,
		},
		{
			name:     "多个参数",
			args:     []any{"key", "value", 100, true},
			expected: 4,
		},
		{
			name:     "混合类型参数",
			args:     []any{"string", 123, 45.67, true, nil},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := RequestLog{
				Command: "TEST",
				Args:    tt.args,
				Ok:      true,
			}
			if len(log.Args) != tt.expected {
				t.Errorf("Args 长度 = %d, 期望 %d", len(log.Args), tt.expected)
			}
		})
	}
}

// BenchmarkFormatZSet 性能测试 FormatZSet
func BenchmarkFormatZSet(b *testing.B) {
	data := make([]redis.Z, 100)
	for i := 0; i < 100; i++ {
		data[i] = redis.Z{Score: float64(i), Member: "member"}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatZSet(data, "member", "score")
	}
}

// BenchmarkFormatHash 性能测试 FormatHash
func BenchmarkFormatHash(b *testing.B) {
	data := make(map[string]string)
	for i := 0; i < 100; i++ {
		data["field"+string(rune(i))] = "value"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatHash(data)
	}
}
