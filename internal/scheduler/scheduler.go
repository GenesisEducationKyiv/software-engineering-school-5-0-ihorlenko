package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/ihorlenko/weather_notifier/internal/models"
	"github.com/ihorlenko/weather_notifier/internal/types"
	"github.com/robfig/cron/v3"
)

type SubscriptionRepositoryManager interface {
	Create(sub *models.Subscription) error
	GetByConfirmationToken(token string) (*models.Subscription, error)
	GetByUnsubscribeToken(token string) (*models.Subscription, error)
	UpdateStatus(id uint, status string) error
	GetActiveSubscriptionsByFrequency(frequency string) ([]models.Subscription, error)
}

type WeatherRetriever interface {
	GetWeather(ctx context.Context, city string) (*types.WeatherData, error)
}

type EmailSender interface {
	SendConfirmationEmail(email, city, token string) error
	SendWeatherUpdate(email, city string, weather *types.WeatherData, unsubscribeToken string) error
}

type WeatherScheduler struct {
	subscriptionRepo SubscriptionRepositoryManager
	weatherService   WeatherRetriever
	emailService     EmailSender
	cron             *cron.Cron
}

func NewWeatherScheduler(
	subscriptionRepo SubscriptionRepositoryManager,
	weatherService WeatherRetriever,
	emailService EmailSender,
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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		weather, err := s.weatherService.GetWeather(ctx, sub.City)
		cancel()

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
