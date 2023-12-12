package event

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jasonzhao47/cuddle/internal/domain/event/article"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"time"
)

type Consumer interface {
	Start(msg *sarama.ConsumerMessage, events article.ReadEvent) error
	BatchStart(msgs []*sarama.ConsumerMessage, events []article.ReadEvent) error
}

type ReadEventConsumer struct {
	repo repository.UserActivityRepository
}

func NewReadEventConsumer(repo repository.UserActivityRepository) Consumer {
	return &ReadEventConsumer{repo: repo}
}

func (r *ReadEventConsumer) Start(msg *sarama.ConsumerMessage, event article.ReadEvent) error {
	// 开始消费
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	err := r.repo.IncrRead(ctx, "article", event.Aid)
	return err
}

func (r *ReadEventConsumer) BatchStart(msgs []*sarama.ConsumerMessage, events []article.ReadEvent) error {
	bizs := make([]string, 0, len(events))
	bizIds := make([]int64, 0, len(events))
	for _, evt := range events {
		bizs = append(bizs, "article")
		bizIds = append(bizIds, evt.Aid)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.repo.BatchIncrRead(ctx, bizs, bizIds)
}
