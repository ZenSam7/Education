package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/tools"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	IdUser int32 `json:"id_user"`
}

func (d *RedisTaskDistributor) DistributeTaskVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	bytePayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("не получилось конвертировать payload в byte: %s", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, bytePayload, opts...)
	info, err := d.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("не получилось создать задачу: %s", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", bytePayload).
		Str("queue", info.Queue).
		Msg("enqueued VerifyEmail")

	return nil
}

// ProcessTaskSendVerifyEmail Обрабатываем в брокере сообщений (в редиске) задачу на верифицирование почты при создании пользователя
func (p *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		// Если у нас не получается достать информацию, то не даём повторить запрос
		return fmt.Errorf("не получилось конвертировать byte в payload: %s", asynq.SkipRetry)
	}

	user, err := p.queries.GetUser(ctx, payload.IdUser)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если пользователя нету, то не даём повторить запрос
			return fmt.Errorf("пользователь не найден")
		}
		return fmt.Errorf("не удалось получить пользователя: %s", err)
	}

	// Если запрос на верификацию уже сделан (т.е. пользователь зачем-то отправяет повторный
	// запрос), у нас несколько вариантов действий:
	lastVerifyRequest, err := p.queries.GetVerifyRequest(ctx, user.IDUser)
	// запрос уже есть
	if err == nil {
		// Срок действия запроса истёк, удаляем
		if time.Now().After(lastVerifyRequest.ExpiredAt.Time) {
			_, err = p.queries.DeleteVerifyRequest(ctx, user.IDUser)
			if err != nil {
				return err
			}
		} else {
			// Не истёк
			return fmt.Errorf("запрос на верификацию почты уже отправлен")
		}
	} else {
		// запроса нету, т.е. это первый
		if err.Error() == "no rows in result set" {
			// продолжаем дальше работать
		} else {
			// ошибка при запросе
			return err
		}
	}

	verifyEmail, err := p.queries.CreateVerifyRequest(ctx, db.CreateVerifyRequestParams{
		IDUser:    user.IDUser,
		SecretKey: tools.GetRandomString(32),
	})
	if err != nil {
		return fmt.Errorf("не получилось создать запрос на подтверждение почты: %s", err)
	}

	// Отправляем сообщение
	err = p.sender.SendMail(user.Email, "email_verify.html", map[string]string{
		"{{ verification_url }}": fmt.Sprintf(
			"http://localhost:%s/verify_email?id_user=%d&secret_key=%s",
			strings.Split(p.config.HttpServerAddress, ":")[1], // Порт
			verifyEmail.IDUser,
			verifyEmail.SecretKey,
		),
	})
	if err != nil {
		return fmt.Errorf("не удалось отправить сообщение: %s", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Int("user", int(user.IDUser)).
		Msg("processed VerifyEmail")

	return nil
}
