package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendGetUser = "task:send_get_user"

type PayloadSendVerifyEmail struct {
	IdUser int32 `json:"id_user"`
}

func (d *RedisTaskDistributor) DistributeTaskVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	bytePayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("не получилось конвертировать payload в byte: %w", err)
	}

	task := asynq.NewTask(TaskSendGetUser, bytePayload, opts...)
	info, err := d.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("не получилось создать задачу: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", bytePayload).
		Str("queue", info.Queue).
		Msg("enqueued GetUser")

	return nil
}

// ProcessTaskSendVerifyEmail Обрабатываем в брокере сообщений (в редиске) задачу на верифицирование почты при создании пользователя
func (p *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		// Если у нас не получается достать информацию, то не даём повторить запрос
		return fmt.Errorf("не получилось конвертировать byte в payload: %w", asynq.SkipRetry)
	}

	user, err := p.queries.GetUser(ctx, payload.IdUser)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если пользователя нету, то не даём повторить запрос
			return fmt.Errorf("пользователь не найден")
		}
		return fmt.Errorf("не удалось получить пользователя: %w", err)
	}

	// TODO: сделать отправку сообщений для верификации почты

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Int("user", int(user.IDUser)).
		Msg("processed GetUser")

	return nil
}
