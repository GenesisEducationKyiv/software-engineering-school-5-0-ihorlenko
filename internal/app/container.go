package app

import (
	"fmt"

	"github.com/ihorlenko/weather_notifier/internal/config"
	"github.com/ihorlenko/weather_notifier/internal/database"
	"github.com/ihorlenko/weather_notifier/internal/repositories"
	"github.com/ihorlenko/weather_notifier/internal/services"
	"gorm.io/gorm"
)

type Container struct {
	config *config.Config
	db     *gorm.DB

	userRepo         *repositories.UserRepository
	subscriptionRepo *repositories.SubscriptionRepository

	weatherService      *services.WeatherService
	emailService        *services.EmailService
	subscriptionService *services.SubscriptionService
}

func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{
		config: cfg,
	}

	if err := container.initializeDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	container.initializeRepositories()
	container.initializeServices()

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

func (c *Container) initializeServices() {
	c.weatherService = services.NewWeatherService(c.config)
	c.emailService = services.NewEmailService(c.config)
	c.subscriptionService = services.NewSubscriptionService(c.userRepo, c.subscriptionRepo)
}

func (c *Container) GetDatabase() *gorm.DB {
	return c.db
}

func (c *Container) GetUserRepository() *repositories.UserRepository {
	return c.userRepo
}

func (c *Container) GetSubscriptionRepository() *repositories.SubscriptionRepository {
	return c.subscriptionRepo
}

func (c *Container) GetWeatherService() *services.WeatherService {
	return c.weatherService
}

func (c *Container) GetEmailService() *services.EmailService {
	return c.emailService
}

func (c *Container) GetSubscriptionService() *services.SubscriptionService {
	return c.subscriptionService
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
