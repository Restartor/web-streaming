package domain

import (
	"time"

	"github.com/lib/pq"
)

type Filem struct {
	ID          uint           `gorm:"primaryKey"`
	Title       string         `gorm:"not null"`
	Description string         `gorm:"not null"`
	Genre       pq.StringArray `gorm:"not null; type:text[]"`
	Year        int            `gorm:"not null"`
	PosterURL   string         `gorm:"not null"`
	Rating      float64        `gorm:"not null; default:0"`
	VideoURL    string         `gorm:"not null"`
}

type UserHistory struct {
	UserID        uint      `gorm:"primaryKey; not null;index"`
	FilmID        uint      `gorm:"primaryKey; not null"`
	LastWatchedAt time.Time `gorm:"not null"`
}

type UserWatchedList struct {
	UserID uint `gorm:"primaryKey; not null"`
	FilmID uint `gorm:"not null"`
}

type FilmRepository interface {
	FindAll() ([]Filem, error)
	FindByTitle(title string) ([]Filem, error)
	Create(filem *Filem) error
	Update(filem *Filem) error
	Delete(id uint) error
}

type FilmService interface {
	GetAllFilms() ([]Filem, error)
	GetFilmByTitle(title string) ([]Filem, error)
	CreateFilm(filem *Filem) error
	UpdateFilm(filem *Filem) error
	DeleteFilm(id uint) error
}

type HistoryRepository interface {
	UserSeeHistory(userID uint) ([]UserHistory, error)
	UserDeleteHistoryID(userID uint, filmID uint) error
	UserDeleteEveryHistory(userID uint) error
	UserRecordWatch(userID uint, filmID uint) error
}

type HistoryService interface {
	GetAllHistory(userID uint) ([]UserHistory, error)
	DeleteHistoryOne(userID uint, filmID uint) error
	DeleteAllHistory(userID uint) error
	RecordWatch(userID uint, filmID uint) error
}

type WatchlistRepository interface {
	UserAddWatchlist(userID uint, filmID uint) error
	RemoveFromWatchlist(userID uint, filmID uint) error
	GetWatchlist(userID uint) ([]UserWatchedList, error)
}

type WatchlistService interface {
	AddToWatchlist(userID uint, filmID uint) error
	RemoveFromWatchlist(userID uint, filmID uint) error
	GetWatchlist(userID uint) ([]UserWatchedList, error)
}
