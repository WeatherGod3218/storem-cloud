package tags

import (
	"fmt"
	"net/http"

	"github.com/WeatherGod3218/storem-cloud-server/internal/auth"
	"github.com/WeatherGod3218/storem-cloud-server/internal/database"
	"github.com/WeatherGod3218/storem-cloud-server/internal/logging"
	"github.com/WeatherGod3218/storem-cloud-server/internal/models"
	"github.com/WeatherGod3218/storem-cloud-server/internal/users"
	"github.com/gin-gonic/gin"
)

// CreateTag godoc
//
// @Summary      Create a new tag
// @Description  Creates a new tag with the provided name, autogenerating a color.
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        request  body      models.CreateTagRequest  true  "Tag Info"
// @Success      204      none
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Router       /api/v2/tags/create [post]
func CreateTag(c *gin.Context) {
	user := users.GetUserByToken(c.Get("User"))

	if user == nil || !user.CanModifyVideoData() {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "You do not have permission for this action.",
		})
		return
	}
	var req models.CreateTagRequest

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

	database.LogAction(user, fmt.Sprintf(`Created a New Tag "%s"`,
		database.ActionTagName(req.Name),
	))
	c.JSON(http.StatusNoContent, gin.H{})
}

// GetAllTags godoc
//
// @Summary      Gets all the tags
// @Description  Gets all of the tags in the system.
// @Tags         tags
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

	c.JSON(http.StatusOK, list)
}

// GetVideoTags godoc
//
// @Summary      Get Video's Tags
// @Description  Gets all of the tags associated with a Video ID
// @Tags         videos, tags
// @Accept       json
// @Produce      json
// @Param        request  path     string  true  "Video id"
// @Success      200      {object}  []models.Tag
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/tags/video/get/{id} [get]
func GetVideoTags(c *gin.Context) {
	_ = users.GetUserByToken(c.Get("User"))

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

// AddVideoTag godoc
//
// @Summary      Add Tag to Video
// @Description Adds the associated TagId onto a VideoId
// @Tags         videos, tags
// @Accept       json
// @Produce      json
// @Param        request  body      models.ModifyVideoTagRequest  true  "Tag / Video Id"
// @Success      204 	  none
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Router       /api/v2/tags/video/add [post]
func AddVideoTag(c *gin.Context) {
	user := users.GetUserByToken(c.Get("User"))

	if user == nil || !user.CanModifyVideoData() {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "You do not have permission for this action.",
		})
		return
	}

	var req models.ModifyVideoTagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.Warnf("Error with parsing request %s", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request",
		})
		return
	}

	if err := database.AddTagToVideo(user, req.VideoID, req.TagID); err != nil {
		logging.Logger.Warnf("Error adding video tag to database %s", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request",
		})
		return
	}

	database.LogAction(user, fmt.Sprintf("Added Tag %s to Video %s",
		database.ActionTagName(req.TagID),
		database.ActionVideoName(req.VideoID),
	))

	c.JSON(http.StatusNoContent, gin.H{})
}

// DeleteVideoTag godoc
//
// @Summary      Remove Video Tag
// @Description  Removes Tag from a provided Video Id
// @Tags         videos, tags
// @Accept       json
// @Produce      json
// @Param        request  body      models.ModifyVideoTagRequest  true  "Tag / Video Id"
// @Success      204	  none
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Router       /api/v2/tags/video/delete [delete]
func DeleteVideoTag(c *gin.Context) {
	user := users.GetUserByToken(c.Get("User"))

	if user == nil || !user.CanModifyVideoData() {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "You do not have permission for this action.",
		})
		return
	}

	var req models.ModifyVideoTagRequest

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

	database.LogAction(user, fmt.Sprintf("Removed Tag %s from Video %s",
		database.ActionTagName(req.TagID),
		database.ActionVideoName(req.VideoID),
	))

	c.JSON(http.StatusNoContent, gin.H{})
}

func Routes(r *gin.RouterGroup) {
	tags := r.Group("/tags", auth.OAuthMiddleware())
	tags.POST("/create", CreateTag)
	tags.GET("/get", GetAllTags)

	videos := tags.Group("/video")

	videos.GET("/get/:id", GetVideoTags)
	videos.DELETE("/remove", DeleteVideoTag)
	videos.POST("/add", AddVideoTag)
}
