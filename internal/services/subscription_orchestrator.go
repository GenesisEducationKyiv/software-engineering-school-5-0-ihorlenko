package services

import (
	"context"
	"fmt"

	"github.com/ihorlenko/weather_notifier/internal/models"
	"github.com/ihorlenko/weather_notifier/internal/weather"
)

type CityValidator interface {
	ValidateCity(ctx context.Context, city string) error
}

type EmailSender interface {
	SendConfirmationEmail(email, city, token string) error
	SendWeatherUpdate(email, city string, weather *weather.Data, unsubscribeToken string) error
}

type SubscriptionProcessor interface {
	CreateSubscription(email, city, frequency string) (*models.Subscription, error)
	ConfirmSubscription(token string) error
	Unsubscribe(token string) error
}

type SubscriptionOrchestrator struct {
	validator           CityValidator
	subscriptionService SubscriptionProcessor
	emailService        EmailSender
}

func NewSubscriptionProcessor(
	validator CityValidator,
	subscriptionService SubscriptionProcessor,
	emailService EmailSender,
) *SubscriptionOrchestrator {
	return &SubscriptionOrchestrator{
		validator:           validator,
		subscriptionService: subscriptionService,
		emailService:        emailService,
	}
}

func (sp *SubscriptionOrchestrator) ProcessSubscription(
	ctx context.Context,
	email, city, frequency string,
) (*models.Subscription, error) {
	if err := sp.validator.ValidateCity(ctx, city); err != nil {
		return nil, fmt.Errorf("city validation failed: %w", err)
	}

	subscription, err := sp.subscriptionService.CreateSubscription(email, city, frequency)
	if err != nil {
		return nil, fmt.Errorf("subscription creation failed: %w", err)
	}

	err = sp.emailService.SendConfirmationEmail(email, city, subscription.ConfirmationToken)
	if err != nil {
		return nil, fmt.Errorf("confirmation email failed: %w", err)
	}

	return subscription, nil
}
