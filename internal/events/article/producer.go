package article

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const TopicReadEvent = "article_read"

type Producer interface {
	ProduceReadEvent(ctx context.Context, event ReadEvent) error
}
type KafkaProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &KafkaProducer{
		producer: producer,
		topic:    TopicReadEvent,
	}
}

// ProduceReadEvent 装饰器扩展重试机制
func (k *KafkaProducer) ProduceReadEvent(ctx context.Context, event ReadEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}

type ReadEvent struct {
	UserId    int64 `json:"user_id"`
	ArticleId int64 `json:"article_id"`
}
