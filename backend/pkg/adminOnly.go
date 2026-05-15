package pkg

import (
	"backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {

	return func(c *gin.Context) {
		role, exist := c.Get("role")
		if !exist || role != "admin" {
			response.Error(c, http.StatusForbidden, "forbidden access")
			c.Abort()
			return
		}
		c.Next()
	}
}
