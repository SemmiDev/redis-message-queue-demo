package rmq

import (
	"github.com/google/uuid"
	"sync"
)

var ArticleSubscribers sync.Map

type ArticleSubscriber struct {
	SubscriberID uuid.UUID `json:"subscriber_id"`
	Email        string    `json:"email"`
	Subscribed   bool      `json:"subscribed"`
}

type ArticleSubscriberRepository interface {
	GetArticleSubscriberByID(id uuid.UUID) (ArticleSubscriber, bool)
	GetArticleSubscribers() []ArticleSubscriber
	CreateArticleSubscriber(subscriber ArticleSubscriber) uuid.UUID
	UnsubscribeArticleSubscriber(id uuid.UUID) bool
}

type InMemoryArticleSubscriberRepository struct{}

func (r *InMemoryArticleSubscriberRepository) GetArticleSubscriberByID(id uuid.UUID) (ArticleSubscriber, bool) {
	subscriber, exists := ArticleSubscribers.Load(id)
	if !exists {
		return ArticleSubscriber{}, false
	}

	return subscriber.(ArticleSubscriber), true
}

func (r *InMemoryArticleSubscriberRepository) CreateArticleSubscriber(subscriber ArticleSubscriber) uuid.UUID {
	subscriber.SubscriberID = uuid.New()

	ArticleSubscribers.Store(subscriber.SubscriberID, subscriber)
	return subscriber.SubscriberID
}

func (r *InMemoryArticleSubscriberRepository) UnsubscribeArticleSubscriber(id uuid.UUID) bool {
	subscriber, exists := ArticleSubscribers.Load(id)
	if !exists {
		return false
	}

	updatedSubscriber := subscriber.(ArticleSubscriber)
	updatedSubscriber.Subscribed = false

	ArticleSubscribers.Store(id, updatedSubscriber)
	return true
}

func (r *InMemoryArticleSubscriberRepository) GetArticleSubscribers() []ArticleSubscriber {
	subscribers := make([]ArticleSubscriber, 0)

	ArticleSubscribers.Range(func(key, value interface{}) bool {
		subscriber := value.(ArticleSubscriber)
		if subscriber.Subscribed {
			subscribers = append(subscribers, subscriber)
		}
		return true
	})

	return subscribers
}
