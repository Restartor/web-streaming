package main

import (
	"log"
	"web-streaming/config"
	"web-streaming/internal/handler"
	"web-streaming/internal/repository"
	"web-streaming/internal/service"
	"web-streaming/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error", err)
	}

	config.DatabaseConnection()

	// masukkan repo,handler,service untuk menggabungkan mereka bertiga
	userRepository := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	filmRepository := repository.NewFilmRepository(config.DB)
	filmService := service.NewFilmService(filmRepository)
	filmHandler := handler.NewFilmHandler(filmService)

	// routes
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(router, userHandler, filmHandler)

	router.Run(":01010")
}
