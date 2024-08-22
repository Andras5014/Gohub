package ioc

import (
	"github.com/Andras5014/webook/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})
}
