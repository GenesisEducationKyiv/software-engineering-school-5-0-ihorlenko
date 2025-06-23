package interfaces

import (
	"context"

	"github.com/ihorlenko/weather_notifier/internal/models"
	"github.com/ihorlenko/weather_notifier/internal/types"
)

type WeatherService interface {
	GetWeather(ctx context.Context, city string) (*types.WeatherData, error)
}

type EmailService interface {
	SendConfirmationEmail(email, city, token string) error
	SendWeatherUpdate(email, city string, weather *types.WeatherData, unsubscribeToken string) error
}

type SubscriptionService interface {
	CreateSubscription(email, city, frequency string) (*models.Subscription, error)
	ConfirmSubscription(token string) error
	Unsubscribe(token string) error
}
