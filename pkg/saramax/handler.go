package saramax

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jasonzhao47/cuddle/internal/logger"
	"time"
)

const batchSize = 10

type Handler struct {
	fn func(msgs []*sarama.ConsumerMessage) error
	l  logger.Logger
}

func (h *Handler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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
		err := h.fn(batch)
		if err != nil {
			h.l.Error("处理消息失败",
				logger.Error(err))
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
