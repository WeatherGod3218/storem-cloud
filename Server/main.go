package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strings"

	v2 "github.com/WeatherGod3218/weather-reels-server/internal/api/v2"
	"github.com/WeatherGod3218/weather-reels-server/internal/api/v2/transfer"
	"github.com/WeatherGod3218/weather-reels-server/internal/auth"
	"github.com/WeatherGod3218/weather-reels-server/internal/database"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/WeatherGod3218/weather-reels-server/internal/redis"
	"github.com/WeatherGod3218/weather-reels-server/internal/s3"
	"github.com/WeatherGod3218/weather-reels-server/internal/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.yaml.in/yaml/v3"

	_ "github.com/WeatherGod3218/weather-reels-server/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed web/dist
var distFS embed.FS

var indexBuffer string = ""

func serveIcon(c *gin.Context, icon []byte) {
	c.Data(http.StatusOK, "image/png", icon)
}

func newIndexFile(metadata string) *string {
	newFile := strings.Replace(indexBuffer, "<!--SSR_META-->", metadata, 1)
	return &newFile
}

func serveFrontend(c *gin.Context) {
	finalIndex := &indexBuffer
	if strings.HasPrefix(c.FullPath(), "/videos/video/") {
		finalIndex = HandleEmbedForVideo(c)
	}

	if finalIndex == nil {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexBuffer))
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(*finalIndex))
}

// @title WeatherReels
// @version 1.0
// @description API for backing up videos.
// @BasePath /api/v2
func main() {

	cfgFile, err := os.ReadFile("config.yaml")
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Error loading config file!")
	}

	var config models.Config

	err = yaml.Unmarshal(cfgFile, &config)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("Error unmarshling config file!")
	}

	logging.Logger.WithFields(logrus.Fields{"config": config}).Info("Loaded Config!")
	users.InitUsers(config)

	if err := redis.InitRedis(); err != nil {
		//its dead, but thats fine because we dont neeeeed it
		logging.Logger.Warnf("failed to initialize redis! continuing without %s", err)
	}

	if err := s3.InitS3(); err != nil {
		logging.Logger.Fatalf("failed to initialize s3 for the app %s", err)
	}

	if err := database.InitDatabase(); err != nil {
		logging.Logger.Fatalf("failed to initialize database for the app %s", err)
	}

	if err := transfer.InitTus(); err != nil {
		logging.Logger.Fatalf("failed to initialize tus protocol for the app %s", err)
	}

	dist, err := fs.Sub(distFS, "web/dist")
	if err != nil {
		logging.Logger.Fatalf("failed to initialize dist for the app %s", err)
	}

	assets, err := fs.Sub(dist, "assets")
	if err != nil {
		logging.Logger.Fatalf("failed to initialize assets for the app %s", err)
	}

	index, err := fs.ReadFile(dist, "index.html")
	if err != nil {
		logging.Logger.Fatalf("failed to load index %s", err)
	}
	indexBuffer = string(index)

	icon, err := fs.ReadFile(dist, "favicon.png")
	if err != nil {
		logging.Logger.Fatal("No favicon?")
	}

	router := gin.Default()
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, models.SuccessResponse{
			Success: true,
		})
	})

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{ //shut the f up
		SkipPaths: []string{"/health"},
	}))

	router.RedirectTrailingSlash = false

	frontend := router.Group("")
	frontend.Use(cors.Default())
	frontend.StaticFS("/assets", http.FS(assets))

	frontend.GET("/favicon.ico", func(ctx *gin.Context) {
		serveIcon(ctx, icon)
	})

	frontend.GET("/favicon.svg", func(ctx *gin.Context) { // just do both and make my life easy
		serveIcon(ctx, icon)
	})

	frontend.GET("/favicon.png", func(ctx *gin.Context) { // theres a third now ig.
		serveIcon(ctx, icon)
	})

	frontend.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	frontend.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	api := router.Group("/api")
	v2.SetRoutes(api)

	if os.Getenv("DEV_MODE") == "true" {
		logging.Logger.Info("RUNNING IN DEV MODE!")
		router.NoRoute(createViteProxy())
	} else {
		router.NoRoute(auth.OAuthMiddleware(), func(c *gin.Context) {
			reqURL := c.Request.URL.String()

			if strings.Contains(reqURL, "api/") {
				c.JSON(http.StatusNotFound, models.ErrorResponse{
					Error: "API not found!",
				})
				return
			}
			serveFrontend(c)
		})
	}

	router.Run(":8080")
}
