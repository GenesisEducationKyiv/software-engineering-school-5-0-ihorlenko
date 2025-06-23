package interfaces

import (
	"context"
)

type WeatherValidator interface {
	ValidateCity(ctx context.Context, city string) error
}
