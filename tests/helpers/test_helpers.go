package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Initializing Postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Initializing file source driver
	"github.com/ihorlenko/weather_notifier/internal/app"
	"github.com/ihorlenko/weather_notifier/internal/config"
	"github.com/ihorlenko/weather_notifier/internal/database"
	"github.com/ihorlenko/weather_notifier/internal/server"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

type TestSuite struct {
	PgContainer  testcontainers.Container
	DB           *gorm.DB
	Config       *config.Config
	Router       *gin.Engine
	TestServer   *httptest.Server
	AppContainer *app.Container
	ctx          context.Context
}

func SetupTestSuite(t *testing.T) *TestSuite {
	ctx := context.Background()

	projectRoot, err := findProjectRoot()
	require.NoError(t, err)

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("test_weather_api"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	cfg := &config.Config{
		AppConfig: config.AppConfig{
			BaseURL: "http://localhost:8080",
			Port:    "8080",
		},
		DBConfig: config.DBConfig{
			Host:     host,
			Port:     port.Port(),
			User:     "test_user",
			Password: "test_password",
			DBName:   "test_weather_api",
			SSLMode:  "disable",
		},
		WeatherAPIConfig: config.WeatherAPIConfig{
			APIKey: "test_api_key", // We'll mock this
		},
		EmailConfig: config.EmailConfig{
			From:     "test@example.com",
			Password: "test_password",
			SMTPHost: "localhost",
			SMTPPort: "587",
		},
	}

	os.Setenv("TEST_MIGRATIONS_PATH", filepath.Join(projectRoot, "migrations"))
	defer os.Unsetenv("TEST_MIGRATIONS_PATH")

	err = runTestMigrations(cfg)
	require.NoError(t, err)

	db, err := database.NewDBConnection(cfg)
	require.NoError(t, err)

	os.Setenv("SKIP_MIGRATIONS", "true")
	defer os.Unsetenv("SKIP_MIGRATIONS")

	appContainer, err := app.NewContainer(cfg)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)

	srv := server.New(
		cfg,
		appContainer.GetWeatherService(),
		appContainer.GetSubscriptionService(),
		appContainer.GetEmailService(),
		appContainer.GetSubscriptionOrchestrator(),
	)

	router := srv.SetupRoutes()

	testServer := httptest.NewServer(router)

	return &TestSuite{
		PgContainer:  pgContainer,
		DB:           db,
		Config:       cfg,
		Router:       router,
		TestServer:   testServer,
		AppContainer: appContainer,
		ctx:          ctx,
	}
}

func (ts *TestSuite) TeardownTestSuite(t *testing.T) {
	if ts.TestServer != nil {
		ts.TestServer.Close()
	}

	if ts.AppContainer != nil {
		err := ts.AppContainer.Close()
		if err != nil {
			log.Printf("Error closing app container: %v", err)
		}
	}

	if ts.PgContainer != nil {
		err := ts.PgContainer.Terminate(ts.ctx)
		require.NoError(t, err)
	}
}

func (ts *TestSuite) CleanupDatabase(t *testing.T) {
	err := ts.DB.Exec("TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE").Error
	require.NoError(t, err)

	err = ts.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error
	require.NoError(t, err)
}

func (ts *TestSuite) GetBaseURL() string {
	return ts.TestServer.URL
}

func runTestMigrations(cfg *config.Config) error {
	var migrationsPath string

	if envPath := os.Getenv("TEST_MIGRATIONS_PATH"); envPath != "" {
		migrationsPath = fmt.Sprintf("file://%s", envPath)
	} else if envPath := os.Getenv("MIGRATIONS_PATH"); envPath != "" {
		migrationsPath = fmt.Sprintf("file://%s", envPath)
	} else {
		projectRoot, err := findProjectRoot()
		if err != nil {
			return fmt.Errorf("failed to find project root: %w", err)
		}

		migrationsPath = fmt.Sprintf("file://%s/migrations", projectRoot)
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

	log.Printf("Successfully ran migrations from: %s", migrationsPath)
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
