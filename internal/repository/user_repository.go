package repository

import (
	"context"
	"errors"

	"github.com/Restartor/web-streaming/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	users map[string]domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: map[string]domain.User{}}
}

func (r *UserRepository) Create(_ context.Context, user domain.User) (domain.User, error) {
	user.ID = int64(len(r.users) + 1)
	r.users[user.Email] = user
	return user, nil
}

func (r *UserRepository) FindByEmail(_ context.Context, email string) (domain.User, error) {
	user, ok := r.users[email]
	if !ok {
		return domain.User{}, ErrUserNotFound
	}
	return user, nil
}
