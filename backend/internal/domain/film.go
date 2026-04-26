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
