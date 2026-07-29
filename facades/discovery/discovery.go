package discovery

import (
	internal "git.intra.weibo.com/search_fe/wbutil-go/internal/discovery"
)

// ConfigGet 对外暴露的配置获取接口，内部实现下沉在 internal/discovery。
//   - group 由组件内部从 ./config/app.yaml 懒加载（discovery.group）
//   - serialization 参数与底层保持一致，用于标识是否需要 JSON 语义（当前底层只返回 raw 字符串，上层可自行处理反序列化）。
//   - defaultValue：在获取失败时作为 raw 的回退值。
//
// 返回值与 internal.ConfigGet 以及 go-vintage LookupKey 对齐，仅返回 (string, error)。
func ConfigGet(key string, serialization bool, defaultValue any) (any, error) {
	return internal.ConfigGet(key, serialization, defaultValue)
}

// DumpCache 打印配置缓存
func DumpCache() map[string]any {
	return internal.DumpCache()
}

// NamingNode 对外暴露的节点结构，直接复用 internal/discovery.NamingNode。
type NamingNode = internal.NamingNode

// NamingGet 对外暴露的命名发现接口，内部实现下沉在 internal/discovery。
func NamingGet(service, cluster string) ([]NamingNode, error) {
	return internal.NamingGet(service, cluster)
}
