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

/*
type WatchedService interface {
	GetAllHistory(userID uint) ([]UserWatchedList, error)
	DeleteHistoryOne(userID uint, filmID uint) error
	DeleteAllHistory(userID uint) error
	AddToWatchlist(userID uint, filmID uint) error
}

*/

func (r *WatchedHandler) GetAllHistory(c *gin.Context) {

	var userwatchedlist domain.UserWatchedList

	history, err := r.service.GetAllHistory(userwatchedlist.UserID)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "error loading history...")
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

	response.Success(c, http.StatusOK, gin.H{"message": "film successfully deleted from history"})
}
