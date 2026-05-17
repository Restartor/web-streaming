package repository

import (
	"backend/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type WatchlistRepository struct {
	db *gorm.DB
}

func (r *WatchlistRepository) UserAddWatchlist(userID uint, filmID uint) error {

	// buat jadi errornya itu ketika film yang ingin ditambahkan ke watchlist tidak ada di database, karena database sudah diisi dengan data film, jadi ketika ingin menambahkan ke watchlist, maka harus memastikan bahwa film yang ingin ditambahkan ke watchlist sudah ada di database, jika belum ada di database, maka akan error, karena database tidak bisa menemukan data film yang ingin ditambahkan ke watchlist, jadi pastikan bahwa film yang ingin ditambahkan ke watchlist sudah ada di database, jika belum ada di database, maka tambahkan data filmnya terlebih dahulu ke database, baru kemudian tambahkan ke watchlist
	var film domain.Filem

	err := r.db.First(&film, filmID).Error
	if err != nil {
		return err
	}
	watchlist := domain.UserWatchList{UserID: userID, FilmID: filmID}
	return r.db.Create(&watchlist).Error

}

func (r *WatchlistRepository) RemoveFromWatchlist(userID uint, filmID uint) error {

	adaFilm := r.db.Where("user_id = ? AND film_id = ?", userID, filmID).Delete(&domain.UserWatchList{})

	if adaFilm.Error != nil {
		return adaFilm.Error
	}
	if adaFilm.RowsAffected == 0 {
		return errors.New("film not found in watchlist")
	}

	return nil
}

func (r *WatchlistRepository) GetWatchlist(userID uint) ([]domain.UserWatchList, error) {
	var watchlist []domain.UserWatchList
	err := r.db.Where("user_id = ?", userID).Find(&watchlist).Error
	return watchlist, err
}

func NewWatchlistRepository(db *gorm.DB) domain.WatchlistRepository {
	return &WatchlistRepository{db: db}
}
