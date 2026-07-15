package transfer

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/WeatherGod3218/weather-reels-server/internal/auth"
	"github.com/WeatherGod3218/weather-reels-server/internal/database"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/WeatherGod3218/weather-reels-server/internal/s3"
	"github.com/WeatherGod3218/weather-reels-server/internal/thumbnails"
	"github.com/WeatherGod3218/weather-reels-server/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	tusd "github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
)

var Handler *tusd.Handler

func OnVideoUpload() {
	for event := range Handler.CompleteUploads {
		var uploadS3Key string
		if event.Upload.Storage != nil {
			uploadS3Key = event.Upload.Storage["Key"]
		} else {
			logging.Logger.Info("Upload did not have any storage metadata!")

			return
		}

		filelength, err := strconv.ParseFloat(event.Upload.MetaData["video_length"], 64)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "tus", "method": "InitTus"}).Warning("Unable to parse filelength")
			return
		}

		if err := thumbnails.GenerateThumbnailFromVideo(uploadS3Key, filelength); err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "tus", "method": "InitTus", "Id": event.Upload.ID}).Warning("Unable to generate thumbnail!")
			return
		}

		fileModMicro, err := strconv.ParseInt(event.Upload.MetaData["file_mod_date"], 10, 64)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "tus", "method": "InitTus"}).Warning("Unable to parse file mod date")
			return
		}

		userId := event.HTTPRequest.Header.Get("X-User-Id")

		fileModDate := time.UnixMicro(int64(fileModMicro))

		entry := models.VideoDatabaseEntry{
			Filename:    event.Upload.MetaData["filename"],
			FileSize:    event.Upload.Size,
			FileLength:  filelength,
			FileModDate: fileModDate,
			VideoId:     uploadS3Key,
			UserId:      userId,
		}

		rowId, err := database.CreateVideoRow(entry)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "tus", "method": "InitTus"}).Warning("Unable to create database entry for video")
			return
		}

		vidUrl, err := s3.CreateGetPresignedVideoURL(uploadS3Key)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "tus", "method": "InitTus"}).Warning("Error getting video URL")
			return
		}

		thumbnailUrl, err := s3.GetThumbnailPresignedURL(uploadS3Key)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "tus", "method": "InitTus"}).Warning("Error getting thumbnail URL")
			return
		}

		logging.Logger.Info(fmt.Sprintf("Created new database entry for video at row %s", rowId))
		logging.Logger.Infof("Video URL %s", vidUrl)
		logging.Logger.Infof("Thumbnail URL %s", thumbnailUrl)
	}
}

func InitTus() error {

	logging.Logger.Info("AWS S3 client initialized")

	fileStore := s3store.New(os.Getenv("AWS_S3_BUCKET"), s3.S3Client)
	composer := tusd.NewStoreComposer()
	fileStore.UseIn(composer)

	config := &tusd.Config{
		BasePath:      "/api/v2/files",
		StoreComposer: composer,

		RespectForwardedHeaders: true,
		NotifyCompleteUploads:   true,
		PreUploadCreateCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, tusd.FileInfoChanges, error) {
			userId := hook.HTTPRequest.Header.Get("X-User-Id")
			logging.Logger.Infof("Got User Id %s", userId)
			exists := users.VerifyUser(userId)
			if !exists {
				logging.Logger.Infof("User %s does not exists", userId)
				return tusd.HTTPResponse{
					StatusCode: 403,
					Body:       "invalid user",
				}, tusd.FileInfoChanges{}, errors.New("Invalid User Id")
			}
			return tusd.HTTPResponse{}, tusd.FileInfoChanges{}, nil
		},
		// PreFinishResponseCallback: func(hook tusd.HookEvent) (tusd.HTTPResponse, error) {
		// 	go func() {
		// 		if err := thumbnails.GenerateThumbnailFromVideo(hook.Upload.ID); err != nil {
		// 			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "tus", "method": "InitTus", "Id": hook.Upload.ID}).Warning("Unable to generate thumbnail!")
		// 			return
		// 		}

		// 		filelength, err := strconv.ParseFloat(hook.Upload.MetaData["videolength"], 64)
		// 		if err != nil {
		// 			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "tus", "method": "InitTus"}).Warning("Unable to parse filelength")
		// 			return
		// 		}

		// 		entry := models.VideoDatabaseEntry{
		// 			FileName:   hook.Upload.MetaData["filename"],
		// 			FileSize:   hook.Upload.Size,
		// 			FileLength: filelength,
		// 			VideoId:    hook.Upload.ID,
		// 		}

		// 		rowId, err := database.CreateVideoRow(entry)
		// 		if err != nil {
		// 			logging.Logger.WithFields(logrus.Fields{"error": err, "module": "tus", "method": "InitTus"}).Warning("Unable to create database entry for video")
		// 			return
		// 		}

		// 		logging.Logger.Info(fmt.Sprintf("Created new database entry for video at row %s", rowId))
		// 	}()
		// 	return tusd.HTTPResponse{}, nil
		// },
	}

	var err error
	Handler, err = tusd.NewHandler(*config)
	if err != nil {
		return err
	}

	go func() {
		OnVideoUpload()
	}()
	return nil
}

func Routes(r *gin.RouterGroup) {
	videos := r.Group("/files", auth.ApiKeyMiddleware())
	strippedHandler := http.StripPrefix("/api/v2/files", Handler)
	videos.Any("/*any", gin.WrapH(strippedHandler))
}
