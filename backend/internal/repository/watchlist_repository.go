package repository

import (
	"backend/internal/domain"

	"gorm.io/gorm"
)

type WatchlistRepository struct {
	db *gorm.DB
}

func (r *WatchlistRepository) UserAddWatchlist(userID uint, filmID uint) error {

	watchlist := domain.UserWatchList{UserID: userID, FilmID: filmID}
	return r.db.Create(&watchlist).Error

}

func (r *WatchlistRepository) RemoveFromWatchlist(userID uint, filmID uint) error {
	return r.db.Where("user_id = ? AND film_id = ?", userID, filmID).Delete(&domain.UserWatchList{}).Error
}

func (r *WatchlistRepository) GetWatchlist(userID uint) ([]domain.UserWatchList, error) {
	var watchlist []domain.UserWatchList
	err := r.db.Where("user_id = ?", userID).Find(&watchlist).Error
	return watchlist, err
}

func NewWatchlistRepository(db *gorm.DB) domain.WatchlistRepository {
	return &WatchlistRepository{db: db}
}
