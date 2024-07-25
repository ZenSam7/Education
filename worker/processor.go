package worker

import (
	"context"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
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
			// Чтобы было удобно парсить логи (и смотреть на них) реализовал их в мой tools.Log
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("ошибка в task")
			}),
			Logger: NewLogger(),
		}),
		queries: queries,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendGetUser, p.ProcessTaskSendVerifyEmail)
	// TODO: добавить какие-нибудь ещё функции

	return p.server.Start(mux)
}
