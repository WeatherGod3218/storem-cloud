package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	v1 "github.com/WeatherGod3218/weather-reels-server/internal/api/v1"
	"github.com/WeatherGod3218/weather-reels-server/internal/auth"
	"github.com/WeatherGod3218/weather-reels-server/internal/database"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/WeatherGod3218/weather-reels-server/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed web/dist
var staticFS embed.FS

func serveFrontend(dist fs.FS) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := fs.ReadFile(dist, "index.html")
		if err != nil {
			logging.Logger.Error(fmt.Sprintf("failed to read index.html: %s", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	}
}

// @title WeatherReels
// @version 1.0
// @description API for backing up videos.
// @BasePath /api/v1
func main() {
	if err := s3.InitS3(); err != nil {
		logging.Logger.Fatal(fmt.Sprintf("failed to initialize s3 for the app %s", err))
	}
	if err := database.InitDatabase(); err != nil {
		logging.Logger.Fatal(fmt.Sprintf("failed to initialize database for the app %s", err))
	}

	dist, err := fs.Sub(staticFS, "web/dist")
	if err != nil {
		logging.Logger.Fatal(fmt.Sprintf("failed to initialize dist for the app %s", err))
	}

	assets, err := fs.Sub(dist, "assets")
	if err != nil {
		logging.Logger.Fatal(fmt.Sprintf("failed to initialize assets for the app %s", err))
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.StaticFS("/assets", http.FS(assets))

	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	router.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	api := router.Group("/api", auth.ApiKeyMiddleware())
	v1.SetRoutes(api)

	router.NoRoute(auth.OAuthMiddleware(), serveFrontend(dist))

	router.Run(":8080")
}
