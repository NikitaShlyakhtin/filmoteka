package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"filmoteka/internal/validator"
	"time"
	"unicode/utf8"
)

type Actor struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"` // RFC3339
	Movies    []int     `json:"movies"`
}

type ActorModel interface {
	Insert(actor *Actor) error
	Delete(actor_id int64) error
	Get(id int64) (*Actor, error)
	GetAll() ([]Actor, error)
	Update(actor *Actor) error
}

type ActorDB struct {
	DB *sql.DB
}

type MockActorDB struct {
	Actors map[int64]*Actor
}

var (
	ErrDuplicateName = errors.New("duplicate full name")
)

func ValidateActor(v *validator.Validator, actor *Actor) {
	v.Check(actor.FullName != "", "full_name", "must be provided")
	v.Check(utf8.RuneCountInString(actor.FullName) <= 200, "full_name", "must be no more than 200 symbols")

	v.Check(actor.Gender != "", "gender", "must be provided")
	v.Check(actor.Gender == "male" || actor.Gender == "female", "gender", "must be either male or female")

	v.Check(!actor.BirthDate.Equal(time.Time{}), "birth_date", "must be provided")
	v.Check(actor.BirthDate.Before(time.Now()), "birth_date", "must be a valid date")
}

func (m ActorDB) Insert(actor *Actor) error {
	query := `
		INSERT INTO Actors (full_name, gender, birth_date)
		VALUES ($1, $2, $3)
		RETURNING actor_id`

	args := []interface{}{actor.FullName, actor.Gender, actor.BirthDate}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&actor.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "actors_full_name_key"`:
			return ErrDuplicateName
		default:
			return err
		}
	}

	return nil
}

/*
Удаляет актера и все его связи с фильмами из таблицы Movies_actors,
но не удаляет сами фильмы из таблицы Movies.
*/
func (m ActorDB) Delete(actor_id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE FROM 
			Actors
		WHERE
			actor_id = $1
	`

	result, err := m.DB.ExecContext(ctx, query, actor_id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m ActorDB) Get(id int64) (*Actor, error) {
	var actor Actor

	query := `
	SELECT
		a.actor_id, a.full_name, a.gender, a.birth_date, json_agg(m.movie_id)
	FROM
		Actors a
	LEFT JOIN
		Movies_actors m ON a.actor_id = m.actor_id
	WHERE
		a.actor_id = $1
	GROUP BY
		a.actor_id, a.full_name, a.gender, a.birth_date	
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var movies json.RawMessage

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&actor.ID,
		&actor.FullName,
		&actor.Gender,
		&actor.BirthDate,
		&movies)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	err = json.Unmarshal(movies, &actor.Movies)
	if err != nil {
		return nil, err
	}

	if len(actor.Movies) == 1 && actor.Movies[0] == 0 {
		actor.Movies = []int{}
	}

	return &actor, nil
}

func (m ActorDB) GetAll() ([]Actor, error) {
	var actors []Actor

	query := `
	SELECT
		a.actor_id, a.full_name, a.gender, a.birth_date, json_agg(m.movie_id)
	FROM
		Actors a
	LEFT JOIN
		Movies_actors m ON a.actor_id = m.actor_id
	GROUP BY
		a.actor_id, a.full_name, a.gender, a.birth_date	
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var actor Actor
		var movies json.RawMessage

		err = rows.Scan(&actor.ID,
			&actor.FullName,
			&actor.Gender,
			&actor.BirthDate,
			&movies)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(movies, &actor.Movies)
		if err != nil {
			return nil, err
		}

		if len(actor.Movies) == 1 && actor.Movies[0] == 0 {
			actor.Movies = []int{}
		}

		actors = append(actors, actor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}

func (m ActorDB) Update(actor *Actor) error {
	query := `
		UPDATE 
			Actors
		SET 
			full_name = $1,
			gender = $2,
			birth_date = $3
		WHERE 
			actor_id = $4
	`

	args := []interface{}{
		actor.FullName,
		actor.Gender,
		actor.BirthDate,
		actor.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		case err.Error() == `pq: duplicate key value violates unique constraint "actors_full_name_key"`:
			return ErrDuplicateName
		default:
			return err
		}
	}

	return nil
}

func (m *MockActorDB) Insert(actor *Actor) error {
	if _, found := m.Actors[actor.ID]; found {
		return ErrDuplicateName
	}

	m.Actors[actor.ID] = actor

	return nil
}

func (m *MockActorDB) Get(id int64) (*Actor, error) {
	actor, ok := m.Actors[id]

	if !ok {
		return nil, ErrRecordNotFound
	}

	return actor, nil
}

func (m *MockActorDB) GetAll() ([]Actor, error) {
	var actors []Actor

	for _, actor := range m.Actors {
		actors = append(actors, *actor)
	}

	return actors, nil
}

func (m *MockActorDB) Update(actor *Actor) error {
	if _, found := m.Actors[actor.ID]; !found {
		return ErrRecordNotFound
	}

	for _, a := range m.Actors {
		if a.FullName == actor.FullName && a.ID != actor.ID {
			return ErrDuplicateName
		}
	}

	m.Actors[actor.ID] = actor

	return nil
}

func (m *MockActorDB) Delete(actor_id int64) error {
	if _, found := m.Actors[actor_id]; !found {
		return ErrRecordNotFound
	}

	delete(m.Actors, actor_id)

	return nil
}
