package routes

import (
	"backend/config"
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
	cfg config.AppConfig,
) {
	user := routes.Group("/api/v1")
	{
		user.POST("/register", middleware.RateLimiter("5-M"), userHandler.Register)
		user.POST("/login", middleware.RateLimiter("10-M"), userHandler.Login)
		user.POST("/refresh-token", middleware.RateLimiter("10-M"), userHandler.RefreshToken)
		user.GET("/films", filmHandler.GetAllFilms)
		user.GET("/films/search", filmHandler.GetFilmByTitle)
	}

	userAuth := routes.Group("/api/v1")
	userAuth.Use(middleware.AuthMiddleware(cfg))
	{
		userAuth.GET("/watchlist", middleware.RateLimiter("10-M"), watchedHandler.GetWatchlist)
		userAuth.DELETE("/watchlist/:id", middleware.RateLimiter("3-M"), watchedHandler.RemoveFromWatchlist)
		userAuth.POST("/watchlist", middleware.RateLimiter("5-M"), watchedHandler.AddToWatchlist)
		userAuth.GET("/history", middleware.RateLimiter("10-M"), historyHandler.GetAllHistory)
		userAuth.POST("/history/:id", middleware.RateLimiter("10-M"), historyHandler.RecordWatch)
		userAuth.DELETE("/history/:id", middleware.RateLimiter("3-M"), historyHandler.DeleteHistoryOne)
		userAuth.DELETE("/history", middleware.RateLimiter("3-M"), historyHandler.DeleteAllHistory)
		userAuth.POST("/logout", middleware.RateLimiter("3-M"), userHandler.Logout)

	}

	protected := routes.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg), pkg.AdminOnly())
	{
		protected.POST("/films", middleware.RateLimiter("5-M"), filmHandler.CreateFilm)
		protected.PUT("/films/:id", middleware.RateLimiter("3-M"), filmHandler.UpdateFilm)
		protected.DELETE("/films/:id", middleware.RateLimiter("3-M"), filmHandler.DeleteFilm)
	}

}
