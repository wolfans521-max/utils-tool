package kwfilter

import (
	"testing"
)

// TestFilterText 文本过滤
func TestFilterText(t *testing.T) {

	tests := []struct {
		text string
	}{
		{
			text: "这是一段测试文本",
		},
		{
			text: "天安门运动和法轮功事件的影响",
		},
		{
			text: "创造营亚洲2全部人背影照",
		},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			statuses, ruleStr := FilterText(tt.text, "")

			t.Logf("statuses: %d, ruleStr: %v", statuses, ruleStr)
		})
	}
}
