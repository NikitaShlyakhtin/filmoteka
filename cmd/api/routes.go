package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/users", app.createUserHandler)

	router.HandlerFunc(http.MethodPost, "/actors", app.requireRoleAdmin(app.addActorHandler))
	router.HandlerFunc(http.MethodGet, "/actors/:id", app.requireAuthenticatedUser(app.getActorHandler))
	router.HandlerFunc(http.MethodPatch, "/actors/:id", app.requireRoleAdmin(app.updateActorHandler))
	router.HandlerFunc(http.MethodDelete, "/actors/:id", app.requireRoleAdmin(app.deleteActorHandler))
	router.HandlerFunc(http.MethodGet, "/actors", app.requireAuthenticatedUser(app.getActorsHandler))

	router.HandlerFunc(http.MethodPost, "/movies", app.requireRoleAdmin(app.addMovieHandler))
	router.HandlerFunc(http.MethodGet, "/movies/:id", app.requireAuthenticatedUser(app.getMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/movies/:id", app.requireRoleAdmin(app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/movies/:id", app.requireRoleAdmin(app.deleteMovieHandler))
	router.HandlerFunc(http.MethodGet, "/movies", app.requireAuthenticatedUser(app.getMoviesHandler))

	router.HandlerFunc(http.MethodGet, "/search", app.requireAuthenticatedUser(app.searchMovieHandler))

	return app.recoverPanic(app.rateLimit(app.logRequest(app.authenticate(router))))
}
