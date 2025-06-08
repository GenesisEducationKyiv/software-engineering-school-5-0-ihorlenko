package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	apperrors "github.com/ihorlenko/weather_notifier/internal/errors"
	"github.com/ihorlenko/weather_notifier/internal/models"
	"github.com/ihorlenko/weather_notifier/internal/repositories"
)

type SubscriptionService struct {
	userRepo         *repositories.UserRepository
	subscriptionRepo *repositories.SubscriptionRepository
}

func NewSubscriptionService(
	userRepo *repositories.UserRepository,
	subscriptionRepo *repositories.SubscriptionRepository,
) *SubscriptionService {
	return &SubscriptionService{
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *SubscriptionService) CreateSubscription(email, city, frequency string) (*models.Subscription, error) {
	user, err := s.userRepo.GetOrCreate(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	confirmToken, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate confirmation token: %w", err)
	}

	unsubToken, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate unsubscribe token: %w", err)
	}

	subscription := &models.Subscription{
		UserID:            user.ID,
		City:              city,
		Frequency:         frequency,
		Status:            "pending",
		ConfirmationToken: confirmToken,
		UnsubscribeToken:  unsubToken,
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return subscription, nil
}

func (s *SubscriptionService) ConfirmSubscription(token string) error {
	subscription, err := s.subscriptionRepo.GetByConfirmationToken(token)
	if err != nil {
		if errors.Is(err, apperrors.ErrSubscriptionNotFound) {
			return fmt.Errorf("subscription not found")
		}
		return fmt.Errorf("failed to find subscription: %w", err)
	}

	if subscription.Status != "pending" {
		return fmt.Errorf("subscription is already confirmed or cancelled")
	}

	if err := s.subscriptionRepo.UpdateStatus(subscription.ID, "active"); err != nil {
		return fmt.Errorf("failed to update subscription status: %w", err)
	}

	return nil
}

func (s *SubscriptionService) Unsubscribe(token string) error {
	subscription, err := s.subscriptionRepo.GetByUnsubscribeToken(token)
	if err != nil {
		if errors.Is(err, apperrors.ErrSubscriptionNotFound) {
			return fmt.Errorf("subscription not found")
		}
		return fmt.Errorf("failed to find subscription: %w", err)
	}

	if subscription.Status == "cancelled" {
		return fmt.Errorf("subscription is already cancelled")
	}

	if err := s.subscriptionRepo.UpdateStatus(subscription.ID, "cancelled"); err != nil {
		return fmt.Errorf("failed to update subscription status: %w", err)
	}

	return nil
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
