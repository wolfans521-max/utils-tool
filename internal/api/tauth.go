// Package api TAuth认证相关功能
package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"golang.org/x/sync/singleflight"
)

// TAuthToken TAuth token结构
type TAuthToken struct {
	Token       string `json:"tauth_token"`
	TokenSecret string `json:"tauth_token_secret"`
}

// tokenCacheEntry Token 缓存条目
type tokenCacheEntry struct {
	token     *TAuthToken
	fetchTime time.Time
}

// TokenCacheStats Token 缓存监控指标
type TokenCacheStats struct {
	Hits         int64                 // 缓存命中次数
	Misses       int64                 // 缓存未命中次数（需重新获取）
	Errors       int64                 // 获取 Token 失败次数
	StaleHits    int64                 // 过期缓存兜底命中次数（获取失败时使用过期缓存）
	RefreshCount int64                 // 缓存刷新次数（成功获取新 Token）
	CacheCount   int                   // 当前缓存条目数
	TTL          string                // 当前缓存 TTL
	Entries      []TokenCacheEntryInfo // 各缓存条目详情
}

// TokenCacheEntryInfo 单个缓存条目信息
type TokenCacheEntryInfo struct {
	Key                string    // 缓存 key（tokenFile 路径）
	FetchTime          time.Time // 获取时间
	Age                string    // 缓存年龄
	ExpiresIn          string    // 距过期剩余时间
	IsExpired          bool      // 是否已过期
	TokenPreview       string    // tauth_token 预览（脱敏，仅显示前4+后4字符）
	TokenSecretPreview string    // tauth_token_secret 预览（脱敏，仅显示前4+后4字符）
}

// Token 缓存：避免每次请求都读文件/发 HTTP，预计节省 ~8% CPU
var (
	tokenCacheMap = make(map[string]*tokenCacheEntry)
	tokenCacheMu  sync.RWMutex
	tokenCacheTTL = 5 * time.Minute // 缓存 TTL，Token 文件极少变化

	// 监控计数器（原子操作，无锁）
	tokenCacheHits      atomic.Int64
	tokenCacheMisses    atomic.Int64
	tokenCacheErrors    atomic.Int64
	tokenCacheStaleHits atomic.Int64
	tokenCacheRefreshes atomic.Int64

	// singleflight 防止缓存过期时多个并发请求同时触发 fetchToken（惊群效应）
	tokenFetchGroup singleflight.Group
)

// SetTokenCacheTTL 设置 Token 缓存 TTL（用于测试或动态配置）
func SetTokenCacheTTL(ttl time.Duration) {
	tokenCacheMu.Lock()
	tokenCacheTTL = ttl
	tokenCacheMu.Unlock()
}

// InvalidateTokenCache 使指定 tokenFile 的缓存失效
func InvalidateTokenCache(tokenFile string) {
	tokenCacheMu.Lock()
	delete(tokenCacheMap, tokenFile)
	tokenCacheMu.Unlock()
}

// GetTAuthToken 获取TAuth认证头
// tokenFile: token文件路径或URL
// uid: 用户ID（可选）
func GetTAuthToken(tokenFile, uid string) (string, error) {
	// 获取token信息（带缓存）
	token, err := fetchTokenWithCache(tokenFile)
	if err != nil {
		return "", err
	}

	if token.Token == "" || token.TokenSecret == "" {
		return "", fmt.Errorf("token信息不完整")
	}

	// 如果没有uid，返回简单格式
	if uid == "" {
		return fmt.Sprintf("TAuth2 token=%s", url.QueryEscape(token.Token)), nil
	}

	// 生成带签名的认证头
	param := fmt.Sprintf("uid=%s", uid)
	sign := generateHmacSha1Sign(token.TokenSecret, param)

	authHeader := fmt.Sprintf(
		`TAuth2 token="%s",param="%s",sign="%s"`,
		url.QueryEscape(token.Token),
		url.QueryEscape(param),
		url.QueryEscape(sign),
	)

	return authHeader, nil
}

