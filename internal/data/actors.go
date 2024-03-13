package data

import (
	"filmoteka/internal/validator"
	"time"
	"unicode/utf8"
)

type Actor struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"` // RFC3339
}

func ValidateActor(v *validator.Validator, actor *Actor) {
	v.Check(actor.FullName != "", "full_name", "must be provided")
	v.Check(utf8.RuneCountInString(actor.FullName) <= 200, "full_name", "must be no more than 200 symbols")

	v.Check(actor.Gender != "", "gender", "must be provided")
	v.Check(actor.Gender == "male" || actor.Gender == "female", "gender", "must be either male or female")

	v.Check(!actor.BirthDate.Equal(time.Time{}), "birth_date", "must be provided")
	v.Check(actor.BirthDate.Before(time.Now()), "birth_date", "must be a valid date")
}
