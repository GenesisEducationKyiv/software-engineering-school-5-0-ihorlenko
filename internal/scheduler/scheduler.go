package scheduler

import (
	"log"
	"time"

	"github.com/ihorlenko/weather_notifier/internal/repositories"
	"github.com/ihorlenko/weather_notifier/internal/services"
	"github.com/robfig/cron/v3"
)

type WeatherScheduler struct {
	subscriptionRepo *repositories.SubscriptionRepository
	weatherService   *services.WeatherService
	emailService     *services.EmailService
	cron             *cron.Cron
}

func NewWeatherScheduler(
	subscriptionRepo *repositories.SubscriptionRepository,
	weatherService *services.WeatherService,
	emailService *services.EmailService,
) *WeatherScheduler {
	c := cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))

	return &WeatherScheduler{
		subscriptionRepo: subscriptionRepo,
		weatherService:   weatherService,
		emailService:     emailService,
		cron:             c,
	}
}

func (s *WeatherScheduler) Start() {
	_, err := s.cron.AddFunc("0 0 * * * *", func() {
		s.sendUpdates("hourly")
	})
	if err != nil {
		log.Printf("Failed to schedule hourly updates: %v", err)
	}

	_, err = s.cron.AddFunc("0 0 9 * * *", func() {
		s.sendUpdates("daily")
	})
	if err != nil {
		log.Printf("Failed to schedule daily updates: %v", err)
	}

	s.cron.Start()
	log.Println("Weather scheduler started")
}

func (s *WeatherScheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		log.Println("Weather scheduler stopped")
	}
}

func (s *WeatherScheduler) sendUpdates(frequency string) {
	log.Printf("Sending %s weather updates", frequency)

	subscriptions, err := s.subscriptionRepo.GetActiveSubscriptionsByFrequency(frequency)
	if err != nil {
		log.Printf("Failed to fetch subscriptions: %v", err)
		return
	}

	log.Printf("Found %d active subscriptions for %s updates", len(subscriptions), frequency)

	for _, sub := range subscriptions {
		weather, err := s.weatherService.GetWeather(sub.City)
		if err != nil {
			log.Printf("Failed to get weather for city %s: %v", sub.City, err)
			continue
		}

		err = s.emailService.SendWeatherUpdate(sub.User.Email, sub.City, weather, sub.UnsubscribeToken)
		if err != nil {
			log.Printf("Failed to send weather update to %s: %v", sub.User.Email, err)
		} else {
			log.Printf("Weather update sent to %s for city %s", sub.User.Email, sub.City)
		}
	}
}
