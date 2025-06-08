package database

import (
	"fmt"
	"log"

	"github.com/ihorlenko/weather_notifier/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDBConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DBConfig.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}
