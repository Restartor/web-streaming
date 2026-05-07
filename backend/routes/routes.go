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
	watchedHandler *handler.WatchlistHandler,
	historyHandler *handler.HistoryHandler,
) {
	user := routes.Group("/api/v1")
	{
		user.POST("/register", middleware.RateLimiter("5-M"), userHandler.Register)
		user.POST("/login", middleware.RateLimiter("10-M"), userHandler.Login)
		user.GET("/films", filmHandler.GetAllFilms)
		user.GET("/films/search", filmHandler.GetFilmByTitle)
	}

	userAuth := routes.Group("/api/v1")
	userAuth.Use(middleware.AuthMiddleware())
	{
		userAuth.GET("/watchlist", watchedHandler.GetWatchlist)
		userAuth.DELETE("/watchlist/:id", middleware.RateLimiter("3-M"), watchedHandler.RemoveFromWatchlist)
		userAuth.POST("/watchlist", middleware.RateLimiter("5-M"), watchedHandler.AddToWatchlist)
		userAuth.GET("/history", historyHandler.GetAllHistory)
		userAuth.DELETE("/history/:id", middleware.RateLimiter("3-M"), historyHandler.DeleteHistoryOne)
		userAuth.DELETE("/history", middleware.RateLimiter("3-M"), historyHandler.DeleteAllHistory)

	}

	protected := routes.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(), pkg.AdminOnly())
	{
		protected.POST("/films", middleware.RateLimiter("5-M"), filmHandler.CreateFilm)
		protected.PUT("/films/:id", middleware.RateLimiter("3-M"), filmHandler.UpdateFilm)
		protected.DELETE("/films/:id", middleware.RateLimiter("3-M"), filmHandler.DeleteFilm)
	}

}
