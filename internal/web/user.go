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

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}
func (u *UserHandler) RegisterRouters(engine *gin.Engine) {
	ug := engine.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit).Use(middleware.NewLoginMiddlewareBuilder().Build())
	ug.GET("/profile", u.Profile).Use(middleware.NewLoginMiddlewareBuilder().Build())
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

func (u *UserHandler) Login(ctx *gin.Context) {
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

	// 设置session
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge: 60 * 60 * 24 * 7,
	})
	sess.Save()

	ctx.String(http.StatusOK, "登录成功")
	return
}

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

	claims := domain.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("secret"))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "登录成功")
	return
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
