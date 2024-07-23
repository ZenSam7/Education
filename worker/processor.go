package worker

import (
	"context"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendGetUser(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server  *asynq.Server
	queries *db.Queries
}

func NewRedisTaskProcessor(opt asynq.RedisClientOpt, queries *db.Queries) TaskProcessor {
	return &RedisTaskProcessor{
		server: asynq.NewServer(opt, asynq.Config{
			// Queues Важные задачи распределяем по отдельным потокам (цифра = степень важности)
			Queues: map[string]int{
				QueueCritical: 9,
				QueueDefault:  3,
				QueueLow:      1,
			},
		}),
		queries: queries,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendGetUser, p.ProcessTaskSendGetUser)

	return p.server.Start(mux)
}
