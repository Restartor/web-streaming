package domain

type Filem struct {
	ID          uint
	Title       string
	Description string
	Genre       string
	Year        int
	PosterURL   string
	Rating      float64
	VideoURL    string
}

type FilmRepository interface {
	FindAll() ([]Filem, error)
	FindByID(id uint) (*Filem, error)
	Create(film *Filem) error
	Update(film *Filem) error
	Delete(id uint) error
}

type FilmService interface {
	GetAllFilms() ([]Filem, error)
	GetFilmByID(id uint) (*Filem, error)
	CreateFilm(film *Filem) error
	UpdateFilm(film *Filem) error
	DeleteFilm(id uint) error
}
