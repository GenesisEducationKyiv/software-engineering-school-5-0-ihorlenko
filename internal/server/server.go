package server

import (
	"context"

	"github.com/gin-gonic/gin"
	// Required for Swagger documentation
	_ "github.com/ihorlenko/weather_notifier/docs"
	"github.com/ihorlenko/weather_notifier/internal/api/handlers"
	"github.com/ihorlenko/weather_notifier/internal/config"
	"github.com/ihorlenko/weather_notifier/internal/models"
	"github.com/ihorlenko/weather_notifier/internal/weather"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Processor interface {
	ProcessSubscription(ctx context.Context, email, city, frequency string) (*models.Subscription, error)
}

type WeatherRetriever interface {
	GetWeather(ctx context.Context, city string) (*weather.Data, error)
}

type SubscriptionManager interface {
	CreateSubscription(email, city, frequency string) (*models.Subscription, error)
	ConfirmSubscription(token string) error
	Unsubscribe(token string) error
}

type EmailSender interface {
	SendConfirmationEmail(email, city, token string) error
	SendWeatherUpdate(email, city string, weather *weather.Data, unsubscribeToken string) error
}

type Server struct {
	config              *config.Config
	weatherService      WeatherRetriever
	subscriptionService SubscriptionManager
	emailService        EmailSender

	subscriptionProcessor Processor
}

func New(
	config *config.Config,
	weatherService WeatherRetriever,
	subscriptionService SubscriptionManager,
	emailService EmailSender,
	subscriptionProcessor Processor,
) *Server {
	return &Server{
		config:                config,
		weatherService:        weatherService,
		subscriptionService:   subscriptionService,
		emailService:          emailService,
		subscriptionProcessor: subscriptionProcessor,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.Default()

	weatherHandler := handlers.NewWeatherHandler(s.weatherService)

	subscriptionHandler := handlers.NewSubscriptionHandler(
		s.subscriptionProcessor,
		s.subscriptionService,
	)

	router.GET("/ping", handlers.PingHandler)

	api := router.Group("/api")
	{
		api.GET("/weather", weatherHandler.GetWeather)
		api.POST("/subscribe", subscriptionHandler.Subscribe)
		api.GET("/confirm/:token", subscriptionHandler.Confirm)
		api.GET("/unsubscribe/:token", subscriptionHandler.Unsubscribe)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Static("/static", "./web/static")
	router.StaticFile("/", "./web/index.html")

	return router
}
