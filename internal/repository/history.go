package repository

import (
	"context"
	"github.com/Andras5014/gohub/internal/domain"
)

type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}
