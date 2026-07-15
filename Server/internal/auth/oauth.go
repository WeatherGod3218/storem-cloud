package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func InitOAuth() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("SERVER_HOST"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state")
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Handle the exchange code to initiate a transport.
	tok, err := conf.Exchange(ctx, "authorization-code")
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(oauth2.NoContext, tok)
	client.Get("...")
}

func OAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
