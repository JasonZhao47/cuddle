package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/pkg/logger"
	"github.com/jasonzhao47/cuddle/pkg/saramax"
	"time"
)

const topic = "user_activity"

type Consumer interface {
	Start() error
	BatchStart() error
	Consume(messages []*sarama.ConsumerMessage, event ReadEvent) error
}

type UserActivityEventConsumer struct {
	repo   repository.UserActivityRepository
	client sarama.Client
	l      logger.Logger
}

func NewReadEventConsumer(repo repository.UserActivityRepository, client sarama.Client, l logger.Logger) Consumer {
	return &UserActivityEventConsumer{
		repo:   repo,
		client: client,
		l:      l,
	}
}

func (u *UserActivityEventConsumer) Consume(messages []*sarama.ConsumerMessage, event ReadEvent) error {
	// 消费特定内容
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return u.repo.IncrRead(ctx, "article", event.Aid)
}

func (u *UserActivityEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(topic, u.client)
	if err != nil {
		return err
	}
	go func() {
		u.l.Info("开始消费", logger.String("topic", TopicReadEvent))
		// 插件和钩子
		err := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewHandler[ReadEvent](u.Consume, u.l))
		if err != nil {
			u.l.Error("消费错误", logger.Error(err))
		}
	}()
	return nil
}

func (u *UserActivityEventConsumer) BatchStart() error {
	cg, err := sarama.NewConsumerGroupFromClient(topic, u.client)
	if err != nil {
		return err
	}
	go func() {
		u.l.Info("开始消费", logger.String("topic", TopicReadEvent))
		// 插件和钩子
		err := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewBatchHandler[ReadEvent](u.Consume, u.l))
		if err != nil {
			u.l.Error("消费错误", logger.Error(err))
		}
	}()
	return nil
}
