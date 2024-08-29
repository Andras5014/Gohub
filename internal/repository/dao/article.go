package dao

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

type GormArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GormArticleDAO{db: db}
}
func (g *GormArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := sql.NullInt64{Int64: time.Now().UnixMilli(), Valid: true}
	article.CreatedAt = now
	article.UpdatedAt = now
	return article.Id, g.db.WithContext(ctx).Create(&article).Error
}

func (g *GormArticleDAO) UpdateById(ctx context.Context, article Article) error {
	article.UpdatedAt = sql.NullInt64{Int64: time.Now().UnixMilli(), Valid: true}
	res := g.db.WithContext(ctx).Model(&article).Where("id = ? And author_id = ?", article.Id, article.AuthorId).Updates(&article)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("更新失败，可能是非法操作")
	}
	return nil
}

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"not null"`
	Content string `gorm:"type=BLOB"`

	AuthorId int64 `gorm:"index"`

	CreatedAt sql.NullInt64
	UpdatedAt sql.NullInt64
	DeletedAt sql.NullInt64
}
