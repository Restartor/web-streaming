package handler

import (
	"backend/internal/domain"
	"backend/internal/dto"
	"backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FilmHandler struct {
	service domain.FilmService
}

func (r *FilmHandler) GetAllFilms(c *gin.Context) {

	var query domain.PaginationQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		query.Page = 1
		query.Limit = 10
	}
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 20 {
		query.Limit = 10
	}

	filems, err := r.service.GetAllFilms(query)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "error loading films..")
		return
	}
	response.Success(c, http.StatusOK, filems)
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

	var createinput dto.CreateFilmInput
	if err := c.ShouldBindJSON(&createinput); err != nil {
		response.Error(c, http.StatusBadRequest, "create input not valid")
		return
	}

	filem := domain.Filem{
		Title:       createinput.Title,
		Description: createinput.Description,
		Genre:       createinput.Genre,
		Year:        createinput.Year,
		PosterURL:   createinput.PosterURL,
		Rating:      createinput.Rating,
		VideoURL:    createinput.VideoURL,
	}

	if err := r.service.CreateFilm(&filem); err != nil {
		response.Error(c, http.StatusBadRequest, "error adding films please try again")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"message": "Film Added"})
}

func (r *FilmHandler) UpdateFilm(c *gin.Context) {

	var updateinput dto.UpdateFilmInput

	if err := c.ShouldBindJSON(&updateinput); err != nil {
		response.Error(c, http.StatusBadRequest, "update input not valid")
		return
	}
	filem := domain.Filem{
		Title:       updateinput.Title,
		Description: updateinput.Description,
		Genre:       updateinput.Genre,
		Year:        updateinput.Year,
		PosterURL:   updateinput.PosterURL,
		Rating:      updateinput.Rating,
		VideoURL:    updateinput.VideoURL,
	}

	id, err := strconv.Atoi(c.Param("id")) // ambil parameter id dari URL dan konversi ke integer
	// karena id di database itu uint, jadi kita harus mengkonversi id yang sudah diambil dari URL ke uint
	if err != nil {
		response.Error(c, http.StatusBadRequest, "error getting parameters")
		return
	}
	filem.ID = uint(id)

	if err := r.service.UpdateFilm(&filem); err != nil {
		response.Error(c, http.StatusBadRequest, "error updating films please try again")
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
