package context

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// TestGetUserId 测试获取用户ID
func TestGetUserId(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		headers        map[string]string
		expectedUserId string
	}{
		{
			name: "从cuid参数获取用户ID",
			queryParams: map[string]string{
				"cuid": "1234567890",
			},
			headers:        map[string]string{},
			expectedUserId: "1234567890",
		},
		{
			name:        "从x-login-uid请求头获取用户ID",
			queryParams: map[string]string{},
			headers: map[string]string{
				"x-login-uid": "9876543210",
			},
			expectedUserId: "9876543210",
		},
		{
			name: "从login_uid参数获取用户ID",
			queryParams: map[string]string{
				"login_uid": "5555555555",
			},
			headers:        map[string]string{},
			expectedUserId: "5555555555",
		},
		{
			name: "优先级：cuid > x-login-uid > login_uid",
			queryParams: map[string]string{
				"cuid":      "1111111111",
				"login_uid": "2222222222",
			},
			headers: map[string]string{
				"x-login-uid": "3333333333",
			},
			expectedUserId: "1111111111",
		},
		{
			name: "优先级：x-login-uid > login_uid (无cuid)",
			queryParams: map[string]string{
				"login_uid": "2222222222",
			},
			headers: map[string]string{
				"x-login-uid": "3333333333",
			},
			expectedUserId: "3333333333",
		},
		{
			name:           "无用户ID时返回空字符串",
			queryParams:    map[string]string{},
			headers:        map[string]string{},
			expectedUserId: "",
		},
		{
			name: "cuid为空字符串时使用x-login-uid",
			queryParams: map[string]string{
				"cuid": "",
			},
			headers: map[string]string{
				"x-login-uid": "4444444444",
			},
			expectedUserId: "4444444444",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试用的gin.Context
			ginCtx := createTestGinContext(tt.queryParams, tt.headers)
			ctx := New(ginCtx)

			// 调用GetUserId方法
			userId := ctx.GetUserId()

			// 验证结果
			if userId != tt.expectedUserId {
				t.Errorf("GetUserId() = %v, 期望 %v", userId, tt.expectedUserId)
			}
		})
	}
}

// TestIsTeenager 测试判断是否为青少年用户
func TestIsTeenager(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		headers        map[string]string
		expectedResult bool
	}{
		{
			name: "从is_teenager参数判断为青少年",
			queryParams: map[string]string{
				"is_teenager": "1",
			},
			headers:        map[string]string{},
			expectedResult: true,
		},
		{
			name:        "从X-Teenager-Flag请求头判断为青少年",
			queryParams: map[string]string{},
			headers: map[string]string{
				"X-Teenager-Flag": "1",
			},
			expectedResult: true,
		},
		{
			name: "优先级：is_teenager参数优先于请求头",
			queryParams: map[string]string{
				"is_teenager": "1",
			},
			headers: map[string]string{
				"X-Teenager-Flag": "0",
			},
			expectedResult: true,
		},
		{
			name: "is_teenager参数任意非空值都返回true",
			queryParams: map[string]string{
				"is_teenager": "yes",
			},
			headers:        map[string]string{},
			expectedResult: true,
		},
		{
			name:        "X-Teenager-Flag请求头任意非空值都返回true",
			queryParams: map[string]string{},
			headers: map[string]string{
				"X-Teenager-Flag": "true",
			},
			expectedResult: true,
		},
		{
			name:           "无青少年标识时返回false",
			queryParams:    map[string]string{},
			headers:        map[string]string{},
			expectedResult: false,
		},
		{
			name: "is_teenager为空字符串时检查请求头",
			queryParams: map[string]string{
				"is_teenager": "",
			},
			headers: map[string]string{
				"X-Teenager-Flag": "1",
			},
			expectedResult: true,
		},
		{
			name: "is_teenager和X-Teenager-Flag都为空时返回false",
			queryParams: map[string]string{
				"is_teenager": "",
			},
			headers: map[string]string{
				"X-Teenager-Flag": "",
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试用的gin.Context
			ginCtx := createTestGinContext(tt.queryParams, tt.headers)
			ctx := New(ginCtx)

			// 调用IsTeenager方法
			result := ctx.IsTeenager()

			// 验证结果
			if result != tt.expectedResult {
				t.Errorf("IsTeenager() = %v, 期望 %v", result, tt.expectedResult)
			}
		})
	}
}

// TestGetUserIdAndIsTeenagerCombined 测试GetUserId和IsTeenager组合使用
func TestGetUserIdAndIsTeenagerCombined(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		headers        map[string]string
		expectedUserId string
		expectedTeen   bool
	}{
		{
			name: "普通成年用户",
			queryParams: map[string]string{
				"cuid": "1234567890",
			},
			headers:        map[string]string{},
			expectedUserId: "1234567890",
			expectedTeen:   false,
		},
		{
			name: "青少年用户",
			queryParams: map[string]string{
				"cuid":        "1234567890",
				"is_teenager": "1",
			},
			headers:        map[string]string{},
			expectedUserId: "1234567890",
			expectedTeen:   true,
		},
		{
			name: "通过请求头识别的青少年用户",
			queryParams: map[string]string{
				"cuid": "1234567890",
			},
			headers: map[string]string{
				"X-Teenager-Flag": "1",
			},
			expectedUserId: "1234567890",
			expectedTeen:   true,
		},
		{
			name: "通过x-login-uid获取ID的青少年用户",
			queryParams: map[string]string{
				"is_teenager": "1",
			},
			headers: map[string]string{
				"x-login-uid": "9876543210",
			},
			expectedUserId: "9876543210",
			expectedTeen:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试用的gin.Context
			ginCtx := createTestGinContext(tt.queryParams, tt.headers)
			ctx := New(ginCtx)

			// 调用方法
			userId := ctx.GetUserId()
			isTeen := ctx.IsTeenager()

			// 验证结果
			if userId != tt.expectedUserId {
				t.Errorf("GetUserId() = %v, 期望 %v", userId, tt.expectedUserId)
			}
			if isTeen != tt.expectedTeen {
				t.Errorf("IsTeenager() = %v, 期望 %v", isTeen, tt.expectedTeen)
			}
		})
	}
}

// createTestGinContext 创建测试用的gin.Context
func createTestGinContext(queryParams, headers map[string]string) *gin.Context {
	// 构建查询参数
	values := url.Values{}
	for k, v := range queryParams {
		values.Set(k, v)
	}

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test?"+values.Encode(), nil)

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 创建测试用的gin.Context
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = req

	return ginCtx
}
