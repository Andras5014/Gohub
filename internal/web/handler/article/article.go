package article

import (
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/web/handler"
	"github.com/Andras5014/webook/internal/web/result"
	"github.com/Andras5014/webook/pkg/logx"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ handler.Handler = &Handler{}

type Handler struct {
	svc    service.ArticleService
	logger logx.Logger
}

func NewArticleHandler(svc service.ArticleService, logger logx.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}
func (a *Handler) RegisterRoutes(engine *gin.Engine) {
	ug := engine.Group("/articles")
	ug.POST("/edit", a.Edit)
	ug.POST("/publish", a.Publish)
	ug.POST("/withdraw", a.Withdraw)
}

func (a *Handler) Edit(ctx *gin.Context) {

	var req articleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {

		return
	}
	// 检测输入
	//aid := ctx.MustGet("userId").(int64)
	authorId := ctx.GetInt64("userId")
	id, err := a.svc.Save(ctx, req.toDomain(authorId))
	if err != nil {
		ctx.JSON(http.StatusOK, result.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("保存文章失败", logx.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, result.Result{
		Msg:  "ok",
		Data: id,
	})
}

func (a *Handler) Publish(ctx *gin.Context) {
	var req articleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}
	// 检测输入
	//aid := ctx.MustGet("userId").(int64)
	authorId := ctx.GetInt64("userId")
	id, err := a.svc.Publish(ctx, req.toDomain(authorId))
	if err != nil {
		ctx.JSON(http.StatusOK, result.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("发表帖子失败", logx.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, result.Result{
		Msg:  "ok",
		Data: id,
	})
}

func (a *Handler) Withdraw(ctx *gin.Context) {
	var req articleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}
	// 检测输入
	//aid := ctx.MustGet("userId").(int64)
	authorId := ctx.GetInt64("userId")
	id, err := a.svc.Withdraw(ctx, req.toDomain(authorId))
	if err != nil {
		ctx.JSON(http.StatusOK, result.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("发表帖子失败", logx.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, result.Result{
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
