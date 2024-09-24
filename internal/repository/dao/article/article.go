package article

import (
	"context"
)

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	SyncV1(ctx context.Context, article Article) (int64, error)
	SyncStatus(ctx context.Context, article Article) (int64, error)
	FindByAuthorId(Dao context.Context, id int64, offset int, limit int) ([]Article, error)
	GetById(ctx context.Context, id int64) (Article, error)
	GetPubById(ctx context.Context, id int64) (PublishedArticle, error)
}
