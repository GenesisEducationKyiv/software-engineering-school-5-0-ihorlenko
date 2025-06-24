package mocks

import (
	"context"

	"github.com/ihorlenko/weather_notifier/internal/interfaces"
	"github.com/ihorlenko/weather_notifier/internal/types"
	"github.com/stretchr/testify/mock"
)

type MockWeatherService struct {
	mock.Mock
}

var _ interfaces.WeatherService = (*MockWeatherService)(nil)

func (m *MockWeatherService) GetWeather(ctx context.Context, city string) (*types.WeatherData, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(*types.WeatherData), args.Error(1)
}

type MockEmailService struct {
	mock.Mock
}

var _ interfaces.EmailService = (*MockEmailService)(nil)

func (m *MockEmailService) SendConfirmationEmail(email, city, token string) error {
	args := m.Called(email, city, token)
	return args.Error(0)
}

func (m *MockEmailService) SendWeatherUpdate(email, city string, weather *types.WeatherData, unsubscribeToken string) error {
	args := m.Called(email, city, weather, unsubscribeToken)
	return args.Error(0)
}

func NewMockWeatherService() *MockWeatherService {
	return &MockWeatherService{}
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}
