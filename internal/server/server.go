package server

import (
	"net/http"
	"strings"

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

	subscriptionOrchestrator interfaces.SubscriptionOrchestrator
}

func New(
	config *config.Config,
	weatherService interfaces.WeatherService,
	subscriptionService interfaces.SubscriptionService,
	emailService interfaces.EmailService,
	subscriptionOrchestrator interfaces.SubscriptionOrchestrator,
) *Server {
	return &Server{
		config:                   config,
		weatherService:           weatherService,
		subscriptionService:      subscriptionService,
		emailService:             emailService,
		subscriptionOrchestrator: subscriptionOrchestrator,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.Default()

	weatherHandler := handlers.NewWeatherHandler(s.weatherService)
	subscriptionHandler := handlers.NewSubscriptionHandler(
		s.subscriptionOrchestrator,
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

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "endpoint not found",
			})
			return
		}

		c.File("./web/index.html")
	})

	return router
}
