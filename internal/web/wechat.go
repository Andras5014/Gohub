package web

import (
	"errors"
	"fmt"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/service/oauth2"
	ijwt "github.com/Andras5014/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type OAuth2WeChatHandler struct {
	svc      oauth2.Service
	userSvc  service.UserService
	stateKey []byte
	ijwt.Handler
}

func NewOAuth2WeChatHandler(svc oauth2.Service, userSvc service.UserService, jwtHdl ijwt.Handler) *OAuth2WeChatHandler {
	return &OAuth2WeChatHandler{
		svc:      svc,
		userSvc:  userSvc,
		stateKey: []byte("secret"),
		Handler:  jwtHdl,
	}
}

func (h *OAuth2WeChatHandler) RegisterRoutes(engine *gin.Engine) {
	g := engine.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.GET("/callback", h.Callback)

}

func (h *OAuth2WeChatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New().String()
	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登陆url失败",
		})
		return
	}
	if err = h.setStateCookie(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Data: url,
	})
}

func (h *OAuth2WeChatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return err
	}
	ctx.SetCookie("jwt-state", tokenStr, 60, "/oauth2/wechat/callback", "", false, true)
	return nil
}

func (h *OAuth2WeChatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {

	}
	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "登陆成功",
	})
}

func (h *OAuth2WeChatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")

	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到cookie,%w", err)
	}

	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token过期,%w", err)
	}

	if sc.State != state {
		return errors.New("state不匹配")
	}
	return nil
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
