package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Andras5014/gohub/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GeFirstPage(ctx context.Context, id int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, articles []domain.Article) error
	DelFirstPage(ctx context.Context, id int64) error

	Set(ctx context.Context, article domain.Article) error
	Get(ctx context.Context, id int64) (domain.Article, error)

	SetPub(ctx context.Context, article domain.Article) error
	DelPub(ctx context.Context, id int64) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func NewRedisArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		client: client,
	}
}
func (r *RedisArticleCache) GeFirstPage(ctx context.Context, id int64) ([]domain.Article, error) {
	data, err := r.client.Get(ctx, r.firstPageKey(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var articles []domain.Article
	return articles, json.Unmarshal(data, &articles)
}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context, articles []domain.Article) error {
	for i := 0; i < len(articles); i++ {
		articles[i].Content = articles[i].Abstract()
	}
	data, err := json.Marshal(articles)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.firstPageKey(articles[0].Id), data, time.Minute*10).Err()
}

func (r *RedisArticleCache) DelFirstPage(ctx context.Context, id int64) error {
	return r.client.Del(ctx, r.firstPageKey(id)).Err()
}

func (r *RedisArticleCache) Set(ctx context.Context, article domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.authorArtKey(article.Id), data, time.Minute*10).Err()
}

func (r *RedisArticleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	data, err := r.client.Get(ctx, r.authorArtKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var article domain.Article
	return article, json.Unmarshal(data, &article)
}

func (r *RedisArticleCache) SetPub(ctx context.Context, article domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.readerArtKey(article.Id), data, time.Minute*30).Err()
}

func (r *RedisArticleCache) DelPub(ctx context.Context, id int64) error {
	return r.client.Del(ctx, r.readerArtKey(id)).Err()
}

func (r *RedisArticleCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	data, err := r.client.Get(ctx, r.readerArtKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var article domain.Article
	err = json.Unmarshal(data, &article)
	return article, err
}

// 创作端的缓存设置
func (r *RedisArticleCache) authorArtKey(id int64) string {
	return fmt.Sprintf("article:author:%d", id)
}

// 读者端的缓存设置
func (r *RedisArticleCache) readerArtKey(id int64) string {
	return fmt.Sprintf("article:reader:%d", id)
}

func (r *RedisArticleCache) firstPageKey(author int64) string {
	return fmt.Sprintf("article:first_page:%d", author)
}
