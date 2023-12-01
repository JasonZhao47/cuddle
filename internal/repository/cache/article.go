package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/redis/go-redis/v9"
)

type ArticleCache interface {
	SetFirstPage(arts []*domain.Article) error
	GetFirstPage(ctx context.Context, authorId int64) ([]*domain.Article, error)
}

type articleCache struct {
	client redis.Cmdable
}

func (a *articleCache) SetFirstPage(arts []*domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func NewArticleCache(cmd redis.Cmdable) ArticleCache {
	return &articleCache{client: cmd}
}

func (a *articleCache) GetFirstPage(ctx context.Context, authorId int64) ([]*domain.Article, error) {
	content, err := a.client.Get(ctx, a.firstPageKey(authorId)).Bytes()
	if err != nil {
		return nil, err
	}
	var res []*domain.Article
	err = json.Unmarshal(content, &res)
	if err != nil {
		return res, nil
	}
	return res, nil
}

func (a *articleCache) firstPageKey(uid int64) string {
	return fmt.Sprintf("article_cache:first_page:%d", uid)
}
