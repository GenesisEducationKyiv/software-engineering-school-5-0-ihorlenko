package services

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/ihorlenko/weather_notifier/internal/errors"
	"github.com/ihorlenko/weather_notifier/internal/weather"
)

type WeatherRetriever interface {
	GetWeather(ctx context.Context, city string) (*weather.Data, error)
}

type WeatherValidator struct {
	weatherService WeatherRetriever
}

func NewWeatherValidator(weatherService WeatherRetriever) *WeatherValidator {
	return &WeatherValidator{
		weatherService: weatherService,
	}
}

func (v *WeatherValidator) ValidateCity(ctx context.Context, city string) error {
	_, err := v.weatherService.GetWeather(ctx, city)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidCity) {
			return fmt.Errorf("city '%s' not found", city)
		}
		return fmt.Errorf("unable to validate city: %w", err)
	}

	return nil
}
