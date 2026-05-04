package handler

import (
	"backend/internal/domain"
	"backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FilmHandler struct {
	service domain.FilmService
}

func (r *FilmHandler) GetAllFilms(c *gin.Context) {

	filems, err := r.service.GetAllFilms()

	if err != nil {
		response.Error(c, http.StatusBadRequest, "error loading films..")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"films": filems})
}

func (r *FilmHandler) GetFilmByTitle(c *gin.Context) {

	title := c.Query("title")

	titlefilems, err := r.service.GetFilmByTitle(title)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "error loading films..")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"films": titlefilems})

}

func (r *FilmHandler) CreateFilm(c *gin.Context) {
	var filem domain.Filem

	if err := c.ShouldBindJSON(&filem); err != nil {
		response.Error(c, http.StatusBadRequest, "create input not valid")
		return
	}

	if err := r.service.CreateFilm(&filem); err != nil {
		response.Error(c, http.StatusBadRequest, "error adding films please try again")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"message": "Film Added"})
}

func (r *FilmHandler) UpdateFilm(c *gin.Context) {

	var filem domain.Filem

	id, err := strconv.Atoi(c.Param("id")) // ambil parameter id dari URL dan konversi ke integer
	filem.ID = uint(id)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "error getting parameters")
	}

	if err := c.ShouldBindJSON(&filem); err != nil {
		response.Error(c, http.StatusBadRequest, "input not valid")
		return
	}
	if err := r.service.UpdateFilm(&filem); err != nil {
		response.Error(c, http.StatusBadRequest, "error updating films")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "Film Updated Successfully!"})

}

func (r *FilmHandler) DeleteFilm(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "error getting parameters id")
		return
	}
	if err := r.service.DeleteFilm(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, "Error deleting films please try again")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "film successfully deleted!"})
}

func NewFilmHandler(service domain.FilmService) *FilmHandler {
	return &FilmHandler{service: service}
}
