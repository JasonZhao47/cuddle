package repository

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
)

type RankingRepository interface {
	Set(ctx context.Context, arts []domain.PublishedArticle) error
}

type CacheRankingRepository struct {
	cache cache.RankingCache
}

func NewCacheRankingRepository(cache cache.RankingCache) *CacheRankingRepository {
	return &CacheRankingRepository{cache: cache}
}

func (c *CacheRankingRepository) Set(ctx context.Context, arts []domain.PublishedArticle) error {
	return c.cache.Set(ctx, arts)
}
