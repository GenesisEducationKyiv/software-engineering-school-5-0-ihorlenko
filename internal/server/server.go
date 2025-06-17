package server

import (
	"github.com/gin-gonic/gin"
	// Required for Swagger documentation
	_ "github.com/ihorlenko/weather_notifier/docs"
	"github.com/ihorlenko/weather_notifier/internal/api/handlers"
	"github.com/ihorlenko/weather_notifier/internal/config"
	"github.com/ihorlenko/weather_notifier/internal/interfaces"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config              *config.Config
	weatherService      interfaces.WeatherService
	subscriptionService interfaces.SubscriptionService
	emailService        interfaces.EmailService
}

func New(
	config *config.Config,
	weatherService interfaces.WeatherService,
	subscriptionService interfaces.SubscriptionService,
	emailService interfaces.EmailService,
) *Server {
	return &Server{
		config:              config,
		weatherService:      weatherService,
		subscriptionService: subscriptionService,
		emailService:        emailService,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.Default()

	weatherHandler := handlers.NewWeatherHandler(s.weatherService)
	subscriptionHandler := handlers.NewSubscriptionHandler(
		s.subscriptionService,
		s.emailService,
		s.weatherService,
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
