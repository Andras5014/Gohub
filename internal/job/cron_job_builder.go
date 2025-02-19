package job

import (
	"context"
	"github.com/Andras5014/gohub/pkg/logx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type CronJobBuilder struct {
	p      *prometheus.SummaryVec
	tracer trace.Tracer
	l      logx.Logger
}

func NewCronJobBuilder(opt prometheus.SummaryOpts, l logx.Logger) *CronJobBuilder {
	p := prometheus.NewSummaryVec(opt, []string{"name"})
	prometheus.MustRegister(p)
	return &CronJobBuilder{
		p:      p,
		l:      l,
		tracer: otel.GetTracerProvider().Tracer("gohub/internal/job"),
	}
}

func (c *CronJobBuilder) Build(j Job) cron.Job {
	name := j.Name()
	return cronJobAdapter(func() {
		ctx, span := c.tracer.Start(context.Background(), name)
		start := time.Now()
		c.l.Info("开始执行任务", logx.String("name", name))

		defer func() {
			span.End()
			duration := time.Since(start)
			c.l.Info("任务执行完毕", logx.String("name", name), logx.Duration("cost", duration))
			c.p.WithLabelValues(name).Observe(float64(duration.Milliseconds()))
		}()

		err := j.Run(ctx)
		if err != nil {
			span.RecordError(err)
			c.l.Error("任务执行失败", logx.String("name", name), logx.Error(err))
		}

	})

}

type cronJobAdapter func()

func (c cronJobAdapter) Run() {
	c()
}
