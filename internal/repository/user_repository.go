package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/Restartor/web-streaming/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	users  map[string]domain.User
	nextID int64
	mu     sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: map[string]domain.User{}, nextID: 1}
}

func (r *UserRepository) Create(_ context.Context, user domain.User) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.nextID++
	r.users[user.Email] = user
	return user, nil
}

func (r *UserRepository) FindByEmail(_ context.Context, email string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[email]
	if !ok {
		return domain.User{}, ErrUserNotFound
	}
	return user, nil
}
