package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

const (
	jobStatusWaiting = iota
	jobStatusRunning
	jobStatusPaused
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, id int64) error
	Stop(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	UpdateUpdatedAt(ctx context.Context, id int64) error
}

type GormJobDAO struct {
	db *gorm.DB
}

func NewJobDAO(db *gorm.DB) JobDAO {
	return &GormJobDAO{db: db}
}
func (g *GormJobDAO) Preempt(ctx context.Context) (Job, error) {
	for {
		now := time.Now()
		var j Job

		// todo 找失败的job执行
		err := g.db.WithContext(ctx).Where("status = ? and next_time <= ?", jobStatusWaiting, now.Unix()).First(&j).Error
		if err != nil {
			return Job{}, err
		}

		res := g.db.WithContext(ctx).Where("id = ? and version = ?", j.Id, j.Version).Updates(map[string]any{
			"status":     jobStatusRunning,
			"updated_at": now.UnixMilli(),
			"version":    j.Version + 1,
		})
		if res.Error != nil {
			return Job{}, res.Error
		}
		if res.RowsAffected == 0 {
			continue
		}
		return j, nil
	}
}

func (g *GormJobDAO) Release(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]any{
		"status":     jobStatusWaiting,
		"updated_at": time.Now().UnixMilli(),
	}).Error
}

func (g *GormJobDAO) Stop(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]any{
		"status":     jobStatusPaused,
		"updated_at": time.Now().UnixMilli(),
	}).Error
}

func (g *GormJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return g.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]any{
		"next_time":  next.UnixMilli(),
		"updated_at": time.Now().UnixMilli(),
	}).Error
}
func (g *GormJobDAO) UpdateUpdatedAt(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]any{
		"updated_at": time.Now().UnixMilli(),
	}).Error
}

type Job struct {
	Id             int64  `gorm:"primaryKey,autoIncrement"`
	Cfg            string `gorm:"type:json"`
	Name           string `gorm:"type:varchar(100),uniqueIndex"`
	Executor       string `gorm:"type:varchar(100)"`
	CronExpression string
	Status         int `gorm:"index:idx_status_next_time"`
	// NextTime 下一次执行时间
	NextTime  int64 `gorm:"index:idx_status_next_time"`
	Version   int
	CreatedAt int64
	UpdatedAt int64
}
