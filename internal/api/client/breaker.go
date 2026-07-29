// Package client 提供下游服务熔断器实现
//
// 熔断器模式：当下游服务连续失败时，自动切断请求（快速失败），
// 避免大量 goroutine 阻塞在超时等待上，防止雪崩效应。
//
// 三种状态：
//   - Closed（关闭）：正常放行请求，统计连续失败次数
//   - Open（打开）：快速失败，不发送实际请求
//   - HalfOpen（半开）：放行少量请求探测下游是否恢复
package client

import (
	"fmt"
	"sync"
	"time"
)

// BreakerState 熔断器状态
type BreakerState int32

const (
	BreakerClosed   BreakerState = iota // 关闭：正常放行
	BreakerOpen                         // 打开：快速失败
	BreakerHalfOpen                     // 半开：探测恢复
)

// String 实现 Stringer 接口
func (s BreakerState) String() string {
	switch s {
	case BreakerClosed:
		return "Closed"
	case BreakerOpen:
		return "Open"
	case BreakerHalfOpen:
		return "HalfOpen"
	default:
		return "Unknown"
	}
}

// BreakerConfig 熔断器配置
type BreakerConfig struct {
	// FailureThreshold 连续失败多少次后打开熔断器，默认 5
	FailureThreshold int `yaml:"failure_threshold" mapstructure:"failure_threshold"`
	// OpenTimeout 熔断器打开后等待多久进入半开状态，默认 10s
	OpenTimeout time.Duration `yaml:"open_timeout" mapstructure:"open_timeout"`
	// HalfOpenMaxProbes 半开状态下最多放行多少个探测请求，默认 1
	HalfOpenMaxProbes int `yaml:"half_open_max_probes" mapstructure:"half_open_max_probes"`
	// SuccessThreshold 半开状态下连续成功多少次后关闭熔断器，默认 2
	SuccessThreshold int `yaml:"success_threshold" mapstructure:"success_threshold"`
}

// DefaultBreakerConfig 默认熔断器配置
// 注意：FailureThreshold 不宜过低，高 QPS 下短暂网络抖动就可能触发连续 5 次失败
// 导致误熔断。建议根据业务 QPS 调整：
//   - 低 QPS (<100): 5-10
//   - 中 QPS (100-1000): 10-30
//   - 高 QPS (>1000): 30-50
func DefaultBreakerConfig() BreakerConfig {
	return BreakerConfig{
		FailureThreshold:  30,             // 连续失败 30 次才熔断（高 QPS 下避免误触发）
		OpenTimeout:       5 * time.Second, // 熔断 5 秒后进入半开（快速探测恢复）
		HalfOpenMaxProbes: 3,              // 半开状态放行 3 个探测请求（提高恢复概率）
		SuccessThreshold:  2,              // 半开状态连续成功 2 次后关闭
	}
}

// CircuitBreaker 熔断器（每个下游服务一个实例）
type CircuitBreaker struct {
	config BreakerConfig
	name   string // 下游服务标识（urlKey）

	mu               sync.Mutex
	state            BreakerState
	consecutiveFails int
	consecutiveOKs   int
	openedAt         time.Time // 进入 Open 状态的时间
	halfOpenProbes   int       // 半开状态下已放行的探测数
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(name string, config BreakerConfig) *CircuitBreaker {
	if config.FailureThreshold <= 0 {
		config.FailureThreshold = 5
	}
	if config.OpenTimeout <= 0 {
		config.OpenTimeout = 10 * time.Second
	}
	if config.HalfOpenMaxProbes <= 0 {
		config.HalfOpenMaxProbes = 1
	}
	if config.SuccessThreshold <= 0 {
		config.SuccessThreshold = 2
	}
	return &CircuitBreaker{
		name:   name,
		config: config,
		state:  BreakerClosed,
	}
}

// Allow 检查是否允许请求通过
// 返回 error 表示熔断器打开，应快速失败
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case BreakerClosed:
		return nil

	case BreakerOpen:
		// 检查是否到了半开时间
		if time.Since(cb.openedAt) >= cb.config.OpenTimeout {
			cb.state = BreakerHalfOpen
			cb.halfOpenProbes = 0
			cb.consecutiveOKs = 0
			// 继续到 HalfOpen 逻辑
		} else {
			return fmt.Errorf("circuit breaker [%s] is open, fast fail (retry after %s)",
				cb.name, cb.config.OpenTimeout-time.Since(cb.openedAt))
		}
		fallthrough

	case BreakerHalfOpen:
		// 半开状态下限制探测请求数
		if cb.halfOpenProbes >= cb.config.HalfOpenMaxProbes {
			return fmt.Errorf("circuit breaker [%s] is half-open, probe limit reached", cb.name)
		}
		cb.halfOpenProbes++
		return nil

	default:
		return nil
	}
}

