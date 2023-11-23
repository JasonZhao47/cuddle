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
	List(ctx context.Context, authorId int64, page int, pageSize int) ([]*domain.Article, error)
	// 操作的表不一样
	Publish(context.Context, *domain.Article) (int64, error)
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
	// upsert, so
	// id is used to identify which case
	id, err := svc.repo.Insert(ctx, art)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (svc *articleService) List(ctx context.Context, authorId int64, page int, pageSize int) ([]*domain.Article, error) {
	return svc.repo.GetByAuthorId(ctx, authorId, page, pageSize)
}

func (svc *articleService) Publish(ctx context.Context, art *domain.Article) (int64, error) {
	art.Status = 1
	return svc.repo.Sync(ctx, art)
}
