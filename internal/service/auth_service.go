package service

import (
	"context"
	"errors"

	"github.com/Restartor/web-streaming/internal/domain"
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
	if user.Password != password {
		return domain.User{}, ErrInvalidCredentials
	}
	return user, nil
}
