package service

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository"
)

type ArticleService interface {
	GetById(context.Context, int64) (*domain.Article, error)
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
