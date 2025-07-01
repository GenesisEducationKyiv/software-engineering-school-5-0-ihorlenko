package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ihorlenko/weather_notifier/internal/interfaces"
)

type WeatherHandler struct {
	weatherService interfaces.WeatherService
}

func NewWeatherHandler(weatherService interfaces.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

// GetWeather godoc
// @Summary      Get current weather for a city
// @Description  Returns current weather for a given city
// @Tags         weather
// @Accept       json
// @Produce      json
// @Param        city query string true "City name"
// @Success      200  {object}  types.WeatherData
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /weather [get]
func (h *WeatherHandler) GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	weatherData, err := h.weatherService.GetWeather(ctx, city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, weatherData)
}
