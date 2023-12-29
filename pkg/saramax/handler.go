package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	logger2 "github.com/jasonzhao47/cuddle/pkg/logger"
)

type Handler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, event T) error
	l  logger2.Logger
}

func NewHandler[T any](
	fn func(msg *sarama.ConsumerMessage, event T) error,
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
	// 循环消费接收到的msg
	messages := claim.Messages()
	select {
	case msg := <-messages:
		// 记录正在处理什么消息
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			h.l.Error("反序列消息体失败",
				logger2.String("topic", msg.Topic),
				logger2.Int32("partition", msg.Partition),
				logger2.Int64("offset", msg.Offset),
				logger2.Error(err),
			)
		}

		// 处理消息，如果有问题，报错
		err = h.fn(msg, t)
		if err != nil {
			h.l.Error("消费消息失败",
				logger2.String("topic", msg.Topic),
				logger2.Int32("partition", msg.Partition),
				logger2.Int64("offset", msg.Offset),
				logger2.Error(err),
			)
		}

		// 标记已经消费 - 这个流程是什么样子的？
		session.MarkMessage(msg, "")
	}
	return nil
}
