package auth

import (
	"log"
	"os"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// func InitOAuth() {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	conf := &oauth2.Config{
// 		ClientID:     os.Getenv("CLIENT_ID"),
// 		ClientSecret: os.Getenv("CLIENT_SECRET"),
// 		RedirectURL:  os.Getenv("SERVER_HOST"),
// 		Scopes: []string{
// 			"openid",
// 			"https://www.googleapis.com/auth/userinfo.email",
// 		},
// 		Endpoint: google.Endpoint,
// 	}

// 	url := conf.AuthCodeURL("state")
// 	fmt.Printf("Visit the URL for the auth dialog: %v", url)

// 	// Handle the exchange code to initiate a transport.
// 	tok, err := conf.Exchange(ctx, "authorization-code")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }

func OAuthMiddleware() gin.HandlerFunc {
	k, err := keyfunc.NewDefault([]string{os.Getenv("SUPABASE_JWKS_ENDPOINT")})
	if err != nil {
		log.Fatalf("failed to create JWKS keyfunc: %s", err)
	}

	return func(c *gin.Context) {
		rawString := c.GetHeader("Authorization")
		if rawString == "" {
			c.Set("User", &models.UserToken{Validated: false})
			c.Next()
			return
		}
		trimmed := strings.TrimPrefix(rawString, "Bearer ")

		token, err := jwt.Parse(trimmed, k.Keyfunc)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "token": token}).Warning("Error handling token")
			c.Set("User", &models.UserToken{Validated: false})
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logging.Logger.Warnf("Error handling claims, not of claims type")
			c.Set("User", &models.UserToken{Validated: false})
			c.Next()
			return
		}

		sub, err := claims.GetSubject()
		if err != nil {
			logging.Logger.Warnf("Error getting claims subject: %s", err)
			c.Set("User", &models.UserToken{Validated: false})
			c.Next()
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			logging.Logger.Warning("Error getting email on user")
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
