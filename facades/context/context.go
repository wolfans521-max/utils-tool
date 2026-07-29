package context

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	string2 "git.intra.weibo.com/search_fe/wbutil-go/tools/utils"
	"github.com/gin-gonic/gin"
)

const (
	defaultAppKey = "WB_SEARCH"
)

type Context struct {
	AbortErr    error
	ginContext  *gin.Context
	User        *User
	queryParams map[string]string
	paramsCache map[string]string
	formParams  map[string]string
	*From
	httpBody   string
	remoteIP   string
	clientIP   string
	seqId      string
	shutdowns  []shutdown
	paramMutex sync.RWMutex
	AppendBuf  map[string]map[string][]string
	AppendMu   sync.Mutex // 保护 AppendBuf 的并发访问
}

type shutdown func(*Context)

var ctxPool = sync.Pool{New: func() any {
	return &Context{
		User: &User{},
	}
}}

func New(ctx *gin.Context) *Context {
	c := ctxPool.Get().(*Context)
	c.reset()

	c.ginContext = ctx

	c.initParamCache()
	c.InitClientIP()
	c.InitRemoteIP()

	c.From = NewFrom(c.RequestParam("from"))

	return c
}

func (ctx *Context) Abort(err error) {
	ctx.AbortErr = err
	ctx.ginContext.Abort()
}

func (ctx *Context) Aborted() (aborted bool, err error) {
	return ctx.ginContext.IsAborted(), ctx.AbortErr
}

func Put(ctx *Context) {
	ctx.reset()
	ctxPool.Put(ctx)
}

func (ctx *Context) Set(key string, value any) {
	ctx.ginContext.Set(key, value)
}

func (ctx *Context) Get(key string) (any, bool) {
	return ctx.ginContext.Get(key)
}

func (ctx *Context) GetString(key string) string {
	return ctx.ginContext.GetString(key)
}

func (ctx *Context) GetHeader(key string) string {
	return ctx.ginContext.GetHeader(key)
}

func (ctx *Context) Header(key, value string) {
	ctx.ginContext.Header(key, value)
}

func (ctx *Context) InitRemoteIP() {
	ctx.remoteIP = ctx.ginContext.RemoteIP()
}

func (ctx *Context) RemoteIP() string {
	if ctx.remoteIP == "" {
		return ctx.ginContext.RemoteIP()
	}

	return ctx.remoteIP
}

func (ctx *Context) SetSeqId(seqId string) {
	ctx.seqId = seqId
}

func (ctx *Context) GetSeqId() string {
	if ctx.seqId == "" {
		return string2.GenerateRequestId()
	}

	return ctx.seqId
}

func (ctx *Context) SetFrom(from string) {
	ctx.From = NewFrom(from)
}

func (ctx *Context) GetFrom() string {
	return ctx.From.raw
}

func (ctx *Context) RequestUrl() *url.URL {
	return ctx.ginContext.Request.URL
}

func (ctx *Context) Request() *http.Request {
	return ctx.ginContext.Request
}

func (ctx *Context) Path() string {
	return ctx.ginContext.Request.URL.Path
}

func (ctx *Context) Method() string {
	return ctx.Request().Method
}

func (ctx *Context) HostPath() string {
	return ctx.ginContext.Request.Host + ctx.ginContext.Request.URL.Path
}

func (ctx *Context) SetRequestStart(appKey string) {
	if appKey == "" {
		appKey = defaultAppKey
	}
	ctx.Set(appKey, time.Now())
}

func (ctx *Context) GetRequestStart(appKey string) time.Time {
	if appKey == "" {
		appKey = defaultAppKey
	}
	t, ok := ctx.Get(appKey)
	if !ok {
		return time.Now()
	}

	return t.(time.Time)
}

func (ctx *Context) DefaultHeader(key, def string) string {
	if h := ctx.ginContext.GetHeader(key); h != "" {
		return h
	}

	return def
}

func (ctx *Context) Sid() string {
	switch ctx.RequestParam("c") {
	case "iphone", "ipad":
		return "t_wap_ios"
	case "android", "harmony":
		return "t_wap_android"
	case "h5":
		return "t_wap_h5"
	}

	if ctx.From != nil {
		switch {
		case ctx.From.IsIPhone(): // platform=3，iPhone
			return "t_wap_ios"
		case ctx.From.IsAndroid(): // platform=5，Android
			return "t_wap_android"
		case ctx.From.IsIPadHD(): // platform=9, minor_platform=01，iPad HD
			return "t_wap_ios"
		case ctx.From.IsHarmonySKF(): // platform=A, minor_platform=02，鸿蒙单框架
			return "t_wap_android"
		}
	}

	return "t_wap"
}

// reset 重置 Context 所有字段，确保对象池复用时不会泄露上一个请求的数据。
// 注意：queryParams/formParams/httpBody 由 initParamCache() 负责重置，
// 此处不重复处理，保持职责清晰。
func (ctx *Context) reset() {
	ctx.AbortErr = nil
	ctx.ginContext = nil
	// User 字段级 reset，避免每次新建对象
	if ctx.User != nil {
		ctx.User.Id = ""
	} else {
		ctx.User = &User{}
	}
	ctx.clientIP = ""
	ctx.remoteIP = ""
	ctx.seqId = ""
	ctx.paramsCache = nil
	ctx.queryParams = nil
	ctx.formParams = nil
	ctx.httpBody = ""
	ctx.shutdowns = nil
	ctx.From = nil
	ctx.AppendBuf = nil
}

func (ctx *Context) Bind(obj any) error {
	return ctx.ginContext.ShouldBind(obj)
}

func (ctx *Context) BindJson(obj any) error {
	return ctx.ginContext.ShouldBindJSON(obj)
}

func (ctx *Context) RegisterShutdown(s shutdown) {
	ctx.shutdowns = append(ctx.shutdowns, s)
}

func (ctx *Context) Shutdown() {
	if len(ctx.shutdowns) > 0 {
		for _, s := range ctx.shutdowns {
			s(ctx)
		}
	}
}

// GinContext 返回底层的gin上下文信息
func (ctx *Context) GinContext() *gin.Context {
	return ctx.ginContext
}
