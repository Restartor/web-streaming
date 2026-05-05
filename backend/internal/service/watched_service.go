package service

import (
	"backend/internal/domain"
	"errors"
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

	if err := r.repo.UserAddWatchlist(userID, filmID); err != nil {
		return errors.New("error data tidak valid")
	}

	return nil
}

func NewHistoryService(repo domain.WatchedRepository) domain.WatchedService {
	return &WatchedService{repo: repo}
}

// finished
