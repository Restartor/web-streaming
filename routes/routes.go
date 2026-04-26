package routes

import (
	"web-streaming/internal/handler"
	"web-streaming/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	routes *gin.Engine,
	userHandler *handler.UserHandler,
	filmHandler *handler.FilmHandler,
) {
	user := routes.Group("/api")
	{
		user.POST("/register", userHandler.Register)
		user.POST("/login", userHandler.Login)
	}

	protected := routes.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/films", filmHandler.GetAllFilms)
		protected.GET("/films/search", filmHandler.GetFilmByTitle)
		protected.POST("/films", filmHandler.CreateFilm)
		protected.PUT("/films/:id", filmHandler.UpdateFilm)
		protected.DELETE("/films/:id", filmHandler.DeleteFilm)
	}

}
