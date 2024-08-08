package worker

import (
	"context"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/tools"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server  *asynq.Server
	queries *db.Queries
	sender  tools.EmailSender
	config  tools.Config
}

func NewRedisTaskProcessor(opt redis.Options, queries *db.Queries) TaskProcessor {
	cnfg := tools.LoadConfig()

	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: opt.Addr,
		},

		asynq.Config{
			// Queues Важные задачи распределяем по отдельным потокам (цифра = степень важности)
			Queues: map[string]int{
				QueueCritical: 9,
				QueueDefault:  3,
				QueueLow:      1,
			},
			// Чтобы было удобно парсить логи (и смотреть на них) реализовал их в tools.Log
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("ошибка в таске")
			}),
			Logger: NewWorkerLogger(),
		},
	)

	return &RedisTaskProcessor{
		queries: queries,
		sender:  &tools.GmailSender{Config: cnfg},
		config:  cnfg,
		server:  server,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)

	return p.server.Start(mux)
}

func (p *RedisTaskProcessor) Shutdown() {
	p.server.Shutdown()
}
