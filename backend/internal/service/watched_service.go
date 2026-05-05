package service

import (
	"backend/internal/domain"
)

type HistoryService struct {
	repo domain.WatchedRepository
}

func (r *HistoryService) GetAllHistory(userID uint) ([]domain.UserWatchedList, error) {
	return r.repo.UserSeeHistory(userID)
}

func (r *HistoryService) DeleteHistoryOne(userID uint, filmID uint) error {
	return r.repo.UserDeleteHistoryID(userID, filmID)
}

func (r *HistoryService) DeleteAllHistory(userID uint) error {
	return r.repo.UserDeleteEveryHistory(userID)
}

func (r *HistoryService) AddToWatchlist(userID uint, filmID uint) error {
	return r.repo.UserAddWatchlist(userID, filmID)
}
