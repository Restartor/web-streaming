package repository

import (
	"backend/internal/domain"

	"gorm.io/gorm"
)

type FilmRepository struct {
	db *gorm.DB
}

func (r *FilmRepository) FindAll(query domain.PaginationQuery) (domain.PaginatedFilms, error) {
	var filems []domain.Filem
	var total int64

	offset := (query.Page - 1) * query.Limit

	r.db.Model(&domain.Filem{}).Count(&total)
	err := r.db.Offset(offset).Limit(query.Limit).Find(&filems).Error

	return domain.PaginatedFilms{
		Data:  filems,
		Total: total,
		Page:  query.Page,
		Limit: query.Limit,
	}, err
}

func (r *FilmRepository) FindByTitle(title string) ([]domain.Filem, error) {
	var filems []domain.Filem

	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&filems).Error

	return filems, err
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
