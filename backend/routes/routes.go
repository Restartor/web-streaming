package routes

import (
	"backend/internal/handler"
	"backend/pkg"
	"backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	routes *gin.Engine,
	userHandler *handler.UserHandler,
	filmHandler *handler.FilmHandler,
) {
	user := routes.Group("/api/v1")
	{
		user.POST("/register", middleware.RateLimiter("5-M"), userHandler.Register)
		user.POST("/login", middleware.RateLimiter("10-M"), userHandler.Login)
		user.GET("/films", filmHandler.GetAllFilms)
		user.GET("/films/search", filmHandler.GetFilmByTitle)
	}

	protected := routes.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(), pkg.AdminOnly())
	{
		protected.POST("/films", filmHandler.CreateFilm)
		protected.PUT("/films/:id", filmHandler.UpdateFilm)
		protected.DELETE("/films/:id", filmHandler.DeleteFilm)
	}

}
