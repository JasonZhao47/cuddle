package service

import (
	"context"
	"fmt"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/internal/repository/cache"
	"github.com/jasonzhao47/cuddle/internal/service/sms"
	"math/rand"
)

var (
	ErrTooManyCodeSend = cache.ErrNextCodeTooSoon
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type SMSCodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

func NewSMSCodeService(repo *repository.CodeRepository, sms sms.Service) *SMSCodeService {
	return &SMSCodeService{repo: repo, sms: sms}
}

func (c *SMSCodeService) Send(ctx context.Context, biz string, phone string) error {
	// 通过biz进行业务区别
	tempCode := c.generateCode()
	err := c.repo.Set(ctx, biz, phone, tempCode)
	if err != nil {
		return err
	}
	const codeTplId = "19381892"
	return c.sms.Send(ctx, codeTplId, []string{tempCode}, []string{phone})
}

func (c *SMSCodeService) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	success, err := c.repo.Verify(ctx, biz, phone, code)
	if err == repository.ErrVerifyTooManyTimes {
		return false, err
	}
	// sms: verify message
	return success, nil
}

func (c *SMSCodeService) generateCode() string {
	code := rand.Intn(999999)
	return fmt.Sprintf("%6d", code)
}
