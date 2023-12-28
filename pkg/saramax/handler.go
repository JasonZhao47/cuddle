package saramax

import (
	"github.com/IBM/sarama"
	logger2 "github.com/jasonzhao47/cuddle/pkg/logger"
)

type Handler struct {
	fn func(msgs []*sarama.ConsumerMessage) error
	l  logger2.Logger
}

func (h *Handler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return nil
}
