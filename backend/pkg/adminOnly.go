package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {

	return func(c *gin.Context) {
		role, exist := c.Get("role")
		if !exist || role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized entry..."})
			c.Abort()
			return
		}
		c.Next()
	}
}
