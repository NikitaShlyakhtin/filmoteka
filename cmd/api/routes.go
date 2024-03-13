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

	router.HandlerFunc(http.MethodPost, "/actors", app.addActorHandler)
	router.HandlerFunc(http.MethodPut, "/actors/:id", app.updateActorHandler)
	router.HandlerFunc(http.MethodPatch, "/actors/:id", app.updateActorPartialHandler)
	router.HandlerFunc(http.MethodDelete, "/actors/:id", app.deleteActorHandler)
	router.HandlerFunc(http.MethodGet, "/actors", app.getActorsHandler)

	router.HandlerFunc(http.MethodPost, "/movies", app.addMovieHandler)
	router.HandlerFunc(http.MethodPut, "/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/movies/:id", app.updateMoviePartialHandler)
	router.HandlerFunc(http.MethodDelete, "/movies/:id", app.deleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/movies", app.getMoviesHandler)

	router.HandlerFunc(http.MethodGet, "/search", app.searchHandler)

	return app.recoverPanic(app.rateLimit(app.logRequest(router)))
}
