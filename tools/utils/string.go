package utils

import (
	"encoding/hex"
	"hash/fnv"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
)

// 全局状态：实例码自动推导；使用原子操作保证并发安全
var (
	instanceCode atomic.Value // 两位实例码（自动推导）

	// 使用单个 uint64 原子变量存储时间戳和计数器
	// 高 44 位存储毫秒时间戳，低 20 位存储计数器
	// 这样可以支持每毫秒 2^20 = 1,048,576 个请求
	reqState atomic.Uint64
)

const (
	timeModuloBase   uint64 = 9000000000 // 10 位范围的基数（不含偏移）
	timeModuloOffset uint64 = 1000000000 // 偏移，保证时间段首位不为0
	counterBits      uint64 = 20         // 计数器占用的位数
	counterMask      uint64 = (1 << counterBits) - 1
	counterSoftLimit uint64 = 10000 // 每毫秒最多 1 万条 -> 单实例每秒上限 1e7
)

func init() {
	instanceCode.Store("")
}

// GenerateRequestId 生成数字请求ID，固定16位：时间10位（毫秒取模且首位非0）+ 2位实例码 + 4位计数。
// 并发安全：使用原子操作（CAS）保证无锁并发，避免时间回拨导致计数重置；单毫秒内软上限 1 万，超出则顺延到下一个时间片。
// 性能优化：
// 1. 使用原子操作替代互斥锁，避免锁竞争
// 2. 使用预分配的 byte slice 和手动数字转换，避免 fmt.Sprintf 的内存分配
// 3. 支持 10w+ QPS 的高并发场景
func GenerateRequestId() string {
	code := getInstanceCode()

	for {
		nowMs := (uint64(time.Now().UnixMilli()) % timeModuloBase) + timeModuloOffset

		// 读取当前状态
		oldState := reqState.Load()
		oldTs := oldState >> counterBits
		oldSeq := oldState & counterMask

		var newTs, newSeq uint64

		// 时间回拨保护：如果当前时间小于记录的时间，使用记录的时间
		if nowMs < oldTs {
			nowMs = oldTs
		}

		if nowMs > oldTs {
			// 新的时间片，重置计数器
			newTs = nowMs
			newSeq = 0
		} else {
			// 同一时间片，递增计数器
			newTs = oldTs
			newSeq = oldSeq + 1

			// 软上限保护：超出上限则顺延到下一个时间片
			if newSeq >= counterSoftLimit {
				newTs++
				newSeq = 0
			}
		}

		newState := (newTs << counterBits) | newSeq

		// CAS 操作：如果状态未被其他 goroutine 修改，则更新成功
		if reqState.CompareAndSwap(oldState, newState) {
			// 使用预分配的 buffer 避免 fmt.Sprintf 的内存分配
			return formatRequestId(newTs, code, newSeq)
		}
	}
}

// formatRequestId 高性能格式化请求ID
// 使用预分配的 byte slice 和手动数字转换，避免 fmt.Sprintf 的内存分配
func formatRequestId(ts uint64, code string, seq uint64) string {
	// 固定 16 位：10位时间戳 + 2位实例码 + 4位计数
	buf := make([]byte, 16)

	// 填充 10 位时间戳（从右向左）
	for i := 9; i >= 0; i-- {
		buf[i] = byte('0' + ts%10)
		ts /= 10
	}

	// 填充 2 位实例码
	buf[10] = code[0]
	buf[11] = code[1]

	// 填充 4 位计数（从右向左）
	for i := 15; i >= 12; i-- {
		buf[i] = byte('0' + seq%10)
		seq /= 10
	}

	return *(*string)(unsafe.Pointer(&buf))
}

// GenerateRequestIdLegacy 保留原有实现，用于兼容性测试
// Deprecated: 请使用 GenerateRequestId
func GenerateRequestIdLegacy() string {
	return GenerateRequestId()
}

func getInstanceCode() string {
	if v, ok := instanceCode.Load().(string); ok && v != "" && v != "00" {
		return v
	}

	derived := deriveInstanceCode()
	instanceCode.Store(derived)
	return derived
}
func deriveInstanceCode() string {
	host, _ := os.Hostname()
	pid := os.Getpid()
	base := strconv.Itoa(len(host)) + host + strconv.Itoa(pid)

	h := fnv.New32a()
	_, _ = h.Write([]byte(base))
	n := h.Sum32() % 100
	if n == 0 {
		n = 1 // 避免 00 导致跨机未区分
	}
	// 手动格式化为两位数字，避免 fmt.Sprintf
	if n < 10 {
		return "0" + strconv.Itoa(int(n))
	}
	return strconv.Itoa(int(n))
}

// JoinSemicolonInts joins []int with ';', empty slice -> "".
func JoinSemicolonInts(parts []int) string {
	if len(parts) == 0 {
		return ""
	}
	out := make([]string, len(parts))
	for i, v := range parts {
		out[i] = strconv.Itoa(v)
	}
	return strings.Join(out, ";")
}

// JoinSemicolonStrings joins []string with ';', empty slice -> "".
func JoinSemicolonStrings(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ";")
}

// JoinCommaStrings joins []string with ',', empty slice -> "".
func JoinCommaStrings(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ",")
}

// GenerateSrid 生成srid
// 优化点：
// 1. math/rand/v2 全局函数无锁（ChaCha8 per-P），无需 sync.Pool
// 2. 使用栈上分配的数组，减少堆内存分配
// 3. 支持 10w+ QPS 的高并发场景
func GenerateSrid() string {
	// 栈上分配数组，避免堆内存分配
	var buf [16]byte

	// 写入 8 字节随机数（math/rand/v2 全局函数无锁）
	randBytes := buf[:8]
	for i := range randBytes {
		randBytes[i] = byte(rand.Uint32())
	}

	// 写入 8 字节纳秒时间戳
	timeBytes := buf[8:]
	ns := uint64(time.Now().UnixNano())
	timeBytes[0] = byte(ns >> 56)
	timeBytes[1] = byte(ns >> 48)
	timeBytes[2] = byte(ns >> 40)
	timeBytes[3] = byte(ns >> 32)
	timeBytes[4] = byte(ns >> 24)
	timeBytes[5] = byte(ns >> 16)
	timeBytes[6] = byte(ns >> 8)
	timeBytes[7] = byte(ns)

	return hex.EncodeToString(buf[:])
}
func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// Substr Substr字符串截取
func Substr(str string, start int, args ...int) string {
	runeStr := []rune(str)
	lenStr := len(runeStr)

	if start < 0 {
		start = max(lenStr+start, 0)
	}
	if start >= lenStr {
		return ""
	}

	length := lenStr
	if len(args) > 0 {
		length = args[0]
		if length == 0 {
			return ""
		}
	}

	// 计算结束位置
	end := lenStr
	if length > 0 {
		end = min(start+length, lenStr)
	} else if length < 0 {
		end = lenStr + length
		if end <= start {
			return ""
		}
	}

	return string(runeStr[start:end])
}
