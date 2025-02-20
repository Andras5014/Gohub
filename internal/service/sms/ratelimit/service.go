package ratelimit

import (
	"context"
	"fmt"
	"github.com/Andras5014/gohub/internal/service/sms"
	"github.com/Andras5014/gohub/pkg/ratelimit"
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
func (s *LimitSmsService) Send(ctx context.Context, tplToken string, args []sms.NamedArg, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:aliyun")
	if err != nil {
		return fmt.Errorf("限流出现问题: %w", err)
	}
	if limited {
		return fmt.Errorf("触发限流")
	}
	err = s.svc.Send(ctx, tplToken, args, numbers...)
	return err
}
