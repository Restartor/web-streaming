package repository

import (
	"backend/internal/domain"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	adaFilm := r.db.Where("user_id = ? AND film_id = ?", userID, filmID).Delete(&domain.UserHistory{})

	if adaFilm.Error != nil {
		return adaFilm.Error
	}
	if adaFilm.RowsAffected == 0 {
		return errors.New("film not found in history")
	}

	return nil
}
func (r *HistoryRepository) UserDeleteEveryHistory(userID uint) error {

	adaFilm := r.db.Where("user_id = ?", userID).Delete(&domain.UserHistory{})
	if adaFilm.Error != nil {
		return adaFilm.Error
	}
	if adaFilm.RowsAffected == 0 {
		return errors.New("history not found")
	}
	return nil
}

func (r *HistoryRepository) UserRecordWatch(userID uint, filmID uint) error {

	var film domain.Filem

	err := r.db.First(&film, filmID).Error
	if err != nil {
		return err
	}

	record := domain.UserHistory{
		UserID:        userID,
		FilmID:        filmID,
		LastWatchedAt: time.Now(),
	}

	dapathistory := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "film_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_watched_at"}),
	}).Create(&record)

	if dapathistory.Error != nil {
		return dapathistory.Error
	}

	return nil
	// penjelasan kode di atas adalah: fungsi UserRecordWatch menerima tiga parameter: userID, filmID, dan waktu
	// saat ini (time.Now()). Fungsi ini membuat sebuah record baru dalam tabel UserHistory
	// dengan nilai userID, filmID, dan lastWatchedAt yang diatur ke waktu saat ini.
	// Jika sudah ada entri dengan kombinasi user_id dan film_id yang sama, maka nilai last_watched_at
	//  akan diperbarui dengan waktu saat ini.
}

func NewHistoryRepository(db *gorm.DB) domain.HistoryRepository {
	return &HistoryRepository{db: db}
}
