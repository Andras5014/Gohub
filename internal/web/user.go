package web

import (
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/service"
	"github.com/gin-gonic/gin"
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
	if err != nil {
		return
	}
	ctx.String(200, "注册成功")

}

func (u *UserHandler) Login(ctx *gin.Context) {
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
