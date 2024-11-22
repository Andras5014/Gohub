package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id             int64
	Name           string
	Executor       string
	CronExpression string
	Cfg            string

	CancelFunc func() error
}

func (j Job) NextTime() time.Time {
	c := cron.NewParser(cron.Second | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := c.Parse(j.CronExpression)
	return s.Next(time.Now())
}
