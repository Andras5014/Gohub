package startup

import (
	"github.com/Andras5014/gohub/internal/service/sms"
	"github.com/Andras5014/gohub/internal/service/sms/aliyun"
	"github.com/Andras5014/gohub/internal/service/sms/memory"
	"github.com/Andras5014/gohub/internal/service/sms/tencent"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSmsService() sms.Service {
	return initSmsMemoryService()
}

func initSmsTencentService() sms.Service {
	secretId, ok := os.LookupEnv("TENCENT_SMS_SECRET_ID")
	if !ok {
		panic("没有找到腾讯云短信的SecretId")
	}
	secretKey, ok := os.LookupEnv("TENCENT_SMS_SECRET_KEY")
	c, err := tencentSms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-chengdu",
		profile.NewClientProfile())
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "123456", "gohub")
}

func initSmsAliService() sms.Service {
	const (
		endpoint = "dysmsapi.aliyuncs.com"
	)
	accessKeyId, ok := os.LookupEnv("ALIYUN_SMS_SECRET_ID")
	if !ok {
		panic("没有找到阿里云短信的SecretId")
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
	return aliyun.NewService(client, "gohub")
}
func initSmsMemoryService() sms.Service {
	return memory.NewService()
}
