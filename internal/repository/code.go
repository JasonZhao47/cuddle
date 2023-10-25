package repository

import "github.com/jasonzhao47/cuddle/internal/repository/cache"

type CodeRepository struct {
	cache *cache.RedisCodeCache
}

func NewCodeRepository(cache *cache.RedisCodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}
}
