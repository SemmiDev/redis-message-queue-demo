package rmq

import "context"

import (
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendNewArticleNotificationEmail(
		ctx context.Context,
		payload *PayloadNewArticleNotificationEmail,
		opts ...asynq.Option,
	) error

	DistributeTaskSendNewSubscriberWelcomeEmail(
		ctx context.Context,
		payload *PayloadNewSubscriberWelcomeEmail,
		opts ...asynq.Option,
	) error

	DistributeTaskTrackReadArticle(
		ctx context.Context,
		payload *PayloadTrackReadArticle,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
