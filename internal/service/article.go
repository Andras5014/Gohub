package service

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	if article.Id > 0 {
		return article.Id, a.repo.Update(ctx, article)
	}
	return a.repo.Create(ctx, article)

}
