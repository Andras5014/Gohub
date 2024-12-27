package events

import (
	"context"
	"github.com/Andras5014/webook/interactive/repository"
	"github.com/Andras5014/webook/pkg/logx"
	"github.com/Andras5014/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

const ReadEventConsumerGroup = "article_read_event"
const TopicReadEvent = "article_read"

type InteractiveReadEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logx.Logger
}

func NewInteractiveReadEventConsumer(client sarama.Client, repo repository.InteractiveRepository, l logx.Logger) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}
func (k *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(ReadEventConsumerGroup, k.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewHandler[ReadEvent](k.l, k.Consume))
		if err != nil {
			k.l.Error("消费消息失败", logx.Error(err))
		}
	}()
	return err
}

func (k *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return k.repo.IncrReadCnt(ctx, "article", t.ArticleId)
}
