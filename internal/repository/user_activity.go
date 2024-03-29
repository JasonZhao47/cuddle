package repository

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
)

type UserActivityRepository interface {
	IncrRead(ctx context.Context, biz string, bizId int64) error
	BatchIncrRead(ctx context.Context, bizs []string, bizIds []int64) error
	GetReadByIds(ctx context.Context, biz string, ids []int64) ([]domain.UserActivity, error)
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

func (repo *CacheUserActivityRepository) BatchIncrRead(ctx context.Context, bizs []string, bizIds []int64) error {
	err := repo.dao.BatchIncrReadCntIfPresent(ctx, bizs, bizIds)
	if err != nil {
		return err
	}
	// 顺序保证要让调用者保证
	n := len(bizs)
	go func() {
		for i := 0; i < n; i++ {
			er := repo.cache.IncrReadCntIfPresent(ctx, bizs[i], bizIds[i])
			if er != nil {
				// log here
			}
		}
	}()
	return nil
}

func (repo *CacheUserActivityRepository) GetReadByIds(ctx context.Context, biz string, ids []int64) ([]domain.UserActivity, error) {
	res := make([]domain.UserActivity, 0)
	userDao, err := repo.dao.GetReadByIds(ctx, biz, ids)
	if err != nil {
		return []domain.UserActivity{}, nil
	}
	for _, d := range userDao {
		res = append(res, repo.toDomain(d))
	}
	return res, nil
}

func (repo *CacheUserActivityRepository) toDomain(dao dao.UserActivity) domain.UserActivity {
	return domain.UserActivity{
		Id:      dao.Id,
		ReadCnt: dao.ReadCnt,
	}

}
