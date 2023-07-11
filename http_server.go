package rmq

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type Server struct {
	Router              *chi.Mux
	ArticleRepo         ArticleRepository
	SubscriberRepo      ArticleSubscriberRepository
	ArticleAnalyticRepo ArticleAnalyticRepository
	TaskDistributor     TaskDistributor
}

func NewServer(
	articleRepo ArticleRepository,
	subscriberRepo ArticleSubscriberRepository,
	articleAnalyticRepo ArticleAnalyticRepository,
	taskDistributor TaskDistributor,
) *Server {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	s := &Server{
		Router:              r,
		ArticleRepo:         articleRepo,
		SubscriberRepo:      subscriberRepo,
		ArticleAnalyticRepo: articleAnalyticRepo,
		TaskDistributor:     taskDistributor,
	}

	s.setupArticleRoutes()
	s.setupArticleSubscriberRoutes()
	s.setupArticleAnalyticRoutes()

	return s
}

func (s *Server) Start(port int) {
	log.Printf("app started on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), s.Router))
}
