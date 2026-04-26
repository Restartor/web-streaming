package main

import (
	"log"
	"web-streaming/config"
	"web-streaming/internal/handler"
	"web-streaming/internal/repository"
	"web-streaming/internal/service"
	"web-streaming/routes"

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

	// routes
	router := gin.Default()

	routes.SetupRoutes(router, userHandler)

	router.Run(":01010")
}
