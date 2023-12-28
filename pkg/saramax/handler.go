package saramax

import (
	"github.com/IBM/sarama"
	logger2 "github.com/jasonzhao47/cuddle/pkg/logger"
)

type Handler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, event T) error
	l  logger2.Logger
}

func NewHandler[T any](
	fn func(msgs []*sarama.ConsumerMessage, event T) error,
	l logger2.Logger) *Handler[T] {
	return &Handler[T]{
		fn: fn,
		l:  l,
	}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return nil
}
