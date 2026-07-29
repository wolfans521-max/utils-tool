package context

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// sprRegexCache 用于缓存 SprParams 方法中的正则表达式
// key 是动态生成的模式，value 是编译后的正则
var sprRegexCache = struct {
	mu    sync.RWMutex
	cache map[string]*regexp.Regexp
}{
	cache: make(map[string]*regexp.Regexp),
}

// getSprRegex 从缓存获取或创建正则表达式
func getSprRegex(pattern string) *regexp.Regexp {
	// 先尝试读缓存
	sprRegexCache.mu.RLock()
	if re, ok := sprRegexCache.cache[pattern]; ok {
		sprRegexCache.mu.RUnlock()
		return re
	}
	sprRegexCache.mu.RUnlock()

	// 创建新的正则表达式
	re := regexp.MustCompile(pattern)

	// 写入缓存
	sprRegexCache.mu.Lock()
	sprRegexCache.cache[pattern] = re
	sprRegexCache.mu.Unlock()

	return re
}

const (
	defaultMultipartMemory = 32 << 20
)

func (ctx *Context) initParamCache() {
	ctx.paramMutex.Lock()
	defer ctx.paramMutex.Unlock()
	ctx.paramsCache = make(map[string]string)
	ctx.queryParams = make(map[string]string)
	ctx.formParams = make(map[string]string)
	ctx.httpBody = ""

	if ctx.ginContext == nil || ctx.ginContext.Request == nil {
		return
	}

	fn := func(k string, v []string) {
		if len(v) == 0 {
			return
		}
		ctx.paramsCache[k] = v[0]
	}

	for k, v := range ctx.ginContext.Request.URL.Query() {
		fn(k, v)
		if len(v) != 0 {
			ctx.queryParams[k] = v[0]
		}
	}

	if ctx.ginContext.Request.Method == http.MethodGet {
		return
	}

	body := ctx.ginContext.Request.Body
	if b, e := ioutil.ReadAll(body); e == nil {
		ctx.httpBody = string(b)
		ctx.ginContext.Request.Body = ioutil.NopCloser(bytes.NewReader(b))
	}

	req := ctx.Request()
	if err := req.ParseMultipartForm(defaultMultipartMemory); err != nil && len(req.PostForm) == 0 {
		return
	}

	for k, v := range req.PostForm {
		fn(k, v)
		if len(v) != 0 {
			ctx.formParams[k] = v[0]
		}
	}
}

func (ctx *Context) SetRequestParam(key, value string) {
	ctx.paramMutex.Lock()
	defer ctx.paramMutex.Unlock()
	if ctx.paramsCache == nil {
		ctx.paramsCache = make(map[string]string)
	}

	ctx.paramsCache[key] = value
}

// RequestParam 获取 request param, 包含Get和POST参数
func (ctx *Context) RequestParam(key string) string {
	v, _ := ctx.GetRequestParam(key)

	return v
}

// RequestParams 获取多个param, 每个key都会返回
func (ctx *Context) RequestParams(keys ...string) map[string]string {
	r := make(map[string]string, len(keys))

	for _, k := range keys {
		r[k] = ctx.RequestParam(k)
	}

	return r
}

func (ctx *Context) GetRequestParam(key string) (string, bool) {
	ctx.paramMutex.RLock()
	defer ctx.paramMutex.RUnlock()

	v, ok := ctx.paramsCache[key]

	return v, ok
}

func (ctx *Context) RequestBody() string {
	return ctx.httpBody
}

func (ctx *Context) RequestInt(key string) int {
	v, ok := ctx.GetRequestParam(key)
	if !ok {
		return 0
	}

	i, _ := strconv.Atoi(v)

	return i
}

func (ctx *Context) RequestAll() map[string]any {
	ctx.paramMutex.RLock()
	defer ctx.paramMutex.RUnlock()

	params := make(map[string]any, len(ctx.paramsCache))
	for k, v := range ctx.paramsCache {
		params[k] = v
	}

	return params
}

func (ctx *Context) QueryParamsAll() map[string]string {
	ctx.paramMutex.RLock()
	defer ctx.paramMutex.RUnlock()

	params := make(map[string]string, len(ctx.queryParams))
	for k, v := range ctx.queryParams {
		params[k] = v
	}

	return params
}

func (ctx *Context) FormParamsAll() map[string]string {
	ctx.paramMutex.RLock()
	defer ctx.paramMutex.RUnlock()

	params := make(map[string]string, len(ctx.formParams))
	for k, v := range ctx.formParams {
		params[k] = v
	}

	return params
}

