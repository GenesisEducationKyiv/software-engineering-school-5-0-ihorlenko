package interfaces

import (
	"github.com/ihorlenko/weather_notifier/internal/models"
)

type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	Create(email string) (*models.User, error)
	GetOrCreate(email string) (*models.User, error)
}

type SubscriptionRepository interface {
	Create(sub *models.Subscription) error
	GetByConfirmationToken(token string) (*models.Subscription, error)
	GetByUnsubscribeToken(token string) (*models.Subscription, error)
	UpdateStatus(id uint, status string) error
	GetActiveSubscriptionsByFrequency(frequency string) ([]models.Subscription, error)
}
