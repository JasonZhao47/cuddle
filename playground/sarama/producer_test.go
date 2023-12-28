package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

var addr = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	// 可以指定分区
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	producer, err := sarama.NewSyncProducer(addr, config)
	assert.NoError(t, err)

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "testing_topic",
		Value: sarama.StringEncoder("测试的消息"),
		Headers: []sarama.RecordHeader{
			sarama.RecordHeader{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
		Metadata: "测试用元数据",
	})
	assert.NoError(t, err)
}

func TestAsyncProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	// 可以指定分区
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	producer, err := sarama.NewAsyncProducer(addr, config)
	assert.NoError(t, err)

	msg := producer.Input()
	msg <- &sarama.ProducerMessage{
		Topic: "testing_topic",
		Value: sarama.StringEncoder("测试的消息"),
		Headers: []sarama.RecordHeader{
			sarama.RecordHeader{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
		Metadata: "测试用元数据",
	}
	select {
	case ms := <-producer.Successes():
		t.Log("发送成功", string(ms.Value.(sarama.StringEncoder)))
	case ms := <-producer.Errors():
		msgs, _ := ms.Msg.Value.Encode()
		t.Log("发送失败", string(msgs))
	}
}
