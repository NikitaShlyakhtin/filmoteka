package main

import (
	"errors"
	"filmoteka/internal/data"
	"filmoteka/internal/validator"
	"net/http"
)

type MovieInput struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	ReleaseDate *string  `json:"release_date"` // RFC3339
	Rating      *float32 `json:"rating"`
	Actors      *[]int   `json:"actors"`
}

type MovieEnvelope struct {
	Movie data.Movie `json:"movie"`
}

type MoviesEnvelope struct {
	Movie []data.Movie `json:"movie"`
}

// @Summary Add a new movie
// @Description Adds a new movie to the database. The request body should include the movie's title, description, release date, rating, and a list of actor IDs.
// @Tags Movies
// @Accept json
// @Produce json
// @Param input body MovieInput true "Movie data"
// @Success 201 {object} MovieEnvelope "Movie successfully created"
// @Failure 400 {object} errorResponse "Client error"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 403 {object} errorResponse "Forbidden"
// @Failure 422 {object} errorResponse "Validation error"
// @Failure 500 {object} errorResponse "Internal server error"
// @Security BasicAuth
// @Router /movies [post]
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

// @Summary Update a movie
// @Description Updates the information of a specific movie in the database. This can be a partial or full update. If a field is not provided in the request body, the current value of that field will be retained.
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param input body MovieInput true "Movie data"
// @Success 200 {object} MovieEnvelope "Movie successfully updated"
// @Failure 400 {object} errorResponse "Client error"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 403 {object} errorResponse "Forbidden"
// @Failure 404 {object} errorResponse "Movie not found"
// @Failure 422 {object} errorResponse "Validation error"
// @Failure 500 {object} errorResponse "Internal server error"
// @Security BasicAuth
// @Router /movies/{id} [patch]
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

		return
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

// @Summary Get a movie
// @Description Retrieves detailed information about a specific movie, including its title, description, release date, rating, and a list of actor IDs.
// @Tags Movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} MovieEnvelope "Movie data"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 404 {object} errorResponse "Movie not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Security BasicAuth
// @Router /movies/{id} [get]
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

// @Summary Get all movies
// @Description Retrieves a list of all movies in the database. Each entry includes the movie's title, description, release date, rating, and a list of actor IDs. The result can be sorted by title, rating, or release date, in ascending or descending order. The default sort order is by rating in descending order.
// @Tags Movies
// @Produce json
// @Param sort query string false "Sort order: title, rating, release_date, -title, -rating, -release_date"
// @Success 200 {object} MoviesEnvelope "List of movies"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 422 {object} errorResponse "Validation error"
// @Failure 500 {object} errorResponse "Internal server error"
// @Security BasicAuth
// @Router /movies [get]
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

// @Summary Delete a movie
// @Description Deletes a specific movie from the database. All information about the movie, including its title, description, release date, rating, and list of actor IDs, will be permanently removed.
// @Tags Movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} MessageEnvelope "Deletion message"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 404 {object} errorResponse "Movie not found"
// @Failure 403 {object} errorResponse "Forbidden"
// @Failure 500 {object} errorResponse "Internal server error"
// @Security BasicAuth
// @Router /movies/{id} [delete]
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

// @Summary Search for movies
// @Description Searches for movies by part of the title or actor name. The query parameters should include the title and/or actor.
// @Tags Search
// @Produce json
// @Param title query string false "Movie title"
// @Param actor query string false "Actor name"
// @Success 200 {object} MoviesEnvelope "List of movies"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 500 {object} errorResponse "Internal server error"
// @Security BasicAuth
// @Router /search [get]
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
