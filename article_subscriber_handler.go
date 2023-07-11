package rmq

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (s *Server) setupArticleSubscriberRoutes() {
	s.Router.Post("/subscriptions", s.CreateArticleSubscriberHandler)
	s.Router.Get("/subscriptions/{id}", s.GetArticleSubscriberByIDHandler)
	s.Router.Put("/subscriptions/{id}/unsubscribe", s.UnsubscribeArticleSubscriberHandler)
}

func (s *Server) GetArticleSubscriberByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid subscriber ID", http.StatusBadRequest)
		return
	}

	subscriber, exists := s.SubscriberRepo.GetArticleSubscriberByID(id)
	if !exists {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(subscriber)

}

func (s *Server) CreateArticleSubscriberHandler(w http.ResponseWriter, r *http.Request) {
	var subscriberReq struct{ Email string }
	err := json.NewDecoder(r.Body).Decode(&subscriberReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	subscriber := ArticleSubscriber{
		Email:      subscriberReq.Email,
		Subscribed: true,
	}

	id := s.SubscriberRepo.CreateArticleSubscriber(subscriber)
	subscriber.SubscriberID = id

	taskPayload := &PayloadNewSubscriberWelcomeEmail{
		Email: subscriber.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(2),
		asynq.Queue(QueueCritical),
	}

	if taskErr := s.TaskDistributor.DistributeTaskSendNewSubscriberWelcomeEmail(context.Background(), taskPayload, opts...); taskErr != nil {
		log.Err(taskErr).Msg("failed to distribute task")
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscriber)
}

func (s *Server) UnsubscribeArticleSubscriberHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid subscriber ID", http.StatusBadRequest)
		return
	}

	success := s.SubscriberRepo.UnsubscribeArticleSubscriber(id)
	if !success {
		http.NotFound(w, r)
		return
	}

	subscriber, _ := s.SubscriberRepo.GetArticleSubscriberByID(id)
	json.NewEncoder(w).Encode(subscriber)
}