// fetchTokenWithCache 带缓存的 Token 获取
// 优先从内存缓存读取，缓存未命中或过期时才读文件/发 HTTP
// 预计节省 ~8% CPU（原来每次请求都读文件）
func fetchTokenWithCache(tokenFile string) (*TAuthToken, error) {
	// 快速路径：读锁检查缓存
	tokenCacheMu.RLock()
	entry, ok := tokenCacheMap[tokenFile]
	ttl := tokenCacheTTL
	tokenCacheMu.RUnlock()

	if ok && time.Since(entry.fetchTime) < ttl {
		tokenCacheHits.Add(1)
		return entry.token, nil
	}

	// 慢速路径：缓存未命中或过期，重新获取（使用 singleflight 合并并发请求，防止惊群效应）
	tokenCacheMisses.Add(1)
	result, err, _ := tokenFetchGroup.Do(tokenFile, func() (interface{}, error) {
		return fetchToken(tokenFile)
	})

	if err != nil {
		tokenCacheErrors.Add(1)
		// 获取失败时，如果缓存存在且未超过 2*TTL，仍可使用过期缓存
		if ok && time.Since(entry.fetchTime) < 2*ttl {
			tokenCacheStaleHits.Add(1)
			return entry.token, nil
		}
		return nil, err
	}

	// 类型断言
	token, typeOk := result.(*TAuthToken)
	if !typeOk {
		tokenCacheErrors.Add(1)
		return nil, fmt.Errorf("token 类型断言失败")
	}

	// 更新缓存
	tokenCacheRefreshes.Add(1)
	tokenCacheMu.Lock()
	tokenCacheMap[tokenFile] = &tokenCacheEntry{
		token:     token,
		fetchTime: time.Now(),
	}
	tokenCacheMu.Unlock()

	return token, nil
}

// GetTokenCacheStats 获取 Token 缓存监控指标
func GetTokenCacheStats() TokenCacheStats {
	tokenCacheMu.RLock()
	ttl := tokenCacheTTL
	entries := make([]TokenCacheEntryInfo, 0, len(tokenCacheMap))
	for key, entry := range tokenCacheMap {
		age := time.Since(entry.fetchTime)
		expiresIn := ttl - age
		if expiresIn < 0 {
			expiresIn = 0
		}
		entries = append(entries, TokenCacheEntryInfo{
			Key:                key,
			FetchTime:          entry.fetchTime,
			Age:                age.Truncate(time.Millisecond).String(),
			ExpiresIn:          expiresIn.Truncate(time.Millisecond).String(),
			IsExpired:          age >= ttl,
			TokenPreview:       maskSecret(entry.token.Token),
			TokenSecretPreview: maskSecret(entry.token.TokenSecret),
		})
	}
	count := len(tokenCacheMap)
	tokenCacheMu.RUnlock()

	return TokenCacheStats{
		Hits:         tokenCacheHits.Load(),
		Misses:       tokenCacheMisses.Load(),
		Errors:       tokenCacheErrors.Load(),
		StaleHits:    tokenCacheStaleHits.Load(),
		RefreshCount: tokenCacheRefreshes.Load(),
		CacheCount:   count,
		TTL:          ttl.String(),
		Entries:      entries,
	}
}

// fetchToken 从文件或URL获取token（无缓存，直接读取）
func fetchToken(tokenFile string) (*TAuthToken, error) {
	// 判断是URL还是文件路径
	if isURL(tokenFile) {
		return fetchTokenFromURL(tokenFile)
	}
	return fetchTokenFromFile(tokenFile)
}

// isURL 判断是否为URL
func isURL(path string) bool {
	return len(path) > 7 && (path[:7] == "http://" || path[:8] == "https://")
}

// fetchTokenFromURL 从URL获取token
func fetchTokenFromURL(tokenURL string) (*TAuthToken, error) {
	resp, err := http.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("请求token失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求token失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取token响应失败: %w", err)
	}

	var token TAuthToken
	if err := sonic.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("解析token失败: %w", err)
	}

	return &token, nil
}

// fetchTokenFromFile 从文件获取token
func fetchTokenFromFile(tokenFile string) (*TAuthToken, error) {
	// 读取本地文件
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("读取token文件失败: %w", err)
	}

	// 解析JSON
	// 解析JSON（使用 sonic 替代 encoding/json）
	var token TAuthToken
	if err := sonic.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("解析token文件失败: %w", err)
	}
	return &token, nil
}

// generateHmacSha1Sign 生成HMAC-SHA1签名
func generateHmacSha1Sign(secret, data string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// BuildTAuthHeader 构建TAuth认证头（便捷方法）
func BuildTAuthHeader(token, tokenSecret, uid string) string {
	if uid == "" {
		return fmt.Sprintf("TAuth2 token=%s", url.QueryEscape(token))
	}

	param := fmt.Sprintf("uid=%s", uid)
	sign := generateHmacSha1Sign(tokenSecret, param)

	return fmt.Sprintf(
		`TAuth2 token="%s",param="%s",sign="%s"`,
		url.QueryEscape(token),
		url.QueryEscape(param),
		url.QueryEscape(sign),
	)
}

// maskSecret 对敏感字符串进行脱敏处理
// 保留前4位和后4位，中间用 **** 替代
// 长度 <= 8 时仅显示前2+后2，中间用 **** 替代
// 空字符串直接返回
func maskSecret(s string) string {
	if s == "" {
		return ""
	}
	n := len(s)
	if n <= 8 {
		if n <= 4 {
			return "****"
		}
		return s[:2] + "****" + s[n-2:]
	}
	return s[:4] + "****" + s[n-4:]
}
