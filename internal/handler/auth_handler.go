package handler

import (
	"encoding/json"
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
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if _, err := h.service.Login(r.Context(), req.Email, req.Password); err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
}
