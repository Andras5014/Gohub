package metrics

import (
	"context"
	"github.com/Andras5014/gohub/internal/service/sms"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type PrometheusDecorator struct {
	svc    sms.Service
	vector *prometheus.SummaryVec
}

func NewPrometheusDecorator(svc sms.Service) sms.Service {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "andras",
		Subsystem: "gohub",
		Help:      "sms send duration in seconds",
		Name:      "sms_send_duration_seconds",
	}, []string{"biz"})
	prometheus.MustRegister(vector)
	return &PrometheusDecorator{
		svc:    svc,
		vector: vector,
	}
}

func (p *PrometheusDecorator) Send(ctx context.Context, tplToken string, args []sms.NamedArg, numbers ...string) error {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Milliseconds()
		p.vector.WithLabelValues(tplToken).Observe(float64(duration))
	}()
	return p.svc.Send(ctx, tplToken, args, numbers...)
}
