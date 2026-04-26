package repository

import (
	"web-streaming/internal/domain"

	"gorm.io/gorm"
)

type FilmRepository struct {
	db *gorm.DB
}

func (r *FilmRepository) FindAll() ([]domain.Filem, error) {
	var filems []domain.Filem
	err := r.db.Find(&filems).Error
	return filems, err
}

func (r *FilmRepository) FindByTitle(title string) (*domain.Filem, error) {
	var filem domain.Filem

	err := r.db.Where("title = ?", title).First(&filem).Error

	return &filem, err

}

func (r *FilmRepository) Create(filem *domain.Filem) error {
	return r.db.Create(filem).Error
}

func (r *FilmRepository) Update(filem *domain.Filem) error {
	return r.db.Save(filem).Error
}

func (r *FilmRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Filem{}, id).Error
}

func NewFilmRepository(db *gorm.DB) domain.FilmRepository {
	return &FilmRepository{db: db}
}
