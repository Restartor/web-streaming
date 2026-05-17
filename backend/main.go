package main

import (
	"backend/config"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/routes"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	logger.Init()
	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Gagal mendapatkan env")
	}

	config.DatabaseConnection()

	// masukkan repo,handler,service untuk menggabungkan mereka bertiga
	appConfig := config.LoadAppConfig()

	userRepository := repository.NewUserRepository(config.DB)
	refreshTokenRepository := repository.NewRefreshTokenRepository(config.DB)
	userService := service.NewUserService(userRepository, refreshTokenRepository, appConfig)
	userHandler := handler.NewUserHandler(userService)

	filmRepository := repository.NewFilmRepository(config.DB)
	filmService := service.NewFilmService(filmRepository)
	filmHandler := handler.NewFilmHandler(filmService)

	watchlistRepository := repository.NewWatchlistRepository(config.DB)
	watchlistService := service.NewWatchlistService(watchlistRepository, filmRepository)
	watchlistHandler := handler.NewWatchlistHandler(watchlistService)

	historyRepository := repository.NewHistoryRepository(config.DB)
	historyService := service.NewHistoryService(historyRepository, filmRepository)
	historyHandler := handler.NewHistoryHandler(historyService)
	// routes
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGINS"))
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	routes.SetupRoutes(router, userHandler, filmHandler, watchlistHandler, historyHandler, appConfig)

	osPort := os.Getenv("PORT")
	if osPort == "" {
		osPort = "1010"
	}

	srver := &http.Server{
		Addr:    ":" + osPort,
		Handler: router,
	}

	go func() {
		if err := srver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info().Msg("shutting down server....")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srver.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("server forced to shutdown")
	}
	logger.Log.Info().Msg("server exited")

}
