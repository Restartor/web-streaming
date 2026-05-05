package repository

import (
	"backend/internal/domain"

	"gorm.io/gorm"
)

type WatchedRepository struct {
	db *gorm.DB
}

func (r *WatchedRepository) UserSeeHistory(userID uint) ([]domain.UserWatchedList, error) {
	var watchlist []domain.UserWatchedList
	err := r.db.Where("user_id = ?", userID).Find(&watchlist).Error
	// line 15 adalah : mencari semua entri dalam tabel UserWatchedList yang memiliki user_id yang sesuai dengan nilai userID yang diberikan,
	// dan menyimpan hasilnya dalam variabel watchlist.
	return watchlist, err
}

func (r *WatchedRepository) UserDeleteHistoryID(userID uint, filmID uint) error {
	return r.db.Where("user_id = ? AND film_id = ?", userID, filmID).Delete(&domain.UserWatchedList{}).Error
}
func (r *WatchedRepository) UserDeleteEveryHistory(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.UserWatchedList{}).Error
}

func (r *WatchedRepository) UserAddWatchlist(userID uint, filmID uint) error {

	watchlist := domain.UserWatchedList{UserID: userID, FilmID: filmID}
	return r.db.Create(&watchlist).Error

}

func NewHistoryRepository(db *gorm.DB) domain.WatchedRepository {
	return &WatchedRepository{db: db}
}
