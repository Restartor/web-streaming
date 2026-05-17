package dto

import "github.com/lib/pq"

type CreateFilmInput struct {
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Genre       pq.StringArray `json:"genre" binding:"required"`
	Year        int            `json:"year" binding:"required"`
	PosterURL   string         `json:"poster_url" binding:"required"`
	Rating      float64        `json:"rating" binding:"required"`
	VideoURL    string         `json:"video_url" binding:"required"`
}

type UpdateFilmInput struct {
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Genre       pq.StringArray `json:"genre" binding:"required"`
	Year        int            `json:"year" binding:"required"`
	PosterURL   string         `json:"poster_url" binding:"required"`
	Rating      float64        `json:"rating" binding:"required"`
	VideoURL    string         `json:"video_url" binding:"required"`
}
