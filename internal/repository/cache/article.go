package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	SetFirstPage(ctx context.Context, arts []domain.Article, authorId int64) error
	GetFirstPage(ctx context.Context, authorId int64) ([]domain.Article, error)
}

type articleCache struct {
	client redis.Cmdable
}

func NewArticleCache(cmd redis.Cmdable) ArticleCache {
	return &articleCache{client: cmd}
}

func (a *articleCache) GetFirstPage(ctx context.Context, authorId int64) ([]domain.Article, error) {
	content, err := a.client.Get(ctx, a.firstPageKey(authorId)).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(content, &res)
	if err != nil {
		return res, nil
	}
	return res, nil
}

func (a *articleCache) SetFirstPage(ctx context.Context, arts []domain.Article, authorId int64) error {
	for idx := range arts {
		arts[idx].Content = arts[idx].Abstract()
	}
	// only set cache on abstracts
	// not contents
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	// SETNX firstPageKey val 600
	err = a.client.Set(ctx, a.firstPageKey(authorId), val, time.Minute*10).Err()
	if err != nil {
		return err
	}
	return nil
}

func (a *articleCache) firstPageKey(uid int64) string {
	return fmt.Sprintf("article_cache:first_page:%d", uid)
}
