package handler

import (
	"backend/internal/domain"
	"backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service domain.UserService
}

func (r *UserHandler) Register(c *gin.Context) {
	var input domain.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Register input Invalid")
		return
	}

	user := domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := r.service.UserRegister(&user); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to register, please try again")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"message": "Berhasil Register!"})
}

func (r *UserHandler) Login(c *gin.Context) {
	var user domain.LoginInput

	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid Input Data")
		return
	}

	token, err := r.service.UserLogin(user.Email, user.Password)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Wrong email or password, please try again")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"token": token})

}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}
