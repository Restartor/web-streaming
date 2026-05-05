package main

import (
	"backend/config"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Gagal mendapatkan env")
	}

	logger.Init()
	config.DatabaseConnection()

	if os.Getenv("JWT_SECRET") == "" {
		logger.Log.Fatal().Err(err).Msg("Jwt_secret belum di set!!")
	}

	// masukkan repo,handler,service untuk menggabungkan mereka bertiga
	userRepository := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	filmRepository := repository.NewFilmRepository(config.DB)
	filmService := service.NewFilmService(filmRepository)
	filmHandler := handler.NewFilmHandler(filmService)

	watchedRepository := repository.NewHistoryRepository(config.DB)
	watchedService := service.NewHistoryService(watchedRepository)
	watchedHandler := handler.NewWatchlistHandler(watchedService)

	// routes
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5174")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	routes.SetupRoutes(router, userHandler, filmHandler, watchedHandler)

	router.Run(":1010")
}
