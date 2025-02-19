package retryable

import (
	"context"
	"errors"
	"github.com/Andras5014/gohub/internal/service/sms"
)

type Service struct {
	svc      sms.Service
	retryMax int
}

func (s *Service) Send(ctx context.Context, tplToken string, args []sms.NamedArg, numbers ...string) error {
	err := s.svc.Send(ctx, tplToken, args, numbers...)
	cnt := 1
	for err != nil && cnt < s.retryMax {
		err = s.svc.Send(ctx, tplToken, args, numbers...)
		//if err == nil {
		//	return nil
		//}
		cnt++
	}
	return errors.New("重试发送短信失败")
}
