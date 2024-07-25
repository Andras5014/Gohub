package aliyun

import (
	"github.com/Andras5014/webook/internal/service/sms"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"os"
	"testing"
)

const (
	endpoint   = "dysmsapi.aliyuncs.com"
	templateId = "SMS_468975799"
)

func TestService_Send(t *testing.T) {
	accessKeyId, ok := os.LookupEnv("ALIYUN_SMS_SECRET_ID")
	if !ok {
		t.Skip("TENCENT_SMS_SECRET_ID not set")
	}
	accessKeySecret, ok := os.LookupEnv("ALIYUN_SMS_SECRET_KEY")
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}
	var err error
	client := &dysmsapi.Client{}
	client, err = dysmsapi.NewClient(config)
	if err != nil {
		panic(err)
	}
	s := NewService(client, "webook")

	testCases := []struct {
		name    string
		tplId   string
		params  []sms.NamedArg
		numbers []string
		wantErr error
	}{
		{
			name:  "test",
			tplId: templateId,
			params: []sms.NamedArg{
				{
					Name:  "code",
					Value: "123456",
				},
			},
			numbers: []string{"17383616844"},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		err := s.Send(nil, tc.tplId, tc.params, tc.numbers...)
		if err != nil {
			t.Errorf("TestService_Send failed, name=%v, err=%v", tc.name, err)
		}
	}
}
