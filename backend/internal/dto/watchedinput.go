package dto

type WatchedInput struct {
	FilmID uint `json:"film_id" binding:"required"`
}
