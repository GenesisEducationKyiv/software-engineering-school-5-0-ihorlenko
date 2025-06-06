package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ihorlenko/weather_notifier/internal/config"
)

type WeatherData struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

type WeatherAPIResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  float64 `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

type WeatherService struct {
	baseURL string
	apiKey  string
}

func NewWeatherService(cfg *config.Config) *WeatherService {
	return &WeatherService{
		baseURL: "https://api.weatherapi.com/v1/",
		apiKey:  cfg.WeatherAPIConfig.APIKey,
	}
}

func (ws *WeatherService) GetWeather(city string) (*WeatherData, error) {
	baseURL, err := url.Parse(ws.baseURL + "current.json")
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	params := url.Values{}
	params.Add("key", ws.apiKey)
	params.Add("q", city)
	baseURL.RawQuery = params.Encode()

	finalURL := baseURL.String()

	//nolint:gosec // G107 false positive
	resp, err := http.Get(finalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to weather API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned non-200 status code: %d", resp.StatusCode)
	}

	var apiResp WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode weather API response: %w", err)
	}

	return &WeatherData{
		City:        apiResp.Location.Name,
		Temperature: apiResp.Current.TempC,
		Humidity:    apiResp.Current.Humidity,
		Description: apiResp.Current.Condition.Text,
	}, nil
}
