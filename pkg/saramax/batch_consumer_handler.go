package saramax

import (
	"context"
	"encoding/json"
	"github.com/Andras5014/webook/pkg/logx"
	"github.com/IBM/sarama"
	"time"
)

const BatchSize = 10

type BatchHandler[T any] struct {
	l         logx.Logger
	fn        func(msg []*sarama.ConsumerMessage, ts []T) error
	batchSize int
}

func NewBatchHandler[T any](l logx.Logger, fn func(msg []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{
		l:         l,
		fn:        fn,
		batchSize: BatchSize,
	}
}
func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for {
		batch := make([]*sarama.ConsumerMessage, 0, b.batchSize)
		ts := make([]T, 0, b.batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		done := false
		for i := 0; i < b.batchSize && !done; i++ {
			select {
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}

				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列化失败", logx.Error(err),
						logx.Any("topic", msg.Topic),
						logx.Any("partition", msg.Partition),
						logx.Any("offset", msg.Offset))
					continue
				}
				batch = append(batch, msg)
				ts = append(ts, t)
			case <-ctx.Done():
				//超时
				done = true
			}
		}
		cancel()
		// 凑够batch size 就进行处理
		err := b.fn(batch, ts)
		if err != nil {
			b.l.Error("处理失败", logx.Error(err))
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
