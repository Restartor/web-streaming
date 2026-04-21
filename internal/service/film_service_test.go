package service

import (
	"context"
	"testing"

	"github.com/Restartor/web-streaming/internal/domain"
)

type filmRepoStub struct {
	films []domain.Film
}

func (s *filmRepoStub) Create(_ context.Context, film domain.Film) (domain.Film, error) {
	s.films = append(s.films, film)
	return film, nil
}

func (s *filmRepoStub) List(_ context.Context) ([]domain.Film, error) {
	return s.films, nil
}

func TestFilmServiceList(t *testing.T) {
	repo := &filmRepoStub{films: []domain.Film{{ID: 1, Title: "Film A"}}}
	svc := NewFilmService(repo)

	films, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(films) != 1 || films[0].Title != "Film A" {
		t.Fatalf("unexpected films result: %+v", films)
	}
}
