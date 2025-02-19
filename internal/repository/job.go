package repository

import (
	"context"
	"github.com/Andras5014/gohub/internal/domain"
	"github.com/Andras5014/gohub/internal/repository/dao"
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

func NewJobRepository(dao dao.JobDAO) JobRepository {
	return &jobRepository{
		dao: dao,
	}
}

func (j *jobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := j.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		Id:             job.Id,
		Cfg:            job.Cfg,
		CronExpression: job.CronExpression,
		Executor:       job.Executor,
		Name:           job.Name,
	}, nil
}

func (j *jobRepository) Release(ctx context.Context, id int64) error {
	return j.dao.Release(ctx, id)
}
func (j *jobRepository) Stop(ctx context.Context, id int64) error {
	return j.dao.Stop(ctx, id)
}

func (j *jobRepository) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return j.dao.UpdateNextTime(ctx, id, next)
}

func (j *jobRepository) UpdateUpdatedAt(ctx context.Context, id int64) error {
	return j.dao.UpdateUpdatedAt(ctx, id)
}
