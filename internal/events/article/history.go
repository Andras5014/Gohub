package article

import (
	"context"
	"github.com/Andras5014/gohub/internal/domain"
	"github.com/Andras5014/gohub/internal/repository"
	"github.com/Andras5014/gohub/pkg/logx"
	"github.com/Andras5014/gohub/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

type HistoryRecordConsumer struct {
	repo   repository.HistoryRecordRepository
	client sarama.Client
	l      logx.Logger
}

func (i *HistoryRecordConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("history_record", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicReadEvent},
			saramax.NewHandler[ReadEvent](i.l, i.Consume))
		if er != nil {
			i.l.Error("退出消费", logx.Error(er))
		}
	}()
	return err
}

func (i *HistoryRecordConsumer) Consume(msg *sarama.ConsumerMessage,
	event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.AddRecord(ctx, domain.HistoryRecord{
		BizId: event.ArticleId,
		Biz:   "article",
		Uid:   event.UserId,
	})
}
