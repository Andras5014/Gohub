package web

import (
	"errors"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
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
	err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或者密码错误")
		return
	}
	if err != nil {
		return
	}

	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
