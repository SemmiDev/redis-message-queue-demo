package main

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"redis-message-queue"
)

func main() {
	articleRepo := &rmq.InMemoryArticleRepository{}

	articleSubscriberRepo := &rmq.InMemoryArticleSubscriberRepository{}

	articleAnalytics := &rmq.InMemoryArticleAnalyticsRepository{
		ArticleRepository: articleRepo,
	}

	redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379"}
	taskDistributor := rmq.NewRedisTaskDistributor(redisOpt)

	mailer := rmq.NewFakeGmailSender("sammi", "sammi@gmail.com", "xxx")
	taskProcessor := rmq.NewRedisTaskProcessor(redisOpt, articleRepo, articleSubscriberRepo, articleAnalytics, mailer)

	go func() {
		log.Info().Msg("start task processor")
		err := taskProcessor.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start task processor")
		}
	}()

	s := rmq.NewServer(articleRepo, articleSubscriberRepo, articleAnalytics, taskDistributor)
	s.Start(8080)
}
