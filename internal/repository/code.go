package repository

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
)

type CodeRepository struct {
	cache *cache.RedisCodeCache
}

var (
	ErrVerifyTooManyTimes = cache.ErrVerifyTooManyTimes
)

func NewCodeRepository(cache *cache.RedisCodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}
}

func (c *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	err := c.cache.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return nil
}

func (c *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	ok, err := c.cache.Verify(ctx, biz, phone, code)
	if err != nil {
		return false, err
	}
	return ok, nil
}
