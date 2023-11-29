package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"time"
)

type ArticleRepository interface {
	GetById(context.Context, int64) (*domain.Article, error)
	Insert(context.Context, *domain.Article) (int64, error)
	GetByAuthorId(context.Context, int64, int, int) ([]*domain.Article, error)
	Sync(context.Context, *domain.Article) (int64, error)
	SyncStatus(context.Context, int64, int64, domain.ArticleStatus) error
}

type articleRepository struct {
	dao dao.ArticleDAO
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &articleRepository{
		dao: dao,
	}
}

func (repo *articleRepository) GetById(ctx context.Context, id int64) (*domain.Article, error) {
	// need to add cache here
	art, err := repo.dao.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return repo.toDomain(art), nil
}

func (repo *articleRepository) Insert(ctx context.Context, article *domain.Article) (int64, error) {
	id, err := repo.dao.Insert(ctx, repo.toEntity(article))
	if err != nil {
		return article.Id, err
	}
	return id, nil
}

func (repo *articleRepository) GetByAuthorId(ctx context.Context, authorId int64, limit int, offset int) ([]*domain.Article, error) {
	arts, err := repo.dao.GetByAuthorId(ctx, authorId, limit, offset)
	if err != nil {
		return []*domain.Article{}, err
	}

	res := slice.Map[*dao.Article, *domain.Article](arts, func(idx int, src *dao.Article) *domain.Article {
		return repo.toDomain(src)
	})

	return res, err
}

func (repo *articleRepository) Sync(ctx context.Context, article *domain.Article) (int64, error) {
	// dao同步数据
	return repo.dao.Sync(ctx, repo.toEntity(article))
}

func (repo *articleRepository) SyncStatus(ctx context.Context, userId int64, artId int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, userId, artId, status)
}

func (repo *articleRepository) toDomain(dao *dao.Article) *domain.Article {
	return &domain.Article{
		Id: dao.Id,
		Author: domain.Author{
			Id: dao.AuthorId,
			// what about author name?
			// join?
		},
		Topic:   dao.Topic,
		Status:  domain.ArticleStatus(dao.Status),
		Content: dao.Content,
		CTime:   time.UnixMilli(dao.CTime),
		UTime:   time.UnixMilli(dao.UTime),
	}
}

func (repo *articleRepository) toEntity(art *domain.Article) *dao.Article {
	return &dao.Article{
		Id:       art.Id,
		AuthorId: art.Author.Id,
		Topic:    art.Topic,
		Status:   art.Status.ToUint8(),
		Content:  art.Content,
		CTime:    art.CTime.UnixMilli(),
		UTime:    art.UTime.UnixMilli(),
	}
}
