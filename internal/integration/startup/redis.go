package startup

import (
	"github.com/Andras5014/gohub/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})
}
