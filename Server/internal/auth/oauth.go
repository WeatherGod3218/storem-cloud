package auth

import (
	"github.com/gin-gonic/gin"
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
	return func(c *gin.Context) {
		c.Next()
	}
}
