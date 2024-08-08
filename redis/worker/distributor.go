package worker

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// TaskDistributor Создаём и (асинхронно) выполняем все эти задачи через брокер сообщений
type TaskDistributor interface {
	DistributeTaskVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error
}

// RedisTaskDistributor Реализуем создание и асинхронное выполнение всех задач (используем редиску как брокер сообщений)
type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(opt redis.Options) TaskDistributor {
	asynqOption := asynq.RedisClientOpt{
		Addr: opt.Addr,
	}

	return &RedisTaskDistributor{client: asynq.NewClient(asynqOption)}
}
