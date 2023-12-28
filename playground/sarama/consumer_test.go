package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	groupID = "testing_group"
)

func TestSyncConsumer(t *testing.T) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	_, err := sarama.NewConsumerGroup(addr, groupID, config)
	assert.NoError(t, err)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	//err = cg.Consume(ctx, "testing_topic")
}
