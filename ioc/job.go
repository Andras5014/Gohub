package ioc

import (
	"github.com/Andras5014/gohub/internal/job"
	"github.com/Andras5014/gohub/internal/service"
	"github.com/Andras5014/gohub/pkg/logx"
	"github.com/go-redsync/redsync/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"time"
)

func InitRankingJob(svc service.RankingService, lock *redsync.Redsync, l logx.Logger) *job.RankingJob {
	mu := lock.NewMutex("ranking_job_lock", redsync.WithExpiry(time.Minute*5))
	return job.NewRankingJob(svc, time.Second*30, mu, l)
}

func InitJobs(rankingJob *job.RankingJob, l logx.Logger) *cron.Cron {
	builder := job.NewCronJobBuilder(prometheus.SummaryOpts{
		Namespace: "echohub",
		Subsystem: "job",
		Name:      "cron_job",
		Help:      "cron job exec info",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, l)
	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 10s", builder.Build(rankingJob))
	if err != nil {
		panic(err)
	}
	return expr
}
