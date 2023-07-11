package rmq

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/techschool/simplebank/mail"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendNewSubscriberWelcomeEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskSendNewArticleNotificationEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskTrackReadArticle(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server                    *asynq.Server
	articleRepo               ArticleRepository
	articleSubscriberRepo     ArticleSubscriberRepository
	articleAnalyticRepository ArticleAnalyticRepository
	mailer                    mail.EmailSender
}

func NewRedisTaskProcessor(
	redisOpt asynq.RedisClientOpt,
	articleRepo ArticleRepository,
	articleSubscriberRepo ArticleSubscriberRepository,
	articleAnalyticRepository ArticleAnalyticRepository,
	mailer mail.EmailSender,
) TaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 25,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: logger,
		},
	)

	return &RedisTaskProcessor{
		server:                    server,
		articleRepo:               articleRepo,
		articleSubscriberRepo:     articleSubscriberRepo,
		articleAnalyticRepository: articleAnalyticRepository,
		mailer:                    mailer,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendNewSubscriberWelcomeEmail, processor.ProcessTaskSendNewSubscriberWelcomeEmail)
	mux.HandleFunc(TaskSendNewArticleNotificationEmail, processor.ProcessTaskSendNewArticleNotificationEmail)
	mux.HandleFunc(TaskTrackReadArticle, processor.ProcessTaskTrackReadArticle)

	return processor.server.Start(mux)
}
