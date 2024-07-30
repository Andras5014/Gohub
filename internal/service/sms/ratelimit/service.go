package ratelimit

import (
	"context"
	"fmt"
	"github.com/Andras5014/webook/internal/service/sms"
	"github.com/Andras5014/webook/pkg/ratelimit"
)

type LimitSmsService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewLimitSmsService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &LimitSmsService{
		svc:     svc,
		limiter: limiter,
	}
}
func (s *LimitSmsService) Send(ctx context.Context, tpl string, args []sms.NamedArg, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:aliyun")
	if err != nil {
		return fmt.Errorf("限流出现问题: %w", err)
	}
	if limited {
		return fmt.Errorf("触发限流")
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	return err
}
