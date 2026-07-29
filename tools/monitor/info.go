package monitor

import (
	"encoding/json"
	"net/http"
	"time"
)

// ServiceInfo 服务信息

var (
	serviceInfo *ServiceInfo
	startTime   = time.Now()
)

// SetServiceInfo 设置服务信息
func SetServiceInfo(info ServiceInfo) {
	info.StartAt = startTime.Format(time.RFC3339)
	serviceInfo = &info
	updateUpTime()
}

// RegisterInfo 注册服务信息路由
func RegisterInfo(mux *http.ServeMux) {
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		updateUpTime()
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(serviceInfo)
		if err != nil {
			return
		}
	})
}

func updateUpTime() {
	if serviceInfo != nil {
		serviceInfo.UpTime = time.Since(startTime).String()
	}
}
