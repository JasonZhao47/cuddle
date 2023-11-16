package repository

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"time"
)

type ArticleRepository interface {
	GetById(context.Context, int64) (*domain.Article, error)
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
