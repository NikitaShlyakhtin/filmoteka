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

type ActorModel struct {
	DB *sql.DB
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

func (m ActorModel) Insert(actor *Actor) error {
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
func (m ActorModel) Delete(actor_id int64) error {
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

func (m ActorModel) Get(id int64) (*Actor, error) {
	var actor Actor

	query := `
	SELECT
		a.actor_id, a.full_name, a.gender, a.birth_date, json_agg(m.movie_id)
	FROM
		Actors a
	JOIN
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

	return &actor, nil
}

func (m ActorModel) GetAll() ([]Actor, error) {
	var actors []Actor

	query := `
	SELECT
		a.actor_id, a.full_name, a.gender, a.birth_date, json_agg(m.movie_id)
	FROM
		Actors a
	JOIN
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

		actors = append(actors, actor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}

func (m ActorModel) Update(actor *Actor) error {
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
