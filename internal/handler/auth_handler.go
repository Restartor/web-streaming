package handler

import (
	"net/http"

	"github.com/Restartor/web-streaming/internal/service"
	"github.com/Restartor/web-streaming/pkg/utils"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")
	user, err := h.service.Login(r.Context(), email, password)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}
