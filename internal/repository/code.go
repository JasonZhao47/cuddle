package repository

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
)

type CodeRepository struct {
	cache cache.CodeCache
}

var (
	ErrVerifyTooManyTimes = cache.ErrVerifyTooManyTimes
)

func NewCodeRepository(cache cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}
}

func (repo *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	err := repo.cache.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return nil
}

func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	ok, err := repo.cache.Verify(ctx, biz, phone, code)
	if err != nil {
		return false, err
	}
	return ok, nil
}
