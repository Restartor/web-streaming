package service

import (
	"backend/internal/domain"
	"errors"
)

type WatchlistService struct {
	repo     domain.WatchlistRepository
	filmRepo domain.FilmRepository
}

func (r *WatchlistService) AddToWatchlist(userID uint, filmID uint) error {

	if _, err := r.filmRepo.FindByID(filmID); err != nil {
		return errors.New("film not found")
	}

	return r.repo.UserAddWatchlist(userID, filmID)
}

func (r *WatchlistService) RemoveFromWatchlist(userID uint, filmID uint) error {
	return r.repo.RemoveFromWatchlist(userID, filmID)
}

func (r *WatchlistService) GetWatchlist(userID uint) ([]domain.UserWatchList, error) {
	return r.repo.GetWatchlist(userID)
}

func NewWatchlistService(repo domain.WatchlistRepository, filmRepo domain.FilmRepository) domain.WatchlistService {
	return &WatchlistService{repo: repo, filmRepo: filmRepo}
}
