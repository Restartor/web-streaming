package service

import (
	"context"
	"errors"

	"github.com/Restartor/web-streaming/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	repository domain.UserRepository
}

func NewAuthService(repository domain.UserRepository) *AuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (domain.User, error) {
	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.User{}, ErrInvalidCredentials
	}
	return user, nil
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
