package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	GetById(context.Context, int64) (*Article, error)
	Insert(context.Context, *Article) (int64, error)
	GetByAuthorId(context.Context, int64, int, int) ([]*Article, error)
	Sync(context.Context, *Article) (int64, error)
}

type ArticleGormDAO struct {
	db *gorm.DB
}

func NewArticleGormDAO(db *gorm.DB) ArticleDAO {
	return &ArticleGormDAO{
		db: db,
	}
}

func (d *ArticleGormDAO) GetById(ctx context.Context, id int64) (*Article, error) {
	var art *Article
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&art).Error
	if err != nil {
		return nil, err
	}
	return art, err
}

func (d *ArticleGormDAO) Insert(ctx context.Context, article *Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.CTime = now
	article.UTime = now
	// insert - not upsert!
	err := d.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

func (d *ArticleGormDAO) GetByAuthorId(ctx context.Context, authorId int64, page int, pageSize int) ([]*Article, error) {
	var arts []*Article
	err := d.db.WithContext(ctx).Where("author_id = ?", authorId).Find(&arts).Error
	return arts, err
}

func (d *ArticleGormDAO) Sync(ctx context.Context, art *Article) (int64, error) {
	// check art's status
	// change the status to "published"
	// transaction
	var id = art.Id
	art.Status = 1
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		dao := NewArticleGormDAO(tx)
		if id > 0 {
			//dao.UpdateById(ctx, id)
		} else {
			id, err := dao.Insert(ctx, art)
			if err != nil {
				return err
			}
			art.Id = id
		}
		return nil
	})
	return id, err
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
