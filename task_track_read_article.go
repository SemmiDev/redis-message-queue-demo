package rmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskTrackReadArticle = "task:track_read_article"

type PayloadTrackReadArticle struct {
	ArticleID   uuid.UUID `json:"article_id"`
	QueryParams string    `json:"query_params"`
}

func (distributor *RedisTaskDistributor) DistributeTaskTrackReadArticle(
	ctx context.Context,
	payload *PayloadTrackReadArticle,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskTrackReadArticle, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskTrackReadArticle(_ context.Context, task *asynq.Task) error {
	var payload PayloadTrackReadArticle
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	event := GetEventFromParams(payload.QueryParams)
	processor.articleAnalyticRepository.Create(payload.ArticleID, event)

	log.Info().Str("type", task.Type()).Msg("processed task")

	return nil
}
