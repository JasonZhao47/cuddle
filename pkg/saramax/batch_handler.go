package saramax

import (
	"context"
	"github.com/IBM/sarama"
	logger2 "github.com/jasonzhao47/cuddle/pkg/logger"
	"time"
)

const batchSize = 10

type BatchHandler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, event T) error
	l  logger2.Logger
}

func NewBatchHandler[T any](
	fn func(msgs []*sarama.ConsumerMessage, event T) error,
	l logger2.Logger) *BatchHandler[T] {
	return &BatchHandler[T]{
		fn: fn,
		l:  l,
	}
}

func (h *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				batch = append(batch, msg)
			}
		}
		cancel()

		var event T

		err := h.fn(batch, event)
		if err != nil {
			h.l.Error("处理消息失败", logger2.Error(err))
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
