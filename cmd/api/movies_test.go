package main

import (
	"context"
	"encoding/json"
	"filmoteka/internal/data"
	"filmoteka/internal/jsonlog"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
)

func TestAddMovieHandler(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			ReleaseDate time.Time `json:"release_date"` // RFC3339
			Rating      float32   `json:"rating"`
			Actors      []int64   `json:"actors"`
		}{
			Title:       "New Movie",
			Description: "New Movie Description",
			ReleaseDate: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			Rating:      8.0,
			Actors:      []int64{1, 2},
		}

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		res := httptest.NewRecorder()

		app.addMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusCreated {
			t.Errorf("expected status code %d, but got %d", http.StatusCreated, res.Code)
		}

		var respBody struct {
			Movie data.Movie `json:"movie"`
		}
		err = json.NewDecoder(res.Body).Decode(&respBody)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if respBody.Movie.Title != input.Title {
			t.Errorf("expected user name %q, but got %q", input.Title, respBody.Movie.Title)
		}
	})

	t.Run("InvalidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			ReleaseDate time.Time `json:"release_date"` // RFC3339
			Rating      float32   `json:"rating"`
			Actors      []int64   `json:"actors"`
		}{
			Title:       "New Movie",
			Description: "New Movie Description",
			Rating:      8.0,
			Actors:      []int64{1, 2},
		}

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		res := httptest.NewRecorder()

		app.addMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, but got %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("InvalidActors", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			ReleaseDate time.Time `json:"release_date"` // RFC3339
			Rating      float32   `json:"rating"`
			Actors      []int64   `json:"actors"`
		}{
			Title:       "New Movie",
			Description: "New Movie Description",
			ReleaseDate: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			Rating:      8.0,
			Actors:      []int64{1, -1},
		}

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		res := httptest.NewRecorder()

		app.addMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, but got %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("TitleExists", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			ReleaseDate time.Time `json:"release_date"` // RFC3339
			Rating      float32   `json:"rating"`
			Actors      []int64   `json:"actors"`
		}{
			Title:       "Mock Movie 1",
			Description: "New Movie Description",
			ReleaseDate: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			Rating:      8.0,
			Actors:      []int64{1, 2},
		}

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		res := httptest.NewRecorder()

		app.addMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, but got %d", http.StatusUnprocessableEntity, res.Code)
		}
	})
}

func TestDeleteMovieHandler(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		id := int64(1)

		req := httptest.NewRequest(http.MethodDelete, "/movies/", nil)
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		params := httprouter.Params{
			httprouter.Param{
				Key:   "id",
				Value: fmt.Sprint(id),
			},
		}
		req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

		res := httptest.NewRecorder()

		app.deleteMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusNoContent, res.Code)
		}
	})
}

func TestUpdateMovieHandler(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			ReleaseDate time.Time `json:"release_date"` // RFC3339
			Rating      float32   `json:"rating"`
			Actors      []int64   `json:"actors"`
		}{
			Title:       "New Movie",
			Description: "New Movie Description",
			ReleaseDate: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			Rating:      8.0,
			Actors:      []int64{1, 2},
		}

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		id := int64(1)

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPatch, "/movies/", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		params := httprouter.Params{
			httprouter.Param{
				Key:   "id",
				Value: fmt.Sprint(id),
			},
		}
		req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

		res := httptest.NewRecorder()

		app.updateMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, res.Code)
		}

		var respBody struct {
			Movie data.Movie `json:"movie"`
		}
		err = json.NewDecoder(res.Body).Decode(&respBody)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if respBody.Movie.Title != input.Title {
			t.Errorf("expected user name %q, but got %q", input.Title, respBody.Movie.Title)
		}
	})

	t.Run("DuplicateTitle", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			ReleaseDate time.Time `json:"release_date"` // RFC3339
			Rating      float32   `json:"rating"`
			Actors      []int64   `json:"actors"`
		}{
			Title: "Mock Movie 1",
		}

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		id := int64(1)

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPatch, "/movies/", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		params := httprouter.Params{
			httprouter.Param{
				Key:   "id",
				Value: fmt.Sprint(id),
			},
		}
		req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

		res := httptest.NewRecorder()

		app.updateMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, but got %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("InvalidActors", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			ReleaseDate time.Time `json:"release_date"` // RFC3339
			Rating      float32   `json:"rating"`
			Actors      []int64   `json:"actors"`
		}{
			Actors: []int64{1, -1},
		}

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		id := int64(1)

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPatch, "/movies/", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		params := httprouter.Params{
			httprouter.Param{
				Key:   "id",
				Value: fmt.Sprint(id),
			},
		}
		req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

		res := httptest.NewRecorder()

		app.updateMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, but got %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("MovieNotFound", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			ReleaseDate time.Time `json:"release_date"` // RFC3339
			Rating      float32   `json:"rating"`
			Actors      []int64   `json:"actors"`
		}{
			Title: "Mock Movie -1 Update",
		}

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		id := int64(-1)

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPatch, "/movies/", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		params := httprouter.Params{
			httprouter.Param{
				Key:   "id",
				Value: fmt.Sprint(id),
			},
		}
		req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

		res := httptest.NewRecorder()

		app.updateMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, but got %d", http.StatusNotFound, res.Code)
		}
	})
}

func TestGetMovieHandler(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		id := int64(1)

		req := httptest.NewRequest(http.MethodGet, "/movies/", nil)

		params := httprouter.Params{
			httprouter.Param{
				Key:   "id",
				Value: fmt.Sprint(id),
			},
		}
		req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

		res := httptest.NewRecorder()

		app.getMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, res.Code)
		}
	})

	t.Run("MovieNotFound", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		id := int64(-1)

		req := httptest.NewRequest(http.MethodGet, "/movies/", nil)

		params := httprouter.Params{
			httprouter.Param{
				Key:   "id",
				Value: fmt.Sprint(id),
			},
		}
		req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

		res := httptest.NewRecorder()

		app.getMovieHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, but got %d", http.StatusNotFound, res.Code)
		}
	})
}

func TestGetMoviesHandler(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		req := httptest.NewRequest(http.MethodGet, "/movies", nil)

		res := httptest.NewRecorder()

		app.getMoviesHandler(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, res.Code)
		}
	})
}

func TestSearchMoviesHandler(t *testing.T) {
	t.Run("ValidTitle", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		req := httptest.NewRequest(http.MethodGet, "/movies/search", nil)

		title := "1"

		params := httprouter.Params{
			httprouter.Param{
				Key:   "title",
				Value: title,
			},
		}
		req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

		res := httptest.NewRecorder()

		app.searchMovieHandler(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, res.Code)
		}
	})

	t.Run("ValidActor", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		req := httptest.NewRequest(http.MethodGet, "/movies/search", nil)

		actor := "1"

		params := httprouter.Params{
			httprouter.Param{
				Key:   "actor",
				Value: actor,
			},
		}
		req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

		res := httptest.NewRecorder()

		app.searchMovieHandler(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, res.Code)
		}
	})
}
