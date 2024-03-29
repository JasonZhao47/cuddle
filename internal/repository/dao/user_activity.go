package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserActivityDAO interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCntIfPresent(ctx context.Context, bizs []string, bizIds []int64) error
	GetReadByIds(ctx context.Context, biz string, bizIds []int64) ([]UserActivity, error)
}

type userActivityDAO struct {
	db *gorm.DB
}

func NewUserActivityDAO(db *gorm.DB) UserActivityDAO {
	return &userActivityDAO{db: db}
}

func (d *userActivityDAO) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	err := d.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"read_cnt": gorm.Expr("`read_cnt` + 1"),
				"utime":    now,
			}),
		}).Create(&UserActivity{
		Biz:     biz,
		BizId:   bizId,
		ReadCnt: 1,
		Utime:   now,
		Ctime:   now,
	}).Error
	return err
}

func (d *userActivityDAO) BatchIncrReadCntIfPresent(ctx context.Context, bizs []string, bizIds []int64) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 调用者应该保证对应关系
		n := len(bizs)
		for i := 0; i < n; i++ {
			err := d.IncrReadCntIfPresent(ctx, bizs[i], bizIds[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *userActivityDAO) GetReadByIds(ctx context.Context, biz string, bizIds []int64) ([]UserActivity, error) {
	var res []UserActivity
	err := d.db.WithContext(ctx).
		Where("biz = ? AND biz_id IN ?", biz, bizIds).
		Find(&res).Error
	return res, err
}

type UserActivity struct {
	// 联合唯一索引：防止并发写入问题
	// 保证任何时刻这两个列都是一致的
	Id          int64  `gorm:"primaryKey,autoIncrement"`
	Biz         string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	BizId       int64  `gorm:"uniqueIndex:biz_type_id"`
	ReadCnt     int64
	LikeCnt     int64
	BookmarkCnt int64
	Utime       int64
	Ctime       int64
}
