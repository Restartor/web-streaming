package service

import (
	"backend/internal/domain"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo             domain.UserRepository // db hanya diketahui pas userRepository
	refreshTokenRepo domain.RefreshTokenRepository
}

func (r *UserService) UserRegister(user *domain.User) error {

	user.Role = "user"

	if _, err := r.repo.FindByEmail(user.Email); err == nil {
		return errors.New("email sudah digunakan")
	}

	if _, err := r.repo.FindByUser(user.Username); err == nil {
		return errors.New("username sudah digunakan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return r.repo.Create(user)
}

func (r *UserService) UserLogin(email, password string) (accessToken string, refreshToken string, err error) {

	user, err := r.repo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("email atau password salah")
	}

	// bandingkan dengan password

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("email atau password salah")
	}

	// generate JWT TOKEN - return tokenstring, nil sama kyk ecommerce repo

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", "", errors.New("token gagal ter generate")
	}

	refreshTokenString := uuid.New().String()

	rt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	if err := r.refreshTokenRepo.Create(rt); err != nil {
		return "", "", errors.New("gagal menyimpan refresh token")
	}

	return tokenString, refreshTokenString, nil
}

func (r *UserService) RefreshAccessToken(refreshToken string) (accessToken string, err error) {

	rt, err := r.refreshTokenRepo.FindByToken(refreshToken)

	if err != nil {
		return "", errors.New("refresh token is not valid")
	}

	if time.Now().After(rt.ExpiresAt) {
		return "", errors.New("refresh token expire")
	}

	user, err := r.repo.FindByID(rt.UserID)

	if err != nil {
		return "", errors.New("user tidak ditemukan")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", errors.New("gagal generate token baru")
	}

	return tokenString, nil
}

func NewUserService(repo domain.UserRepository, refreshTokenRepo domain.RefreshTokenRepository) domain.UserService {
	return &UserService{repo: repo, refreshTokenRepo: refreshTokenRepo}
}
