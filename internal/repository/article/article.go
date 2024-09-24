package article

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository"
	"github.com/Andras5014/webook/internal/repository/cache"
	dao "github.com/Andras5014/webook/internal/repository/dao/article"
	"github.com/Andras5014/webook/pkg/logx"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type Repository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error

	// Sync 存储同步数据
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncV1(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx context.Context, article domain.Article) (int64, error)

	// List 制作库
	List(ctx context.Context, id int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)

	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}
type CacheArticleRepository struct {
	dao dao.ArticleDAO

	// v1 操作二个dao
	authorDAO dao.AuthorDAO
	readerDAO dao.ReaderDAO

	userRepo repository.UserRepository

	cache cache.ArticleCache
	l     logx.Logger
}

func (c *CacheArticleRepository) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	pubArt, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	user, err := c.userRepo.FindById(ctx, pubArt.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	return domain.Article{
		Id: pubArt.Id,
		Author: domain.Author{
			Id:   user.Id,
			Name: user.NickName,
		},
		Content:   pubArt.Content,
		Title:     pubArt.Title,
		Status:    domain.ArticleStatus(pubArt.Status),
		CreatedAt: time.UnixMilli(pubArt.CreatedAt),
		UpdatedAt: time.UnixMilli(pubArt.UpdatedAt),
	}, nil
}

func (c *CacheArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	data, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.toDomain(data), nil

}

func NewArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache, l logx.Logger) Repository {
	return &CacheArticleRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}
func (c *CacheArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	defer func() {
		c.cache.DelFirstPage(ctx, article.Author.Id)
	}()
	return c.dao.Insert(ctx, dao.Article{
		AuthorId: article.Author.Id,
		Content:  article.Content,
		Title:    article.Title,
		Status:   article.Status.ToUint8(),
	})
}

func (c *CacheArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		AuthorId: article.Author.Id,
		Content:  article.Content,
		Title:    article.Title,
		Status:   article.Status.ToUint8(),
	})
}
func (c *CacheArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.toEntity(article))
}
func (c *CacheArticleRepository) SyncV1(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	articleEntity := c.toEntity(article)
	if article.Id > 0 {
		err = c.authorDAO.UpdateById(ctx, articleEntity)
	} else {
		id, err = c.authorDAO.Insert(ctx, articleEntity)
	}
	if err != nil {
		return id, err
	}
	return id, c.readerDAO.Upsert(ctx, articleEntity)
}
func (c *CacheArticleRepository) SyncStatus(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.SyncStatus(ctx, c.toEntity(article))
}

func (c *CacheArticleRepository) List(ctx context.Context, id int64, offset int, limit int) ([]domain.Article, error) {
	// 缓存方案
	if offset == 0 && limit <= 100 {
		data, err := c.cache.GeFirstPage(ctx, id)
		if err == nil {
			return data, nil
		}
	}

	res, err := c.dao.FindByAuthorId(ctx, id, offset, limit)
	if err != nil {
		return nil, err
	}

	data := slice.Map[dao.Article, domain.Article](res, func(idx int, src dao.Article) domain.Article {
		return c.toDomain(src)
	})
	go func() {
		err := c.cache.SetFirstPage(ctx, data)
		if err != nil {
			c.l.Error("缓存失败", logx.Any("err", err))
		}
	}()
	return data, err
}
func (c *CacheArticleRepository) preCache(ctx context.Context, articles []domain.Article) {
	const contentSizeThreshold = 1024 * 1024
	if len(articles) > 0 && len(articles[0].Content) <= contentSizeThreshold {
		if err := c.cache.Set(ctx, articles[0]); err != nil {
			c.l.Error("缓存失败", logx.Error(err))
		}
	}
}

func (c *CacheArticleRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		Id:       article.Id,
		AuthorId: article.Author.Id,
		Content:  article.Content,
		Title:    article.Title,
		Status:   article.Status.ToUint8(),
	}
}
func (c *CacheArticleRepository) toDomain(article dao.Article) domain.Article {
	return domain.Article{
		Id:        article.Id,
		Author:    domain.Author{Id: article.AuthorId},
		Content:   article.Content,
		Title:     article.Title,
		Status:    domain.ArticleStatus(article.Status),
		CreatedAt: time.UnixMilli(article.CreatedAt),
		UpdatedAt: time.UnixMilli(article.UpdatedAt),
	}
}
