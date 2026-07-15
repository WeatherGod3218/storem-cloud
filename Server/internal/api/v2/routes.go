package v2

import (
	"github.com/WeatherGod3218/weather-reels-server/internal/api/v2/transfer"
	"github.com/WeatherGod3218/weather-reels-server/internal/api/v2/videos"
	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.RouterGroup) {
	v2Group := router.Group("/v2")
	videos.Routes(v2Group)
	transfer.Routes(v2Group)
}
