package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const (
	readCntField     = "read_cnt"
	likeCntField     = "like_cnt"
	bookmarkCntField = "bookmark_cnt"
)

type UserActivityCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
}

type userActivityCache struct {
	client redis.Cmdable
}

func NewUserActivityCache(client redis.Cmdable) UserActivityCache {
	return &userActivityCache{client: client}
}

func (c *userActivityCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	// check if cnt present
	// then updates
	// concurrency problem
	// lua script
	actKey := c.key(biz, bizId)
	return c.client.Eval(ctx, luaIncrCnt, []string{actKey}, readCntField, 1).Err()
}

func (c *userActivityCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("user_activity:%s:%d", biz, bizId)
}
