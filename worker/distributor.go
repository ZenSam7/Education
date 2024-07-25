package worker

import (
	"context"
	"github.com/hibiken/asynq"
)

// TaskDistributor Создаём и (асинхронно) выполняем все эти задачи через брокер сообщений
type TaskDistributor interface {
	DistributeTaskVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error
}

// RedisTaskDistributor Реализуем создание и асинхронное выполнение всех задач (используем редиску как брокер сообщений)
type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOptions asynq.RedisClientOpt) TaskDistributor {
	return &RedisTaskDistributor{client: asynq.NewClient(redisOptions)}
}
