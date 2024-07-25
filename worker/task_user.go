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

type PayloadSendGetUser struct {
	IdUser int32 `json:"id_user"`
}

func (d *RedisTaskDistributor) DistributeTaskGetUser(ctx context.Context, payload *PayloadSendGetUser, opts ...asynq.Option) error {
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

func (p *RedisTaskProcessor) ProcessTaskSendGetUser(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendGetUser
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("не получилось конвертировать byte в payload: %w", asynq.SkipRetry)
	}

	user, err := p.queries.GetUser(ctx, payload.IdUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("пользователь не найден")
		}
		return fmt.Errorf("не удалось получить пользователя: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Int("user", int(user.IDUser)).
		Msg("processed GetUser")

	return nil
}

// TODO: добавить все остальные функции для user
