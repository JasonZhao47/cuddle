package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

const TopicReadEvent = "article_read"

type Producer interface {
	ProduceReadEvent(event ReadEvent) error
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{producer: producer}
}

func (r *SaramaSyncProducer) ProduceReadEvent(event ReadEvent) error {
	val, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, _, err = r.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicReadEvent,
		Key:   sarama.ByteEncoder(val),
	})
	return err
}
