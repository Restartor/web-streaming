package repository

import (
	"web-streaming/internal/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	return &user, err

}

func (r *userRepository) FindByEmail() {

}

func (r *userRepository) Create() {

}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}
