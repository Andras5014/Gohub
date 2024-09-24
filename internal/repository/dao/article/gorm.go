package article

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GormArticleDAO struct {
	db *gorm.DB
}

func (g *GormArticleDAO) FindByAuthorId(Dao context.Context, id int64, offset int, limit int) ([]Article, error) {
	var arts []Article
	// orderby 命中索引
	err := g.db.WithContext(Dao).Where("author_id = ?", id).Limit(limit).Offset(offset).Order("updated_at DESC").Find(&arts).Error
	return arts, err
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GormArticleDAO{db: db}
}
func (g *GormArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.CreatedAt = now
	article.UpdatedAt = now
	return article.Id, g.db.WithContext(ctx).Create(&article).Error
}

func (g *GormArticleDAO) UpdateById(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.UpdatedAt = now
	res := g.db.WithContext(ctx).Model(&article).Where("id = ? And author_id = ?", article.Id, article.AuthorId).Updates(&article)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("更新失败，可能是非法操作")
	}
	return nil
}
func (g *GormArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {
	var id = article.Id
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)
		dao := NewArticleDAO(tx)
		if id > 0 {
			err = dao.UpdateById(ctx, article)
		} else {
			id, err = dao.Insert(ctx, article)
		}
		if err != nil {
			return err
		}
		article.Id = id
		now := time.Now().UnixMilli()
		pubArt := PublishedArticle(article)
		pubArt.CreatedAt = now
		pubArt.UpdatedAt = now
		err = tx.Clauses(clause.OnConflict{
			// 对MySQL不起效，但是可以兼容别的方言
			// INSERT xxx ON DUPLICATE KEY SET `title`=?
			// 别的方言：
			// sqlite INSERT XXX ON CONFLICT DO UPDATES WHERE
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":      pubArt.Title,
				"content":    pubArt.Content,
				"updated_at": now,
				"status":     pubArt.Status,
			}),
		}).Create(&pubArt).Error
		return err
	})
	return id, err
}

func (g *GormArticleDAO) SyncV1(ctx context.Context, art Article) (int64, error) {
	tx := g.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 防止后面业务panic
	defer tx.Rollback()

	var (
		id  = art.Id
		err error
	)
	dao := NewArticleDAO(tx)
	if id > 0 {
		err = dao.UpdateById(ctx, art)
	} else {
		id, err = dao.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	now := time.Now().UnixMilli()
	pubArt := PublishedArticle(art)
	pubArt.CreatedAt = now
	pubArt.UpdatedAt = now
	err = tx.Clauses(clause.OnConflict{
		// 对MySQL不起效，但是可以兼容别的方言
		// INSERT xxx ON DUPLICATE KEY SET `title`=?
		// 别的方言：
		// sqlite INSERT XXX ON CONFLICT DO UPDATES WHERE
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":      pubArt.Title,
			"content":    pubArt.Content,
			"updated_at": now,
		}),
	}).Create(&pubArt).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (g *GormArticleDAO) SyncStatus(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.UpdatedAt = now
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&article).Where("id = ? And author_id = ?", article.Id, article.AuthorId).Updates(&article)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("更新失败，可能是非法操作")
		}
		return tx.Model(&PublishedArticle{}).Where("id = ?", article.Id).Updates(&article).Error
	})
	return article.Id, err
}
func (g *GormArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	var article Article
	err := g.db.WithContext(ctx).Where("id = ?", id).First(&article).Error
	return article, err
}

func (g *GormArticleDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	var pubArt PublishedArticle
	err := g.db.WithContext(ctx).Where("id = ?", id).First(&pubArt).Error
	return pubArt, err
}
