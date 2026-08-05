package main

import (
	"fmt"
	"html"
	"os"

	"github.com/WeatherGod3218/weather-reels-server/internal/database"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/redis"
	"github.com/WeatherGod3218/weather-reels-server/internal/s3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var DEFAULT_DESCRIPTION string = "No description has been applied to this video."

func HandleEmbedForVideo(c *gin.Context) *string {
	id := c.Param("id")
	serverHost := os.Getenv("SERVER_HOST")

	if id == "" {
		logging.Logger.WithFields(logrus.Fields{"module": "main", "method": "HandleEmbedForVideo"}).Warning("failed to get video ID")
		return nil
	}

	video, err := database.GetVideoData(id)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "HandleEmbedForVideo"}).Warning("failed get video data")
		return nil
	}

	videoURL, err := redis.GetPresignedURL(video.S3ID, redis.Video)
	if err != nil {
		videoURL, err = s3.CreateGetPresignedVideoURL(video.S3ID)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "HandleEmbedForVideo"}).Warning("failed get video url")
			return nil
		}
	}

	thumbnailURL, err := redis.GetPresignedURL(video.S3ID, redis.Thumbnail)
	if err != nil {
		thumbnailURL, err = s3.GetThumbnailPresignedURL(video.S3ID)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "HandleEmbedForVideo"}).Warning("failed get video thumnail url")
			return nil
		}
	}

	if video.CustomTitle == nil {
		video.CustomTitle = &video.Filename
	}
	if video.CustomDescription == nil {
		video.CustomDescription = &DEFAULT_DESCRIPTION
	}

	meta := fmt.Sprintf(`
		<meta property="og:type" content="video.other" />
		<meta property="og:title" content="%s" />
		<meta property="og:description" content="%s" />
		<meta property="og:url" content="%s/videos/video/%s" />
		<meta property="og:image" content="%s" />
		<meta property="og:video" content="%s" />
		<meta property="og:video:secure_url" content="%s" />
		<meta property="og:video:type" content="video/mp4" />
		<meta property="og:video:width" content="%d" />
		<meta property="og:video:height" content="%d" />
	`,
		html.EscapeString(*video.CustomTitle),
		html.EscapeString(*video.CustomDescription),
		serverHost, id,
		thumbnailURL,
		videoURL,
		videoURL,
		1920, 1080,
	)

	return &meta
}
