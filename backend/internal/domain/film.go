package domain

import (
	"time"

	"github.com/lib/pq"
)

type Filem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"not null" json:"description"`
	Genre       pq.StringArray `gorm:"not null;type:text[]" json:"genre"`
	Year        int            `gorm:"not null" json:"year"`
	PosterURL   string         `gorm:"not null" json:"poster_url"`
	Rating      float64        `gorm:"default:0" json:"rating"`
	VideoURL    string         `gorm:"not null" json:"video_url"`
}
type PaginationQuery struct {
	Page  int `form:"page" binding:"required,min=1"`
	Limit int `form:"limit" binding:"required,min=1"`
}
type PaginatedFilms struct {
	Data  []Filem `json:"films"`
	Total int64   `json:"total"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
}
type UserHistory struct {
	UserID        uint      `gorm:"primaryKey; not null;index"`
	FilmID        uint      `gorm:"primaryKey; not null"`
	LastWatchedAt time.Time `gorm:"not null"`
}

type UserWatchList struct {
	UserID uint `gorm:"primaryKey; not null;index"`
	FilmID uint `gorm:"primaryKey;not null;index"`
}

type FilmRepository interface {
	FindAll(query PaginationQuery) (PaginatedFilms, error)
	FindByTitle(title string) ([]Filem, error)
	Create(filem *Filem) error
	Update(filem *Filem) error
	Delete(id uint) error
}

type FilmService interface {
	GetAllFilms(query PaginationQuery) (PaginatedFilms, error)
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
	GetWatchlist(userID uint) ([]UserWatchList, error)
}

type WatchlistService interface {
	AddToWatchlist(userID uint, filmID uint) error
	RemoveFromWatchlist(userID uint, filmID uint) error
	GetWatchlist(userID uint) ([]UserWatchList, error)
}
