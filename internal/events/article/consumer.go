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

type InteractiveReadEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logx.Logger
}

func NewInteractiveReadEventConsumer(client sarama.Client, repo repository.InteractiveRepository, l logx.Logger) events.Consumer {
	return &InteractiveReadEventConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}
func (k *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", k.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(), []string{"article_read"}, saramax.NewHandler[ReadEvent](k.l, k.Consume))
		if err != nil {
			k.l.Error("消费消息失败", logx.Error(err))
		}
	}()
	return err
}

func (k *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return k.repo.IncrReadCnt(ctx, "article", t.ArticleId, t.UserId)
}
