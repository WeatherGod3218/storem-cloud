package actions

import (
	"net/http"

	"github.com/WeatherGod3218/storem-cloud-server/internal/auth"
	"github.com/WeatherGod3218/storem-cloud-server/internal/database"
	"github.com/WeatherGod3218/storem-cloud-server/internal/logging"
	"github.com/WeatherGod3218/storem-cloud-server/internal/models"
	"github.com/WeatherGod3218/storem-cloud-server/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const MODULE = "v2/api/actions" //stfu sonar

// GetActionGroup godoc
//
// @Summary      Gets a group of logged actions
// @Description  Gets a group of logged actions for auditing
// @Tags actions
// @Accept       json
// @Produce      json
// @Param        request  body      models.GetActionGroupRequest  true  "Action Group"
// @Success      200      {object}  models.GetActionsGroupResponse
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Router       /api/v2/actions/group [post]
func GetActionGroup(c *gin.Context) {
	user := users.GetUserByToken(c.Get("User"))

	if user == nil || !user.CanViewAuditLogs() {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "You are unauthorized to view this content.",
		})
		return
	}

	var req models.GetActionGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": MODULE, "method": "GetActionGroup"}).Warning("failed to bind JSON")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	actions, hasMore, err := database.GetActionGroup(&req, user)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": MODULE, "method": "GetActionGroup"}).Warning("failed to get videos from database")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Unable to process request!",
		})
		return
	}

	resp := models.GetActionsGroupResponse{
		Actions: actions,
	}

	if hasMore {
		last := actions[len(actions)-1]
		resp.Cursor = &models.GetActionGroupCursor{Timestamp: last.Timestamp, RowID: last.RowID}
	}

	c.JSON(http.StatusOK, resp)
}

func Routes(r *gin.RouterGroup) {
	actions := r.Group("/actions", auth.OAuthMiddleware())

	actions.POST("/group", GetActionGroup) //replace with query later
}
