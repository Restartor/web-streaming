package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint
	Username string
	Email    string
	Password string
	Role     int
}
