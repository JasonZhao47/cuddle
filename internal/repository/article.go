package repository

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"time"
)

type ArticleRepository interface {
	GetById(context.Context, int64) (*domain.Article, error)
	Insert(context.Context, *domain.Article) (int64, error)
	GetByAuthorId(context.Context, int64, int, int) ([]*domain.Article, error)
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

func (repo *articleRepository) GetByAuthorId(ctx context.Context, authorId int64, page int, pageSize int) ([]*domain.Article, error) {
	arts, err := repo.dao.GetByAuthorId(ctx, authorId, page, pageSize)
	if err != nil {
		return []*domain.Article{}, err
	}

	res := make([]*domain.Article, len(arts))
	for i := range arts {
		res[i] = repo.toDomain(arts[i])
	}

	return res, err
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
		Status:   uint8(art.Status),
		Content:  art.Content,
		CTime:    art.CTime.UnixMilli(),
		UTime:    art.UTime.UnixMilli(),
	}
}
