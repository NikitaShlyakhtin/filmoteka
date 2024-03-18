package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModel
	Actors ActorModel
	Users  UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieDB{DB: db},
		Actors: ActorDB{DB: db},
		Users:  UserDB{DB: db},
	}
}

func NewMockModels() Models {
	movies := make(map[int64]*Movie)
	actors := make(map[int64]*Actor)
	users := make(map[string]*User)

	hash, _ := GeneratePasswordHash("password123")
	users["user"] = &User{ID: 1, Name: "user", Password: password{hash: hash}, Role: "user"}
	users["admin"] = &User{ID: 2, Name: "admin", Password: password{hash: hash}, Role: "admin"}

	actors[1] = &Actor{
		ID:        1,
		FullName:  "Mock Actor 1",
		Gender:    "male",
		BirthDate: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
		Movies:    []int{1},
	}
	actors[2] = &Actor{
		ID:        2,
		FullName:  "Mock Actor 2",
		Gender:    "female",
		BirthDate: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
		Movies:    []int{1},
	}

	movies[1] = &Movie{
		ID:          1,
		Title:       "Mock Movie 1",
		ReleaseDate: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		Rating:      7.0,
		Actors:      []int64{1, 2},
	}

	return Models{
		Movies: &MockMovieDB{Movies: movies, Actors: actors},
		Actors: &MockActorDB{Actors: actors},
		Users:  &MockUserDB{Users: users},
	}
}
