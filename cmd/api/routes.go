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
	router.HandlerFunc(http.MethodPut, "/actors/{actor_id}", app.updateActorHandler)
	router.HandlerFunc(http.MethodPatch, "/actors/{actor_id}", app.updateActorPartialHandler)
	router.HandlerFunc(http.MethodDelete, "/actors/{actor_id}", app.deleteActorHandler)
	router.HandlerFunc(http.MethodGet, "/actors", app.getActorsHandler)

	router.HandlerFunc(http.MethodPost, "/movies", app.addMovieHandler)
	router.HandlerFunc(http.MethodPut, "/movies/{movie_id}", app.updateMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/movies/{movie_id}", app.updateMoviePartialHandler)
	router.HandlerFunc(http.MethodDelete, "/movies/{movie_id}", app.deleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/movies", app.getMoviesHandler)

	router.HandlerFunc(http.MethodGet, "/search", app.searchHandler)

	return router
}
