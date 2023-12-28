package repository

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/slice"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"time"
)

var (
	ErrIllegalOffsetOrLimit = errors.New("非法偏移量")
)

type ArticleRepository interface {
	GetById(context.Context, int64) (domain.Article, error)
	Insert(context.Context, domain.Article) (int64, error)
	GetByAuthor(context.Context, int64, int, int) ([]domain.Article, error)
	Sync(context.Context, domain.Article) (int64, error)
	SyncStatus(context.Context, int64, int64, domain.ArticleStatus) error
	GetPubById(context.Context, int64) (domain.PublishedArticle, error)
}

type CachedArticleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache
}

func NewArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache) ArticleRepository {
	return &CachedArticleRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CachedArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	// need to add cache here
	res, err := repo.cache.Get(ctx, id)
	if err == nil {
		return res, nil
	}
	art, err := repo.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	s := repo.toDomain(art)
	// 设置缓存不是主要逻辑
	// 可以放在线程里面做
	go func() {
		err = repo.cache.Set(ctx, s)
		if err != nil {
			// log here
		}
	}()
	return s, nil
}

func (repo *CachedArticleRepository) Insert(ctx context.Context, article domain.Article) (int64, error) {
	id, err := repo.dao.Insert(ctx, repo.toEntity(article))
	if err == nil {
		err = repo.cache.EraseFirstPage(ctx, article.Author.Id)
		// need to add cache here
		if err != nil {
			// log here?
		}
	}
	return id, err
}

func (repo *CachedArticleRepository) GetByAuthor(ctx context.Context, authorId int64, limit int, offset int) ([]domain.Article, error) {
	// handle when offset is not legal
	var res []domain.Article
	if limit < 0 || offset < 0 {
		return res, ErrIllegalOffsetOrLimit
	}
	if limit == 0 && offset <= 100 {
		// cached range
		res, err := repo.cache.GetFirstPage(ctx, authorId)
		if err == nil {
			return res, nil
		} else {
			// log here
		}
	}
	arts, err := repo.dao.GetByAuthorId(ctx, authorId, limit, offset)
	if err != nil {
		return res, err
	}
	res = slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return repo.toDomain(src)
	})

	// 不能简单同步处理
	// 过段时间要取消掉
	// 回写缓存错误不是重要的错误
	// 那么我们异步处理也可以
	go func() {
		if limit == 0 && offset <= 100 {
			err = repo.cache.SetFirstPage(ctx, res, authorId)
			if err != nil {
				// log here
			}
		}
	}()
	go func() {
		repo.preCache(ctx, res)
	}()
	return res, nil
}

func (repo *CachedArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	// dao同步数据
	// need to erase cache here
	id, err := repo.dao.Sync(ctx, repo.toEntity(article))
	if err == nil {
		// clear cache here
		err = repo.cache.EraseFirstPage(ctx, article.Author.Id)
		if err != nil {
			// log here?
		}
	}
	// 异步设置缓存
	//go func() {
	//
	//}
	return id, err
}

func (repo *CachedArticleRepository) SyncStatus(ctx context.Context, userId int64, artId int64, status domain.ArticleStatus) error {
	// need to erase cache here
	err := repo.dao.SyncStatus(ctx, userId, artId, status)
	if err == nil {
		// clear cache here
		err = repo.cache.EraseFirstPage(ctx, artId)
		if err != nil {
			// log here?
		}
	}
	return err
}

func (repo *CachedArticleRepository) GetPubById(ctx context.Context, id int64) (domain.PublishedArticle, error) {
	// get cache of first page
	res, err := repo.cache.GetPub(ctx, id)
	if err == nil {
		return res, nil
	}
	// log here to indicate cache miss
	art, err := repo.dao.GetByPublishedId(ctx, id)
	if err != nil {
		return domain.PublishedArticle{}, nil
	}
	s := repo.toPublishedDomain(art)
	go func() {
		err = repo.cache.SetPub(ctx, s)
		if err != nil {
			// log here
		}
	}()
	return s, nil
}

func (repo *CachedArticleRepository) toDomain(dao dao.Article) domain.Article {
	return domain.Article{
		Id: dao.Id,
		Author: domain.Author{
			Id: dao.AuthorId,
			// what about author name?
			// join?
		},
		Topic:   dao.Topic,
		Status:  domain.ArticleStatus(dao.Status),
		Content: dao.Content,
		Ctime:   time.UnixMilli(dao.Ctime),
		Utime:   time.UnixMilli(dao.Utime),
	}
}

func (repo *CachedArticleRepository) toPublishedDomain(dao dao.PublishedArticle) domain.PublishedArticle {
	return domain.PublishedArticle{
		Id: dao.Id,
		Author: domain.Author{
			Id: dao.AuthorId,
			// what about author name?
			// join?
		},
		Topic:   dao.Topic,
		Status:  domain.ArticleStatus(dao.Status),
		Content: dao.Content,
		Ctime:   time.UnixMilli(dao.Ctime),
		Utime:   time.UnixMilli(dao.Utime),
	}
}

func (repo *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		AuthorId: art.Author.Id,
		Topic:    art.Topic,
		Status:   art.Status.ToUint8(),
		Content:  art.Content,
		Ctime:    art.Ctime.UnixMilli(),
		Utime:    art.Utime.UnixMilli(),
	}
}

func (repo *CachedArticleRepository) preCache(ctx context.Context, arts []domain.Article) {
	const contentSizeThreshold = 1024 * 1024
	// 方案：提前缓存好某几个条目
	if len(arts) > 0 && len(arts[0].Content) <= contentSizeThreshold {
		if err := repo.cache.Set(ctx, arts[0]); err != nil {
			// log here
		}
	}
}
