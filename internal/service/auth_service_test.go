package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Restartor/web-streaming/internal/domain"
)

type userRepoStub struct {
	user domain.User
	err  error
}

func (s *userRepoStub) Create(_ context.Context, user domain.User) (domain.User, error) {
	return user, nil
}

func (s *userRepoStub) FindByEmail(_ context.Context, _ string) (domain.User, error) {
	if s.err != nil {
		return domain.User{}, s.err
	}
	return s.user, nil
}

func TestAuthServiceLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := NewAuthService(&userRepoStub{user: domain.User{Email: "demo@example.com", Password: "secret"}})
		user, err := svc.Login(context.Background(), "demo@example.com", "secret")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Email != "demo@example.com" {
			t.Fatalf("unexpected user: %+v", user)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		svc := NewAuthService(&userRepoStub{user: domain.User{Email: "demo@example.com", Password: "secret"}})
		_, err := svc.Login(context.Background(), "demo@example.com", "wrong")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
		}
	})
}
