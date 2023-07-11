package rmq

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (s *Server) setupArticleAnalyticRoutes() {
	s.Router.Post("/analytics", s.CreateArticleAnalytic)
	s.Router.Get("/analytics/top-most-read", s.GetTopNMostReadArticles)
	s.Router.Get("/analytics/top-clicked-via-email", s.GetTopNClickedArticlesViaEmail)
}

func (s *Server) CreateArticleAnalytic(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ArticleID uuid.UUID `json:"article_id"`
		EventType EventType `json:"event_type"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	s.ArticleAnalyticRepo.Create(requestData.ArticleID, requestData.EventType)

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) GetTopNMostReadArticles(w http.ResponseWriter, r *http.Request) {
	nParam := r.URL.Query().Get("n")
	n, err := strconv.Atoi(nParam)
	if err != nil {
		http.Error(w, "Invalid 'n' parameter", http.StatusBadRequest)
		return
	}

	articles := s.ArticleAnalyticRepo.FindTopNMostReadArticles(n)

	json.NewEncoder(w).Encode(articles)
}

func (s *Server) GetTopNClickedArticlesViaEmail(w http.ResponseWriter, r *http.Request) {
	nParam := r.URL.Query().Get("n")
	n, err := strconv.Atoi(nParam)
	if err != nil {
		http.Error(w, "Invalid 'n' parameter", http.StatusBadRequest)
		return
	}

	articles := s.ArticleAnalyticRepo.FindTopNClickedArticlesViaEmail(n)

	json.NewEncoder(w).Encode(articles)
}
