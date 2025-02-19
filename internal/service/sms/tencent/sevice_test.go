package tencent

import (
	"errors"
	mysms "github.com/Andras5014/gohub/internal/service/sms"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"testing"
)

func TestService_Send(t *testing.T) {
	secretId, ok := os.LookupEnv("TENCENT_SMS_SECRET_ID")
	if !ok {
		t.Skip("TENCENT_SMS_SECRET_ID not set")
	}
	secretKey, ok := os.LookupEnv("TENCENT_SMS_SECRET_KEY")
	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-chengdu",
		profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
	}
	s := NewService(c, "1400603310", "腾讯云短信测试")

	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:    "test",
			tplId:   "1400603310",
			params:  []string{"123456"},
			numbers: []string{"+86123456789"},
			wantErr: nil,
		},
		{
			name:    "test",
			tplId:   "1400603310",
			params:  []string{"123456"},
			numbers: []string{"+86123456789"},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		args := make([]mysms.NamedArg, len(tc.params))
		for _, param := range tc.params {
			args = append(args, mysms.NamedArg{
				Name:  "code",
				Value: param,
			})
		}
		t.Run(tc.name, func(t *testing.T) {
			err := s.Send(nil, tc.tplId, args, tc.numbers...)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("want %v, got %v", tc.wantErr, err)
			}
		})
	}
}
