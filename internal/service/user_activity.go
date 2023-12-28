package service

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/repository"
)

type UserActivityService interface {
	IncrRead(ctx context.Context, biz string, bizId int64) error
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
