package handler

import (
	"backend/internal/domain"
	"backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	service domain.HistoryService
}

func (r *HistoryHandler) GetAllHistory(c *gin.Context) {

	val, _ := c.Get("user_id")
	userID := val.(uint)
	history, err := r.service.GetAllHistory(userID)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "error loading history..")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"history": history})

}

func (r *HistoryHandler) DeleteHistoryOne(c *gin.Context) {

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

func (r *HistoryHandler) DeleteAllHistory(c *gin.Context) {

	val, _ := c.Get("user_id") // karena user_id itu disimpan di context dengan tipe data uint, jadi kita harus melakukan type assertion untuk mengubahnya menjadi uint
	userID := val.(uint)       // jelasinnya userID itu uint, jadi kita type assertion ke uint

	err := r.service.DeleteAllHistory(userID)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "bad request try again")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "all history deleted successfully"})

}

func NewHistoryHandler(service domain.HistoryService) *HistoryHandler {
	return &HistoryHandler{service: service}
}
