package article

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
)

type AuthorRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}
