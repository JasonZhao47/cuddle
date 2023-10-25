package cache

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"time"
)

type LocalCodeCache struct {
	cache      *lru.Cache
	expiration time.Duration
}

func NewLocalCodeCache(cache *lru.Cache, expiry time.Duration) *LocalCodeCache {
	return &LocalCodeCache{
		cache:      cache,
		expiration: expiry,
	}
}

func (c *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	return nil
}

func (c *LocalCodeCache) Verify(ctx context.Context, biz string, phone string, input string) (bool, error) {
	return false, nil
}

type codeItem struct {
	code   string
	cnt    int
	expire time.Time
}
