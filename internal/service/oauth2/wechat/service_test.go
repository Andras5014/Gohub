//go:build manual

package wechat

import (
	"context"
	"github.com/Andras5014/webook/ioc"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestOAuth2WeChatService_e2e_VerifyCode(t *testing.T) {
	logger := ioc.InitLogger()
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID not set")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("WECHAT_APP_KEY not set")
	}
	svc := NewOAuth2WeChatService(appId, appKey, logger)
	res, err := svc.VerifyCode(context.Background(), "")
	require.NoError(t, err)
	t.Log(res)
}
