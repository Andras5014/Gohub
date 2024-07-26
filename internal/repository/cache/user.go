package cache

import (
	"context"
	"encoding/json"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var (
	ErrKeyNotExist = redis.Nil
)

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: 10 * time.Minute,
	}
}

func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	val, err := cache.client.Get(ctx, cache.key(id)).Result()
	// 数据不存在 err = redis.Nil
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(val), &u)
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return cache.client.Set(ctx, cache.key(u.Id), val, cache.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return "user:info:" + strconv.FormatInt(id, 10)
}
