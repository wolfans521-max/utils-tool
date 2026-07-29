package kwfilter

/*
#cgo amd64 LDFLAGS: -L${SRCDIR}/lib/amd64 -lkwords_filter -lsensitive_words -lstdc++ -pthread
#cgo arm64 LDFLAGS: -L${SRCDIR}/lib/arm64 -lkwords_filter -lsensitive_words -lstdc++ -pthread
#include <stdlib.h>

// C函数声明
char* filter_exec(const char *dict_path, const char* log_path, const char* str, char* json_str, int* match_ret, int rand);
*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

// 默认配置常量
const (
	DefaultDictPath      = "/data1/apache2/config/sensitive_word_list_new_I.dict"
	DefaultLogPath       = "/dev/null"
	DefaultRandomPercent = 10
)

// FilterConfig 过滤配置
type FilterConfig struct {
	DictPath      string // 词典路径
	LogPath       string // 日志路径
	RandomPercent int    // 词典更新检测概率(0-100)
}

// Kwfilter 敏感词过滤器
type Kwfilter struct {
	config      *FilterConfig
	mu          sync.RWMutex
	initialized bool
}

// New 创建新的过滤器实例
// 可以使用以下方式调用:
//   - New()         使用所有默认配置
//   - New(config)   使用提供的配置，空字段自动填充默认值
func New(config ...*FilterConfig) *Kwfilter {
	var cfg *FilterConfig
	if len(config) == 0 || config[0] == nil {
		cfg = &FilterConfig{}
	} else {
		cfg = config[0]
	}

	if cfg.DictPath == "" {
		cfg.DictPath = DefaultDictPath
	}
	if cfg.LogPath == "" {
		cfg.LogPath = DefaultLogPath
	}
	if cfg.RandomPercent <= 0 || cfg.RandomPercent > 100 {
		cfg.RandomPercent = DefaultRandomPercent
	}

	kf := &Kwfilter{
		config: cfg,
	}

	kf.initialized = true

	return kf
}

// Filter 执行敏感词过滤
// text: 待过滤文本
// attrJSON: JSON格式的属性配置
// 返回: 过滤结果JSON字符串, 匹配等级, 错误
func (k *Kwfilter) FilterExec(text string, attrJSON string) (string, int, error) {
	if !k.initialized {
		return "", -1, errors.New("kwfilter not initialized")
	}

	if text == "" {
		return "", -1, errors.New("text cannot be empty")
	}

	if attrJSON == "" {
		return "", -1, errors.New("attrJSON cannot be empty")
	}

	// 转换为C字符串
	cDictPath := C.CString(k.config.DictPath)
	defer C.free(unsafe.Pointer(cDictPath))

	cLogPath := C.CString(k.config.LogPath)
	defer C.free(unsafe.Pointer(cLogPath))

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cAttr := C.CString(attrJSON)
	defer C.free(unsafe.Pointer(cAttr))

	// 使用 Go 变量来接收匹配结果
	var matchRet C.int

	cResult := C.filter_exec(cDictPath, cLogPath, cText, cAttr, &matchRet, C.int(k.config.RandomPercent))

	if cResult == nil {
		return "", int(matchRet), nil
	}

	// 转换结果
	result := C.GoString(cResult)
	C.free(unsafe.Pointer(cResult))

	return result, int(matchRet), nil
}

// Close 释放资源
func (k *Kwfilter) Close() {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.initialized = false
}
