package repository

import "github.com/SaranHiruthikM/newsletter-system/internal/domain"

type SubscriberRepository interface {
	FindByEmail(email string) (*domain.Subscriber, error)
	Create(subscriber *domain.Subscriber) error
	FindByToken(token string) (*domain.Subscriber, error)
	UpdateConfirmed(id string, confirmed bool) error
	FindAllConfirmed() ([]*domain.Subscriber, error)
}

type NewsletterRepository interface {
	Create(newsletter *domain.NewsletterSend) error
	IncrementSentCount(id string) error
	IncrementFailCount(id string) error
}
