package monitor

import (
	"net/http"
	"sync"
)

var (
	healthChecks []func() bool
	healthMu     sync.Mutex
)

// AddHealthCheck 添加健康检查
func AddHealthCheck(check func() bool) {
	healthMu.Lock()
	defer healthMu.Unlock()
	healthChecks = append(healthChecks, check)
}

// RegisterHealthz 注册健康检查路由
func RegisterHealthz(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		healthMu.Lock()
		defer healthMu.Unlock()

		status := http.StatusOK
		for _, check := range healthChecks {
			if !check() {
				status = http.StatusServiceUnavailable
				break
			}
		}

		w.WriteHeader(status)
		_, err := w.Write([]byte(http.StatusText(status)))
		if err != nil {
			return
		}
	})
}
