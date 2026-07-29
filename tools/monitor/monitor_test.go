package monitor_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.intra.weibo.com/search_fe/wbutil-go/tools/monitor"
)

func TestHealthCheck(t *testing.T) {
	mux := http.NewServeMux()
	monitor.RegisterHealthz(mux)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestInfoHandler(t *testing.T) {
	mux := http.NewServeMux()
	monitor.SetServiceInfo(monitor.ServiceInfo{
		Name:    "test",
		Version: "1.0.0",
		BuildAt: "2025-07-02",
	})
	monitor.RegisterInfo(mux)

	req := httptest.NewRequest("GET", "/info", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRegisterAll(t *testing.T) {
	mux := http.NewServeMux()
	monitor.RegisterAll(mux, "test-service", "1.0.0", "2025-07-02")

	// 测试/metrics端点
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for /metrics, got %d", http.StatusOK, w.Code)
	}

	// 测试/info端点
	req = httptest.NewRequest("GET", "/info", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for /info, got %d", http.StatusOK, w.Code)
	}
}

func TestMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	wrapped := monitor.Middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
