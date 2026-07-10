package v1

import (
	"github.com/WeatherGod3218/weather-reels-server/internal/api/v1/videos"
	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.RouterGroup) {
	v1Group := router.Group("/v1")
	videos.Routes(v1Group)
}
