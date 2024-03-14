package main

import (
	"errors"
	"filmoteka/internal/data"
	"filmoteka/internal/validator"
	"net/http"
)

// Добавляет информацию о фильме в базу данных.
func (app *application) addMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		ReleaseDate string  `json:"release_date"` // RFC3339
		Rating      float32 `json:"rating"`
		Actors      []int   `json:"actors"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:       input.Title,
		Description: input.Description,
		ReleaseDate: input.ReleaseDate,
		Rating:      input.Rating,
		Actors:      input.Actors,
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateName):
			v.AddError("title", "movie with this title already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrActorsNotFound):
			v.AddError("actors", "one or more actor IDs do not exist")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*
Обновляет информацию о фильме в базе данных частично или полностью через метод PATCH.
Если поле не передано, оно сохраняет свое значение.
*/
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		ReleaseDate *string  `json:"release_date"` // RFC3339
		Rating      *float32 `json:"rating"`
		Actors      []int    `json:"actors"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Description != nil {
		movie.Description = *input.Description
	}

	if input.ReleaseDate != nil {
		movie.ReleaseDate = *input.ReleaseDate
	}

	if input.Rating != nil {
		movie.Rating = *input.Rating
	}

	if input.Actors != nil {
		movie.Actors = input.Actors
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(*movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateName):
			v.AddError("title", "movie with this title already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrActorsNotFound):
			v.AddError("actors", "one or more actor IDs do not exist")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Возвращает список фильмов с возможностью сортировки и фильтрации.
func (app *application) getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Sort = app.readString(qs, "sort", "-rating")
	input.Filters.SortSafelist = []string{"title", "rating", "release_date", "-title", "-rating", "-release_date"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, err := app.models.Movies.GetAll(input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Удаляет информацию о фильме из базы данных.
func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) searchMovieHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	title := app.readString(qs, "title", "")
	actor := app.readString(qs, "actor", "")

	movies, err := app.models.Movies.Search(title, actor)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
