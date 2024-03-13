package data

import (
	"database/sql"
	"filmoteka/internal/validator"
	"unicode/utf8"
)

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	ReleaseDate string  `json:"release_date"` // RFC3339
	Rating      float32 `json:"rating"`
	Actors      []int   `json:"actors"`
}

type MovieModel struct {
	DB *sql.DB
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(utf8.RuneCountInString(movie.Title) <= 150, "title", "must be no more than 150 symbols")

	v.Check(movie.Description != "", "description", "must be provided")
	v.Check(utf8.RuneCountInString(movie.Description) < 1000, "description", "must be no more than 1000 symbols")

	v.Check(movie.ReleaseDate != "", "release_date", "must be provided")

	v.Check(movie.Rating >= 0 && movie.Rating <= 10, "rating", "must be between 0 and 10")

	v.Check(len(movie.Actors) >= 1, "actors", "must contain at least one actor")
}