func (ctx *Context) SprParams(key string) string {
	if !ctx.RequestHas("spr") {
		return ""
	}

	spr := ctx.RequestParam("spr")
	if !strings.Contains(spr, key) {
		return ""
	}

	// 使用缓存的正则表达式
	pattern := key + ":(.*?);"
	reg := getSprRegex(pattern)
	match := reg.FindStringSubmatch(spr)

	if len(match) > 1 {
		return match[1]
	}

	return ""
}

func (ctx *Context) RequestHas(key string) bool {
	_, ok := ctx.GetRequestParam(key)

	return ok
}

// DefaultRequestParam 获取request param, 没有则返回默认值
func (ctx *Context) DefaultRequestParam(key, def string) string {
	v, ok := ctx.GetRequestParam(key)
	if ok {
		return v
	}

	return def
}

func (ctx *Context) Timestamp() string {
	return fmt.Sprintf("%d", ctx.Time())
}

func (ctx *Context) Time() int64 {
	return time.Now().Unix()
}

func (ctx *Context) InitClientIP() {
	if ctx.DefaultHeader("X-Shanhai-Flag", "") == "1" {
		ApiRemoteIP := ctx.DefaultHeader("API-RemoteIPV6", "")
		if ApiRemoteIP != "" {
			ctx.clientIP = ApiRemoteIP
			return
		}
	}

	ApiRemoteIP := ctx.DefaultHeader("API-RemoteIP", "")
	if ApiRemoteIP != "" {
		ctx.clientIP = ApiRemoteIP
		return
	}

	if ctx.clientIP == "" {
		ctx.clientIP = ctx.ginContext.ClientIP()
	}
}

func (ctx *Context) ClientIP() string {
	if ctx.clientIP == "" {
		if ctx.ginContext == nil {
			log.Printf("[PANIC_DEBUG] Context.ginContext is nil in ClientIP() method")
			return "127.0.0.1"
		}

		if ctx.ginContext.Request == nil {
			log.Printf("[PANIC_DEBUG] Context.ginContext.Request is nil in ClientIP() method")
			return "127.0.0.1"
		}

		return ctx.ginContext.ClientIP()
	}

	return ctx.clientIP
}

func (ctx *Context) Ua() string {
	ua := ctx.RequestParam("ua")
	if ua != "" {
		return ua
	}

	return ctx.GetHeader("User-Agent")
}

func (ctx *Context) C() string {
	return ctx.RequestParam("c")
}

func (ctx *Context) Wm() string {
	return ctx.RequestParam("wm")
}

func (ctx *Context) Lang() string {
	return ctx.DefaultRequestParam("lang", "zh_CN")
}

// GetSpr 获取spr
func (ctx *Context) GetSpr() string {
	if spr := ctx.DefaultRequestParam("spr", ""); spr != "" {
		return spr
	}

	from := ctx.GetFrom()
	wm := ctx.Wm()

	// 使用 strings.Builder 优化字符串拼接
	var sb strings.Builder
	sb.Grow(256) // 预分配足够的空间

	// 基础 spr
	sb.WriteString("from:")
	sb.WriteString(from)
	sb.WriteString(";wm:")
	sb.WriteString(wm)

	if v := ctx.DefaultRequestParam("uicode", ""); v != "" {
		sb.WriteString(";uicode:")
		sb.WriteString(v)
	}
	if v := ctx.DefaultRequestParam("containerid", ""); v != "" {
		sb.WriteString(";fid:")
		sb.WriteString(v)
	}
	if v := ctx.DefaultRequestParam("aid", ""); v != "" {
		sb.WriteString(";aid:")
		sb.WriteString(v)
		// lang 使用 ctx.Lang() 获取，但只有非默认值时才输出
		if lang := ctx.Lang(); lang != "" {
			sb.WriteString(";lang:")
			sb.WriteString(lang)
		}
	}
	if v := ctx.DefaultRequestParam("networktype", ""); v != "" {
		sb.WriteString(";networktype:")
		sb.WriteString(v)
	}
	if v := ctx.DefaultRequestParam("launchid", ""); v != "" {
		sb.WriteString(";launchid:")
		sb.WriteString(v)
	}
	if v := ctx.DefaultRequestParam("ul_hid", ""); v != "" {
		sb.WriteString(";ul_hid:")
		sb.WriteString(v)
	}
	if v := ctx.DefaultRequestParam("ul_sid", ""); v != "" {
		sb.WriteString(";ul_sid:")
		sb.WriteString(v)
	}
	if v := ctx.Ua(); v != "" {
		sb.WriteString(";ua:")
		sb.WriteString(v)
	}

	return sb.String()
}
