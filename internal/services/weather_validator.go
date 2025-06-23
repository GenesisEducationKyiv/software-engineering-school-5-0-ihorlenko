package services

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/ihorlenko/weather_notifier/internal/errors"
	"github.com/ihorlenko/weather_notifier/internal/interfaces"
)

var _ interfaces.WeatherValidator = (*WeatherValidator)(nil)

type WeatherValidator struct {
	weatherService interfaces.WeatherService
}

func NewWeatherValidator(weatherService interfaces.WeatherService) interfaces.WeatherValidator {
	return &WeatherValidator{
		weatherService: weatherService,
	}
}

func (v *WeatherValidator) ValidateCity(ctx context.Context, city string) error {
	if city == "" {
		return errors.New("city cannot be empty")
	}

	_, err := v.weatherService.GetWeather(ctx, city)
	if err != nil {
		if errors.Is(err, apperrors.ErrWeatherServiceUnavailable) {
			return fmt.Errorf("weather service unavailable, please try again later")
		}
		if errors.Is(err, apperrors.ErrInvalidCity) {
			return fmt.Errorf("city '%s' not found", city)
		}
		return fmt.Errorf("unable to validate city: %w", err)
	}

	return nil
}
