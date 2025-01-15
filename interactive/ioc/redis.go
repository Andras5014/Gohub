package ioc

import (
	"github.com/Andras5014/webook/interactive/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})
}

//func InitRedisUniversalClient(cfg *config.Config) redis.UniversalClient {
//
//	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
//		Addrs: []string{cfg.Redis.Addr},
//	})
//	if err := redisClient.Ping(context.Background()).Err(); err != nil {
//		panic(err)
//	}
//
//	return redisClient
//}
