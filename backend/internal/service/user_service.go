package service

import (
	"backend/internal/domain"
	"backend/pkg/logger"
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
	accessDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		accessDuration = time.Minute * 15
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(accessDuration).Unix(),
	}

	jwtsecret := os.Getenv("JWT_SECRET")
	if jwtsecret == "" {
		logger.Log.Fatal().Msg("jwt secret blm di set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtsecret))

	if err != nil {
		return "", "", errors.New("token gagal ter generate")
	}

	refreshTokenString := uuid.New().String()

	refreshDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		refreshDuration = time.Hour * 24 * 7
	}

	rt := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(refreshDuration),
	}

	if err := r.refreshTokenRepo.Create(rt); err != nil {
		return "", "", errors.New("gagal menyimpan refresh token")
	}

	return tokenString, refreshTokenString, nil
}

func (r *UserService) RefreshAccessToken(refreshToken string) (accessToken string, newRefreshToken string, err error) {

	rt, err := r.refreshTokenRepo.FindByToken(refreshToken)

	if err != nil {
		return "", "", errors.New("refresh token is not valid")
	}

	if time.Now().After(rt.ExpiresAt) {
		return "", "", errors.New("refresh token expire")
	}

	user, err := r.repo.FindByID(rt.UserID)

	if err != nil {
		return "", "", errors.New("user tidak ditemukan")
	}

	duration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		duration = time.Minute * 15
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(duration).Unix(),
	}

	jwtsecret := os.Getenv("JWT_SECRET")
	if jwtsecret == "" {
		logger.Log.Fatal().Msg("jwt secret blm di set")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtsecret))

	if err != nil {
		return "", "", errors.New("gagal generate token baru")
	}
	// rotate refresh token
	if err := r.refreshTokenRepo.DeleteByToken(refreshToken); err != nil {
		return "", "", errors.New("gagal rotate refresh token")
	}

	refreshDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		refreshDuration = time.Hour * 24 * 7
	}
	newRT := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(refreshDuration),
	}
	if err := r.refreshTokenRepo.Create(newRT); err != nil {
		return "", "", errors.New("gagal simpan refresh token baru")
	}

	return tokenString, newRT.Token, nil
}

func (r *UserService) UserLogout(userID uint) error {
	return r.refreshTokenRepo.DeleteByUserID(userID)
}

func NewUserService(repo domain.UserRepository, refreshTokenRepo domain.RefreshTokenRepository) domain.UserService {
	return &UserService{repo: repo, refreshTokenRepo: refreshTokenRepo}
}
