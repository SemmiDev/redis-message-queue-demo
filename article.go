package rmq

import (
	"strings"
	"sync"

	"github.com/google/uuid"
)

var Articles sync.Map

type Article struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Slug    string    `json:"slug"`
}

func (a *Article) GenerateSlug() {
	title := strings.ToLower(a.Title)
	title = strings.ReplaceAll(title, " ", "-")
	a.Slug = title
}

type ArticleRepository interface {
	GetAllArticles() []Article
	GetArticleByID(id uuid.UUID) (Article, bool)
	GetArticleBySlug(slug string) (Article, bool)
	CreateArticle(article Article) Article
	UpdateArticle(id uuid.UUID, article Article) bool
	DeleteArticle(id uuid.UUID) bool
}

type InMemoryArticleRepository struct{}

func (r *InMemoryArticleRepository) GetAllArticles() []Article {
	articles := make([]Article, 0)

	Articles.Range(func(key, value interface{}) bool {
		articles = append(articles, value.(Article))
		return true
	})

	return articles
}

func (r *InMemoryArticleRepository) GetArticleByID(id uuid.UUID) (Article, bool) {
	article, exists := Articles.Load(id)
	if !exists {
		return Article{}, false
	}

	return article.(Article), true
}

func (r *InMemoryArticleRepository) GetArticleBySlug(slug string) (Article, bool) {
	articles := r.GetAllArticles()
	for _, article := range articles {
		if article.Slug == slug {
			return article, true
		}
	}
	return Article{}, false
}

func (r *InMemoryArticleRepository) CreateArticle(article Article) Article {
	article.ID = uuid.New()
	article.GenerateSlug()

	Articles.Store(article.ID, article)
	return article
}

func (r *InMemoryArticleRepository) UpdateArticle(id uuid.UUID, article Article) bool {
	_, exists := Articles.Load(id)
	if !exists {
		return false
	}

	article.ID = id
	article.GenerateSlug()

	Articles.Store(id, article)
	return true
}

func (r *InMemoryArticleRepository) DeleteArticle(id uuid.UUID) bool {
	_, exists := Articles.Load(id)
	if !exists {
		return false
	}

	Articles.Delete(id)
	return true
}
