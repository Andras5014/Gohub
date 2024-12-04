package ioc

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

func InitRedSync(client redis.UniversalClient) *redsync.Redsync {
	pool := goredis.NewPool(client)
	return redsync.New(pool)
}
