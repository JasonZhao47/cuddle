package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	groupID = "testing_group"
)

type TestingHandler struct{}

func TestSyncConsumer(t *testing.T) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	_, err := sarama.NewConsumerGroup(addr, groupID, config)
	assert.NoError(t, err)
	//cg.Consume(context.Background(), "testing_topics", handler)
}
