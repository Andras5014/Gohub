package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
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

type Job struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Cfg      string `gorm:"type:json"`
	Name     string `gorm:"type:varchar(100),uniqueIndex"`
	Executor string `gorm:"type:varchar(100)"`
	Status   int    `gorm:"index:idx_status_next_time"`
	// NextTime 下一次执行时间
	NextTime  int64 `gorm:"index:idx_status_next_time"`
	Version   int
	CreatedAt int64
	UpdatedAt int64
}
