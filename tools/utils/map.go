package utils

import "strings"

// MapValue defaultValue 为可变参数，不传时使用 T 的零值作为默认值
//
// 用法示例：
// utils.MapValue[map[string]any](data, "key")        // 不传默认值，未找到返回 nil
// utils.MapValue[int](data, "flag", 0)               // 传默认值，未找到返回 0
// utils.MapValue[string](data, "name", "default")    // 传默认值，未找到返回 "default"
// utils.MapValue[string](data, "key")        // 不传默认值，未找到返回 ""
func MapValue[T any](root any, key string, defaultValue ...T) T {
	var def T
	if len(defaultValue) > 0 {
		def = defaultValue[0]
	}

	if root == nil || key == "" {
		return def
	}

	current := root
	start := 0
	keyLen := len(key)

	for start < keyLen {
		end := strings.IndexByte(key[start:], '.')
		var k string
		if end == -1 {
			k = key[start:]
			start = keyLen
		} else {
			k = key[start : start+end]
			start = start + end + 1
		}

		m, ok := current.(map[string]any)
		if !ok {
			return def
		}

		val, exists := m[k]
		if !exists || val == nil {
			return def
		}
		current = val
	}

	// 最终将结果断言为泛型 T
	if actualVal, ok := current.(T); ok {
		return actualVal
	}

	return def
}
