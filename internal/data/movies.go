package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"filmoteka/internal/validator"
	"fmt"
	"strings"
	"time"
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

var (
	ErrActorsNotFound = errors.New("one or more actor IDs do not exist")
)

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(utf8.RuneCountInString(movie.Title) <= 150, "title", "must be no more than 150 symbols")

	v.Check(movie.Description != "", "description", "must be provided")
	v.Check(utf8.RuneCountInString(movie.Description) < 1000, "description", "must be no more than 1000 symbols")

	v.Check(movie.ReleaseDate != "", "release_date", "must be provided")

	v.Check(movie.Rating >= 0 && movie.Rating <= 10, "rating", "must be between 0 and 10")

	v.Check(len(movie.Actors) >= 1, "actors", "must contain at least one actor")
}

func (m MovieModel) Insert(movie *Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := checkActorsExistence(m.DB, movie.Actors); err != nil {
		return err
	}

	query := `
		INSERT INTO movies (title, description, release_date, rating)
		VALUES ($1, $2, $3, $4)
		RETURNING movie_id`

	args := []interface{}{movie.Title, movie.Description, movie.ReleaseDate, movie.Rating}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "movies_title_key"`:
			return ErrDuplicateName
		default:
			return err
		}
	}

	query = `
		INSERT INTO movies_actors (movie_id, actor_id)
		VALUES ($1, $2)`

	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, actorID := range movie.Actors {
		_, err := stmt.ExecContext(ctx, movie.ID, actorID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m MovieModel) Delete(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	DELETE FROM 
		Movies
	WHERE
		movie_id = $1`

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	if err = checkAffectedRows(result); err != nil {
		return err
	}

	return nil
}

func (m MovieModel) GetAll(filters Filters) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
			SELECT
				m.movie_id,
				title,
				description,
				release_date,
				rating,
				json_agg(actor_id)
			FROM
				Movies m
			JOIN
				Movies_actors ma ON m.movie_id = ma.movie_id
			GROUP BY
					m.movie_id,
					title,
					description,
					release_date,
					rating
			ORDER BY
				%s %s`, filters.sortColumn(), filters.sortDirection())

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	movies, err := parseMoviesRows(rows)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT
			m.movie_id,
			title,
			description,
			release_date,
			rating,
			json_agg(actor_id)
		FROM
			Movies m
		JOIN
			Movies_actors ma ON m.movie_id = ma.movie_id
		GROUP BY
				m.movie_id,
				title,
				description,
				release_date,
				rating
		HAVING
			m.movie_id = $1`

	var movie Movie
	var actors json.RawMessage

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.ReleaseDate,
		&movie.Rating,
		&actors,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	err = json.Unmarshal(actors, &movie.Actors)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (m MovieModel) Update(movie Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := checkActorsExistence(m.DB, movie.Actors); err != nil {
		return err
	}

	query := `
		UPDATE movies
		SET title = $1, description = $2, release_date = $3, rating = $4
		WHERE movie_id = $5`

	result, err := m.DB.ExecContext(ctx, query, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating, movie.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "movies_title_key"`:
			return ErrDuplicateName
		default:
			return err
		}
	}

	if err = checkAffectedRows(result); err != nil {
		return err
	}

	query = `
		DELETE FROM movies_actors
		WHERE movie_id = $1`

	result, err = m.DB.ExecContext(ctx, query, movie.ID)
	if err != nil {
		return err
	}

	if err = checkAffectedRows(result); err != nil {
		return err
	}

	query = `
		INSERT INTO movies_actors (movie_id, actor_id)
		VALUES ($1, $2)`

	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, actorID := range movie.Actors {
		_, err := stmt.ExecContext(ctx, movie.ID, actorID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m MovieModel) Search(title, actor string) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT
			m.movie_id
		FROM
			Movies m
		JOIN
			Movies_actors ma ON m.movie_id = ma.movie_id
		JOIN
			Actors a ON ma.actor_id = a.actor_id
		WHERE
			title ILIKE '%' || $1 || '%'
		AND
			a.full_name ILIKE '%' || $2 || '%'
		GROUP BY
				m.movie_id`

	rows, err := m.DB.QueryContext(ctx, query, title, actor)
	if err != nil {
		return nil, err
	}

	var moviesIDs []string

	for rows.Next() {
		var movieID string
		err = rows.Scan(&movieID)

		if err != nil {
			return nil, err
		}

		moviesIDs = append(moviesIDs, movieID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	rows.Close()

	query = fmt.Sprintf(`
		SELECT
			m.movie_id,
			title,
			description,
			release_date,
			rating,
			json_agg(actor_id)
		FROM
			Movies m
		JOIN
			Movies_actors ma ON m.movie_id = ma.movie_id
		GROUP BY
				m.movie_id,
				title,
				description,
				release_date,
				rating
		HAVING
			m.movie_id IN (%s)`, strings.Join(moviesIDs, ","))

	rows, err = m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	movies, err := parseMoviesRows(rows)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func parseMoviesRows(rows *sql.Rows) ([]*Movie, error) {
	var movies []*Movie

	for rows.Next() {
		var movie Movie
		var actors json.RawMessage

		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.ReleaseDate,
			&movie.Rating,
			&actors,
		)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(actors, &movie.Actors)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func checkActorsExistence(db *sql.DB, actors []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, actorID := range actors {
		query := `
			SELECT actor_id
			FROM actors
			WHERE actor_id = $1
			`

		var result int
		err := db.QueryRowContext(ctx, query, actorID).Scan(&result)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return ErrActorsNotFound
			default:
				return err
			}
		}
	}

	return nil
}

func checkAffectedRows(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
