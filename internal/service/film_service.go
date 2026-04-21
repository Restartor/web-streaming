package service

import (
	"context"

	"github.com/Restartor/web-streaming/internal/domain"
)

type FilmService struct {
	repository domain.FilmRepository
}

func NewFilmService(repository domain.FilmRepository) *FilmService {
	return &FilmService{repository: repository}
}

func (s *FilmService) Create(ctx context.Context, film domain.Film) (domain.Film, error) {
	return s.repository.Create(ctx, film)
}

func (s *FilmService) List(ctx context.Context) ([]domain.Film, error) {
	return s.repository.List(ctx)
}
