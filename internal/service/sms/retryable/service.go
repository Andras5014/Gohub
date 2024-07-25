package retryable

import (
	"context"
	"github.com/Andras5014/webook/internal/service/sms"
)

type Service struct {
	svc        sms.Service
	retryCount int
}

func (s *Service) Send(ctx context.Context, tpl string, args []sms.NamedArg, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)
	for err != nil && s.retryCount < 10 {
		err = s.svc.Send(ctx, tpl, args, numbers...)
		s.retryCount++
	}
	return err
}
