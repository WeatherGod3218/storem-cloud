package videos

import (
	"fmt"
	"net/http"

	"github.com/WeatherGod3218/storem-cloud-server/internal/auth"
	"github.com/WeatherGod3218/storem-cloud-server/internal/database"
	"github.com/WeatherGod3218/storem-cloud-server/internal/logging"
	"github.com/WeatherGod3218/storem-cloud-server/internal/models"
	"github.com/WeatherGod3218/storem-cloud-server/internal/redis"
	"github.com/WeatherGod3218/storem-cloud-server/internal/s3"
	"github.com/WeatherGod3218/storem-cloud-server/internal/users"
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
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "VerifyVideos"}).Warning("failed to bind video list")
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

// GetVideoGroup godoc
//
// @Summary      Gets a group of video/thumbnail urls
// @Description  Gets a group of video and thumbnail urls with the given offset. Max of 10
// @Accept       json
// @Produce      json
// @Param        request  body      models.GetVideoGroupRequest  true  "Video information"
// @Success      200      {object}  models.GetVideoGroupResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/group [post]
func GetVideoGroup(c *gin.Context) {
	user := users.GetUserByToken(c.Get("User"))

	var req models.GetVideoGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "GetVideoGroup"}).Warning("failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	videoEntries, hasMore, err := database.GetVideoGroup(&req, user)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "GetVideoGroup"}).Warning("failed to get videos from database")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	videos := make([]models.GetVideoGroupPartResponse, len(videoEntries))
	for i, video := range videoEntries {
		var thumbnailURL string

		thumbnailURL, err = redis.GetPresignedURL(video.S3ID, redis.Thumbnail)
		if err != nil {
			thumbnailURL, err = s3.GetThumbnailPresignedURL(video.S3ID)
			if err != nil {
				logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "GetVideoData"}).Warning("failed get video data")
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Error: "Unable to process request!",
				})
				return
			}
		}

		username := users.GetUsername(video.UserId)

		videos[i] = models.GetVideoGroupPartResponse{
			RowID:             video.RowID,
			S3ID:              video.S3ID,
			CustomTitle:       video.CustomTitle,
			CustomDescription: video.CustomDescription,
			Visibility:        video.Visibility,
			Username:          username,
			Filename:          video.Filename,
			ThumbnailURL:      thumbnailURL,
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
// @Success      204      none
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/title [put]
func ChangeVideoTitle(c *gin.Context) {
	user := users.GetUserByToken(c.Get("User"))
	if user == nil || !user.CanModifyVideoData() {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "You do not have permission for this action.",
		})
		return
	}

	var req models.ChangeVideoTitleRequest

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

	database.LogAction(user, fmt.Sprintf(`Changed Video %s Title to "%s"`,
		database.ActionVideoName(req.RowID),
		req.Title,
	))

	c.JSON(http.StatusNoContent, gin.H{})
}

// ChangeVideoVisibility godoc
//
// @Summary      Updates a videos visiblity
// @Description  Updates a given videoID's visiblity
// @Accept       json
// @Produce      json
// @Param        request  body      models.ChangeVideoVisibilityRequest  true  "New Visibility"
// @Success      204      none
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/visbility [put]
func ChangeVideoVisibility(c *gin.Context) {
	user := users.GetUserByToken(c.Get("User"))
	if user == nil || !user.CanModifyVideoData() || !user.CanViewPrivateVideos() {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "You do not have permission for this action.",
		})
		return
	}

	var req models.ChangeVideoVisibilityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "ChangeVideoDescription"}).Warning("failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	if err := database.ChangeVideoVisibility(req.RowID, req.Visibility); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "ChangeVideoDescription"}).Warning("failed to change description in database")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	database.LogAction(user, fmt.Sprintf(`Changed Video %s Visibility to "%s"`,
		database.ActionVideoName(req.RowID),
		req.Visibility,
	))

	c.JSON(http.StatusNoContent, gin.H{})
}

// ChangeVideoDescription godoc
//
// @Summary      Updates a videos description
// @Description  Updates a given videoID's description
// @Accept       json
// @Produce      json
// @Param        request  body      models.ChangeVideoDescriptionRequest  true  "New TDesription"
// @Success      204      none
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/description [put]
func ChangeVideoDescription(c *gin.Context) {
	user := users.GetUserByToken(c.Get("User"))
	if user == nil || !user.CanModifyVideoData() {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "You do not have permission for this action.",
		})
		return
	}

	var req models.ChangeVideoDescriptionRequest

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

	database.LogAction(user, fmt.Sprintf(`Changed Video %s Description to "%s"`,
		database.ActionVideoName(req.RowID),
		req.Description,
	))

	c.JSON(http.StatusNoContent, gin.H{})
}

// GetRandomVideo godoc
//
// @Summary      Get random data for a video
// @Description  Gets all the required data for a random video display page
// @Produce      json
// @Success      200      {object}  models.GetRandomVideoResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/videos/random [get]
func GetRandomVideo(c *gin.Context) {
	user := users.GetUserByToken(c.Get("User"))

	rowId, err := database.GetRandomVideoData(user)
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

// GetVideoData godoc
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
	user := users.GetUserByToken(c.Get("User"))

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

	if user == nil {
		if data.Visibility != "Public" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "You do not have permission for this action."})
			return
		}
	} else if !user.CanViewVideo(data) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "You do not have permission for this action."})
		return
	}

	var videoURL string
	var thumbnailURL string

	videoURL, err = redis.GetPresignedURL(data.S3ID, redis.Video)
	if err != nil {
		videoURL, err = s3.CreateGetPresignedVideoURL(data.S3ID)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "GetVideoData"}).Warning("failed get video url")
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Unable to process request!",
			})
			return
		}
	}

	thumbnailURL, err = redis.GetPresignedURL(data.S3ID, redis.Thumbnail)
	if err != nil {
		thumbnailURL, err = s3.GetThumbnailPresignedURL(data.S3ID)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v2/api/videos", "method": "GetVideoData"}).Warning("failed get video thumbnail url")
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Unable to process request!",
			})
			return
		}
	}

	resp := models.GetVideoDataResponse{
		RowID: data.RowID,
		S3ID:  data.S3ID,

		CustomTitle:       data.CustomTitle,
		CustomDescription: data.CustomDescription,
		Visibility:        data.Visibility,
		Username:          users.GetUsername(data.UserID),
		Filename:          data.Filename,

		ThumbnailURL: thumbnailURL,
		VideoURL:     videoURL,

		Timestamp: data.Timestamp,
		CanModify: ((user != nil) && user.CanModifyVideoData()),
	}

	c.JSON(http.StatusOK, resp)
}

func Routes(r *gin.RouterGroup) {
	videos := r.Group("/videos", auth.OAuthMiddleware())

	videos.PUT("/verify", VerifyVideos)
	videos.GET("/video/:id", GetVideoData)
	videos.GET("/random", GetRandomVideo)

	videos.POST("/group", GetVideoGroup) //TODO: REPLACE THIS WITH QUERY WHEN AVAILABLE

	videos.PUT("/visibility", ChangeVideoVisibility)
	videos.PUT("/title", ChangeVideoTitle)
	videos.PUT("/description", ChangeVideoDescription)
}
