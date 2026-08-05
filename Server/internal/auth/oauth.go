package auth

import (
	"os"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func OAuthMiddleware() gin.HandlerFunc {
	k, err := keyfunc.NewDefault([]string{os.Getenv("SUPABASE_JWKS_ENDPOINT")})
	if err != nil {
		logging.Logger.Fatalf("failed to create JWKS keyfunc: %s", err)
	}

	return func(c *gin.Context) {
		rawString := c.GetHeader("Authorization")
		if rawString == "" {
			c.Set("User", &models.UserToken{Validated: false})
			c.Next()
			return
		}
		trimmed := strings.TrimPrefix(rawString, "Bearer ")

		if trimmed == "" {
			c.Set("User", &models.UserToken{Validated: false})
			c.Next()
		}

		token, err := jwt.Parse(trimmed, k.Keyfunc)
		if err != nil {
			c.Set("User", &models.UserToken{Validated: false})
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Set("User", &models.UserToken{Validated: false})
			c.Next()
			return
		}

		sub, err := claims.GetSubject()
		if err != nil {
			c.Set("User", &models.UserToken{Validated: false})
			c.Next()
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			c.Set("User", &models.UserToken{Validated: false})
			c.Next()
			return
		}

		c.Set("User", &models.UserToken{
			Validated: true,
			SubjectId: &sub,
			Email:     &email,
		})
		c.Next()
	}
}
