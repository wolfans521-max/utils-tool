package monitor

import (
	"net/http"
)

// RegisterAll 注册所有监控端点
// config参数可选，用于配置pprof
func RegisterAll(mux *http.ServeMux, serviceName, version, buildTime string, config ...PProfConfig) {
	// 注册基础监控
	RegisterHealthz(mux)
	RegisterMetrics(mux)

	// 注册pprof，支持配置
	if len(config) > 0 {
		RegisterPProf(mux, config[0])
	} else {
		RegisterPProf(mux)
	}

	// 设置服务信息
	SetServiceInfo(ServiceInfo{
		Name:    serviceName,
		Version: version,
		BuildAt: buildTime,
	})
	RegisterInfo(mux)
}
