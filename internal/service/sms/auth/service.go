package auth

import (
	"context"
	"github.com/Andras5014/gohub/internal/service/sms"
	"github.com/golang-jwt/jwt/v5"
)

type SmsService struct {
	svc sms.Service
	key string
}

// Send tplToken 是线下申请的一个业务方的token
func (s *SmsService) Send(ctx context.Context, tplToken string, args []sms.NamedArg, numbers ...string) error {
	var tc TokenClaims
	claims, err := jwt.ParseWithClaims(tplToken, &tc, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !claims.Valid {
		return err
	}
	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

type TokenClaims struct {
	jwt.RegisteredClaims
	Tpl string
}
