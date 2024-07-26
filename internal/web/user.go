package web

import (
	"errors"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/web/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const biz = "login"

type UserHandler struct {
	svc     *service.UserService
	codeSvc *service.CodeService
}

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	return &UserHandler{
		svc:     svc,
		codeSvc: codeSvc,
	}
}
func (u *UserHandler) RegisterRouters(engine *gin.Engine) {
	ug := engine.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit).Use(middleware.NewLoginJWTMiddlewareBuilder().Build())
	ug.GET("/profile", u.Profile).Use(middleware.NewLoginJWTMiddlewareBuilder().Build())
	ug.POST("login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSms)
}

func (u *UserHandler) LoginSms(ctx *gin.Context) {
	type LoginReq struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	if ok, err := u.codeSvc.Verify(ctx, biz, req.Code, req.Phone); err != nil || !ok {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "验证码错误",
		})
		return
	}

	// 新建或者查找用户
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}

	if err := u.setJWTToken(ctx, user.Id); err != nil {
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
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusOK,
			Msg:  "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusOK,
			Msg:  "发送太频繁",
		})
	default:
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
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "两次密码不一致",
		})
		return
	}

	err := u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.JSON(400, gin.H{
			"code": 400,
			"msg":  "邮箱冲突",
		})
		return
	}
	ctx.String(200, "注册成功")

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

	if err := u.setJWTToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := middleware.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("secret"))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
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

}

func (u *UserHandler) Profile(ctx *gin.Context) {
	//c, _ := ctx.Get("claims")
	//
	//claims, ok := c.(*UserClaims)
	//if !ok {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
}
