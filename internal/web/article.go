package web

import (
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ handler = &ArticleHandler{}

type ArticleHandler struct {
	svc    service.ArticleService
	logger logger.Logger
}

func NewArticleHandler(svc service.ArticleService, logger logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:    svc,
		logger: logger,
	}
}
func (a *ArticleHandler) RegisterRoutes(engine *gin.Engine) {
	ug := engine.Group("/articles")
	ug.POST("/edit", a.Edit)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {

		return
	}
	// 检测输入
	aid := ctx.MustGet("userId").(int64)

	id, err := a.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: aid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("保存文章失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "ok",
		Data: id,
	})

}
