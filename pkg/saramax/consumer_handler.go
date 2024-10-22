package saramax

import (
	"encoding/json"
	"errors"
	"github.com/Andras5014/webook/pkg/logx"
	"github.com/IBM/sarama"
)

type Handler[T any] struct {
	l  logx.Logger
	fn func(msg *sarama.ConsumerMessage, t T) error
}

func NewHandler[T any](l logx.Logger, fn func(msg *sarama.ConsumerMessage, t T) error) *Handler[T] {
	return &Handler[T]{
		l:  l,
		fn: fn,
	}
}
func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			h.l.Error("反序列化失败", logx.Error(err),
				logx.Any("topic", msg.Topic),
				logx.Any("partition", msg.Partition),
				logx.Any("offset", msg.Offset))
			continue
		}
		for i := 0; i < 3; i++ {
			err = h.fn(msg, t)
			if errors.Is(err, nil) {
				break
			}
			h.l.Error("消费失败", logx.Any("topic", msg.Topic),
				logx.Any("partition", msg.Partition),
				logx.Any("offset", msg.Offset),
				logx.Any("retry", i),
				logx.Error(err))
		}
		if err != nil {
			h.l.Error("重试失败", logx.Any("topic", msg.Topic),
				logx.Any("partition", msg.Partition),
				logx.Any("offset", msg.Offset),
				logx.Error(err))
		} else {
			session.MarkMessage(msg, "")
		}

	}
	return nil
}
