package videos

import (
	"net/http"

	"github.com/WeatherGod3218/weather-reels-server/internal/auth"
	"github.com/WeatherGod3218/weather-reels-server/internal/database"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/WeatherGod3218/weather-reels-server/internal/s3"
	"github.com/WeatherGod3218/weather-reels-server/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// VerifyVideos godoc
//
// @Summary      Verify uploaded Videos
// @Description  Verifies videos that are already uploaded, returning a list of ones that are not verified
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        request  body      []string  true  "List of video filenames"
// @Success      200      {array}  []string
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/verify [put]
func VerifyVideos(c *gin.Context) {
	logging.Logger.WithFields(logrus.Fields{"module": "v1/api/videos", "method": "VerifyVideos"}).Info("starting verification!")
	var videoList []string

	if err := c.ShouldBindJSON(&videoList); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "VerifyVideos"}).Warning("failed to bind video list")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	nonVerified, err := database.VerifyVideoList(videoList)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "VerifyVideos"}).Warning("failed to verify videos in the database")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	c.JSON(http.StatusOK, nonVerified)
}

// AccessVideo godoc
//
// @Summary      Get URL for a video
// @Description  Generate a presigned URL for a video file
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        request  path     string  true  "Video information"
// @Success      200      {object}  models.SuccessResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/fetch/{video} [get]
func AccessVideo(c *gin.Context) {
	video := c.Param("video")

	if video == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	url, err := s3.CreateGetPresignedVideoURL(video)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "AccessVideo"}).Warning("Could not get video url")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}
	c.JSON(http.StatusOK, models.AccessVideoResponse{
		URL:   url,
		Video: video,
	})
}

// GetVideoGroup godoc
//
// @Summary      Gets a group of video/thumbnail urls
// @Description  Gets a group of video and thumbnail urls with the given offset. Max of 10
// @Accept       json
// @Produce      json
// @Param        request  body      models.GetVideoGroupRequest  true  "Video information"
// @Success      200      {object}  models.GetVideoGroupResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/abort [post]
func GetVideoGroup(c *gin.Context) {
	var req *models.GetVideoGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "GetVideoGroup"}).Warning("failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	videoEntries, hasMore, err := database.GetVideoGroup(req.Timestamp, req.RowID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "GetVideoGroup"}).Warning("failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	videos := make([]models.GetVideoGroupPartResponse, len(videoEntries))
	for i, video := range videoEntries {
		thumbnailUrl, err := s3.GetThumbnailPresignedURL(video.S3Id)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "GetVideoGroup"}).Warning("failed create presigned url for thumbnail")
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Unable to process request!",
			})
			return
		}

		username := users.GetUsername(video.UserId)

		videos[i] = models.GetVideoGroupPartResponse{
			RowID:             video.RowID,
			S3Id:              video.S3Id,
			CustomTitle:       video.CustomTitle,
			CustomDescription: video.CustomDescription,
			Username:          username,
			Filename:          video.Filename,
			ThumbnailURL:      thumbnailUrl,
			Timestamp:         video.Timestamp,
		}
	}

	resp := models.GetVideoGroupResponse{
		Videos: videos,
	}

	if hasMore {
		last := videos[len(videos)-1]
		resp.Cursor = &models.GetVideoGroupCursor{Timestamp: last.Timestamp, RowID: last.RowID}
	}

	c.JSON(http.StatusOK, resp)
}

// ChangeVideoTitle godoc
//
// @Summary      Updates a videos title
// @Description  Updates a given videoID's title
// @Accept       json
// @Produce      json
// @Param        request  body      models.ChangeVideoTitleRequest  true  "New Title"
// @Success      200      {object}  models.SuccessResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/title [put]
func ChangeVideoTitle(c *gin.Context) {
	var req *models.ChangeVideoTitleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "ChangeVideoTitle"}).Warning("failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	if err := database.ChangeVideoTitle(req.RowID, req.Title); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "ChangeVideoTitle"}).Warning("failed to change title in database")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}

// ChangeVideoDescription godoc
//
// @Summary      Updates a videos description
// @Description  Updates a given videoID's description
// @Accept       json
// @Produce      json
// @Param        request  body      models.ChangeVideoDescriptionRequest  true  "New Title"
// @Success      200      {object}  models.SuccessResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/description [put]
func ChangeVideoDescription(c *gin.Context) {
	var req *models.ChangeVideoDescriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "ChangeVideoDescription"}).Warning("failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	if err := database.ChangeVideoDescription(req.RowID, req.Description); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "ChangeVideoDescription"}).Warning("failed to change description in database")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}

// GetRandomVideo godoc
//
// @Summary      Get random data for a video
// @Description  Gets all the required data for a random video display page
// @Accept       json
// @Produce      json
// @Param        request  path     string  true  "Video information"
// @Success      200      {object}  models.GetRandomVideoResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/random [get]
func GetRandomVideo(c *gin.Context) {
	rowId, err := database.GetRandomVideoData()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "GetVideoData"}).Warning("failed get video data")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	resp := models.GetRandomVideoResponse{
		RowID: rowId,
	}

	c.JSON(http.StatusOK, resp)
}

// GetVideoGroup godoc
//
// @Summary      Get data for a video
// @Description  Gets all the required data for a video display page
// @Accept       json
// @Produce      json
// @Param        request  path     string  true  "Video information"
// @Success      200      {object}  models.GetVideoDataResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/video [get]
func GetVideoData(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		logging.Logger.WithFields(logrus.Fields{"module": "v2/api/videos", "method": "GetVideoData"}).Warning("failed to get video ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	data, err := database.GetVideoData(id)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "GetVideoData"}).Warning("failed get video data")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	videoURL, err := s3.CreateGetPresignedVideoURL(data.S3ID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "GetVideoData"}).Warning("failed get video data")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	logging.Logger.Info(videoURL)

	resp := models.GetVideoDataResponse{
		RowID: data.RowID,
		S3ID:  data.S3ID,

		CustomTitle:       data.CustomTitle,
		CustomDescription: data.CustomDescription,

		Username: users.GetUsername(data.UserId),
		Filename: data.Filename,
		VideoURL: videoURL,
	}

	c.JSON(http.StatusOK, resp)
}

func Routes(r *gin.RouterGroup) {
	videos := r.Group("/videos", auth.OAuthMiddleware())

	videos.PUT("/verify", VerifyVideos)
	videos.GET("/video/:id", GetVideoData)
	videos.GET("/random", GetRandomVideo)
	videos.POST("/group", GetVideoGroup) //TODO: REPLACE THIS WITH QUERY WHEN AVAILABLE

	videos.PUT("/title", ChangeVideoTitle)
	videos.PUT("/description", ChangeVideoDescription)
}
