// Package owner 提供日志文件所有者管理功能
// 从 log.yaml 读取 user/group 配置，支持创建文件时自动修改所有者
package owner

import (
	"os"
	"os/user"
	"strconv"
	"sync"

	"git.intra.weibo.com/search_fe/wbutil-go/config"
	"github.com/goccy/go-yaml"
)

// ChownToConfigOwner 将文件/目录所有者修改为 log.yaml 中配置的用户/组
// 如果配置中未设置 user/group，则不执行任何操作
func ChownToConfigOwner(path string) {
	cfg := GetConfig()
	if cfg.UID == -1 && cfg.GID == -1 {
		return
	}
	_ = os.Chown(path, cfg.UID, cfg.GID)
}

// Config 文件所有者配置
type Config struct {
	User  string
	Group string
	UID   int // 解析后的 UID，-1 表示未设置
	GID   int // 解析后的 GID，-1 表示未设置
}

var (
	cfgOnce sync.Once
	cfg     *Config
)

// GetConfig 获取日志文件所有者配置（只解析一次）
func GetConfig() *Config {
	cfgOnce.Do(func() {
		cfg = &Config{UID: -1, GID: -1}
		parseConfig()
	})
	return cfg
}

// parseConfig 解析 log.yaml 中的 user/group 配置
func parseConfig() {
	data := config.LogYaml
	if len(data) == 0 {
		return
	}

	raw := make(map[string]any)
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return
	}

	// 解析 user
	if userVal, ok := raw["user"]; ok {
		cfg.User, cfg.UID = parseUserGroup(userVal)
	}

	// 解析 group
	if groupVal, ok := raw["group"]; ok {
		cfg.Group, cfg.GID = parseUserGroup(groupVal)
	}

	// 如果解析出了用户名但没有 UID，尝试从系统查找
	if cfg.User != "" && cfg.UID == -1 {
		if u, err := user.Lookup(cfg.User); err == nil {
			if uid, err := strconv.Atoi(u.Uid); err == nil {
				cfg.UID = uid
			}
			// 同时获取用户的默认 GID
			if cfg.GID == -1 {
				if gid, err := strconv.Atoi(u.Gid); err == nil {
					cfg.GID = gid
				}
			}
		}
	}

	// 如果解析出了组名但没有 GID，尝试从系统查找
	if cfg.Group != "" && cfg.GID == -1 {
		if g, err := user.LookupGroup(cfg.Group); err == nil {
			if gid, err := strconv.Atoi(g.Gid); err == nil {
				cfg.GID = gid
			}
		}
	}
}

// parseUserGroup 解析用户/组值，支持 string、int、[]interface{} 格式
func parseUserGroup(val any) (string, int) {
	switch v := val.(type) {
	case string:
		return v, parseID(v)
	case int:
		return strconv.Itoa(v), v
	case []interface{}:
		// 处理 user:\n  daemon 这种格式
		if len(v) > 0 {
			if str, ok := v[0].(string); ok {
				return str, parseID(str)
			}
		}
	}
	return "", -1
}

// parseID 将字符串解析为数字 ID（支持数字字符串）
func parseID(s string) int {
	if s == "" {
		return -1
	}
	if id, err := strconv.Atoi(s); err == nil {
		return id
	}
	return -1
}
