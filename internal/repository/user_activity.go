package repository

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
)

type UserActivityRepository interface {
	IncrRead(ctx context.Context, biz string, bizId int64) error
}

type CacheUserActivityRepository struct {
	cache cache.UserActivityCache
	dao   dao.UserActivityDAO
}

func NewCacheUserActivityRepository(cache cache.UserActivityCache, dao dao.UserActivityDAO) UserActivityRepository {
	return &CacheUserActivityRepository{cache: cache, dao: dao}
}

func (repo *CacheUserActivityRepository) IncrRead(ctx context.Context, biz string, bizId int64) error {
	// update db first
	err := repo.dao.IncrReadCntIfPresent(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return repo.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}
