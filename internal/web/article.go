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
	ug.POST("/publish", a.Publish)
	ug.POST("/withdraw", a.Withdraw)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {

	var req articleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {

		return
	}
	// 检测输入
	//aid := ctx.MustGet("userId").(int64)
	authorId := ctx.GetInt64("userId")
	id, err := a.svc.Save(ctx, req.toDomain(authorId))
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

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req articleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}
	// 检测输入
	//aid := ctx.MustGet("userId").(int64)
	authorId := ctx.GetInt64("userId")
	id, err := a.svc.Publish(ctx, req.toDomain(authorId))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "ok",
		Data: id,
	})
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	var req articleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}
	// 检测输入
	//aid := ctx.MustGet("userId").(int64)
	authorId := ctx.GetInt64("userId")
	id, err := a.svc.Withdraw(ctx, req.toDomain(authorId))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "ok",
		Data: id,
	})
}

type articleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (a *articleReq) toDomain(authorId int64) domain.Article {
	return domain.Article{
		Id:      a.Id,
		Title:   a.Title,
		Content: a.Content,
		Author: domain.Author{
			Id: authorId,
		},
	}
}
