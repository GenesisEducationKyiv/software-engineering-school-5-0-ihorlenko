package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ihorlenko/weather_notifier/internal/api/handlers"
)

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", handlers.PingHandler)
	}
	router.Static("/static", "./web")
}
