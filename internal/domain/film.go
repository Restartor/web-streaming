package domain

import "context"

type Film struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type FilmRepository interface {
	Create(ctx context.Context, film Film) (Film, error)
	List(ctx context.Context) ([]Film, error)
}
