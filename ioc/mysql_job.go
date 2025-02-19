package ioc

import (
	"context"
	"github.com/Andras5014/gohub/internal/domain"
	"github.com/Andras5014/gohub/internal/job"
	"github.com/Andras5014/gohub/internal/service"
	"github.com/Andras5014/gohub/pkg/logx"
	"time"
)

func InitScheduler(local *job.LocalFuncExecutor, svc service.JobService, l logx.Logger) *job.Scheduler {
	res := job.NewScheduler(svc, l)
	res.RegisterExecutor(local)
	return res
}

func InitLocalFuncExecutor(svc service.RankingService) *job.LocalFuncExecutor {
	res := job.NewLocalFuncExecutor()
	res.RegisterFunc("ranking", func(ctx context.Context, j domain.Job) error {
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		return svc.TopN(ctx, 100)

	})
	return res
}
