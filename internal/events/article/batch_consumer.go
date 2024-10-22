package article

import (
	"context"
	"github.com/Andras5014/webook/internal/events"
	"github.com/Andras5014/webook/internal/repository"
	"github.com/Andras5014/webook/pkg/logx"
	"github.com/Andras5014/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

type InteractiveReadEventBatchConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logx.Logger
}

func NewInteractiveReadEventBatchConsumer(client sarama.Client, repo repository.InteractiveRepository, l logx.Logger) events.Consumer {
	return &InteractiveReadEventBatchConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}

func (i *InteractiveReadEventBatchConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{"article_read"}, saramax.NewBatchHandler[ReadEvent](i.l, i.BatchConsume))
		if err != nil {
			i.l.Error("消费消息失败", logx.Error(err))
		}
	}()
	return err
}

func (i *InteractiveReadEventBatchConsumer) BatchConsume(msg []*sarama.ConsumerMessage, events []ReadEvent) error {
	bizs := make([]string, 0, len(events))
	bizIds := make([]int64, 0, len(events))
	for _, event := range events {
		bizs = append(bizs, "article")
		bizIds = append(bizIds, event.ArticleId)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.BatchIncrReadCnt(ctx, bizs, bizIds)
}
