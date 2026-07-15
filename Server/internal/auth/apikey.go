package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logging.Logger.Info("Checking API Key!")
		serverApiKey := os.Getenv("SERVER_API_KEY")

		clientToken := c.GetHeader("Authorization")

		clientApiKey := strings.TrimPrefix(clientToken, "Bearer ")
		if serverApiKey == clientApiKey {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "API Access Key was not Authorized",
		})
	}
}
