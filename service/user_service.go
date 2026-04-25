package service

import (
	"errors"
	"web-streaming/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo domain.UserRepository // db hanya diketahui pas userRepository
}

func (r *UserService) UserRegister(user *domain.User) error {
	return r.repo.Create(user)
}

func (r *UserService) UserLogin(email, password string) (string, error) {

	user, err := r.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("email not found!")
	}

	// bandingkan dengan password

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("password yang dimasukkan salah!")
	}

	// generate JWT TOKEN - return tokenstring, nil sama kyk ecommerce repo
}
