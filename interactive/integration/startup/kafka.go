package startup

import (
	"github.com/Andras5014/gohub/config"
	"github.com/Andras5014/gohub/internal/events"
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

func InitConsumers(c events.Consumer) []events.Consumer {
	return []events.Consumer{c}
}
