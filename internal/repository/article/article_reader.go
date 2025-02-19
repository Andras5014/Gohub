package article

import (
	"context"
	"github.com/Andras5014/gohub/internal/domain"
)

type ReaderRepository interface {
	Save(ctx context.Context, article domain.Article) error
}
