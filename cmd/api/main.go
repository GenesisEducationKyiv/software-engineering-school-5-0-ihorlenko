package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/ihorlenko/weather_notifier/docs"
	"github.com/ihorlenko/weather_notifier/internal/api/handlers"
	"github.com/ihorlenko/weather_notifier/internal/config"
	"github.com/ihorlenko/weather_notifier/internal/database"
	"github.com/ihorlenko/weather_notifier/internal/repositories"
	"github.com/ihorlenko/weather_notifier/internal/scheduler"
	"github.com/ihorlenko/weather_notifier/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Weather Notifier API
// @version         1.0
// @description     API for weather notification subscription
// @host            localhost:8080
// @BasePath        /api
func main() {
	cfg := config.LoadConfig()

	if err := database.RunMigrations(cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	db, err := database.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	subscriptionRepo := repositories.NewSubscriptionRepository(db)

	weatherService := services.NewWeatherService(cfg)
	emailService := services.NewEmailService(cfg)
	subscriptionService := services.NewSubscriptionService(userRepo, subscriptionRepo)

	weatherHandler := handlers.NewWeatherHandler(weatherService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, emailService, weatherService)

	weatherScheduler := scheduler.NewWeatherScheduler(subscriptionRepo, weatherService, emailService)

	weatherScheduler.Start()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	router := gin.Default()

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

	go func() {
		if err := router.Run(":" + cfg.AppConfig.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	weatherScheduler.Stop()
}
