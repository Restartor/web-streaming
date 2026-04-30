package handler

import (
	"backend/internal/domain"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "error loading films.."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"films": filems})

}

func (r *FilmHandler) GetFilmByTitle(c *gin.Context) {

	title := c.Query("title")

	titlefilems, err := r.service.GetFilmByTitle(title)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "films your looking for is not available"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"films": titlefilems})

}

func (r *FilmHandler) CreateFilm(c *gin.Context) {
	var filem domain.Filem

	if err := c.ShouldBindJSON(&filem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "input not valid"})
		return
	}

	if err := r.service.CreateFilm(&filem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Film Added Succesfully!!"})

}

func (r *FilmHandler) UpdateFilm(c *gin.Context) {

	var filem domain.Filem

	id, err := strconv.Atoi(c.Param("id"))
	filem.ID = uint(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error getting parameters"})
	}

	if err := c.ShouldBindJSON(&filem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "input not valid"})
		return
	}
	if err := r.service.UpdateFilm(&filem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Film Updated Succesfully!!"})

}

func (r *FilmHandler) DeleteFilm(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error getting parameters"})
		return
	}
	if err := r.service.DeleteFilm(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Film Deleted Successfully!"})
}

func NewFilmHandler(service domain.FilmService) *FilmHandler {
	return &FilmHandler{service: service}
}
