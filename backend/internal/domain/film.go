package domain

type Filem struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Genre       string  `gorm:"not null"`
	Year        int     `gorm:"not null"`
	PosterURL   string  `gorm:"not null"`
	Rating      float64 `gorm:"not null; default:0"`
	VideoURL    string  `gorm:"not null"`
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
