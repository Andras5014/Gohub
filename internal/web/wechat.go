package web

import (
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/service/oauth2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WeChatHandler struct {
	svc     oauth2.Service
	userSvc service.UserService
	JwtHandler
}

func NewOAuth2WeChatHandler(svc oauth2.Service, userSvc service.UserService) *OAuth2WeChatHandler {
	return &OAuth2WeChatHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}

func (h *OAuth2WeChatHandler) RegisterRoutes(engine *gin.Engine) {
	g := engine.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.GET("/callback", h.Callback)

}

func (h *OAuth2WeChatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "微信登录失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Data: url,
	})
}

func (h *OAuth2WeChatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
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
	err = h.setJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Data: gin.H{
			"redirect": state,
		},
	})

}
