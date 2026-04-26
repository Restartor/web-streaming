package routes

import (
	"web-streaming/internal/handler"
	"web-streaming/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(routes *gin.Engine, userHandler *handler.UserHandler) {
	user := routes.Group("/api")
	{
		user.POST("/register", userHandler.Register)
		user.POST("/login", userHandler.Login)
	}

	protected := routes.Group("/api")
	protected.Use(middleware.AuthMiddleware())

}
