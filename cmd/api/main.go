package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ihorlenko/weather_notifier/internal/app"
	"github.com/ihorlenko/weather_notifier/internal/config"
)

// @title           Weather Notifier API
// @version         1.0
// @description     API for weather notification subscription
// @host            localhost:8080
// @BasePath        /api
func main() {
	cfg := config.LoadConfig()

	application := app.New(cfg)

	if err := application.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("Failed to start application: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down application...")

	if err := application.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Application stopped")
}
