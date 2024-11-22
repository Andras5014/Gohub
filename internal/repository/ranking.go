package repository

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
}
type rankingRepository struct {
	cache cache.RankingCache
}

func (r *rankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return r.cache.Set(ctx, arts)
}

func NewRankingRepository() RankingRepository {
	return &rankingRepository{}
}
