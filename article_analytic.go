package rmq

import (
	"github.com/google/uuid"
	"sort"
	"time"
)

type EventType string

const (
	Read      EventType = "read"
	LinkClick EventType = "link_click"
)

type ArticleAnalytic struct {
	ArticleID uuid.UUID `json:"article_id"`
	EventType EventType `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

type ArticleAnalyticRepository interface {
	Create(articleId uuid.UUID, eventType EventType)
	FindTopNMostReadArticles(n int) []ArticleAnalyticResponse
	FindTopNClickedArticlesViaEmail(n int) []ArticleAnalyticResponse
}

var ArticleAnalytics []ArticleAnalytic

type InMemoryArticleAnalyticsRepository struct {
	ArticleRepository ArticleRepository
}

func (i *InMemoryArticleAnalyticsRepository) Create(articleId uuid.UUID, eventType EventType) {
	ArticleAnalytics = append(ArticleAnalytics, ArticleAnalytic{
		ArticleID: articleId,
		EventType: eventType,
		Timestamp: time.Now(),
	})
}

type ArticleAnalyticResponse struct {
	Count   int     `json:"count"`
	Article Article `json:"article"`
}

func (i *InMemoryArticleAnalyticsRepository) FindTopNMostReadArticles(n int) []ArticleAnalyticResponse {
	articleCount := make(map[uuid.UUID]int)

	for _, event := range ArticleAnalytics {
		if event.EventType == Read {
			articleCount[event.ArticleID]++
		}
	}

	sortedArticles := make([]uuid.UUID, 0, len(articleCount))
	for articleID := range articleCount {
		sortedArticles = append(sortedArticles, articleID)
	}
	sort.Slice(sortedArticles, func(i, j int) bool {
		return articleCount[sortedArticles[i]] > articleCount[sortedArticles[j]]
	})

	topNArticles := make([]ArticleAnalyticResponse, 0, n)
	for x := 0; x < n && x < len(sortedArticles); x++ {
		article, exists := i.ArticleRepository.GetArticleByID(sortedArticles[x])
		if !exists {
			continue
		}

		topNArticles = append(topNArticles, ArticleAnalyticResponse{
			Count:   articleCount[sortedArticles[x]],
			Article: article,
		})
	}

	return topNArticles
}

func (i *InMemoryArticleAnalyticsRepository) FindTopNClickedArticlesViaEmail(n int) []ArticleAnalyticResponse {
	articleCount := make(map[uuid.UUID]int)

	for _, event := range ArticleAnalytics {
		if event.EventType == LinkClick {
			articleCount[event.ArticleID]++
		}
	}

	sortedArticles := make([]uuid.UUID, 0, len(articleCount))
	for articleID := range articleCount {
		sortedArticles = append(sortedArticles, articleID)
	}
	sort.Slice(sortedArticles, func(i, j int) bool {
		return articleCount[sortedArticles[i]] > articleCount[sortedArticles[j]]
	})

	topNArticles := make([]ArticleAnalyticResponse, 0, n)
	for x := 0; x < n && x < len(sortedArticles); x++ {
		article, exists := i.ArticleRepository.GetArticleByID(sortedArticles[x])
		if !exists {
			continue
		}

		topNArticles = append(topNArticles, ArticleAnalyticResponse{
			Count:   articleCount[sortedArticles[x]],
			Article: article,
		})
	}

	return topNArticles
}
