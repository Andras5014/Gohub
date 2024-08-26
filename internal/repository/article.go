package repository

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}
type CacheArticleRepository struct {
	dao dao.ArticleDAO
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
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
