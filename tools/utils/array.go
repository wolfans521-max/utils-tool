package utils

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"sort"
	"strings"

	"github.com/tidwall/gjson"
)

// ArrayValue 获取数组中指定key的value值，深度检索多个key用.分割
func ArrayValue(root any, key string, defaultValue ...any) any {
	if root == nil || key == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	current := root
	remainingKey := key

	for {
		// 寻找下一个点号
		dotIdx := strings.IndexByte(remainingKey, '.')

		// 确定当前的查找键
		var k string
		if dotIdx == -1 {
			k = remainingKey
		} else {
			k = remainingKey[:dotIdx]
		}

		// 执行查找
		m, ok := current.(map[string]any)
		if !ok {
			goto NotFound
		}

		val, exists := m[k]
		if !exists {
			goto NotFound
		}
		current = val

		// 如果没有更多层级了，返回结果
		if dotIdx == -1 {
			if current == nil {
				goto NotFound
			}
			return current
		}

		// 准备下一层循环
		remainingKey = remainingKey[dotIdx+1:]
	}

NotFound:
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

// ArrayColumn 获取数组中某一列的值
func ArrayColumn(array any, columnKey string) []any {
	switch a := array.(type) {
	case []any:
		// 预分配内存，减少扩容开销
		elems := make([]any, 0, len(a))
		for _, item := range a {
			if row, ok := item.(map[string]any); ok {
				if val, ok := row[columnKey]; ok {
					elems = append(elems, val)
				}
			}
		}
		return elems

	case []map[string]any:
		elems := make([]any, 0, len(a))
		for _, item := range a {
			if val, ok := item[columnKey]; ok {
				elems = append(elems, val)
			}
		}
		return elems

	case []string:
		elems := make([]any, 0, len(a))
		for _, item := range a {
			// 使用 Value() 保持原始类型，减少 String() 的内存分配
			result := gjson.Get(item, columnKey)
			if result.Exists() {
				elems = append(elems, result.Value())
			}
		}
		return elems

	case string:
		// 针对大字符串，提前探测长度可以更高效
		// 这里 gjson.Parse 已经是 O(n) 了，ForEach 也是 O(n)
		var elems []any
		parsed := gjson.Parse(a)
		if !parsed.IsArray() {
			return []any{}
		}

		// 只有当我们需要频繁操作解析后的对象时，Parse 才显现优势
		parsed.ForEach(func(_, v gjson.Result) bool {
			res := v.Get(columnKey)
			if res.Exists() {
				elems = append(elems, res.Value())
			}
			return true
		})
		return elems
	}
	return []any{}
}

// ArrayKeys array_keys 函数
func ArrayKeys(array any) []string {
	var keys []string
	switch a := array.(type) {
	case string:
		gjson.Parse(a).ForEach(func(k, _ gjson.Result) bool {
			keys = append(keys, k.String())
			return true
		})
	case map[string]any:
		for k, _ := range a {
			keys = append(keys, k)
		}
	default:
		if reflect.TypeOf(array).Kind() == reflect.Map {
			for _, k := range reflect.ValueOf(array).MapKeys() {
				keys = append(keys, k.String())
			}
		}
	}
	return keys
}

// ArraySearch array_search 函数
func ArraySearch(needle, haystack any) (key string, exist bool) {
	switch array := haystack.(type) {
	case map[string]string:
		for k, v := range array {
			if value, ok := needle.(string); ok && v == value {
				return k, true
			}
		}
	case map[string]int:
		for k, v := range array {
			if value, ok := needle.(int); ok && v == value {
				return k, true
			}
		}
	case map[string]int64:
		for k, v := range array {
			if value, ok := needle.(int64); ok && v == value {
				return k, true
			}
		}
	case []int:
		for i, v := range array {
			if value, ok := needle.(int); ok && v == value {
				return fmt.Sprintf("%d", i), true
			}
		}
	case []int64:
		for i, v := range array {
			if value, ok := needle.(int64); ok && v == value {
				return fmt.Sprintf("%d", i), true
			}
		}
	case []string:
		for i, v := range array {
			if value, ok := needle.(string); ok && v == value {
				return fmt.Sprintf("%d", i), true
			}
		}
	}
	return "", false
}

// ArrayShift array_shift 函数
func ArrayShift[T any](array *[]T) (frist T) {
	if len(*array) == 0 {
		return frist
	}
	frist = (*array)[0]
	*array = (*array)[1:]
	return frist
}

// ArrayPop array_pop 函数
func ArrayPop[T any](array *[]T) (last T) {
	if len(*array) == 0 {
		return last
	}
	last = (*array)[len(*array)-1]
	*array = (*array)[:len(*array)-1]
	return last
}

// ArrayFilter array_filter 函数
func ArrayFilter[V any, T map[string]V | []V](array T, callback ...func(value V) bool) T {
	// 如果没有提供 callback 回调函数，将删除数组中 array 的所有“空”元素
	fn := func(value V) bool {
		switch value := any(value).(type) {
		case nil:
			return false
		case bool:
			return value
		case string:
			return value != "" && value != "0"
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return fmt.Sprintf("%d", value) != "0"
		case float32, float64:
			return value != float32(0) && value != float64(0)
		default:
			t := reflect.TypeOf(value).Kind()
			if t == reflect.Map || t == reflect.Slice || t == reflect.Array {
				return reflect.ValueOf(value).Len() > 0
			}
		}
		return true
	}
	if len(callback) > 0 {
		fn = callback[0]
	}

	switch array := any(array).(type) {
	case map[string]V:
		result := make(map[string]V)
		for k, v := range array {
			if fn(v) {
				result[k] = v
			}
		}
		return any(result).(T)
	case []V:
		result := []V{}
		for _, v := range array {
			if fn(v) {
				result = append(result, v)
			}
		}
		return any(result).(T)
	}

	return array
}

// ArrayRandom array_random 函数
// math/rand/v2 全局函数无锁（ChaCha8 per-P），无需 sync.Pool
func ArrayRandom[V any, T map[string]V | []V](array T, num ...int) T {
	n := 1
	if len(num) > 0 {
		n = num[0]
	}
	// 指定要取出的单元数量。 必须大于零，且小于或等于 array 的长度
	if len(array) == 0 || n <= 0 || n > len(array) {
		return array
	}

	// math/rand/v2 全局函数无锁，直接使用
	casual := rand.Perm(len(array))[:n]
	sort.Ints(casual)

	switch array := any(array).(type) {
	case map[string]V:
		result := make(map[string]V, n)
		i := 0
		for k, v := range array {
			for j := 0; j < n; j++ {
				if i == casual[j] {
					result[k] = v
					break
				}
			}
			i++
		}
		return any(result).(T)
	case []V:
		result := make([]V, 0, n)
		for i, v := range array {
			for j := 0; j < n; j++ {
				if i == casual[j] {
					result = append(result, v)
					break
				}
			}
		}
		return any(result).(T)
	}

	return array
}

// Shuffle shuffle函数
// math/rand/v2 全局函数无锁（ChaCha8 per-P），无需 sync.Pool
func Shuffle[T any](array *[]T) bool {
	if len(*array) == 0 {
		return false
	}

	for i := range *array {
		j := rand.IntN(i + 1)
		(*array)[i], (*array)[j] = (*array)[j], (*array)[i]
	}
	return true
}

func Krsort[T any](array *[]T) bool {
	if len(*array) == 0 {
		return false
	}
	for i, j := 0, len(*array)-1; i < j; i, j = i+1, j-1 {
		(*array)[i], (*array)[j] = (*array)[j], (*array)[i]
	}
	return true
}

func KsortMap[T any](array map[string]T) []map[string]T {
	var slice []map[string]T
	var keys []string
	for k, _ := range array {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		slice = append(slice, map[string]T{k: array[k]})
	}
	return slice
}

func KrsortMap[T any](array map[string]T) []map[string]T {
	var slice []map[string]T
	var keys []string
	for k, _ := range array {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	for _, k := range keys {
		slice = append(slice, map[string]T{k: array[k]})
	}
	return slice
}

func ArrayFill[T any](startIndex, count int, value T) []T {
	array := make([]T, startIndex+count)
	for i := startIndex; i < startIndex+count; i++ {
		array[i] = value
	}
	return array
}

func ArrayValues[T any, A ~[]T | ~map[string]T](array A) []T {
	var r []T
	switch array := any(array).(type) {
	case []T:
		for _, v := range array {
			r = append(r, v)
		}
	case map[string]T:
		for _, v := range array {
			r = append(r, v)
		}
	}

	return r
}

func MapKeyStr[T any](array map[string]T, args ...string) string {
	var keys []string
	for k, _ := range array {
		keys = append(keys, k)
	}
	delimiter := ","
	if len(args) > 0 {
		delimiter = args[0]
	}
	return strings.Join(keys, delimiter)
}

func ArrayOnly[T any, A ~[]T | ~map[string]T, K int | string](array A, keys []K) A {
	var r A
	switch array := any(array).(type) {
	case []T:
		var t []T
		for i, v := range array {
			for _, key := range keys {
				if idx, ok := any(key).(int); ok && i == idx {
					t = append(t, v)
				}
			}
		}
		r, _ = any(t).(A)
	case map[string]T:
		t := map[string]T{}
		for k, v := range array {
			for _, key := range keys {
				if ky, ok := any(key).(string); ok && k == ky {
					t[k] = v
				}
			}
		}
		r, _ = any(t).(A)
	}
	return r
}

func ToStringArray(array any) []string {
	var a []string
	switch v := array.(type) {
	case []string:
		return v
	case string:
		if gjson.Parse(v).IsArray() {
			gjson.Parse(v).ForEach(func(_, r gjson.Result) bool {
				a = append(a, r.String())
				return true
			})
		} else {
			a = append(a, v)
		}
	case []any:
		for _, r := range v {
			a = append(a, r.(string))
		}
	}
	return a
}
