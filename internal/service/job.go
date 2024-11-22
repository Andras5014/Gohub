package service

import (
	"context"
	"github.com/Andras5014/webook/internal/domain"
	"github.com/Andras5014/webook/internal/repository"
	"github.com/Andras5014/webook/pkg/logx"
	"time"
)

type JobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, j domain.Job) error
	ResetNextTime(ctx context.Context, j domain.Job) error
	Stop(ctx context.Context, j domain.Job) error
}

type CronJobService struct {
	repo            repository.JobRepository
	refreshInterval time.Duration
	l               logx.Logger
}
