package service

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository"
)

type UserActivityService interface {
	IncrRead(ctx context.Context, biz string, bizId int64) error
	GetReadByIds(ctx context.Context, biz string, ids []int64) ([]domain.UserActivity, error)
}

type activityService struct {
	repo repository.UserActivityRepository
}

func NewUserActivityService(repo repository.UserActivityRepository) UserActivityService {
	return &activityService{repo: repo}
}

func (svc *activityService) IncrRead(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrRead(ctx, biz, bizId)
}

func (svc *activityService) GetReadByIds(ctx context.Context, biz string, ids []int64) ([]domain.UserActivity, error) {
	return svc.repo.GetReadByIds(ctx, biz, ids)
}
