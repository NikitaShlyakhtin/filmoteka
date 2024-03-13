package main

import (
	"errors"
	"filmoteka/internal/data"
	"filmoteka/internal/validator"
	"net/http"
	"strings"
	"time"
)

// Добавляет информацию об актёре в базу данных.
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

	err = app.writeJSON(w, http.StatusCreated, envelope{"actor": actor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*
Обновляет информацию об актёре в базе данных частично или полностью через метод PATCH.
Если поле не передано, оно сохраняет свое значение.
*/
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

	if input.FullName != nil {
		actor.FullName = *input.FullName
	}

	if input.Gender != nil {
		actor.Gender = strings.ToLower(*input.Gender)
	}

	if input.BirthDate != nil {
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

// Возвращает информацию об актере и о фильмах, в которых они снимались.
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

// Возвращает список актёров с информацией о фильмах, в которых они снимались.
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

// Удаляет информацию об актёре из базы данных.
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
