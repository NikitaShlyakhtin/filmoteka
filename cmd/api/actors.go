package main

import (
	"errors"
	"filmoteka/internal/data"
	"filmoteka/internal/validator"
	"net/http"
	"strings"
	"time"
)

type ActorInput struct {
	FullName  *string    `json:"full_name"`
	Gender    *string    `json:"gender"`
	BirthDate *time.Time `json:"birth_date"` // RFC3339
}

type ActorEnvelope struct {
	Actor data.Actor `json:"actor"`
}

type ActorsEnvelope struct {
	Actor []data.Actor `json:"actor"`
}

type MessageEnvelope struct {
	Message string `json:"message"`
}

// @Summary Add new actor
// @Description Adds a new actor to the database. The request body should include the actor's full name, gender, and birth date. Once the actor is added, he can be associated with movies.
// @Tags Actors
// @Accept json
// @Produce json
// @Param input body ActorInput true "Actor data"
// @Success 201 {object} ActorEnvelope "Actor successfully created"
// @Failure 400 {object} errorResponse "Client error"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 403 {object} errorResponse "Forbidden"
// @Failure 422 {object} errorResponse "Validation error"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /actors [post]
// @Security BasicAuth
func (app *application) addActorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName  string    `json:"full_name"`
		Gender    string    `json:"gender"`
		BirthDate time.Time `json:"birth_date"` // RFC3339
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	actor := &data.Actor{
		FullName:  input.FullName,
		Gender:    strings.ToLower(input.Gender),
		BirthDate: input.BirthDate,
	}

	v := validator.New()
	if data.ValidateActor(v, actor); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Actors.Insert(actor)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateName):
			v.AddError("full_name", "actor with this full name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	if actor.Movies == nil {
		actor.Movies = []int{}
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"actor": actor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Update actor
// @Description Updates the information of a specific actor in the database. This can be a partial or full update. If a field is not provided in the request body, the current value of that field will be retained.
// @Tags Actors
// @Accept json
// @Produce json
// @Param id path int true "Actor ID"
// @Param input body ActorInput true "Actor data"
// @Success 200 {object} ActorEnvelope "Actor successfully updated"
// @Failure 400 {object} errorResponse "Client error"
// @Failure 403 {object} errorResponse "Forbidden"
// @Failure 404 {object} errorResponse "Actor not found"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 422 {object} errorResponse "Validation error"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /actors/{id} [patch]
// @Security BasicAuth
func (app *application) updateActorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		FullName  *string    `json:"full_name"`
		Gender    *string    `json:"gender"`
		BirthDate *time.Time `json:"birth_date"` // RFC3339
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	actor, err := app.models.Actors.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	if *input.FullName != "" {
		actor.FullName = *input.FullName
	}

	if *input.Gender != "" {
		actor.Gender = strings.ToLower(*input.Gender)
	}

	if !(*input.BirthDate).Equal(time.Time{}) {
		actor.BirthDate = *input.BirthDate
	}

	v := validator.New()
	if data.ValidateActor(v, actor); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Actors.Update(actor)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateName):
			v.AddError("full_name", "actor with this full name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"actor": actor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get actor by ID
// @Description Retrieves information about specific actor from the database, including actor's full name, gender, birth date, and a list of movies IDs he have appeared in. If the actor doesn't appear in any movies, the list will be empty.
// @Tags Actors
// @Accept json
// @Produce json
// @Param id path int true "Actor ID"
// @Success 200 {object} ActorEnvelope "Actor data"
// @Failure 400 {object} errorResponse "Client error"
// @Failure 403 {object} errorResponse "Forbidden"
// @Failure 404 {object} errorResponse "Actor not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Router /actors/{id} [get]
// @Security BasicAuth
func (app *application) getActorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	actor, err := app.models.Actors.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"actor": actor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Get actors
// @Description Retrieves a list of all actors in the database. Each entry includes the actor's full name, gender, birth date, and a list of movies they have appeared in. If the actor doesn't appear in any movies, the list will be empty.
// @Tags Actors
// @Accept json
// @Produce json
// @Success 200 {object} ActorsEnvelope "Actors data"
// @Failure 400 {object} errorResponse "Client error"
// @Failure 500 {object} errorResponse "Internal server error"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Router /actors [get]
// @Security BasicAuth
func (app *application) getActorsHandler(w http.ResponseWriter, r *http.Request) {
	actors, err := app.models.Actors.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"actors": actors}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary Delete actor
// @Description Deletes a specific actor from the database. All information about the actor, including their full name, gender, birth date, and list of movies, will be permanently removed.
// @Tags Actors
// @Accept json
// @Produce json
// @Param id path int true "Actor ID"
// @Success 204 {object} MessageEnvelope "Actor successfully deleted"
// @Failure 400 {object} errorResponse "Client error"
// @Failure 403 {object} errorResponse "Forbidden"
// @Failure 404 {object} errorResponse "Actor not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Router /actors/{id} [delete]
// @Security BasicAuth
func (app *application) deleteActorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Actors.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "actor successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
