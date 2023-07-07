package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

// distribute for client
func (distribute *RedisTaskDistributor) DistributeTastSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opt ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opt...)
	info, err := distribute.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w ", err)
	}
	log.Info().
		Str("type", task.Type()).
		Int("max_retry", info.MaxRetry).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Msg("Enqueue Task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to un marhsal payload from task %w ", err)
	}
	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed is while getting user %w", err)
	}
	// send email
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: utils.RandomString(32),
	})
	if err != nil {
		return err
	}
	subject := "i send mail from task send email"
	verifyURL := fmt.Sprintf("http://127.0.0.1:8888/v1/verify_email?secret_code=%s&email_id=%d", verifyEmail.SecretCode, verifyEmail.ID)

	to := []string{user.Email}
	content := fmt.Sprintf(`
		Hello %s, </br>	
		Thanks for registrering with us! </br>
		Please <a href="%s">click here </a> to verify email</br>
		<h1>test head 1</h1></br><a href=""> text</a>
	`, user.FullName, verifyURL)
	err = processor.mailer.SendEmail(
		subject,
		content,
		to,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("process Task")
	return nil
}
