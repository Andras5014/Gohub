package article

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	dao "github.com/Andras5014/webook/internal/repository/dao/article"
)

type Repository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	// Sync 存储同步数据
	Sync(ctx context.Context, article domain.Article) (int64, error)
}
type CacheArticleRepository struct {
	dao dao.ArticleDAO

	// v1 操作二个dao
	authorDAO dao.AuthorDAO
	readerDAO dao.ReaderDAO
}

func NewArticleRepository(dao dao.ArticleDAO) Repository {
	return &CacheArticleRepository{
		dao: dao,
	}
}
func (c *CacheArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		AuthorId: article.Author.Id,
		Content:  article.Content,
		Title:    article.Title,
	})
}

func (c *CacheArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		AuthorId: article.Author.Id,
		Content:  article.Content,
		Title:    article.Title,
	})
}
func (c *CacheArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
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

func (c *CacheArticleRepository) toEntity(article domain.Article) dao.Article {
	return dao.Article{
		Id:       article.Id,
		AuthorId: article.Author.Id,
		Content:  article.Content,
		Title:    article.Title,
	}
}
