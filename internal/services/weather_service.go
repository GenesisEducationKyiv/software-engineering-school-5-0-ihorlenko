package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ihorlenko/weather_notifier/internal/config"
	apperrors "github.com/ihorlenko/weather_notifier/internal/errors"
	"github.com/ihorlenko/weather_notifier/internal/interfaces"
	"github.com/ihorlenko/weather_notifier/internal/types"
)

var _ interfaces.WeatherService = (*WeatherService)(nil)

type WeatherService struct {
	baseURL string
	apiKey  string
}

func NewWeatherService(cfg *config.Config) interfaces.WeatherService {
	return &WeatherService{
		baseURL: "https://api.weatherapi.com/v1/",
		apiKey:  cfg.WeatherAPIConfig.APIKey,
	}
}

func (ws *WeatherService) GetWeather(ctx context.Context, city string) (*types.WeatherData, error) {
	baseURL, err := url.Parse(ws.baseURL + "current.json")
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	params := url.Values{}
	params.Add("key", ws.apiKey)
	params.Add("q", city)
	baseURL.RawQuery = params.Encode()
	finalURL := baseURL.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, finalURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to make request to weather API: %v", err)
		return nil, apperrors.ErrWeatherServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Weather API returned non 200 status code: %d", resp.StatusCode)
		return nil, apperrors.ErrWeatherServiceUnavailable
	}

	var apiResp types.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Printf("Failed to decode weather API response: %v", err)
		return nil, apperrors.ErrWeatherServiceUnavailable
	}

	return &types.WeatherData{
		City:        apiResp.Location.Name,
		Temperature: apiResp.Current.TempC,
		Humidity:    apiResp.Current.Humidity,
		Description: apiResp.Current.Condition.Text,
	}, nil
}
