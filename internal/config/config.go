package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppConfig        AppConfig
	DBConfig         DBConfig
	WeatherAPIConfig WeatherAPIConfig
	EmailConfig      EmailConfig
}

type AppConfig struct {
	BaseURL string
	Port    string
}

type EmailConfig struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (db *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode)
}

type WeatherAPIConfig struct {
	APIKey string
}

func LoadConfig() *Config {
	return &Config{
		AppConfig: AppConfig{
			BaseURL: getEnv("BASE_URL", "http://localhost:8080"),
			Port:    "8080",
		},

		DBConfig: DBConfig{
			Host:     getEnv("POSTGRES_HOST", "postgres"),
			Port:     "5432",
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_DB", "weather_notifier"),
			SSLMode:  "disable",
		},

		WeatherAPIConfig: WeatherAPIConfig{
			APIKey: getEnv("WEATHER_API_KEY", ""),
		},

		EmailConfig: EmailConfig{
			From:     getEnv("EMAIL_FROM", ""),
			Password: getEnv("EMAIL_PASSWORD", ""),
			SMTPHost: getEnv("EMAIL_SMTP_HOST", "smtp.gmail.com"),
			SMTPPort: getEnv("EMAIL_SMTP_PORT", "587"),
		},
	}
}

func getEnv(key string, defaultValue string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultValue
	}
	return env
}
