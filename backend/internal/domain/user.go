package domain

import "time"

type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
}

type UserRepository interface {
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(user *User) error
}

type UserService interface {
	UserRegister(user *User) error
	UserLogin(email, password string) (string, error)
}
