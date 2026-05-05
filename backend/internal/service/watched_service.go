package service

import (
	"backend/internal/domain"
)

type WatchedService struct {
	repo domain.WatchedRepository
}

func (r *WatchedService) GetAllHistory(userID uint) ([]domain.UserWatchedList, error) {
	return r.repo.UserSeeHistory(userID)
}

func (r *WatchedService) DeleteHistoryOne(userID uint, filmID uint) error {
	return r.repo.UserDeleteHistoryID(userID, filmID)
}

func (r *WatchedService) DeleteAllHistory(userID uint) error {
	return r.repo.UserDeleteEveryHistory(userID)
}

func (r *WatchedService) AddToWatchlist(userID uint, filmID uint) error {
	return r.repo.UserAddWatchlist(userID, filmID)
}

func NewHistoryService(repo domain.WatchedRepository) domain.WatchedService {
	return &WatchedService{repo: repo}
}

// finished
