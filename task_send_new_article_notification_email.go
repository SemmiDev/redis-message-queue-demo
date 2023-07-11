package rmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendNewArticleNotificationEmail = "task:send_new_article_notification_email"

type PayloadNewArticleNotificationEmail struct {
	ArticleTitle string `json:"article_title"`
	ArticleSlug  string `json:"article_slug"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendNewArticleNotificationEmail(
	ctx context.Context,
	payload *PayloadNewArticleNotificationEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendNewArticleNotificationEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendNewArticleNotificationEmail(_ context.Context, task *asynq.Task) error {
	var payload PayloadNewArticleNotificationEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	subscribers := processor.articleSubscriberRepo.GetArticleSubscribers()
	if len(subscribers) == 0 {
		return nil
	}

	to := make([]string, 0, len(subscribers))

	for _, v := range subscribers {
		to = append(to, v.Email)
	}

	subject := "New Article Published"
	baseURL := "http://localhost:8080/articles/" + payload.ArticleSlug
	source := "newsletter"
	medium := "email"
	campaign := "newsletter"

	link, err := GenerateUTMURL(baseURL, source, medium, campaign)
	if err != nil {
		return fmt.Errorf("failed to generate utm url: %w", err)
	}

	content := fmt.Sprintf(`Hello, <br/>
	We have published a new article: <a href="%s">%s</a><br/>
	`, link, payload.ArticleTitle)

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("processed task")

	return nil
}
