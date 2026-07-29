package target

// FileConfig 文件 target 配置（用于构建 zap core）。
// - Path:      日志目录
// - Prefix:    公共前缀
// - Delimiter: 分隔符（默认 "_"）
// - Postfix:   日期后缀（Go time layout，默认 "20060102"）
// - Extension: 扩展名（默认 ".log"）
// - Tags/Levels: 过滤条件；为空表示不过滤
//
// 注意：具体写入行为由 [`BuildZapLogger()`](internal/log/target/zap.go:19) 内的 fileCore 实现。
type FileConfig struct {
	Enable    bool     `json:"enable" yaml:"enable"`
	Path      string   `json:"path" yaml:"path"`
	Prefix    string   `json:"prefix" yaml:"prefix"`
	Delimiter string   `json:"delimiter" yaml:"delimiter"`
	Postfix   string   `json:"postfix" yaml:"postfix"`
	Extension string   `json:"extension" yaml:"extension"`
	Tags      []string `json:"tags" yaml:"tags"`
	Levels    []string `json:"levels" yaml:"levels"`
}

func sliceToSet(xs []string) map[string]struct{} {
	if len(xs) == 0 {
		return map[string]struct{}{}
	}
	m := make(map[string]struct{}, len(xs))
	for _, x := range xs {
		if x == "" {
			continue
		}
		m[x] = struct{}{}
	}
	return m
}
