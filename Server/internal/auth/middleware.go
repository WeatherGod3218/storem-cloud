package auth

import (
	"github.com/gin-gonic/gin"
)

func auth_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
