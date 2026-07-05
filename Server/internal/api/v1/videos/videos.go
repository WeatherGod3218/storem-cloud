package videos

import (
	"net/http"

	"github.com/WeatherGod3218/weather-reels-server/internal/database"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func VerifyVideos(c *gin.Context) {
	logging.Logger.WithFields(logrus.Fields{"module": "v1/api/videos", "method": "VerifyVideos"}).Info("starting verification!")
	var videoList []string

	if err := c.ShouldBindJSON(&videoList); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "VerifyVideos"}).Warning("failed to bind video list")
		c.JSON(400, gin.H{
			"error": "Unable to process request!",
		})
		return
	}

	nonVerified, err := database.VerifyVideoList(videoList)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "VerifyVideos"}).Warning("failed to verify videos in the database")
		c.JSON(400, gin.H{
			"error": "Unable to process request!",
		})
		return
	}

	c.JSON(http.StatusOK, nonVerified)
}

func BackupVideo(c *gin.Context) {

}

func Routes(r *gin.RouterGroup) {
	videos := r.Group("/videos")
	videos.PUT("/verify", VerifyVideos)
}
