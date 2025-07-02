package app

import (
	"context"
	"fmt"

	"github.com/ihorlenko/weather_notifier/internal/config"
	"github.com/ihorlenko/weather_notifier/internal/database"
	"github.com/ihorlenko/weather_notifier/internal/models"
	"github.com/ihorlenko/weather_notifier/internal/repositories"
	"github.com/ihorlenko/weather_notifier/internal/services"
	"github.com/ihorlenko/weather_notifier/internal/weather"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	Create(email string) (*models.User, error)
	GetOrCreate(email string) (*models.User, error)
}

type SubscriptionRepository interface {
	Create(sub *models.Subscription) error
	GetByConfirmationToken(token string) (*models.Subscription, error)
	GetByUnsubscribeToken(token string) (*models.Subscription, error)
	UpdateStatus(id uint, status string) error
	GetActiveSubscriptionsByFrequency(frequency string) ([]models.Subscription, error)
}

type WeatherService interface {
	GetWeather(ctx context.Context, city string) (*weather.Data, error)
}

type EmailService interface {
	SendConfirmationEmail(email, city, token string) error
	SendWeatherUpdate(email, city string, weather *weather.Data, unsubscribeToken string) error
}

type SubscriptionService interface {
	CreateSubscription(email, city, frequency string) (*models.Subscription, error)
	ConfirmSubscription(token string) error
	Unsubscribe(token string) error
}

type SubscriptionProcessor interface {
	ProcessSubscription(ctx context.Context, email, city, frequency string) (*models.Subscription, error)
}

type WeatherValidator interface {
	ValidateCity(ctx context.Context, city string) error
}

type Container struct {
	config *config.Config
	db     *gorm.DB

	userRepo         UserRepository
	subscriptionRepo SubscriptionRepository

	weatherService      WeatherService
	emailService        EmailService
	subscriptionService SubscriptionService

	weatherValidator      WeatherValidator
	subscriptionProcessor SubscriptionProcessor
}

func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{
		config: cfg,
	}

	if err := container.initializeDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	container.initializeRepositories()
	container.initializeCoreServices()
	container.initializeCompositeServices()

	return container, nil
}

func (c *Container) initializeDatabase() error {
	if err := database.RunMigrations(c.config); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	db, err := database.NewDBConnection(c.config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	c.db = db
	return nil
}

func (c *Container) initializeRepositories() {
	c.userRepo = repositories.NewUserRepository(c.db)
	c.subscriptionRepo = repositories.NewSubscriptionRepository(c.db)
}

func (c *Container) initializeCoreServices() {
	c.weatherService = services.NewWeatherService(c.config)
	c.emailService = services.NewEmailService(c.config)
	c.subscriptionService = services.NewSubscriptionService(c.userRepo, c.subscriptionRepo)
}

func (c *Container) initializeCompositeServices() {
	c.weatherValidator = services.NewWeatherValidator(c.weatherService)

	c.subscriptionProcessor = services.NewSubscriptionProcessor(
		c.weatherValidator,
		c.subscriptionService,
		c.emailService,
	)
}

func (c *Container) GetDatabase() *gorm.DB {
	return c.db
}

func (c *Container) GetUserRepository() UserRepository {
	return c.userRepo
}

func (c *Container) GetSubscriptionRepository() SubscriptionRepository {
	return c.subscriptionRepo
}

func (c *Container) GetWeatherService() WeatherService {
	return c.weatherService
}

func (c *Container) GetEmailService() EmailService {
	return c.emailService
}

func (c *Container) GetSubscriptionService() SubscriptionService {
	return c.subscriptionService
}

func (c *Container) GetWeatherValidator() WeatherValidator {
	return c.weatherValidator
}

func (c *Container) GetSubscriptionOrchestrator() SubscriptionProcessor {
	return c.subscriptionProcessor
}

func (c *Container) Close() error {
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
