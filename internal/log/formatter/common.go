package formatter

import (
	"fmt"
	"strings"
	"time"

	"git.intra.weibo.com/search_fe/wbutil-go/facades/context"
	"github.com/bytedance/sonic"
)

func Common(msg any, ctx *context.Context) string {
	// 从 ctx 获取格式化选项
	noFormat := false
	jsonLog := false
	delimiter := "|"
	noTimestamp := false

	if ctx != nil {
		noFormat = ctx.DefaultRequestParam("no_format", "") == "true"
		jsonLog = ctx.DefaultRequestParam("json_log", "") == "true"
		if d := ctx.DefaultRequestParam("delimiter", ""); d != "" {
			delimiter = d
		}
		noTimestamp = ctx.DefaultRequestParam("no_timestamp", "") == "true"
	}

	// no_format：完全不做格式化，也不附加换行/时间戳
	if noFormat {
		return toString(msg)
	}

	// json_log：使用制表符连接各字段
	// 根据 no_timestamp 配置决定是否添加时间戳
	if jsonLog {
		if arr, ok := toInterfaceSlice(msg); ok {
			out := make([]string, 0, len(arr)+1)
			// 如果 no_timestamp 不为 true，则添加时间戳
			if !noTimestamp {
				out = append(out, time.Now().Format("2006-01-02 15:04:05"))
			}
			for _, v := range arr {
				out = append(out, toString(v))
			}
			return strings.Join(out, "\t") + "\n"
		}
		if !noTimestamp {
			return time.Now().Format("2006-01-02 15:04:05") + "\t" + toString(msg) + "\n"
		}
		return toString(msg) + "\n"
	}

	addTs := !noTimestamp

	var ts string
	if addTs {
		ts = time.Now().Format("2006-01-02 15:04:05")
	}

	msgVal := msg
	if s, ok := msg.(string); ok {
		var decoded any
		if err := sonic.Unmarshal([]byte(s), &decoded); err == nil {
			msgVal = decoded
		} else {
			msgVal = s
		}
	}

	// 1) 最终仍然是 string
	if s, ok := msgVal.(string); ok {
		if ts != "" {
			return ts + delimiter + s + "\n"
		}
		return s + "\n"
	}

	// 2) slice 类型
	if arr, ok := toInterfaceSlice(msgVal); ok {
		out := make([]string, 0, len(arr)+1)
		if ts != "" {
			out = append(out, ts)
		}
		for _, v := range arr {
			out = append(out, toString(v))
		}
		return strings.Join(out, delimiter) + "\n"
	}

	// 3) map[string]any 类型
	if m, ok := msgVal.(map[string]any); ok {
		out := make([]string, 0, len(m)+1)
		if ts != "" {
			out = append(out, ts)
		}
		for k, v := range m {
			out = append(out, fmt.Sprintf("%s:%v", k, v))
		}
		return strings.Join(out, delimiter) + "\n"
	}

	// 4) 其他类型，退回到 sprintf
	if ts != "" {
		return ts + delimiter + fmt.Sprint(msgVal) + "\n"
	}
	return fmt.Sprint(msgVal) + "\n"
}

// toInterfaceSlice 尝试将各种常见切片类型转为 []any，用于数组场景处理。
func toInterfaceSlice(v any) ([]any, bool) {
	if v == nil {
		return nil, false
	}
	switch t := v.(type) {
	case []any:
		return t, true
	case []string:
		out := make([]any, len(t))
		for i, s := range t {
			out[i] = s
		}
		return out, true
	case []int:
		out := make([]any, len(t))
		for i, n := range t {
			out[i] = n
		}
		return out, true
	case []float64:
		out := make([]any, len(t))
		for i, n := range t {
			out[i] = n
		}
		return out, true
	}
	return nil, false
}

// toStringOnly：仅处理最常见的 string/[]byte/Stringer
func toStringOnly(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	case fmt.Stringer:
		return t.String()
	}
	return ""
}

// toString：通用转字符串，优先走 toStringOnly。
func toString(v any) string {
	if v == nil {
		return ""
	}
	if s := toStringOnly(v); s != "" {
		return s
	}
	return fmt.Sprint(v)
}
