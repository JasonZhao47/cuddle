package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	GetById(context.Context, int64) (*Article, error)
	Insert(context.Context, *Article) (int64, error)
}

type ArticleGormDAO struct {
	db *gorm.DB
}

func NewArticleGormDAO(db *gorm.DB) ArticleDAO {
	return &ArticleGormDAO{
		db: db,
	}
}

func (dao *ArticleGormDAO) GetById(ctx context.Context, id int64) (*Article, error) {
	var art *Article
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&art).Error
	if err != nil {
		return nil, err
	}
	return art, err
}

func (dao *ArticleGormDAO) Insert(ctx context.Context, article *Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.CTime = now
	article.UTime = now
	err := dao.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// create an index for frequently queried column
	AuthorId int64  `gorm:"index"`
	Topic    string `gorm:"type=varchar(4096)"`
	Content  string `gorm:"type=BLOB"`
	// why use int64 here
	// bson doesn't support time.Time?
	Status uint8
	CTime  int64
	UTime  int64
}
