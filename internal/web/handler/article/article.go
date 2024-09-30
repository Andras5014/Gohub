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
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

var _ handler.Handler = &Handler{}

type Handler struct {
	svc    service.ArticleService
	logger logx.Logger

	intrSvc service.InteractiveService
	biz     string
}

func NewArticleHandler(svc service.ArticleService, intrSvc service.InteractiveService, logger logx.Logger) *Handler {
	return &Handler{
		svc:     svc,
		intrSvc: intrSvc,
		logger:  logger,
		biz:     "article",
	}
}
func (h *Handler) RegisterRoutes(engine *gin.Engine) {
	// 用户对自己操作
	ug := engine.Group("/articles")
	{
		ug.POST("/edit", h.Edit)
		ug.POST("/publish", h.Publish)
		ug.POST("/withdraw", h.Withdraw)
		ug.POST("/list", ginx.WrapBody(h.logger, h.List))
		ug.GET("/detail/:id", ginx.Wrap(h.logger, h.Detail))
	}

	pub := engine.Group("/pub")
	{
		pub.GET("/:id", ginx.Wrap(h.logger, h.PubDetail))
		pub.POST("/like", ginx.WrapBody(h.logger, h.Like))
	}
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

// PubDetail 提供发布详情的接口。
// 通过文章ID获取发布的详细信息，并增加阅读计数。
// 参数:
//   - ctx: gin的上下文，包含请求信息及各种实用方法。
//
// 返回值:
//   - ginx.Result: 包含请求结果的数据结构。
//   - error: 错误信息，如果执行过程中出现错误。
func (h *Handler) PubDetail(ctx *gin.Context) (ginx.Result, error) {
	var (
		id          int64
		article     domain.Article
		interactive domain.Interactive
		eg          errgroup.Group
		err         error
	)

	// 从URL参数中提取文章ID字符串
	idStr := ctx.Param("id")
	// 将ID字符串转换为int64类型
	id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// 如果转换失败，返回参数无效错误
		return ginx.InvalidParam(), err
	}

	eg.Go(func() error {
		// 根据ID获取发布文章详情
		article, err = h.svc.GetPubById(ctx, id)
		return err
	})

	eg.Go(func() error {
		uid := ctx.GetInt64("userId")
		interactive, err = h.intrSvc.Get(ctx, h.biz, id, uid)
		return err
	})

	err = eg.Wait()
	if err != nil {
		//查询出错
		return ginx.SystemError(), err
	}

	// 异步增加阅读计数
	go func() {
		// 调用内部服务增加阅读计数
		er := h.intrSvc.IncrReadCnt(ctx, h.biz, id)
		if er != nil {
			// 如果增加阅读计数失败，记录错误日志
			h.logger.Error("增加阅读计数失败", logx.Error(er), logx.Any("aid", id))
			return
		}
	}()

	// 返回文章详情结果
	return ginx.Result{
		Data: ArticleVO{
			Id:         article.Id,
			Title:      article.Title,
			CreatedAt:  article.CreatedAt.String(),
			UpdatedAt:  article.UpdatedAt.String(),
			Status:     article.Status.ToUint8(),
			Abstract:   article.Abstract(),
			AuthorId:   article.Author.Id,
			AuthorName: article.Author.Name,
			LikeCnt:    interactive.LikeCnt,
			CollectCnt: interactive.CollectCnt,
			ReadCnt:    interactive.ReadCnt,
			Liked:      interactive.Liked,
			Collected:  interactive.Collected,
		},
	}, nil
}

func (h *Handler) Like(ctx *gin.Context, req LikeReq) (ginx.Result, error) {
	Uid := ctx.GetInt64("userId")
	var err error
	if req.Like {
		// 点赞
		err = h.intrSvc.Like(ctx, h.biz, req.Id, Uid)
	} else {
		// 取消点赞
		err = h.intrSvc.CancelLike(ctx, h.biz, req.Id, Uid)
	}
	if err != nil {
		return ginx.SystemError(), err
	}
	return ginx.Success(), nil
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
