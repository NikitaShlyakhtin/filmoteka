package data

import (
	"filmoteka/internal/validator"
	"time"
)

type Actor struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"` // RFC3339
}

func ValidateActor(v *validator.Validator, actor *Actor) {
	v.Check(actor.FullName != "", "full_name", "must be provided")
	v.Check(len(actor.FullName) < 500, "full_name", "must be less than 500 bytes")

	v.Check(actor.Gender != "", "gender", "must be provided")
	v.Check(actor.Gender == "male" || actor.Gender == "female", "gender", "must be either male or female")

	v.Check(!actor.BirthDate.Equal(time.Time{}), "birth_date", "must be provided")
	v.Check(actor.BirthDate.Before(time.Now()), "birth_date", "must be a valid date")
}
