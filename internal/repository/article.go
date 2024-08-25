package repository

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
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
