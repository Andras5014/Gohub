package article

import (
	"context"
	"github.com/IBM/sarama"
	"testing"
)

func TestSend(t *testing.T) {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9094"}, nil)
	if err != nil {
		t.Error(err)
	}
	KafkaPro := NewSaramaSyncProducer(producer)
	err = KafkaPro.ProduceReadEvent(context.Background(), ReadEvent{
		ArticleId: 1,
		UserId:    1,
	})
	if err != nil {
		t.Error(err)
	}
}
func TestBatchSend(t *testing.T) {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9094"}, nil)
	if err != nil {
		t.Error(err)
	}
	KafkaPro := NewSaramaSyncProducer(producer)
	for i := 0; i < 10; i++ {
		err = KafkaPro.ProduceReadEvent(context.Background(), ReadEvent{
			ArticleId: 1,
			UserId:    1,
		})
		if err != nil {
			t.Error(err)
		}
	}
	if err != nil {
		t.Error(err)
	}
}
