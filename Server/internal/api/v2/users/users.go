package users

import (
	"net/http"

	"github.com/WeatherGod3218/storem-cloud-server/internal/auth"
	"github.com/WeatherGod3218/storem-cloud-server/internal/models"
	"github.com/WeatherGod3218/storem-cloud-server/internal/users"
	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	user := users.GetUserByToken(c.Get("User"))
	if user == nil {
		c.JSON(http.StatusOK, models.GetUserInfoResponse{
			Role:    "None",
			Actions: false,
		})
		return
	}

	c.JSON(http.StatusOK, models.GetUserInfoResponse{
		Role:    user.Role,
		Actions: user.CanViewAuditLogs(),
	})
}

func Routes(r *gin.RouterGroup) {
	userGroup := r.Group("/users", auth.OAuthMiddleware())

	userGroup.GET("/info", GetUserInfo) //replace with query later
}
