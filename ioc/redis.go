package ioc

import (
	"github.com/jasonzhao47/cuddle/configs"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: configs.Config.Redis.Addr,
	})
	return redisClient
}
