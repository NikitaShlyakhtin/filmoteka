package main

import (
	"errors"
	"filmoteka/internal/data"
	"filmoteka/internal/validator"
	"net/http"
	"strings"
)

type UserEnvelope struct {
	User data.User `json:"user"`
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Create a new user
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param input body CreateUserInput true "User data"
// @Success 201 {object} UserEnvelope "User created successfully"
// @Failure 400 {object} errorResponse "Bad request"
// @Failure 422 {object} errorResponse "Validation failed"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /users [post]
func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name: input.Name,
		Role: strings.ToLower("user"),
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateName):
			v.AddError("name", "a user with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
