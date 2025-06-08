package errors

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidInput         = errors.New("invalid input")
	ErrInternalError        = errors.New("internal error")

	ErrWeatherServiceUnavailable = errors.New("weather service is currently unavailable")
	ErrInvalidCity               = errors.New("city not found or invalid")
)
