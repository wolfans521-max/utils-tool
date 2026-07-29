package monitor

// Metric 监控指标
type Metric struct {
	Name   string
	Value  float64
	Labels map[string]string
}

// AlertRule 告警规则
type AlertRule struct {
	Name        string
	Expression  string
	Severity    string
	Description string
}

// PProfConfig pprof配置
type PProfConfig struct {
	Enabled   bool
	Prefix    string
	BasicAuth struct {
		Username string
		Password string
	}
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	BuildAt string `json:"build_at"`
	StartAt string `json:"start_at"`
	UpTime  string `json:"up_time"`
}
