package redsyncx

import (
	"context"
	"errors"
	"github.com/go-redsync/redsync/v4"
	"sync"
	"time"
)

var ErrExtendFailed = errors.New("extend lock failed")

type AutoExtendMutex struct {
	unlock chan struct{}
	*redsync.Mutex
	once  sync.Once
	until time.Time
}

func NewAutoExtendMutex(mu *redsync.Mutex) *AutoExtendMutex {
	return &AutoExtendMutex{
		Mutex:  mu,
		unlock: make(chan struct{}),
	}
}

func (m *AutoExtendMutex) stop() {
	m.once.Do(func() {
		close(m.unlock)
	})
}

func (m *AutoExtendMutex) Unlock() (bool, error) {
	m.stop()
	return m.Mutex.Unlock()
}

func (m *AutoExtendMutex) UnlockContext(ctx context.Context) (bool, error) {
	m.stop()
	return m.Mutex.UnlockContext(ctx)
}
func (m *AutoExtendMutex) extendWithTimeout(timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ok, err := m.ExtendContext(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		// ctx 超时可以尝试重试
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !ok {
		return false, ErrExtendFailed
	}
	return true, nil
}

func (m *AutoExtendMutex) AutoExtend(interval time.Duration, timeout time.Duration) error {
	if m.until.Before(time.Now()) {
		return ErrExtendFailed
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop() // 确保在函数退出时停止定时器

	for {
		select {
		case <-ticker.C:
			_, err := m.extendWithTimeout(timeout)
			if err != nil {
				return err
			}
		case <-m.unlock:
			// 主动释放了锁, 退出自动续期
			return nil
		}
	}
}
