package kwfilter

import (
	"encoding/json"
	"os"
	"testing"
)

// TestFilter 敏感词过滤测试
// 注意: 此测试需要有效的词典文件才能正常运行
func TestFilter(t *testing.T) {
	// 检查默认词典文件是否存在
	defaultDictPath := DefaultDictPath
	if _, err := os.Stat(defaultDictPath); os.IsNotExist(err) {
		t.Skipf("Default dictionary file not found at %s, skipping test. Please provide a valid dictionary file.", defaultDictPath)
		return
	}

	kf := New()
	defer kf.Close()

	// 3. 准备属性配置
	attr := map[string]interface{}{
		"log_level":        3,
		"filter_level":     6203,
		"has_user":         true,
		"user_type":        1,
		"has_resrc":        true,
		"resrc_type":       1,
		"has_specialstats": true,
		"special_stats":    1,
	}
	attrJSON, _ := json.Marshal(attr)

	// 4. 执行过滤测试
	tests := []struct {
		text string
	}{
		{
			text: "这是一段测试文本",
		},
		{
			text: "这是一段正常的文本内容，啦啦啦",
		},
		{
			text: "天安门运动事件具体是什么",
		},
		{
			text: "法轮功组织的详细架构和成员情况",
		},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result, matchLevel, err := kf.FilterExec(tt.text, string(attrJSON))

			if err != nil {
				t.Errorf("Filter() error = %v", err)
			}

			// A-F级
			sensitive := 0
			if matchLevel&0x3b > 0 {
				sensitive = 1
			}

			t.Logf("MatchLevel: %d, Sensitive: %d, Result: %s", matchLevel, sensitive, result)
		})
	}
}
