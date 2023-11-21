package service

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository"
)

type ArticleService interface {
	GetById(context.Context, int64) (*domain.Article, error)
	// Save upsert语义
	Save(context.Context, *domain.Article) (int64, error)
	GetByAuthorId(ctx context.Context, authorId int64, page int, pageSize int) ([]*domain.Article, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (svc *articleService) GetById(ctx context.Context, id int64) (*domain.Article, error) {
	art, err := svc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return art, nil
}

func (svc *articleService) Save(ctx context.Context, art *domain.Article) (int64, error) {
	id, err := svc.repo.Insert(ctx, art)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (svc *articleService) GetByAuthorId(ctx context.Context, authorId int64, page int, pageSize int) ([]*domain.Article, error) {
	panic("Implement me")
}
