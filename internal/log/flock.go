package log

import (
	"os"
	"syscall"
	"time"
)

// Flock 对文件加锁
// exclusive: true 表示独占锁(LOCK_EX)，false 表示共享锁(LOCK_SH)
// blocking: true 表示阻塞等待，false 表示非阻塞
func Flock(file *os.File, exclusive bool, blocking bool) error {
	var how int
	if exclusive {
		how = syscall.LOCK_EX
	} else {
		how = syscall.LOCK_SH
	}
	if !blocking {
		how |= syscall.LOCK_NB
	}
	return syscall.Flock(int(file.Fd()), how)
}

// Funlock 释放文件锁
func Funlock(file *os.File) error {
	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}

// FlockWithTimeout 尝试加锁，带超时机制
// 阻塞等待锁，但最多等待 timeout 时长
func FlockWithTimeout(file *os.File, exclusive bool, timeout time.Duration) error {
	// 使用非阻塞锁 + 轮询实现超时机制
	start := time.Now()
	for {
		err := Flock(file, exclusive, false) // 非阻塞尝试
		if err == nil {
			return nil // 获取锁成功
		}
		// 检查是否超时
		if time.Since(start) >= timeout {
			return err // 超时返回错误
		}
		// 短暂等待后重试
		time.Sleep(10 * time.Millisecond)
	}
}

// FlockWithRetry 尝试加锁，带重试机制（已废弃，请使用 FlockWithTimeout）
func FlockWithRetry(file *os.File, exclusive bool, maxRetries int, delay time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		err := Flock(file, exclusive, false) // 非阻塞
		if err == nil {
			return nil
		}
		if i < maxRetries-1 {
			time.Sleep(delay)
		}
	}
	return syscall.EWOULDBLOCK
}
