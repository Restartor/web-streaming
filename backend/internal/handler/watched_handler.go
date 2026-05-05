package handler

import (
	"backend/internal/domain"
	"backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WatchedHandler struct {
	service domain.WatchedService
}

func (r *WatchedHandler) GetAllHistory(c *gin.Context) {

	val, _ := c.Get("user_id")
	userID := val.(uint)
	history, err := r.service.GetAllHistory(userID)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "error loading history..")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"history": history})

}

func (r *WatchedHandler) DeleteHistoryOne(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid film id")
		return
	}
	filmID := uint(id)
	val, _ := c.Get("user_id")

	userID := val.(uint)

	if err := r.service.DeleteHistoryOne(userID, filmID); err != nil {
		response.Error(c, http.StatusBadRequest, "error deleting history, please try again")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "film successfully deleted from history "})
}

func (r *WatchedHandler) DeleteAllHistory(c *gin.Context) {

	var watchlist domain.UserWatchedList

	err := r.service.DeleteAllHistory(watchlist.UserID)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad request try again")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "all history deleted successfully"})

}

func (r *WatchedHandler) AddToWatchlist(c *gin.Context) {

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

func NewWatchlistHandler(service domain.WatchedService) *WatchedHandler {
	return &WatchedHandler{service: service}
}
