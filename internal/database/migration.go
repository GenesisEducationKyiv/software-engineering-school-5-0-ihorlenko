package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	// Initializing Postgres driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// Initializing file source driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ihorlenko/weather_notifier/internal/config"
)

func RunMigrations(cfg *config.Config) error {
	if os.Getenv("SKIP_MIGRATIONS") == "true" {
		log.Println("Skipping migrations due to SKIP_MIGRATIONS environment variable")
		return nil
	}

	var migrationsPath string

	if envPath := os.Getenv("TEST_MIGRATIONS_PATH"); envPath != "" {
		migrationsPath = fmt.Sprintf("file://%s", envPath)
	} else if envPath := os.Getenv("MIGRATIONS_PATH"); envPath != "" {
		migrationsPath = fmt.Sprintf("file://%s", envPath)
	} else {
		if projectRoot, err := findProjectRoot(); err == nil {
			migrationsPath = fmt.Sprintf("file://%s/migrations", projectRoot)
		} else {
			migrationsPath = "file://migrations"
		}
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Host,
		cfg.DBConfig.Port, cfg.DBConfig.DBName, cfg.DBConfig.SSLMode)

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance with path %s: %w", migrationsPath, err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Printf("Migrations successfully executed from: %s", migrationsPath)
	return nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find project root (go.mod not found)")
}
