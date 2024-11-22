package cache

import (
	"context"
	"encoding/json"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

type RedisRankingCache struct {
	client redis.Cmdable
	key    string
}

func NewRedisRankingCache(client redis.Cmdable) RankingCache {
	return &RedisRankingCache{
		client: client,
	}
}
func (r *RedisRankingCache) Set(ctx context.Context, arts []domain.Article) error {
	for i := 0; i < len(arts); i++ {
		arts[i].Content = ""
	}
	data, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key, data, time.Minute*10).Err()
}

func (r *RedisRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	data, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	return arts, json.Unmarshal(data, &arts)
}
