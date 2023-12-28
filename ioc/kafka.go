package ioc

import (
	"github.com/IBM/sarama"
	"github.com/jasonzhao47/cuddle/internal/domain/event"
	"github.com/spf13/viper"
)

func InitSaramaClient() sarama.Client {
	type Config struct {
		Addr []string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, scfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return p
}

func InitConsumers(c1 *event.UserActivityEventConsumer) []event.Consumer {
	return []event.Consumer{c1}
}
