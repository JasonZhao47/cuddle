package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	GetById(context.Context, int64) (*Article, error)
	Insert(context.Context, *Article) (int64, error)
	GetByAuthorId(context.Context, int64, int, int) ([]*Article, error)
	Sync(context.Context, *Article) (int64, error)
	UpdateById(context.Context, *Article) error
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
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// upsert material db first
		dao := NewArticleGormDAO(tx)
		if id > 0 {
			err := dao.UpdateById(ctx, art)
			if err != nil {
				return err
			}
		} else {
			id, err := dao.Insert(ctx, art)
			if err != nil {
				return err
			}
			art.Id = id
		}
		now := time.Now().UnixMilli()
		pubArt := PublishedArticle{
			Article: *art,
		}
		pubArt.CTime = now
		pubArt.UTime = now
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"topic":   pubArt.Topic,
				"content": pubArt.Content,
				"status":  pubArt.Status,
				"utime":   now,
			}),
		}).Create(&pubArt).Error
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (d *ArticleGormDAO) UpdateById(ctx context.Context, art *Article) error {
	// update
	// do we need a lock?
	// pessimistic lock
	// locks everything during operation
	// FOR UPDATE
	// or checks if there's a concurrency happening
	// how?
	now := time.Now().UnixMilli()
	res := d.db.WithContext(ctx).Model(&art).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"topic":   art.Topic,
			"content": art.Content,
			"status":  art.Status,
			"utime":   now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		// possible attacker
		return errors.New("ID不对或者作者不对")
	}
	return nil
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

type PublishedArticle struct {
	Article
}
