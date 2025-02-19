package failover

import (
	"context"
	"github.com/Andras5014/gohub/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSmsService struct {
	svcs []sms.Service
	// 连续超时个数
	cnt int32

	idx int32

	//阈值
	threshold int32
}

func NewTimeoutFailoverSmsService(svcs []sms.Service) *TimeoutFailoverSmsService {
	return &TimeoutFailoverSmsService{
		svcs:      svcs,
		threshold: 3,
		idx:       0,
		cnt:       0,
	}
}

func (f *TimeoutFailoverSmsService) Send(ctx context.Context, tplToken string, args []sms.NamedArg, numbers ...string) error {
	idx := atomic.LoadInt32(&f.idx)
	cnt := atomic.LoadInt32(&f.cnt)

	if cnt > f.threshold {
		newIdx := (idx + 1) % int32(len(f.svcs))
		if atomic.CompareAndSwapInt32(&f.idx, idx, newIdx) {
			// 切换到newIdx
			atomic.StoreInt32(&f.cnt, 0)
		}
		idx = atomic.LoadInt32(&f.idx)
	}

	svc := f.svcs[idx]
	err := svc.Send(ctx, tplToken, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&f.cnt, 1)
	case nil:
		atomic.StoreInt32(&f.cnt, 0)
	default:
		return err
	}
	return nil
}
