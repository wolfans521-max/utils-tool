package monitor

import (
	"fmt"
	"net/http"
	"time"
)

// Middleware HTTP监控中间件
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		method := r.Method
		reqInFlight.Inc()

		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		// 记录指标
		status := fmt.Sprintf("%d", rw.statusCode)
		duration := time.Since(start).Seconds()
		RecordRequest(path, method, status, duration)

		// 记录错误指标
		if rw.statusCode >= 400 {
			reqErrors.WithLabelValues(path, method).Inc()
		}
		reqInFlight.Dec()
	})
}

// responseWriter 包装http.ResponseWriter以获取状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
