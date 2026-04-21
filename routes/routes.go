package routes

import (
	"net/http"

	"github.com/Restartor/web-streaming/internal/handler"
	"github.com/Restartor/web-streaming/internal/repository"
	"github.com/Restartor/web-streaming/internal/service"
)

func NewRouter() http.Handler {
	filmRepo := repository.NewFilmRepository()
	filmService := service.NewFilmService(filmRepo)
	filmHandler := handler.NewFilmHandler(filmService)

	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /films", filmHandler.List)
	mux.HandleFunc("POST /login", authHandler.Login)
	return mux
}
