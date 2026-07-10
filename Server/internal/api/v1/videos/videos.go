package videos

import (
	"net/http"

	"github.com/WeatherGod3218/weather-reels-server/internal/database"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/WeatherGod3218/weather-reels-server/internal/s3"
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
// @Router       /api/v1/videos/verify [put]
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
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "VerifyVideos"}).Warning("failed to verify videos in the database")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	c.JSON(http.StatusOK, nonVerified)
}

// StartVideoUpload godoc
//
// @Summary      Starts a video upload
// @Description  Starts a multipart video upload, creating S3 URLS and marks the record as pending
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        request  body      models.VideoStartBackupRequest  true  "Video information"
// @Success      200      {object}  models.VideoStartBackupResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v1/videos/backup [post]
func StartVideoUpload(c *gin.Context) {
	var req models.VideoStartBackupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "StartVideoUpload"}).Warning("failed to start video upload")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	videoURL, err := s3.CreatePermanentVideoURL(req.FileName)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "StartVideoUpload"}).Warning("failed to start video upload")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	thumbnailURL, err := s3.CreatePermanentThumbnailURL(req.FileName)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "StartVideoUpload"}).Warning("failed to start video upload")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	fullVideoRequest := models.VideoBackupProcessed{
		FileName:       req.FileName,
		FileLength:     req.FileLength,
		FileSize:       req.FileSize,
		VideoS3URL:     videoURL,
		ThumbnailS3URL: thumbnailURL,
	}

	rowid, err := database.StartVideoUpload(fullVideoRequest)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "StartVideoUpload"}).Warning("failed to start video upload")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	video_urls, uploadID, err := s3.CreatePutPresignedVideoURL(req.FileName, req.FileSize)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "StartVideoUpload"}).Warning("failed to start video upload")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	temp_thumbnail_url, err := s3.CreatePresignedThumbnailURL(req.FileName)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "StartVideoUpload"}).Warning("failed to start video upload")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	c.JSON(http.StatusOK, models.VideoStartBackupResponse{
		RowID:           rowid,
		VideoS3URLs:     video_urls,
		VideoS3UploadID: uploadID,
		ThumbnailS3URL:  temp_thumbnail_url,
	})
}

// CompleteVideoUpload godoc
//
// @Summary      Completes a Video Upload
// @Description  Completes a multipart video upload and marks the record as completed
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        request  body      models.VideoCompleteBackupRequest  true  "Video information"
// @Success      200      {object}  models.SuccessResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v1/videos/complete [post]
func CompleteVideoUpload(c *gin.Context) {
	var req models.VideoCompleteBackupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "CompleteVideoUpload"}).Warning("failed to bind json")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	if err := s3.CompleteMultipartUpload(req.Filename, req.VideoS3UploadID, req.CompletedParts); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "CompleteVideoUpload"}).Warning("failed to complete upload")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	if err := database.CompleteVideoUpload(req.RowID); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "CompleteVideoUpload"}).Warning("failed to update database")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}

// AbortVideoUpload godoc
//
// @Summary      Abort video upload
// @Description  Aborts a multipart video upload and marks the failed record
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        request  body      models.VideoAbortBackupRequest  true  "Video information"
// @Success      200      {object}  models.SuccessResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v1/videos/abort [post]
func AbortVideoUpload(c *gin.Context) {
	var req models.VideoAbortBackupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "AbortVideoUpload"}).Warning("failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	if err := s3.AbortMultipartUpload(req.Filename, req.VideoS3UploadID); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "AbortVideoUpload"}).Warning("failed to abort multiupload")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	if err := database.AbortVideoUpload(req.RowID); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "AbortVideoUpload"}).Warning("failed to update database")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
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
// @Router       /api/v1/videos/group [post]
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
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "AccessVideo"}).Warning("Could not get video url")
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
// @Router       /api/v1/videos/abort [post]
func GetVideoGroup(c *gin.Context) {
	var req models.GetVideoGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "GetVideoGroup"}).Warning("failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	videoEntries, err := database.GetVideoGroup(req.Timestamp, req.RowId)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "GetVideoGroup"}).Warning("failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	videos := make([]models.GetVideoGroupPart, len(videoEntries))
	for i, video := range videoEntries {
		videoUrl, err := s3.CreateGetPresignedVideoURL(video.FileName)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "GetVideoGroup"}).Warning("failed create presigned url for video")
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Unable to process request!",
			})
			return
		}

		thumbnailUrl, err := s3.CreatePresignedThumbnailURL(video.FileName)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "v1/api/videos", "method": "GetVideoGroup"}).Warning("failed create presigned url for thumbnail")
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Unable to process request!",
			})
			return
		}

		videos[i] = models.GetVideoGroupPart{
			RowId:        video.RowID,
			Timestamp:    video.Timestamp,
			VideoURL:     videoUrl,
			ThumbnailURL: thumbnailUrl,
		}
	}

	c.JSON(http.StatusOK, models.GetVideoGroupResponse{
		Videos: videos,
	})
}

func Routes(r *gin.RouterGroup) {
	videos := r.Group("/videos")

	videos.PUT("/verify", VerifyVideos)
	videos.POST("/backup", StartVideoUpload)
	videos.POST("/complete", CompleteVideoUpload)
	videos.POST("/abort", AbortVideoUpload)

	videos.GET("/fetch/*video", AccessVideo)
	videos.POST("/group", GetVideoGroup) //TODO: REPLACE THIS WITH QUERY WHEN AVAILABLE
}
