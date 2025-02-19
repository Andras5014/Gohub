package oauth2

import (
	"context"
	"github.com/Andras5014/gohub/internal/domain"
)

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WeChatInfo, error)
}
