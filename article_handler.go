package rmq

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *Server) setupArticleRoutes() {
	s.Router.Post("/articles", s.CreateArticleHandler)
	s.Router.Get("/articles", s.GetAllArticlesHandler)
	s.Router.Get("/articles/{id}", s.GetArticleByIDHandler)
	s.Router.Put("/articles/{id}", s.UpdateArticleHandler)
	s.Router.Delete("/articles/{id}", s.DeleteArticleHandler)
}

func (s *Server) GetAllArticlesHandler(w http.ResponseWriter, r *http.Request) {
	articles := s.ArticleRepo.GetAllArticles()
	json.NewEncoder(w).Encode(articles)
}

func (s *Server) GetArticleByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if strings.TrimSpace(idParam) == "" {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	var article Article

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		articleData, exists := s.ArticleRepo.GetArticleBySlug(idParam)
		if !exists {
			http.NotFound(w, r)
			return
		}

		article = articleData
	} else {
		articleData, exists := s.ArticleRepo.GetArticleByID(id)
		if !exists {
			http.NotFound(w, r)
			return
		}

		article = articleData
	}

	taskPayload := &PayloadTrackReadArticle{
		ArticleID:   article.ID,
		QueryParams: r.URL.Query().Encode(),
	}

	opts := []asynq.Option{
		asynq.MaxRetry(2),
		asynq.Queue(QueueCritical),
	}

	if taskErr := s.TaskDistributor.DistributeTaskTrackReadArticle(context.Background(), taskPayload, opts...); taskErr != nil {
		log.Err(taskErr).Msg("Failed to distribute task")
	}

	json.NewEncoder(w).Encode(article)
}

func (s *Server) CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	var article Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	article = s.ArticleRepo.CreateArticle(article)

	taskPayload := &PayloadNewArticleNotificationEmail{
		ArticleTitle: article.Title,
		ArticleSlug:  article.Slug,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(0),
		asynq.Queue(QueueCritical),
	}

	if taskErr := s.TaskDistributor.DistributeTaskSendNewArticleNotificationEmail(context.Background(), taskPayload, opts...); taskErr != nil {
		log.Err(taskErr).Msg("Failed to distribute task")
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

func (s *Server) UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	var updatedArticle Article
	err = json.NewDecoder(r.Body).Decode(&updatedArticle)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	success := s.ArticleRepo.UpdateArticle(id, updatedArticle)
	if !success {
		http.NotFound(w, r)
		return
	}

	updatedArticle.ID = id
	json.NewEncoder(w).Encode(updatedArticle)
}

func (s *Server) DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	success := s.ArticleRepo.DeleteArticle(id)
	if !success {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
