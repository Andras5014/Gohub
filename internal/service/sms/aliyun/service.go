package aliyun

import (
	"context"
	"encoding/json"
	"github.com/Andras5014/gohub/internal/service/sms"
	dysmapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/ecodeclub/ekit/slice"
	"strings"
)

type Service struct {
	signName string
	client   *dysmsapi.Client
}

func NewService(client *dysmsapi.Client, signName string) sms.Service {
	return &Service{
		client:   client,
		signName: signName,
	}
}
func (s *Service) Send(ctx context.Context, tplToken string, args []sms.NamedArg, numbers ...string) error {

	request := &dysmapi.SendSmsRequest{}
	request.SetPhoneNumbers(s.formatPhoneNumbers(numbers))
	request.SetSignName(s.signName)
	request.SetTemplateCode(tplToken)
	paramMap := make(map[string]string, len(args))
	for _, arg := range args {
		paramMap[arg.Name] = arg.Value
	}
	paramJson, err := json.Marshal(paramMap)
	if err != nil {
		return err
	}
	request.SetTemplateParam(string(paramJson))
	_, err = s.client.SendSms(request)
	return err
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}

func (s *Service) formatPhoneNumbers(numbers []string) string {
	sb := strings.Builder{}
	n := len(numbers)
	for i, number := range numbers {
		sb.WriteString(number)
		if i != n-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
