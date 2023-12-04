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
	EraseFirstPage(ctx context.Context, authorId int64) error
	Set(ctx context.Context, art domain.Article) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	SetPub(ctx context.Context, art domain.PublishedArticle) error
	GetPub(ctx context.Context, id int64) (domain.PublishedArticle, error)
}

type articleCache struct {
	client redis.Cmdable
}

func NewArticleCache(cmd redis.Cmdable) ArticleCache {
	return &articleCache{client: cmd}
}

func (a *articleCache) SetPub(ctx context.Context, art domain.PublishedArticle) error {
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, a.publicArtKey(art.Id), data, time.Minute*10).Err()
}

func (a *articleCache) GetPub(ctx context.Context, artId int64) (domain.PublishedArticle, error) {
	data, err := a.client.Get(ctx, a.publicArtKey(artId)).Bytes()
	if err != nil {
		return domain.PublishedArticle{}, err
	}
	var art domain.PublishedArticle
	err = json.Unmarshal(data, &art)
	return art, err
}

func (a *articleCache) Set(ctx context.Context, art domain.Article) error {
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, a.artKey(art.Id), data, time.Minute*10).Err()
}

func (a *articleCache) GetFirstPage(ctx context.Context, authorId int64) ([]domain.Article, error) {
	content, err := a.client.Get(ctx, a.firstPageKey(authorId)).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(content, &res)

	return res, err
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
	return a.client.Set(ctx, a.firstPageKey(authorId), val, time.Minute*10).Err()
}

func (a *articleCache) EraseFirstPage(ctx context.Context, authorId int64) error {
	return a.client.Del(ctx, a.firstPageKey(authorId)).Err()
}

func (a *articleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	data, err := a.client.Get(ctx, a.artKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(data, &res)
	return res, err
}

func (a *articleCache) firstPageKey(uid int64) string {
	return fmt.Sprintf("article_cache:first_page:%d", uid)
}

func (a *articleCache) artKey(uid int64) string {
	return fmt.Sprintf("article_cache:id:%d", uid)
}

func (a *articleCache) publicArtKey(uid int64) string {
	return fmt.Sprintf("article_cache:pub:id:%d", uid)
}
