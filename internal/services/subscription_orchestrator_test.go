package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ihorlenko/weather_notifier/internal/models"
	"github.com/ihorlenko/weather_notifier/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubscriptionOrchestrator_ProcessSubscription_Success(t *testing.T) {
	mockValidator := &MockWeatherValidator{}
	mockSubscriptionService := &MockSubscriptionService{}
	mockEmailService := mocks.NewMockEmailService()

	orchestrator := NewSubscriptionOrchestrator(
		mockValidator,
		mockSubscriptionService,
		mockEmailService,
	)

	ctx := context.Background()
	email := "test@example.com"
	city := "Odesa"
	frequency := "daily"

	expectedSubscription := &models.Subscription{
		ID:                1,
		City:              city,
		Frequency:         frequency,
		Status:            "pending",
		ConfirmationToken: "test-confirmation-token",
		UnsubscribeToken:  "test-unsubscribe-token",
	}

	mockValidator.On("ValidateCity", ctx, city).Return(nil)
	mockSubscriptionService.On("CreateSubscription", email, city, frequency).Return(expectedSubscription, nil)
	mockEmailService.On("SendConfirmationEmail", email, city, expectedSubscription.ConfirmationToken).Return(nil)

	result, err := orchestrator.ProcessSubscription(ctx, email, city, frequency)

	assert.NoError(t, err)
	assert.Equal(t, expectedSubscription, result)

	mockValidator.AssertExpectations(t)
	mockSubscriptionService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestSubscriptionOrchestrator_ProcessSubscription_CityValidationFails(t *testing.T) {
	mockValidator := &MockWeatherValidator{}
	mockSubscriptionService := &MockSubscriptionService{}
	mockEmailService := mocks.NewMockEmailService()

	orchestrator := NewSubscriptionOrchestrator(
		mockValidator,
		mockSubscriptionService,
		mockEmailService,
	)

	ctx := context.Background()
	email := "test@example.com"
	city := "InvalidCity"
	frequency := "daily"

	validationError := errors.New("city not found")
	mockValidator.On("ValidateCity", ctx, city).Return(validationError)

	result, err := orchestrator.ProcessSubscription(ctx, email, city, frequency)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "city validation failed")

	mockValidator.AssertExpectations(t)
	mockSubscriptionService.AssertNotCalled(t, "CreateSubscription")
	mockEmailService.AssertNotCalled(t, "SendConfirmationEmail")
}

func TestSubscriptionOrchestrator_ProcessSubscription_SubscriptionCreationFails(t *testing.T) {
	mockValidator := &MockWeatherValidator{}
	mockSubscriptionService := &MockSubscriptionService{}
	mockEmailService := mocks.NewMockEmailService()

	orchestrator := NewSubscriptionOrchestrator(
		mockValidator,
		mockSubscriptionService,
		mockEmailService,
	)

	ctx := context.Background()
	email := "test@example.com"
	city := "Odesa"
	frequency := "daily"

	subscriptionError := errors.New("database error")
	mockValidator.On("ValidateCity", ctx, city).Return(nil)
	mockSubscriptionService.On("CreateSubscription", email, city, frequency).Return(nil, subscriptionError)

	result, err := orchestrator.ProcessSubscription(ctx, email, city, frequency)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "subscription creation failed")

	mockValidator.AssertExpectations(t)
	mockSubscriptionService.AssertExpectations(t)
	mockEmailService.AssertNotCalled(t, "SendConfirmationEmail")
}

func TestSubscriptionOrchestrator_ProcessSubscription_EmailSendingFails(t *testing.T) {
	mockValidator := &MockWeatherValidator{}
	mockSubscriptionService := &MockSubscriptionService{}
	mockEmailService := mocks.NewMockEmailService()

	orchestrator := NewSubscriptionOrchestrator(
		mockValidator,
		mockSubscriptionService,
		mockEmailService,
	)

	ctx := context.Background()
	email := "test@example.com"
	city := "Odesa"
	frequency := "daily"

	expectedSubscription := &models.Subscription{
		ID:                1,
		City:              city,
		Frequency:         frequency,
		Status:            "pending",
		ConfirmationToken: "test-confirmation-token",
	}

	emailError := errors.New("SMTP server error")
	mockValidator.On("ValidateCity", ctx, city).Return(nil)
	mockSubscriptionService.On("CreateSubscription", email, city, frequency).Return(expectedSubscription, nil)
	mockEmailService.On("SendConfirmationEmail", email, city, expectedSubscription.ConfirmationToken).Return(emailError)

	result, err := orchestrator.ProcessSubscription(ctx, email, city, frequency)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "confirmation email failed")

	mockValidator.AssertExpectations(t)
	mockSubscriptionService.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

type MockWeatherValidator struct {
	mock.Mock
}

func (m *MockWeatherValidator) ValidateCity(ctx context.Context, city string) error {
	args := m.Called(ctx, city)
	return args.Error(0)
}

type MockSubscriptionService struct {
	mock.Mock
}

func (m *MockSubscriptionService) CreateSubscription(email, city, frequency string) (*models.Subscription, error) {
	args := m.Called(email, city, frequency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionService) ConfirmSubscription(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockSubscriptionService) Unsubscribe(token string) error {
	args := m.Called(token)
	return args.Error(0)
}
