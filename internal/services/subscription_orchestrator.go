package services

import (
	"context"
	"fmt"

	"github.com/ihorlenko/weather_notifier/internal/interfaces"
	"github.com/ihorlenko/weather_notifier/internal/models"
)

type subscriptionOrchestrator struct {
	validator           interfaces.WeatherValidator
	subscriptionService interfaces.SubscriptionService
	emailService        interfaces.EmailService
}

var _ interfaces.SubscriptionOrchestrator = (*subscriptionOrchestrator)(nil)

func NewSubscriptionOrchestrator(
	validator interfaces.WeatherValidator,
	subscriptionService interfaces.SubscriptionService,
	emailService interfaces.EmailService,
) interfaces.SubscriptionOrchestrator {
	return &subscriptionOrchestrator{
		validator:           validator,
		subscriptionService: subscriptionService,
		emailService:        emailService,
	}
}

func (o *subscriptionOrchestrator) ProcessSubscription(
	ctx context.Context,
	email, city, frequency string,
) (*models.Subscription, error) {
	if err := o.validator.ValidateCity(ctx, city); err != nil {
		return nil, fmt.Errorf("city validation failed: %w", err)
	}

	subscription, err := o.subscriptionService.CreateSubscription(email, city, frequency)
	if err != nil {
		return nil, fmt.Errorf("subscription creation failed: %w", err)
	}

	err = o.emailService.SendConfirmationEmail(email, city, subscription.ConfirmationToken)
	if err != nil {
		return nil, fmt.Errorf("confirmation email failed: %w", err)
	}

	return subscription, nil
}
