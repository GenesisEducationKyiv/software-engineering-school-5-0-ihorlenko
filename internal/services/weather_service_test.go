package services

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/ihorlenko/weather_notifier/internal/config"
	apperrors "github.com/ihorlenko/weather_notifier/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestWeatherService_GetWeather_Success(t *testing.T) {
	defer gock.Off()

	cfg := &config.Config{
		WeatherAPIConfig: config.WeatherAPIConfig{
			APIKey: "test-api-key",
		},
	}

	service := NewWeatherService(cfg)
	ctx := context.Background()
	city := "Odesa"

	gock.New("https://api.weatherapi.com").
		Get("/v1/current.json").
		MatchParam("key", "test-api-key").
		MatchParam("q", city).
		Reply(200).
		JSON(map[string]interface{}{
			"location": map[string]interface{}{
				"name": "Odesa",
			},
			"current": map[string]interface{}{
				"temp_c":   15.5,
				"humidity": 65.0,
				"condition": map[string]interface{}{
					"text": "Partly cloudy",
				},
			},
		})

	result, err := service.GetWeather(ctx, city)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Odesa", result.City)
	assert.Equal(t, 15.5, result.Temperature)
	assert.Equal(t, 65.0, result.Humidity)
	assert.Equal(t, "Partly cloudy", result.Description)

	assert.True(t, gock.IsDone())
}

func TestWeatherService_GetWeather_APIReturnsError(t *testing.T) {
	defer gock.Off()

	cfg := &config.Config{
		WeatherAPIConfig: config.WeatherAPIConfig{
			APIKey: "test-api-key",
		},
	}

	service := NewWeatherService(cfg)
	ctx := context.Background()
	city := "InvalidCity"

	gock.New("https://api.weatherapi.com").
		Get("/v1/current.json").
		MatchParam("key", "test-api-key").
		MatchParam("q", city).
		Reply(400).
		JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    1006,
				"message": "No matching location found.",
			},
		})

	result, err := service.GetWeather(ctx, city)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, apperrors.ErrWeatherServiceUnavailable, err)

	assert.True(t, gock.IsDone())
}

func TestWeatherService_GetWeather_NetworkError(t *testing.T) {
	defer gock.Off()

	cfg := &config.Config{
		WeatherAPIConfig: config.WeatherAPIConfig{
			APIKey: "test-api-key",
		},
	}

	service := NewWeatherService(cfg)
	ctx := context.Background()
	city := "Odesa"

	result, err := service.GetWeather(ctx, city)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, apperrors.ErrWeatherServiceUnavailable, err)
}

func TestWeatherService_GetWeather_MalformedJSON(t *testing.T) {
	defer gock.Off()

	cfg := &config.Config{
		WeatherAPIConfig: config.WeatherAPIConfig{
			APIKey: "test-api-key",
		},
	}

	service := NewWeatherService(cfg)
	ctx := context.Background()
	city := "Odesa"

	gock.New("https://api.weatherapi.com").
		Get("/v1/current.json").
		MatchParam("key", "test-api-key").
		MatchParam("q", city).
		Reply(200).
		BodyString("invalid json response")

	result, err := service.GetWeather(ctx, city)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, apperrors.ErrWeatherServiceUnavailable, err)

	assert.True(t, gock.IsDone())
}

func TestWeatherService_GetWeather_ContextCancelled(t *testing.T) {
	defer gock.Off()

	cfg := &config.Config{
		WeatherAPIConfig: config.WeatherAPIConfig{
			APIKey: "test-api-key",
		},
	}

	service := NewWeatherService(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	city := "Odesa"

	result, err := service.GetWeather(ctx, city)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, apperrors.ErrWeatherServiceUnavailable, err)
}

func TestWeatherService_GetWeather_SpecialCharactersInCity(t *testing.T) {
	defer gock.Off()

	cfg := &config.Config{
		WeatherAPIConfig: config.WeatherAPIConfig{
			APIKey: "test-api-key",
		},
	}

	service := NewWeatherService(cfg)
	ctx := context.Background()
	city := "São Paulo"

	gock.New("https://api.weatherapi.com").
		Get("/v1/current.json").
		MatchParam("key", "test-api-key").
		MatchParam("q", city).
		Reply(200).
		JSON(map[string]interface{}{
			"location": map[string]interface{}{
				"name": "São Paulo",
			},
			"current": map[string]interface{}{
				"temp_c":   25.0,
				"humidity": 70.0,
				"condition": map[string]interface{}{
					"text": "Sunny",
				},
			},
		})

	result, err := service.GetWeather(ctx, city)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "São Paulo", result.City)
	assert.Equal(t, 25.0, result.Temperature)

	assert.True(t, gock.IsDone())
}
