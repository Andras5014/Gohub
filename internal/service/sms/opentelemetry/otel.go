package opentelemetry

import (
	"context"
	"github.com/Andras5014/gohub/internal/service/sms"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	svc    sms.Service
	tracer trace.Tracer
}

func NewService(svc sms.Service) sms.Service {
	tp := otel.GetTracerProvider()
	tracer := tp.Tracer("gohub/internal/service/sms/opentelemetry")
	return &Service{
		svc:    svc,
		tracer: tracer}
}
func (s *Service) Send(ctx context.Context, tplId string, args []sms.NamedArg, numbers ...string) error {
	ctx, span := s.tracer.Start(ctx, "sms_send_"+tplId, trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	err := s.svc.Send(ctx, tplId, args, numbers...)
	if err != nil {
		span.RecordError(ctx.Err())
	}
	return nil
}
