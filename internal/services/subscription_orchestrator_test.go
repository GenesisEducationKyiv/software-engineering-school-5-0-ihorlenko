package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ihorlenko/weather_notifier/internal/models"
	"github.com/ihorlenko/weather_notifier/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testEmail     = "test@example.com"
	testCity      = "Odesa"
	testFrequency = "daily"
)

var ErrMockReturnedNil = errors.New("mock returned nil subscription")

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
	email := testEmail
	city := testCity
	frequency := testFrequency

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
	email := testEmail
	city := "InvalidCity"
	frequency := testFrequency

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
	email := testEmail
	city := testCity
	frequency := testFrequency

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
	email := testEmail
	city := testCity
	frequency := testFrequency

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

	if err := args.Error(1); err != nil {
		return nil, err
	}

	if args.Get(0) == nil {
		return nil, ErrMockReturnedNil
	}

	subscription, ok := args.Get(0).(*models.Subscription)
	if !ok {
		return nil, fmt.Errorf("mock returned unexpected type")
	}

	return subscription, nil
}

func (m *MockSubscriptionService) ConfirmSubscription(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockSubscriptionService) Unsubscribe(token string) error {
	args := m.Called(token)
	return args.Error(0)
}
