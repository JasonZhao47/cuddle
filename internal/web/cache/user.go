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

type UserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	// 松耦合，传进来一个接口，让调用方去实现
	// 这样接口保证了双方的一种规则，能够直接被使用
	// 如果写成传入一个address这样具体的参数，还是紧耦合
	// 紧耦合的问题在于如果想更改底层的实现，必须要更改所有的调用方
	// 除了边缘系统以外，一般不用紧耦合
	return &UserCache{
		cmd:        cmd,
		expiration: 15 * time.Minute,
	}
}

func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
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

func (cache *UserCache) Set(ctx context.Context, user domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(user.Id)
	return cache.cmd.Set(ctx, key, []byte(data), cache.expiration).Err()
}
