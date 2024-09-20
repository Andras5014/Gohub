package logx

import "context"

type Logger interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)

	WithCtx(ctx context.Context) Logger
	With(field Field) Logger
}
type Field struct {
	Key   string
	Value any
}

func Any(key string, value any) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

func Error(err error) Field {
	return Any("error", err)
}
