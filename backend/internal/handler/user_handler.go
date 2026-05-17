package handler

import (
	"backend/internal/domain"
	"backend/internal/dto"
	"backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service domain.UserService
}

func (r *UserHandler) Register(c *gin.Context) {
	var input dto.RegisterInput
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
		response.Error(c, http.StatusInternalServerError, "Failed to register, Username or Email already exists, please try again")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"message": "Berhasil Register!"})
}

func (r *UserHandler) Login(c *gin.Context) {
	var user dto.LoginInput

	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid Input Data")
		return
	}

	accessToken, refreshToken, err := r.service.UserLogin(user.Email, user.Password)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Wrong email or password, please try again")
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("access_token", accessToken, 900, "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/", "", true, true)

	response.Success(c, http.StatusOK, gin.H{"message": "Login successful"})
}

func (r *UserHandler) RefreshToken(c *gin.Context) {
	// ambil refresh token dari cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "refresh token not found")
		return
	}

	accessToken, newRefreshToken, err := r.service.RefreshAccessToken(refreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("access_token", accessToken, 900, "/", "", true, true)
	c.SetCookie("refresh_token", newRefreshToken, 604800, "/", "", true, true)

	response.Success(c, http.StatusOK, gin.H{"message": "Token refreshed successfully"})
}

func (r *UserHandler) Logout(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if err := r.service.UserLogout(userID.(uint)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to logout")
		return
	}
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	response.Success(c, http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}
