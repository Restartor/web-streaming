package handler

import (
	"net/http"
	"web-streaming/internal/domain"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service domain.UserService
}

func (r *UserHandler) Register(c *gin.Context) {
	var input domain.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := r.service.UserRegister(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Register Success!"})
}

func (r *UserHandler) Login(c *gin.Context) {
	var user domain.LoginInput

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong User/Password"})
	}

	token, err := r.service.UserLogin(user.Email, user.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}
