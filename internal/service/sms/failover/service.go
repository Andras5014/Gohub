package failover

import (
	"context"
	"errors"
	"github.com/Andras5014/gohub/internal/service/sms"
	"sync/atomic"
)

type FailoverSmsService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSmsService(svcs []sms.Service) sms.Service {
	return &FailoverSmsService{
		svcs: svcs,
	}
}
func (f *FailoverSmsService) Send(ctx context.Context, tplToken string, args []sms.NamedArg, numbers ...string) error {

	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplToken, args, numbers...)
		// 发送成功
		if err == nil {
			return nil
		}
	}
	return errors.New("所有服务商都发送失败")
}

func (f *FailoverSmsService) SendV1(ctx context.Context, tplToken string, args []sms.NamedArg, numbers ...string) error {
	// 去下一个节点作为起始节点
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplToken, args, numbers...)
		// 发送成功
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		}
	}
	return errors.New("所有服务商都发送失败")
}
