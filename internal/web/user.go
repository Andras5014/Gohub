package web

import (
	"errors"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/service"
	ijwt "github.com/Andras5014/webook/internal/web/jwt"
	"github.com/Andras5014/webook/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

const biz = "login"

var _ handler = &UserHandler{}

type UserHandler struct {
	svc     service.UserService
	codeSvc service.CodeService
	cmd     redis.Cmdable
	ijwt.Handler
	Logger logger.Logger
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, jwtHdl ijwt.Handler, logger logger.Logger) *UserHandler {
	return &UserHandler{
		svc:     svc,
		codeSvc: codeSvc,
		Handler: jwtHdl,
		Logger:  logger,
	}
}
func (u *UserHandler) RegisterRoutes(engine *gin.Engine) {
	ug := engine.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/logout", u.LogoutJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSms)
	ug.POST("/refresh_token", u.RefreshToken)
}
func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	refreshToken := u.ExtractToken(ctx)

	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// redis问题 或者退出登录
		ctx.AbortWithStatus(http.StatusUnauthorized)
		u.Logger.Error("redis错误", logger.Any("err", err))
		return
	}
	err = u.SetJwtToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})
}
func (u *UserHandler) LoginSms(ctx *gin.Context) {
	type LoginReq struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "参数错误",
		})
		return
	}
	if ok, err := u.codeSvc.Verify(ctx, biz, req.Code, req.Phone); err != nil || !ok {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		u.Logger.Error("验证码错误", logger.Any("err", err))
		return
	}

	// 新建或者查找用户
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}

	if err := u.SetLoginToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "校验成功")

}
func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type LoginReq struct {
		Phone string `json:"phone" binding:"required"`
	}

	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "参数错误",
		})
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusOK,
			Msg:  "发送成功",
		})
	case errors.Is(err, service.ErrCodeSendTooMany):
		u.Logger.Warn("发送短信过于频繁", logger.Any("err", err))
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusOK,
			Msg:  "发送太频繁",
		})
	default:
		u.Logger.Error("发送短信失败", logger.Any("err", err))
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusOK,
			Msg:  "系统错误",
		})
	}

	ctx.String(http.StatusOK, "发送成功")
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpForm struct {
		Email           string `json:"email" binding:"required,email"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
		Password        string `json:"password" binding:"required"`
	}
	var req SignUpForm
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "参数错误",
		})
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "两次密码不一致",
		})
		return
	}

	err := u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "邮箱冲突",
		})
		return
	}
	ctx.String(http.StatusOK, "注册成功")

}

//func (u *UserHandler) Login(ctx *gin.Context) {
//	type LoginReq struct {
//		Email    string `json:"email" binding:"required,email"`
//		Password string `json:"password" binding:"required"`
//	}
//
//	var req LoginReq
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		return
//	}
//	user, err := u.svc.Login(ctx, req.Email, req.Password)
//	if errors.Is(err, service.ErrInvalidUserOrPassword) {
//		ctx.String(http.StatusOK, "用户名或者密码错误")
//		return
//	}
//	if err != nil {
//		return
//	}
//
//	// 设置session
//	sess := sessions.Default(ctx)
//	sess.Set("userId", user.Id)
//	sess.Options(sessions.Options{
//		MaxAge: 60 * 60 * 24 * 7,
//	})
//	sess.Save()
//
//	ctx.String(http.StatusOK, "登录成功")
//	return
//}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或者密码错误")
		return
	}
	if err != nil {
		return
	}

	// JWT

	if err := u.SetLoginToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "登录成功")
	return
}
func (u *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "退出成功",
	})
}

func (u *UserHandler) Logout(ctx *gin.Context) {

	// 设置session
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()

	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		NickName string `json:"nickName" binding:"required"`
		Birthday string `json:"birthday" binding:"required"`
		AboutMe  string `json:"aboutMe" binding:"required"`
	}

	var req EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	uc := ctx.MustGet("claims").(*ijwt.UserClaims)
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	err = u.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uc.Uid,
		NickName: req.NickName,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "修改成功")

}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	type ProfileResp struct {
		Email    string
		NickName string
		AboutMe  string
		Birthday string
		Phone    string
	}

	uc := ctx.MustGet("claims").(*ijwt.UserClaims)
	user, err := u.svc.Profile(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, ProfileResp{
		Email:    user.Email,
		NickName: user.NickName,
		AboutMe:  user.AboutMe,
		Birthday: user.Birthday.Format("2006-01-02"),
		Phone:    user.Phone,
	})
}
