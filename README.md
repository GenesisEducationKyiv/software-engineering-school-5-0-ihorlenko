# Weather API Application

A RESTful API service that allows users to subscribe to weather updates for cities of their choice.

## Live Demo

The application is deployed and available at:
- **Subscription UI**: [https://weather-notifier-krb7.onrender.com/](https://weather-notifier-krb7.onrender.com/)
- **API Documentation**: [https://weather-notifier-krb7.onrender.com/swagger/index.html](https://weather-notifier-krb7.onrender.com/swagger/index.html)

## Features

- Get current weather information for any city
- Subscribe to weather updates via email (hourly or daily at 9:00)
- Confirmation email for subscriptions
- Unsubscribe functionality
- Simple web interface for subscription management

## Tech Stack

- **Backend**: Go 1.24.3
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Scheduler**: cron
- **Containerization**: Docker, Docker Compose
- **Weather Data**: WeatherAPI.com

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/weather?city={city}` | Get current weather for a given city |
| POST | `/api/subscribe` | Subscribe an email to weather updates |
| GET | `/api/confirm/{token}` | Confirm email subscription |
| GET | `/api/unsubscribe/{token}` | Unsubscribe from weather updates |

## Prerequisites

Before running this application, you'll need:

1. [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/)
2. Free API key from [WeatherAPI.com](https://www.weatherapi.com/)
3. Email account for sending notifications (if using Gmail, you may need an app password)

## Local Setup

### 1. Clone the repository

```bash
git clone https://github.com/ihorlenko/weather_notifier.git
cd weather-api
```

### 2. Configure environment variables

Rename `.env.example` to a `.env`:

```
BASE_URL=http://localhost:8080
POSTGRES_HOST=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=weather_api
WEATHER_API_KEY=your_weather_api_key
EMAIL_FROM=your_email@example.com
EMAIL_PASSWORD=your_email_password
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
```

Replace `your_weather_api_key`, `your_email@example.com`, and `your_email_password` with your actual values.

### 3. Run with Docker Compose

There are two ways to start the application:

#### Using docker-compose directly:

```bash
docker-compose --env-file .env up -d
```

#### Using Make commands:

```bash
# Build the Docker image
make build

# Start containers
make up
```

This command will:
- Build the Go application
- Start a PostgreSQL database
- Set up the necessary network
- Run database migrations
- Start the application on port 8080

### 4. Verify the application is running

Open your browser and navigate to:

```
http://localhost:8080
```

You should see the weather subscription interface.

Test the API endpoint:

```bash
curl http://localhost:8080/api/weather?city=London
```

### 5. Stop the application

When you're done, you can stop the application with:

```bash
# Using docker-compose
docker-compose down

# Or using Make
make down
```

To remove all data volumes:

```bash
docker-compose down -v
```

## Make Commands

The project includes a Makefile with several useful commands:

| Command | Description |
|---------|-------------|
| `make build` | Build the Docker image |
| `make rebuild` | Rebuild the Docker image |
| `make up` | Start the containers |
| `make down` | Stop the containers |
| `make restart` | Restart the containers |
| `make postgres` | Connect to the PostgreSQL database |
| `make help` | Show all available commands |
