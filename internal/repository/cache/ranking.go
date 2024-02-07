package cache

import (
	"context"
	"encoding/json"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	ReplaceTopN(ctx context.Context, arts []domain.PublishedArticle) error
}

type rankingCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}

func NewRankingCache(client redis.Cmdable, key string, expiration time.Duration) RankingCache {
	return &rankingCache{client: client, key: key, expiration: expiration}
}

func (r *rankingCache) ReplaceTopN(ctx context.Context, arts []domain.PublishedArticle) error {
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key, val, r.expiration).Err()
}
