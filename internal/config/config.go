package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

// LoadSectionFromViper 组件配置信息获取
func LoadSectionFromViper(v *viper.Viper, keyPrefix string, out any) error {
	if v == nil {
		return fmt.Errorf("config: viper is nil")
	}
	if keyPrefix == "" {
		return fmt.Errorf("config: keyPrefix is empty")
	}
	if out == nil {
		return fmt.Errorf("config: out is nil")
	}

	sub := v.Sub(keyPrefix)
	if sub == nil {
		return fmt.Errorf("config: viper sub(%q) not found", keyPrefix)
	}

	if err := sub.Unmarshal(out); err != nil {
		return fmt.Errorf("config: unmarshal %q failed: %w", keyPrefix, err)
	}
	return nil
}

const (
	defaultComponentConfigFile = "config/app.yaml" // 业务在使用时需要约定定义app.yaml配置文件
	defaultFeedConfigFile      = "config/feed.yaml" // feed 配置文件
)

var (
	componentOnce sync.Once
	componentV    *viper.Viper
	componentErr  error

	feedOnce sync.Once
	feedV    *viper.Viper
	feedErr  error
)

// ComponentViper 懒加载读取使用方固定配置文件：config/app.yaml
func ComponentViper() (*viper.Viper, error) {
	componentOnce.Do(func() {
		resolved, err := findConfigUpward(defaultComponentConfigFile)
		if err != nil {
			componentErr = err
			return
		}
		v := viper.New()
		v.SetConfigFile(resolved)
		v.SetConfigType("yaml")
		if err := v.ReadInConfig(); err != nil {
			componentErr = fmt.Errorf("config: read %s failed: %w", resolved, err)
			return
		}
		componentV = v
	})
	if componentErr != nil {
		return nil, componentErr
	}
	if componentV == nil {
		return nil, fmt.Errorf("config: component viper is nil")
	}
	return componentV, nil
}

// FeedViper 懒加载读取 feed 配置文件：config/feed.yaml
func FeedViper() (*viper.Viper, error) {
	feedOnce.Do(func() {
		resolved, err := findConfigUpward(defaultFeedConfigFile)
		if err != nil {
			feedErr = err
			return
		}
		v := viper.New()
		v.SetConfigFile(resolved)
		v.SetConfigType("yaml")
		if err := v.ReadInConfig(); err != nil {
			feedErr = fmt.Errorf("config: read %s failed: %w", resolved, err)
			return
		}
		feedV = v
	})
	if feedErr != nil {
		return nil, feedErr
	}
	if feedV == nil {
		return nil, fmt.Errorf("config: feed viper is nil")
	}
	return feedV, nil
}

func findConfigUpward(relPath string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("config: getwd failed: %w", err)
	}

	dir := wd
	for {
		candidate := filepath.Join(dir, relPath)
		if st, e := os.Stat(candidate); e == nil && !st.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("config: read %s failed: open %s: no such file or directory", relPath, relPath)
}

// AbsPathFromConfigFile 将相对路径转为“相对配置文件目录”的绝对路径
func AbsPathFromConfigFile(v *viper.Viper, p string) string {
	if p == "" {
		return p
	}
	if filepath.IsAbs(p) {
		return p
	}
	if v == nil {
		return p
	}
	cfgFile := v.ConfigFileUsed()
	if cfgFile == "" {
		return p
	}
	return filepath.Join(filepath.Dir(cfgFile), p)
}
