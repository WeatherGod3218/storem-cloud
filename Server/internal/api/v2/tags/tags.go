package tags

import (
	"net/http"

	"github.com/WeatherGod3218/weather-reels-server/internal/database"
	"github.com/WeatherGod3218/weather-reels-server/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateTag godoc
//
// @Summary      Verify uploaded Videos
// @Description  Verifies videos that are already uploaded, returning a list of ones that are not verified
// @Tags         videos
// @Accept       json
// @Produce      json
// @Param        request  body      string  true  "Tag to Create"
// @Success      200      {object}  models.SuccessResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /api/v2/tags/create [post]
func CreateTag(c *gin.Context) {
	var req *models.CreateTagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	if err := database.CreateTagRow(req.Name, "testUser"); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
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
}
