package monitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
	reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
	reqInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of in-flight requests",
		},
	)
	reqErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Number of HTTP request errors",
		},
		[]string{"path", "method"},
	)
)

func init() {
	prometheus.MustRegister(reqCounter, reqDuration, reqInFlight, reqErrors)
}

// RegisterMetrics 注册Prometheus指标路由
func RegisterMetrics(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}

// RecordRequest 记录请求指标
func RecordRequest(path, method, status string, duration float64) {
	reqCounter.WithLabelValues(path, method, status).Inc()
	reqDuration.WithLabelValues(path, method).Observe(duration)
}
