package worker

import (
	"context"

	"github.com/rs/zerolog/log"

	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/email"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(
		ctx context.Context,
		task *asynq.Task,
	) error
}
type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer email.Emailsender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer email.Emailsender) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("proccess task failed")
			}),
			Logger: NewLogger(),
		},
	)
	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
