package ioc

import (
	"github.com/Andras5014/webook/interactive/config"
	"github.com/Andras5014/webook/interactive/events"
	"github.com/Andras5014/webook/pkg/saramax"
	"github.com/IBM/sarama"
)

func InitKafka(cfg *config.Config) sarama.Client {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Kafka.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSyncProducer(client sarama.Client) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

func InitConsumers(c *events.InteractiveReadEventBatchConsumer) []saramax.Consumer {
	return []saramax.Consumer{c}
}
