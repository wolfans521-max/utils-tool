package target

// ElkHost ELK UDP 接收端。
type ElkHost struct {
	Host      string `json:"host" yaml:"host"`
	Port      int    `json:"port" yaml:"port"`
	TimeoutMS int    `json:"timeout_ms" yaml:"timeout_ms"`
}

// ElkConfig Elk target 配置。
//
// white_list 为空则不过滤（默认关闭白名单过滤）。
//
// 注意：具体发送行为由 [`BuildZapLogger()`](internal/log/target/zap.go:19) 内的 elkCore 实现。
type ElkConfig struct {
	Enable    bool      `json:"enable" yaml:"enable"`
	Prefix    string    `json:"prefix" yaml:"prefix"`
	Hosts     []ElkHost `json:"hosts" yaml:"hosts"`
	Tags      []string  `json:"tags" yaml:"tags"`
	Levels    []string  `json:"levels" yaml:"levels"`
	WhiteList []int64   `json:"white_list" yaml:"white_list"`
}
