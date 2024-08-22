package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/pkg/logger"
	"net/http"
	"net/url"
)

const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=SCOPE&state=%s#wechat_redirect"
const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"

var redirectURI = url.PathEscape("http://localhost:8080/oauth2/wechat/callback")

type OAuth2WeChatService struct {
	appId     string
	appSecret string
	client    *http.Client
	logger    logger.Logger
}

func NewOAuth2WeChatService(appId string, appSecret string, logger logger.Logger) *OAuth2WeChatService {
	return &OAuth2WeChatService{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
		logger:    logger,
	}
}
func (o *OAuth2WeChatService) AuthURL(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf(urlPattern, o.appId, redirectURI, state), nil
}

func (o *OAuth2WeChatService) VerifyCode(ctx context.Context, code string) (domain.WeChatInfo, error) {
	target := fmt.Sprintf(targetPattern, o.appId, o.appSecret, code)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WeChatInfo{}, err
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return domain.WeChatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)

	if err != nil {
		return domain.WeChatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WeChatInfo{}, fmt.Errorf("微信返回错误码：%d,错误信息：%s", res.ErrCode, res.ErrMsg)
	}

	o.logger.Info("微信返回用户信息", logger.Any("openId", res.OpenId), logger.Any("unionId", res.UnionId))
	return domain.WeChatInfo{
		OpenId:  res.OpenId,
		UnionId: res.UnionId,
	}, nil

}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
}
