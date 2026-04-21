package repository

import (
	"context"
	"sync"

	"github.com/Restartor/web-streaming/internal/domain"
)

type FilmRepository struct {
	films  []domain.Film
	nextID int64
	mu     sync.RWMutex
}

func NewFilmRepository() *FilmRepository {
	return &FilmRepository{films: []domain.Film{}, nextID: 1}
}

func (r *FilmRepository) Create(_ context.Context, film domain.Film) (domain.Film, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	film.ID = r.nextID
	r.nextID++
	r.films = append(r.films, film)
	return film, nil
}

func (r *FilmRepository) List(_ context.Context) ([]domain.Film, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	films := make([]domain.Film, len(r.films))
	copy(films, r.films)
	return films, nil
}
