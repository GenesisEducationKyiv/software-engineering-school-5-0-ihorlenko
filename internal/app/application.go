package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ihorlenko/weather_notifier/internal/config"
	"github.com/ihorlenko/weather_notifier/internal/interfaces"
	"github.com/ihorlenko/weather_notifier/internal/scheduler"
	"github.com/ihorlenko/weather_notifier/internal/server"
)

type Application struct {
	config     *config.Config
	container  *Container
	server     *server.Server
	scheduler  interfaces.WeatherScheduler
	httpServer *http.Server
}

func New(cfg *config.Config) *Application {
	return &Application{
		config: cfg,
	}
}

func (app *Application) Initialize() error {
	container, err := NewContainer(app.config)
	if err != nil {
		return fmt.Errorf("failed to initialize container: %w", err)
	}
	app.container = container

	app.scheduler = scheduler.NewWeatherScheduler(
		container.GetSubscriptionRepository(),
		container.GetWeatherService(),
		container.GetEmailService(),
	)

	app.server = server.New(
		app.config,
		container.GetWeatherService(),
		container.GetSubscriptionService(),
		container.GetEmailService(),
	)

	log.Println("Application initialized successfully")
	return nil
}

func (app *Application) Run() error {
	app.scheduler.Start()
	log.Println("Scheduler started")

	router := app.server.SetupRoutes()
	app.httpServer = &http.Server{
		Addr:         ":" + app.config.AppConfig.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on port %s", app.config.AppConfig.Port)
	return app.httpServer.ListenAndServe()
}

func (app *Application) Shutdown() error {
	if app.scheduler != nil {
		app.scheduler.Stop()
		log.Println("Scheduler stopped")
	}

	if app.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := app.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown HTTP server: %w", err)
		}
		log.Println("HTTP server stopped")
	}

	if app.container != nil {
		if err := app.container.Close(); err != nil {
			return fmt.Errorf("failed to close container: %w", err)
		}
		log.Println("Container resources closed")
	}

	return nil
}
