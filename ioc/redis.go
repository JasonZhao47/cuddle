package ioc

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var config Config
	err := viper.UnmarshalKey("data.redis", &config)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败%v", err.Error()))
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Addr,
	})
	return redisClient
}
