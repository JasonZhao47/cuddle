package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type AsyncSmsDAO interface {
	Insert(ctx context.Context, sms AsyncSms) error
	GetPreemptiveSms(ctx context.Context) (AsyncSms, error)
	MarkSuccess(ctx context.Context, smsId int64) error
	MarkFailure(ctx context.Context, smsId int64) error
}

const (
	// 因为本身状态没有暴露出去，所以不需要在 domain 里面定义
	asyncStatusWaiting = iota
	// 失败了，并且超过了重试次数
	asyncStatusFailed
	asyncStatusSuccess
)

type AsyncGormSmsDAO struct {
	db *gorm.DB
}

func NewAsyncSmsDao(db *gorm.DB) AsyncSmsDAO {
	return &AsyncGormSmsDAO{db: db}
}

func (a *AsyncGormSmsDAO) Insert(ctx context.Context, sms AsyncSms) error {
	return a.db.Create(&sms).Error
}

func (a *AsyncGormSmsDAO) GetPreemptiveSms(ctx context.Context) (AsyncSms, error) {
	// 加锁，抢占一分钟之前的
	var sms AsyncSms
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		Utime := time.Now().UnixMilli() - time.Minute.Milliseconds()
		er := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Where("utime < ? AND status = ?", Utime, asyncStatusWaiting).First(&sms).Error
		if er != nil {
			return er
		}
		// 抢占到了，直接先更新，后边人不能抢到了
		// 防止后面的并发条件导致死锁

		er = tx.Model(&AsyncSms{}).
			Where("id = ?", sms.Id).
			Updates(map[string]any{
				"retry_cnt": gorm.Expr("retry_cnt + 1"),
				// 保证没有人可以抢到了
				"utime": now,
			}).Error
		return er
	})
	return sms, err
}

func (a *AsyncGormSmsDAO) MarkSuccess(ctx context.Context, smsId int64) error {
	now := time.Now().UnixMilli()
	return a.db.WithContext(ctx).Where("id = ?", smsId).Updates(map[string]any{
		"status": asyncStatusSuccess,
		"utime":  now,
	}).Error
}

func (a *AsyncGormSmsDAO) MarkFailure(ctx context.Context, smsId int64) error {
	now := time.Now().UnixMilli()
	return a.db.WithContext(ctx).Where("id = ?", smsId).Updates(map[string]any{
		"status": asyncStatusFailed,
		"utime":  now,
	}).Error
}

type AsyncSms struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Config   []byte `gorm:"type:json"`
	Status   uint8
	RetryCnt int
	RetryMax int
	Utime    int64 `gorm:"index"`
	Ctime    int64
}

type SmsConfig struct {
	TplId     string
	Args      []string
	PhoneNums []string
}
