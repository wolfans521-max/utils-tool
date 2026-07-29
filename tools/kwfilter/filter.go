package kwfilter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bytedance/sonic"
)

// FilterParams 敏感词过滤参数
type FilterParams struct {
	LogLevel     int  `json:"log_level"`     // LogLevel 日志级别
	FilterLevel  int  `json:"filter_level"`  // FilterLevel 过滤级别
	HasUser      bool `json:"has_user"`      // HasUser 是否包含用户信息
	UserType     int  `json:"user_type"`     // UserType 用户类型: 1=非V, 2=蓝V, 3=黄V
	UserVType    int  `json:"user_vtype"`    // UserVType 用户V类型
	UserVLevel   int  `json:"user_vlevel"`   // UserVLevel 用户V等级
	UserStatus   int  `json:"user_status"`   // UserStatus 用户状态
	HasResrc     bool `json:"has_resrc"`     // HasResrc 是否包含资源
	ResrcType    int  `json:"resrc_type"`    // ResrcType 资源类型
	SpecialStats int  `json:"special_stats"` // SpecialStats 特殊状态
}

// FilterResult 过滤结果
type FilterResult struct {
	Rule       string
	MatchLevel int
	Error      error
}

// MatchedWord 敏感词匹配结果
type MatchedWord struct {
	Word  string
	Level int
}

// FilterText 对普通文本进行敏感词过滤
// 参数:
//   - text: 待过滤文本
//   - prefix: 可选前缀，用于标识敏感词来源
//
// 返回:
//   - int: 过滤结果，大于0表示命中敏感词
//   - string: 敏感信息记录（格式: 前缀_文本_词\t等级）
func FilterText(text, prefix string) (int, string) {
	result := CallKWordFilter(text, -1, -1, -1, -1, false, 0x1B, 0)

	// 检查是否命中敏感词
	if result.MatchLevel > 0 {
		sensInfo := FormatFilterLog(result)
		if prefix != "" {
			sensInfo = fmt.Sprintf("%s_%s_rule(%s)", prefix, text, sensInfo)
		} else {
			sensInfo = fmt.Sprintf("%s_rule(%s)", text, sensInfo)
		}

		return result.MatchLevel, sensInfo
	}

	return 0, ""
}

// NewDefaultFilterParams 创建默认的过滤参数
// 对应PHP代码中的默认值:
// log_level: 3
// filter_level: 6203
// has_user: true
func NewDefaultFilterParams() FilterParams {
	return FilterParams{
		LogLevel:     3,
		FilterLevel:  6203,
		HasUser:      true,
		UserType:     1,  // 默认非V
		UserVType:    -1, // 默认无V类型
		UserVLevel:   -1, // 默认无V等级
		UserStatus:   1,  // 默认正常用户
		HasResrc:     false,
		ResrcType:    0,
		SpecialStats: 0,
	}
}

// CallKWordFilter 封装过滤方法
func CallKWordFilter(text string, verified int, verifiedType int, userVLevel int, userType int, hasUser bool, digitAttr int, dispatchCtrl int) FilterResult {
	params := NewDefaultFilterParams()

	// 设置是否包含用户信息
	params.HasUser = hasUser

	// 设置用户类型
	// verified: true=认证用户, false=非认证用户
	// verifiedType: 0=黄V, 1-3=蓝V
	if verified > 0 {
		if verifiedType == 0 {
			params.UserType = 3 // 黄V
		} else {
			params.UserType = 2 // 蓝V
		}
		params.UserVType = verifiedType
	} else {
		params.UserType = 1 // 非V
		params.UserVType = -1
	}

	// 设置用户V等级
	params.UserVLevel = userVLevel

	// 设置用户状态
	params.UserStatus = userType

	// 设置资源信息
	if digitAttr != 0 {
		params.HasResrc = true
		params.ResrcType = digitAttr
	}

	// 设置特殊状态
	if dispatchCtrl != 0 {
		params.SpecialStats = dispatchCtrl
	}

	return CallFilterExec(text, params)
}

// CallFilterExec 基础敏感词过滤函数，用于调用 C 库进行文本过滤。
// text: 待过滤文本
// params: 过滤参数结构
// 返回: 过滤结果
func CallFilterExec(text string, params FilterParams) FilterResult {
	// 移除 t.cn 短链接，编译正则表达式，匹配 http://t.cn/ 开头的短链接
	if text != "" {
		re := regexp.MustCompile(`http://t\.cn/\w+`)
		text = re.ReplaceAllString(text, "")
	}

	// 检查文本是否为空（移除后）
	if text == "" {
		return FilterResult{
			Rule:       "",
			MatchLevel: 0,
			Error:      nil,
		}
	}

	// 将参数转换为JSON字符串
	attrJSON, err := sonic.Marshal(params)
	if err != nil {
		return FilterResult{
			Rule:       "",
			MatchLevel: 0,
			Error:      fmt.Errorf("json转换失败: %w", err),
		}
	}

	kwfilter := New()
	defer kwfilter.Close()

	rule, matchLevel, err := kwfilter.FilterExec(text, string(attrJSON))
	if err != nil {
		return FilterResult{
			Rule:       "",
			MatchLevel: 0,
			Error:      err,
		}
	}

	return FilterResult{
		Rule:       rule,
		MatchLevel: matchLevel,
		Error:      nil,
	}
}

// GetMatchedWordsForLog 从规则 JSON 中提取用于日志的敏感词列表
// 过滤规则:
//  1. 排除 level=13 的词
//  2. 只取 Words 中第一个词（|分隔）
//
// 返回: 敏感词和等级的切片
func GetMatchedWordsForLog(ruleJSON string) []MatchedWord {
	if ruleJSON == "" {
		return nil
	}

	var rule struct {
		MatchedInfo []struct {
			Words string `json:"words"`
			Level int    `json:"level"`
		} `json:"matched_info"`
	}

	err := sonic.Unmarshal([]byte(ruleJSON), &rule)
	if err != nil {
		return nil
	}

	var result []MatchedWord
	for _, item := range rule.MatchedInfo {
		// 过滤 level=13
		if item.Level == 13 {
			continue
		}

		// 只取第一个词（|分隔）
		word := item.Words
		if idx := strings.Index(word, "|"); idx != -1 {
			word = word[:idx]
		}

		result = append(result, MatchedWord{
			Word:  word,
			Level: item.Level,
		})
	}

	return result
}

// FormatFilterLog 格式化敏感词过滤日志
// 返回格式: word1:level1,word2:level2
// 如果没有命中的敏感词（level=13 被过滤），返回空字符串
func FormatFilterLog(result FilterResult) string {
	if result.MatchLevel == 0 {
		return ""
	}

	words := GetMatchedWordsForLog(result.Rule)
	if len(words) == 0 {
		return ""
	}

	parts := make([]string, len(words))
	for i, w := range words {
		parts[i] = fmt.Sprintf("%s:%d", w.Word, w.Level)
	}

	return fmt.Sprintf("%s", strings.Join(parts, ","))
}
