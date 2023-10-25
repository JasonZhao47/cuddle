package service

import (
	"context"
	"fmt"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"math/rand"
)

var (
	ErrTooManyCodeSend = cache.ErrNextCodeTooSoon
)

type CodeService struct {
	repo *repository.CodeRepository
}

func (c *CodeService) Send(ctx context.Context, biz string, phone string) error {
	// 通过biz进行业务区别
	tempCode := generateCode()
	err := c.repo.Set(ctx, biz, phone, tempCode)
	if err != nil {
		return err
	}
	// sms: send message
	return nil
}

func (c *CodeService) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	success, err := c.repo.Verify(ctx, biz, phone, code)
	if err == repository.ErrVerifyTooManyTimes {
		return false, err
	}
	// sms: verify message
	return success, nil
}

func generateCode() string {
	code := rand.Intn(999999)
	return fmt.Sprintf("%6d", code)
}
