package tencent

import (
	"context"
	"fmt"
	sms "github.com/Andras5014/gohub/internal/service/sms"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	tencentSms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *tencentSms.Client
}

func NewService(client *tencentSms.Client, appId string, signName string) sms.Service {
	return &Service{
		client:   client,
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
	}
}

func (s *Service) Send(ctx context.Context, tplToken string, args []sms.NamedArg, numbers ...string) error {
	templateParamSet := make([]*string, len(args))
	for _, arg := range args {
		templateParamSet = append(templateParamSet, ekit.ToPtr[string](arg.Value))
	}

	req := tencentSms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tplToken)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = templateParamSet
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}

	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败, %s,%s", *(status.Code), *(status.Message))
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
