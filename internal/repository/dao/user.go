package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate  = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(context.Context, User) error
	Create(context.Context, User) error
	FindByEmail(context.Context, string) (User, error)
	FindByPhone(context.Context, string) (User, error)
	FindById(context.Context, int64) (User, error)
	UpdateById(context.Context, User) error
}

type UserGormDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &UserGormDAO{
		db: db,
	}
}

func (dao *UserGormDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime, user.Utime = now, now

	err := dao.db.WithContext(ctx).Create(&user).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		//duplicate key violation for mysql
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *UserGormDAO) Create(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime, user.Utime = now, now

	err := dao.db.WithContext(ctx).Create(&user).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		//duplicate key violation for mysql
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *UserGormDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&user).Error
	return user, err
}

func (dao *UserGormDAO) FindById(ctx context.Context, uid int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id=?", uid).First(&user).Error
	return user, err
}

func (dao *UserGormDAO) UpdateById(ctx context.Context, user User) error {
	err := dao.db.WithContext(ctx).Model(&user).Where("id = ?", user.Id).
		Updates(
			map[string]any{
				"utime":    time.Now().UnixMilli(),
				"nickname": user.Nickname,
				"birthday": user.Birthday,
				"about_me": user.AboutMe,
			},
		).Error
	return err
}

func (dao *UserGormDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).First(&user).Error
	return user, err
}

type User struct {
	Id       int64          `gorm:"primaryKey, autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`

	Nickname string `gorm:"type=varchar(128)"`

	// YYYY-MM-DD
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`

	// 数据库溯源
	Ctime int64
	Utime int64
}
