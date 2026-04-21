package repository

import (
	"context"

	"github.com/Restartor/web-streaming/internal/domain"
)

type FilmRepository struct {
	films []domain.Film
}

func NewFilmRepository() *FilmRepository {
	return &FilmRepository{films: []domain.Film{}}
}

func (r *FilmRepository) Create(_ context.Context, film domain.Film) (domain.Film, error) {
	film.ID = int64(len(r.films) + 1)
	r.films = append(r.films, film)
	return film, nil
}

func (r *FilmRepository) List(_ context.Context) ([]domain.Film, error) {
	films := make([]domain.Film, len(r.films))
	copy(films, r.films)
	return films, nil
}
