package main

import (
	"net/http"
)

// Добавляет информацию о фильме в базу данных.
func (app *application) addMovieHandler(w http.ResponseWriter, r *http.Request) {
}

// Обновляет информацию о фильме в базе данных полностью через метод PUT.
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
}

// Обновляет информацию о фильме в базе данных частично через метод PATCH.
func (app *application) updateMoviePartialHandler(w http.ResponseWriter, r *http.Request) {
}

// getMoviesHandler возвращает список фильмов с возможностью сортировки и фильтрации.
func (app *application) getMoviesHandler(w http.ResponseWriter, r *http.Request) {
}

// Удаляет информацию о фильме из базы данных.
func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
}
