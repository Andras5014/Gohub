package repository

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	"time"
)

type JobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUpdatedAt(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Stop(ctx context.Context, id int64) error
}

type jobRepository struct {
	dao dao.JobDAO
}
