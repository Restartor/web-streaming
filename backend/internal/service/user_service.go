package service

import (
	"backend/internal/domain"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo domain.UserRepository // db hanya diketahui pas userRepository
}

func (r *UserService) UserRegister(user *domain.User) error {

	user.Role = "user"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	existing, err := r.repo.FindByEmail(user.Email)
	if existing != nil {
		return errors.New("email sudah digunakan")
	}
	existinguser, err := r.repo.FindByUser(user.Username)
	if existinguser != nil {
		return errors.New("username sudah digunakan")
	}

	return r.repo.Create(user)
}

func (r *UserService) UserLogin(email, password string) (string, error) {

	user, err := r.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// bandingkan dengan password

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// generate JWT TOKEN - return tokenstring, nil sama kyk ecommerce repo

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", errors.New("token gagal ter generate")
	}

	return tokenString, nil
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &UserService{repo: repo}
}
