package handler

import (
	"backend/internal/domain"
	"backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WatchlistHandler struct {
	service domain.WatchlistService
}

func (r *WatchlistHandler) AddToWatchlist(c *gin.Context) {

	var watchlist domain.UserWatchedList

	if err := c.ShouldBindJSON(&watchlist); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid, please try again")
		return
	}

	userID, _ := c.Get("user_id")
	if err := r.service.AddToWatchlist(userID.(uint), watchlist.FilmID); err != nil {
		response.Error(c, http.StatusBadRequest, "error adding to watchlist, please try again")
		return
	}

	response.Success(c, http.StatusOK, "film successfully added to watchlist")

}

func (r *WatchlistHandler) RemoveFromWatchlist(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid film id")
		return
	}
	filmID := uint(id)
	val, _ := c.Get("user_id")
	userID := val.(uint)

	if err := r.service.RemoveFromWatchlist(userID, filmID); err != nil {
		response.Error(c, http.StatusBadRequest, "error removing from watchlist please try again..")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "film successfully deleted from the watchlist!"})

}

func (r *WatchlistHandler) GetWatchlist(c *gin.Context) {
	val, _ := c.Get("user_id")
	userID := val.(uint)

	watchlist, err := r.service.GetWatchlist(userID)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "error loading watchlist")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"watchlist": watchlist})

}

func NewWatchlistHandler(service domain.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{service: service}
}
