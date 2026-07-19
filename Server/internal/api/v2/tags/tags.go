package tags

import (
	"net/http"

	"github.com/WeatherGod3218/weather-reels-server/internal/database"
	"github.com/WeatherGod3218/weather-reels-server/internal/logging"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateTag godoc
//
// @Summary      Verify uploaded Videos
// @Description  Verifies videos that are already uploaded, returning a list of ones that are not verified
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        request  body      string  true  "Tag to Create"
// @Success      200      {object}  models.SuccessResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/tags/create [post]
func CreateTag(c *gin.Context) {
	var req *models.CreateTagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.Warnf("Error unmarshling request %s", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	unique, err := database.CreateTagRow(req.Name, "testUser")
	if err != nil {
		logging.Logger.Warnf("Error creating database entry %s", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	if !unique {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}

// GetAllTags godoc
//
// @Summary      Gets all the tags
// @Description  Gets all of the tags in the system
// @Tags         tags
// @Accept       json
// @Produce      json
// @Success      200      {object}  []*models.Tag
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/tags/get [get]
func GetAllTags(c *gin.Context) {
	list, err := database.GetAllTags()
	if err != nil {
		logging.Logger.Warnf("Error in getting all the tags %s", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	logging.Logger.Info(list)

	c.JSON(http.StatusOK, list)
}

// GetVideoTags godoc
//
// @Summary      Get Video's Tags
// @Description  Gets all of the tags associated with a video
// @Tags         videos, tags
// @Accept       json
// @Produce      json
// @Param        request  path     string  true  "Video id"
// @Success      200      {object}  []models.Tag
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/tags/video/get/{id} [get]
func GetVideoTags(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusOK, models.ErrorResponse{
			Error: "Unable to process request",
		})
		return
	}
	list, err := database.GetAllTagsOnVideo(id)

	if err != nil {
		logging.Logger.Warnf("Error with database getting tags on a video %s", err)
		c.JSON(http.StatusOK, models.ErrorResponse{
			Error: "Unable to process request",
		})
		return
	}

	c.JSON(http.StatusOK, list)
}

// GetVideoTags godoc
//
// @Summary      Get Video's Tags
// @Description  Gets all of the tags associated with a video
// @Tags         videos, tags
// @Accept       json
// @Produce      json
// @Param        request  path     string  true  "Video id"
// @Success      200      {object}  models.SuccessResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/tags/video/add [post]
func AddVideoTag(c *gin.Context) {
	var req *models.ModifyVideoTagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.Warnf("Error with parsing request %s", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request",
		})
		return
	}

	logging.Logger.Infof("%+v", req)

	if err := database.AddTagToVideo(req.VideoID, req.TagID); err != nil {
		logging.Logger.Warnf("Error adding video tag to database %s", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}

// GetVideoTags godoc
//
// @Summary      Get Video's Tags
// @Description  Gets all of the tags associated with a video
// @Tags         videos, tags
// @Accept       json
// @Produce      json
// @Param        request  path     string  true  "Video id"
// @Success      200      {object}  models.SuccessResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/tags/video/delete [delete]
func DeleteVideoTag(c *gin.Context) {
	var req *models.ModifyVideoTagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.Warnf("Error with parsing request %s", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request",
		})
		return
	}

	if err := database.RemoveTagFromVideo(req.VideoID, req.TagID); err != nil {
		logging.Logger.Warnf("Error removing video tag to database %s", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
	})
}

func Routes(r *gin.RouterGroup) {
	tags := r.Group("/tags")
	tags.POST("/create", CreateTag)
	tags.GET("/get", GetAllTags)

	videos := tags.Group("/video")

	videos.GET("/get/:id", GetVideoTags)
	videos.DELETE("/remove", DeleteVideoTag)
	videos.POST("/add", AddVideoTag)
}
