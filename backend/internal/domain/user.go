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

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
