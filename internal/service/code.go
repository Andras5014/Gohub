package service

import (
	"context"
	"fmt"
	"github.com/Andras5014/webook/internal/repository"
	"github.com/Andras5014/webook/internal/service/sms"
	"math/rand"
)

const codeTplId = "SMS_468975799"

var ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, code string, number string) (bool, error)
}
type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}
func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	// 生成验证码
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTplId, []sms.NamedArg{
		{
			Name:  "code",
			Value: code,
		},
	}, phone)

	return err
}

func (svc *codeService) Verify(ctx context.Context, biz string, code string, number string) (bool, error) {
	return svc.repo.Verify(ctx, biz, number, code)
}

func (svc *codeService) generateCode() string {

	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
