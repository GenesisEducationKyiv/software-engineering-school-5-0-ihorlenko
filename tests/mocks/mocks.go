package mocks

import (
	"context"
	"errors"
	"fmt"

	"github.com/ihorlenko/weather_notifier/internal/weather"
	"github.com/stretchr/testify/mock"
)

var ErrMockReturnedNil = errors.New("mock returned nil")

type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetWeather(ctx context.Context, city string) (*weather.Data, error) {
	args := m.Called(ctx, city)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	if args.Get(0) == nil {
		return nil, ErrMockReturnedNil
	}

	weatherData, ok := args.Get(0).(*weather.Data)
	if !ok {
		return nil, fmt.Errorf("mock returned unexpected type")
	}

	return weatherData, nil
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendConfirmationEmail(email, city, token string) error {
	args := m.Called(email, city, token)
	return args.Error(0)
}

func (m *MockEmailService) SendWeatherUpdate(
	email, city string, weather *weather.Data,
	unsubscribeToken string,
) error {
	args := m.Called(email, city, weather, unsubscribeToken)
	return args.Error(0)
}

func NewMockWeatherService() *MockWeatherService {
	return &MockWeatherService{}
}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}
