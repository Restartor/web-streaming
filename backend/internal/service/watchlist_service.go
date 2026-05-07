package service

import (
	"backend/internal/domain"
)

type WatchlistService struct {
	repo domain.WatchlistRepository
}

func (r *WatchlistService) AddToWatchlist(userID uint, filmID uint) error {
	return r.repo.UserAddWatchlist(userID, filmID)
}

func (r *WatchlistService) RemoveFromWatchlist(userID uint, filmID uint) error {
	return r.repo.RemoveFromWatchlist(userID, filmID)
}

func (r *WatchlistService) GetWatchlist(userID uint) ([]domain.UserWatchedList, error) {
	return r.repo.GetWatchlist(userID)
}

func NewWatchlistService(repo domain.WatchlistRepository) domain.WatchlistService {
	return &WatchlistService{repo: repo}
}
