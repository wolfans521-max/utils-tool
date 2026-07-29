package utils

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/bytedance/sonic"
)

// ToMapper 接口允许类型自定义转换为 map[string]any
// 实现此接口的类型将跳过 JSON 往返转换，零反射零分配
type ToMapper interface {
	ToMap() map[string]any
}

// StructToMap 将结构体或任意类型转换为 map[string]any
// 转换优先级：
// 1. nil → nil
// 2. 已经是 map[string]any → 直接返回
// 3. 实现 ToMapper 接口 → 调用 ToMap()（零反射）
// 4. 兜底：通过 sonic JSON 序列化/反序列化转换
func StructToMap(v any) map[string]any {
	if v == nil {
		return nil
	}

	// 如果已经是 map[string]any，直接返回
	if m, ok := v.(map[string]any); ok {
		return m
	}

	// 优先使用 ToMapper 接口（零反射零分配）
	if m, ok := v.(ToMapper); ok {
		return m.ToMap()
	}

	// 兜底：通过 JSON 序列化/反序列化转换
	result := make(map[string]any)
	data, err := sonic.Marshal(v)
	if err != nil {
		return result
	}
	_ = sonic.Unmarshal(data, &result)
	return result
}

// ShuffleSlice 随机打乱切片顺序
// 使用 Fisher-Yates 洗牌算法
// math/rand/v2 全局函数无锁（ChaCha8 per-P），无需 sync.Pool
func ShuffleSlice(slice []any) {
	n := len(slice)
	if n <= 1 {
		return
	}

	for i := n - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// DebugDumpJson 调试打印
func DebugDumpJson(v map[string]any, pretty bool) {
	if pretty {
		jsonPretty, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Println("json MarshalIndent error:", err)
			return
		}
		fmt.Println(string(jsonPretty))
		return
	}
	jsonData, err := json.Marshal(v)
	if err != nil {
		fmt.Println("json Marshal error:", err)
		return
	}
	fmt.Println(string(jsonData))
	return
}

// ConvertToAnySlice 将任意切片类型转换为 []any
// 使用 JSON 序列化/反序列化进行类型转换
func ConvertToAnySlice(items any) []any {
	if items == nil {
		return nil
	}
	// 使用 JSON 序列化/反序列化进行类型转换
	jsonBytes, err := sonic.Marshal(items)
	if err != nil {
		return nil
	}
	var result []any
	if err := sonic.Unmarshal(jsonBytes, &result); err != nil {
		return nil
	}
	return result
}
