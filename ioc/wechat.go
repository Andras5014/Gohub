package ioc

import (
	"github.com/Andras5014/webook/internal/service/oauth2"
	"github.com/Andras5014/webook/internal/service/oauth2/wechat"
	"github.com/Andras5014/webook/pkg/logx"
)

func InitOAuth2WeChatService(l logx.Logger) oauth2.Service {
	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("WECHAT_APP_ID not set")
	//}
	//appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("WECHAT_APP_KEY not set")
	//}
	appId, appKey := "wx0f0b0f0f0f0f0f0f", "0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f"
	return wechat.NewOAuth2WeChatService(appId, appKey, l)
}
