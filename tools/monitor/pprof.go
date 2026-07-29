package monitor

import (
	"net/http"
	"net/http/pprof"
)

var defaultPProfConfig = PProfConfig{
	Enabled: true,
	Prefix:  "/debug/pprof",
}

// RegisterPProf 注册pprof性能分析路由
func RegisterPProf(mux *http.ServeMux, config ...PProfConfig) {
	cfg := defaultPProfConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if !cfg.Enabled {
		return
	}

	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "/debug/pprof"
	}

	// 注册pprof路由
	mux.HandleFunc(prefix+"/", pprof.Index)
	mux.HandleFunc(prefix+"/cmdline", pprof.Cmdline)
	mux.HandleFunc(prefix+"/profile", pprof.Profile)
	mux.HandleFunc(prefix+"/symbol", pprof.Symbol)
	mux.HandleFunc(prefix+"/trace", pprof.Trace)

	// 添加BasicAuth中间件
	if cfg.BasicAuth.Username != "" && cfg.BasicAuth.Password != "" {
		auth := func(handler http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				user, pass, ok := r.BasicAuth()
				if !ok || user != cfg.BasicAuth.Username || pass != cfg.BasicAuth.Password {
					w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				handler(w, r)
			}
		}

		mux.HandleFunc(prefix+"/", auth(pprof.Index))
		mux.HandleFunc(prefix+"/cmdline", auth(pprof.Cmdline))
		mux.HandleFunc(prefix+"/profile", auth(pprof.Profile))
		mux.HandleFunc(prefix+"/symbol", auth(pprof.Symbol))
		mux.HandleFunc(prefix+"/trace", auth(pprof.Trace))
	}
}
