package repository

import (
	"backend/internal/domain"

	"gorm.io/gorm"
)

type HistoryRepository struct {
	db *gorm.DB
}

func (r *HistoryRepository) UserSeeHistory(userID uint) ([]domain.UserWatchedList, error) {
	var watchlist []domain.UserWatchedList
	err := r.db.Where("user_id = ?", userID).Find(&watchlist).Error
	// line 15 adalah : mencari semua entri dalam tabel UserWatchedList yang memiliki user_id yang sesuai dengan nilai userID yang diberikan,
	// dan menyimpan hasilnya dalam variabel watchlist.
	return watchlist, err
}

func (r *HistoryRepository) UserDeleteHistoryID(userID uint, filmID uint) error {
	return r.db.Where("user_id = ? AND film_id = ?", userID, filmID).Delete(&domain.UserWatchedList{}).Error
}
func (r *HistoryRepository) UserDeleteEveryHistory(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.UserWatchedList{}).Error
}

func (r *HistoryRepository) UserAddWatchlist(userID uint, filmID uint) error {

	watchlist := domain.UserWatchedList{UserID: userID, FilmID: filmID}
	return r.db.Create(&watchlist).Error

}

func NewHistoryRepository(db *gorm.DB) domain.WatchedRepository {
	return &HistoryRepository{db: db}
}
