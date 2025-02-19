package service

import (
	"context"
	"github.com/Andras5014/gohub/internal/domain"
	"github.com/Andras5014/gohub/internal/repository"
	"github.com/Andras5014/gohub/pkg/logx"
	"time"
)

type JobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, id int64) error
	ResetNextTime(ctx context.Context, j domain.Job) error
	Stop(ctx context.Context, j domain.Job) error
}

type CronJobService struct {
	repo            repository.JobRepository
	refreshInterval time.Duration
	l               logx.Logger
}

func (c *CronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}

	ticker := time.NewTicker(c.refreshInterval)

	go func() {
		// 续约
		for range ticker.C {
			c.refresh(j.Id)
		}
	}()

	j.CancelFunc = func() error {
		ticker.Stop()
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		return c.repo.Release(ctx, j.Id)
	}
	return j, nil
}

func (c *CronJobService) Release(ctx context.Context, id int64) error {
	return c.repo.Release(ctx, id)
}
func (c *CronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.repo.UpdateUpdatedAt(ctx, id)
	if err != nil {
		c.l.Error("续约失败", logx.Error(err), logx.Int64("jid", id))
	}
}

func (c *CronJobService) ResetNextTime(ctx context.Context, j domain.Job) error {
	next := j.NextTime()
	if next.IsZero() {
		return nil
	}
	return c.repo.UpdateNextTime(ctx, j.Id, next)
}

func (c *CronJobService) Stop(ctx context.Context, j domain.Job) error {
	return c.repo.Stop(ctx, j.Id)
}
