package handler

import (
	"net/http"

	"github.com/Restartor/web-streaming/internal/service"
	"github.com/Restartor/web-streaming/pkg/utils"
)

type FilmHandler struct {
	service *service.FilmService
}

func NewFilmHandler(service *service.FilmService) *FilmHandler {
	return &FilmHandler{service: service}
}

func (h *FilmHandler) List(w http.ResponseWriter, r *http.Request) {
	films, err := h.service.List(r.Context())
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, films)
}
