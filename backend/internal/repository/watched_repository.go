package repository

import (
	"backend/internal/domain"

	"gorm.io/gorm"
)

type HistoryRepository struct {
	db *gorm.DB
}

func (r *HistoryRepository) UserSeeHistory(userID uint) ([]domain.UserHistory, error) {
	var history []domain.UserHistory
	err := r.db.Where("user_id = ?", userID).Find(&history).Error
	// line 15 adalah : mencari semua entri dalam tabel UserHistory yang memiliki user_id yang sesuai dengan nilai userID yang diberikan,
	// dan menyimpan hasilnya dalam variabel history.
	return history, err
}

func (r *HistoryRepository) UserDeleteHistoryID(userID uint, filmID uint) error {
	return r.db.Where("user_id = ? AND film_id = ?", userID, filmID).Delete(&domain.UserHistory{}).Error
}
func (r *HistoryRepository) UserDeleteEveryHistory(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.UserHistory{}).Error
}

func (r *HistoryRepository) UserRecordWatch(userID uint, filmID uint) error {
	record := domain.UserHistory{UserID: userID, FilmID: filmID}
	return r.db.Create(&record).Error
}

func NewHistoryRepository(db *gorm.DB) domain.HistoryRepository {
	return &HistoryRepository{db: db}
}
