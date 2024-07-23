package worker

import (
	"context"
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskGetUser(ctx context.Context, payload *PayloadSendGetUser, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOptions asynq.RedisClientOpt) TaskDistributor {
	return &RedisTaskDistributor{client: asynq.NewClient(redisOptions)}
}
