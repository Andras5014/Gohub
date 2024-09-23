package main

import (
	"fmt"
	"time"
)

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	ch   chan struct{} // 通道，用于控制并发量
	size int           // 滑动窗口的大小
}

// NewSlidingWindowLimiter 创建一个新的滑动窗口限流器
func NewSlidingWindowLimiter(size int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		ch:   make(chan struct{}, size),
		size: size,
	}
}

// Allow 尝试获取一个许可，如果获取不到则阻塞
func (l *SlidingWindowLimiter) Allow() bool {
	select {
	case l.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release 释放一个许可
func (l *SlidingWindowLimiter) Release() {
	<-l.ch
}

func main() {
	// 创建一个窗口大小为3的限流器
	limiter := NewSlidingWindowLimiter(3)

	for i := 0; i < 5; i++ {
		go func(i int) {
			if limiter.Allow() {
				fmt.Printf("Goroutine %d is running\n", i)
				// 模拟任务执行时间
				time.Sleep(2 * time.Second)
				limiter.Release()
				fmt.Printf("Goroutine %d is done\n", i)
			} else {
				fmt.Printf("Goroutine %d is blocked\n", i)
			}
		}(i)
	}

	// 给主goroutine足够的时间来启动所有goroutine
	time.Sleep(10 * time.Second)
}