// RecordSuccess 记录请求成功
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.consecutiveFails = 0

	switch cb.state {
	case BreakerHalfOpen:
		cb.consecutiveOKs++
		if cb.consecutiveOKs >= cb.config.SuccessThreshold {
			cb.state = BreakerClosed
			cb.consecutiveOKs = 0
			cb.halfOpenProbes = 0
		}
	case BreakerClosed:
		// 正常状态，无需额外操作
	}
}

// RecordFailure 记录请求失败
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.consecutiveOKs = 0

	switch cb.state {
	case BreakerClosed:
		cb.consecutiveFails++
		if cb.consecutiveFails >= cb.config.FailureThreshold {
			cb.state = BreakerOpen
			cb.openedAt = time.Now()
		}
	case BreakerHalfOpen:
		// 半开状态下失败，立即回到 Open
		cb.state = BreakerOpen
		cb.openedAt = time.Now()
		cb.halfOpenProbes = 0
	}
}

// State 获取当前状态
func (cb *CircuitBreaker) State() BreakerState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Stats 获取熔断器统计信息（用于监控）
func (cb *CircuitBreaker) Stats() (state BreakerState, consecutiveFails int, consecutiveOKs int) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state, cb.consecutiveFails, cb.consecutiveOKs
}

// Reset 重置熔断器为 Closed 状态
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = BreakerClosed
	cb.consecutiveFails = 0
	cb.consecutiveOKs = 0
	cb.halfOpenProbes = 0
}

// BreakerGroup 熔断器组（管理多个下游服务的熔断器）
type BreakerGroup struct {
	mu       sync.RWMutex
	breakers map[string]*CircuitBreaker
	config   BreakerConfig
}

// NewBreakerGroup 创建熔断器组
func NewBreakerGroup(config BreakerConfig) *BreakerGroup {
	return &BreakerGroup{
		breakers: make(map[string]*CircuitBreaker),
		config:   config,
	}
}

// Get 获取指定服务的熔断器（懒创建）
func (g *BreakerGroup) Get(name string) *CircuitBreaker {
	g.mu.RLock()
	cb, ok := g.breakers[name]
	g.mu.RUnlock()

	if ok {
		return cb
	}

	// 懒创建
	g.mu.Lock()
	defer g.mu.Unlock()

	// double-check
	if cb, ok = g.breakers[name]; ok {
		return cb
	}

	cb = NewCircuitBreaker(name, g.config)
	g.breakers[name] = cb
	return cb
}

// AllStats 获取所有熔断器状态（用于监控接口）
func (g *BreakerGroup) AllStats() map[string]BreakerStats {
	g.mu.RLock()
	defer g.mu.RUnlock()

	result := make(map[string]BreakerStats, len(g.breakers))
	for name, cb := range g.breakers {
		state, fails, oks := cb.Stats()
		result[name] = BreakerStats{
			Name:             name,
			State:            state.String(),
			ConsecutiveFails: fails,
			ConsecutiveOKs:   oks,
		}
	}
	return result
}

// BreakerStats 熔断器统计信息
type BreakerStats struct {
	Name             string `json:"name"`
	State            string `json:"state"`
	ConsecutiveFails int    `json:"consecutive_fails"`
	ConsecutiveOKs   int    `json:"consecutive_oks"`
}
