package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jasonzhao47/cuddle/internal/domain"
	redis "github.com/redis/go-redis/v9"
	"time"
)

// 业务层cache

type UserCache interface {
	Get(context.Context, int64) (domain.User, error)
	Set(context.Context, domain.User) error
}

type UserRedisCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &UserRedisCache{
		cmd:        cmd,
		expiration: 15 * time.Minute,
	}
}

func (cache *UserRedisCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (cache *UserRedisCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	data, err := cache.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (cache *UserRedisCache) Set(ctx context.Context, user domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(user.Id)
	return cache.cmd.Set(ctx, key, data, cache.expiration).Err()
}
