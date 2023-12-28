package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository/dao"
	"time"
)

type AsyncSmsRepository interface {
	Add(ctx context.Context, s domain.AsyncSms) error
	PreemptiveGetSms(ctx context.Context) (domain.AsyncSms, error)
	ReportScheduleResult(ctx context.Context, smsId int64, success bool) error
}

var ErrSmsNotFound = errors.New("没有找到相应短信")

type asyncSmsRepository struct {
	dao dao.AsyncSmsDAO
}

func (a *asyncSmsRepository) Add(ctx context.Context, s domain.AsyncSms) error {
	return a.dao.Insert(ctx, a.toEntity(s))
}

func (a *asyncSmsRepository) PreemptiveGetSms(ctx context.Context) (domain.AsyncSms, error) {
	smsDao, err := a.dao.GetPreemptiveSms(ctx)
	if err != nil {
		return domain.AsyncSms{}, err
	}
	return a.toDomain(smsDao), nil
}

func (a *asyncSmsRepository) ReportScheduleResult(ctx context.Context, smsId int64, success bool) error {
	if !success {
		return a.dao.MarkFailure(ctx, smsId)
	}
	return a.dao.MarkSuccess(ctx, smsId)
}

func (a *asyncSmsRepository) toDomain(s dao.AsyncSms) domain.AsyncSms {
	var cfg dao.SmsConfig
	_ = json.Unmarshal(s.Config, &cfg)
	return domain.AsyncSms{
		Id:        s.Id,
		TplId:     cfg.TplId,
		RetryMax:  s.RetryMax,
		PhoneNums: cfg.PhoneNums,
		Args:      cfg.Args,
	}
}

func (a *asyncSmsRepository) toEntity(s domain.AsyncSms) dao.AsyncSms {
	now := time.Now().UnixMilli()
	smsConfig, _ := json.Marshal(dao.SmsConfig{
		TplId:     s.TplId,
		Args:      s.Args,
		PhoneNums: s.PhoneNums,
	})
	return dao.AsyncSms{
		Id:     s.Id,
		Config: smsConfig,
		// TODO: add status defs
		Status:   0,
		RetryCnt: 0,
		RetryMax: s.RetryMax,
		Utime:    now,
		Ctime:    now,
	}
}
