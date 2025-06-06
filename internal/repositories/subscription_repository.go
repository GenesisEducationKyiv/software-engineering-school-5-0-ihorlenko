package repositories

import (
	"github.com/ihorlenko/weather_notifier/internal/models"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (sr *SubscriptionRepository) Create(sub *models.Subscription) error {
	return sr.db.Create(sub).Error
}

func (sr *SubscriptionRepository) GetByConfirmationToken(token string) (*models.Subscription, error) {
	var sub models.Subscription
	result := sr.db.Where("confirmation_token = ?", token).Find(&sub)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sub, nil
}

func (sr *SubscriptionRepository) GetByUnsubscribeToken(token string) (*models.Subscription, error) {
	var sub models.Subscription
	result := sr.db.Where("unsubscribe_token = ?", token).Find(&sub)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sub, nil
}

func (sr *SubscriptionRepository) UpdateStatus(id uint, status string) error {
	return sr.db.Model(&models.Subscription{}).Where("id = ?", id).Update("status", status).Error
}

func (sr *SubscriptionRepository) GetActiveSubscriptionsByFrequency(frequency string) ([]models.Subscription, error) {
	var subs []models.Subscription
	result := sr.db.Where("status = ? AND frequency = ?", "active", frequency).Preload("User").Find(&subs)
	if result.Error != nil {
		return nil, result.Error
	}
	return subs, nil
}
