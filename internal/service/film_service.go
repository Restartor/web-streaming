package service

import "web-streaming/internal/domain"

type FilmService struct {
	repo domain.FilmRepository
}

func (r *FilmService) GetAllFilms() ([]domain.Filem, error) {
	return r.repo.FindAll()
}

func (r *FilmService) GetFilmByTitle(title string) ([]domain.Filem, error) {
	return r.repo.FindByTitle(title)
}

func (r *FilmService) CreateFilm(filem *domain.Filem) error {
	return r.repo.Create(filem)
}

func (r *FilmService) UpdateFilm(filem *domain.Filem) error {
	return r.repo.Update(filem)
}

func (r *FilmService) DeleteFilm(id uint) error {
	return r.repo.Delete(id)
}

func NewFilmService(repo domain.FilmRepository) domain.FilmService {
	return &FilmService{repo: repo}
}
