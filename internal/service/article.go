package service

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository"
)

type ArticleService interface {
	GetById(context.Context, int64) (domain.Article, error)
	// Save upsert语义
	Save(context.Context, domain.Article) (int64, error)
	List(ctx context.Context, authorId int64, limit int, offset int) ([]domain.Article, error)
	// 操作的表不一样
	Publish(context.Context, domain.Article) (int64, error)
	WithDraw(ctx context.Context, userId int64, artId int64) error
	GetPubById(ctx context.Context, artId int64) (domain.PublishedArticle, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (svc *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := svc.repo.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return art, nil
}

func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	// upsert, so
	// id is used to identify which case
	id, err := svc.repo.Insert(ctx, art)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (svc *articleService) List(ctx context.Context, authorId int64, limit int, offset int) ([]domain.Article, error) {
	return svc.repo.GetByAuthor(ctx, authorId, limit, offset)
}

func (svc *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, art)
}

func (svc *articleService) WithDraw(ctx context.Context, userId int64, artId int64) error {
	return svc.repo.SyncStatus(ctx, userId, artId, domain.ArticleStatusPrivate)
}

func (svc *articleService) GetPubById(ctx context.Context, artId int64) (domain.PublishedArticle, error) {
	return svc.repo.GetPubById(ctx, artId)
}
