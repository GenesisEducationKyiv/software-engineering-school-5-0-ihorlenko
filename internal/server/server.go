package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ihorlenko/weather_notifier/internal/api/handlers"
	"github.com/ihorlenko/weather_notifier/internal/config"
	"github.com/ihorlenko/weather_notifier/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Required for Swagger documentation
	_ "github.com/ihorlenko/weather_notifier/docs"
)

type Server struct {
	config              *config.Config
	weatherService      *services.WeatherService
	subscriptionService *services.SubscriptionService
	emailService        *services.EmailService
}

func New(
	config *config.Config,
	weatherService *services.WeatherService,
	subscriptionService *services.SubscriptionService,
	emailService *services.EmailService,
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
	router.Static("/css", "./web/css")
	router.Static("/js", "./web/js")
	router.StaticFile("/", "./web/index.html")

	return router
}
