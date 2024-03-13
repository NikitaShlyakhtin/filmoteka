package main

import "net/http"

// Добавляет информацию об актёре в базу данных.
func (app *application) addActorHandler(w http.ResponseWriter, r *http.Request) {
}

// Обновляет информацию об актёре в базе данных полностью через метод PUT.
func (app *application) updateActorHandler(w http.ResponseWriter, r *http.Request) {
}

// Обновляет информацию об актёре в базе данных частично через метод PATCH.
func (app *application) updateActorPartialHandler(w http.ResponseWriter, r *http.Request) {
}

// Возвращает список актёров с информацией о фильмах, в которых они снимались.
func (app *application) getActorsHandler(w http.ResponseWriter, r *http.Request) {
}

// Удаляет информацию об актёре из базы данных.
func (app *application) deleteActorHandler(w http.ResponseWriter, r *http.Request) {
}
