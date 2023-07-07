package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTastSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opt ...asynq.Option,
	) error
}
type RedisTaskDistributor struct {
	client *asynq.Client
}

// declared interface but return struct  to implement func from interface for struct
func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
