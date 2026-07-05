package main

import (
	"fmt"

	v1 "github.com/WeatherGod3218/weather-reels-server/internal/api/v1"
	"github.com/WeatherGod3218/weather-reels-server/internal/database"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/s3"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := s3.InitS3(); err != nil {
		logging.Logger.Fatal(fmt.Sprintf("failed to initialize the app %s", err))
	}
	if err := database.InitDatabase(); err != nil {
		logging.Logger.Fatal(fmt.Sprintf("failed to initialize the app %s", err))
	}

	router := gin.Default()

	authGroup := router.Group("/")
	v1.SetRoutes(authGroup)

	router.Run(":8080")
}
