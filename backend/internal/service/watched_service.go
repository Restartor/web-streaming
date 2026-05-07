package service

import (
	"backend/internal/domain"
)

type HistoryService struct {
	repo domain.HistoryRepository
}

func (r *HistoryService) GetAllHistory(userID uint) ([]domain.UserHistory, error) {
	return r.repo.UserSeeHistory(userID)
}

func (r *HistoryService) DeleteHistoryOne(userID uint, filmID uint) error {
	return r.repo.UserDeleteHistoryID(userID, filmID)
}

func (r *HistoryService) DeleteAllHistory(userID uint) error {
	return r.repo.UserDeleteEveryHistory(userID)
}

func (r *HistoryService) RecordWatch(userID uint, filmID uint) error {
	return r.repo.UserRecordWatch(userID, filmID)
}

func NewHistoryService(repo domain.HistoryRepository) domain.HistoryService {
	return &HistoryService{repo: repo}
}

// finished
