package article

import (
	"errors"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/web/handler"
	"github.com/Andras5014/webook/internal/web/result"
	"github.com/Andras5014/webook/pkg/ginx"
	"github.com/Andras5014/webook/pkg/logx"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
func (h *Handler) RegisterRoutes(engine *gin.Engine) {
	// 用户对自己操作
	ug := engine.Group("/articles")
	ug.POST("/edit", h.Edit)
	ug.POST("/publish", h.Publish)
	ug.POST("/withdraw", h.Withdraw)
	ug.POST("/list", ginx.WrapBody(h.logger, h.List))
	ug.GET("/detail/:id", ginx.Wrap(h.logger, h.Detail))
}

func (h *Handler) Edit(ctx *gin.Context) {

	var req articleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {

		return
	}
	// 检测输入
	//aid := ctx.MustGet("userId").(int64)
	authorId := ctx.GetInt64("userId")
	id, err := h.svc.Save(ctx, req.toDomain(authorId))
	if err != nil {
		ctx.JSON(http.StatusOK, result.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.logger.Error("保存文章失败", logx.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, result.Result{
		Msg:  "ok",
		Data: id,
	})
}

func (h *Handler) Publish(ctx *gin.Context) {
	var req articleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}
	// 检测输入
	//aid := ctx.MustGet("userId").(int64)
	authorId := ctx.GetInt64("userId")
	id, err := h.svc.Publish(ctx, req.toDomain(authorId))
	if err != nil {
		ctx.JSON(http.StatusOK, result.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.logger.Error("发表帖子失败", logx.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, result.Result{
		Msg:  "ok",
		Data: id,
	})
}

func (h *Handler) Withdraw(ctx *gin.Context) {
	var req articleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	// 检测输入
	//aid := ctx.MustGet("userId").(int64)
	authorId := ctx.GetInt64("userId")
	id, err := h.svc.Withdraw(ctx, req.toDomain(authorId))
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.SystemError())
		h.logger.Error("发表帖子失败", logx.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, result.Result{
		Msg:  "ok",
		Data: id,
	})
}

func (h *Handler) List(ctx *gin.Context, req ListReq) (ginx.Result, error) {
	id := ctx.MustGet("userId").(int64)
	res, err := h.svc.List(ctx, id, req.Offset, req.Limit)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, nil
	}
	return ginx.Result{
		Data: slice.Map[domain.Article, ArticleVO](res, func(idx int, src domain.Article) ArticleVO {
			return ArticleVO{
				Id:        src.Id,
				Title:     src.Title,
				CreatedAt: src.CreatedAt.String(),
				UpdatedAt: src.UpdatedAt.String(),
				Status:    src.Status.ToUint8(),
				Abstract:  src.Abstract(),
			}
		}),
	}, nil
}

func (h *Handler) Detail(ctx *gin.Context) (ginx.Result, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ginx.InvalidParam(), err
	}

	article, err := h.svc.GetById(ctx, id)
	if err != nil {
		return ginx.SystemError(), err
	}

	userId := ctx.GetInt64("userId")
	if article.Author.Id != userId {
		return ginx.Result{
			Code: 4,
			Msg:  "无权限",
		}, errors.New("无权限")
	}

	return ginx.Result{
		Data: ArticleVO{
			Id:        article.Id,
			Title:     article.Title,
			CreatedAt: article.CreatedAt.String(),
			UpdatedAt: article.UpdatedAt.String(),
			Status:    article.Status.ToUint8(),
			Content:   article.Content,
		},
	}, nil
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
