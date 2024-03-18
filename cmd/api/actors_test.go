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

func TestAddActorHandler(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			FullName  string    `json:"full_name"`
			Gender    string    `json:"gender"`
			BirthDate time.Time `json:"birth_date"` // RFC3339
		}{
			FullName:  "John Doe",
			Gender:    "male",
			BirthDate: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
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

		req := httptest.NewRequest(http.MethodPost, "/actors", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		res := httptest.NewRecorder()

		app.addActorHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusCreated {
			t.Errorf("expected status code %d, but got %d", http.StatusCreated, res.Code)
		}

		var respBody struct {
			Actor data.Actor `json:"actor"`
		}
		err = json.NewDecoder(res.Body).Decode(&respBody)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if respBody.Actor.FullName != input.FullName {
			t.Errorf("expected user name %q, but got %q", input.FullName, respBody.Actor.FullName)
		}
	})

	t.Run("InvalidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			FullName  string    `json:"full_name"`
			Gender    string    `json:"gender"`
			BirthDate time.Time `json:"birth_date"` // RFC3339
		}{
			FullName:  "",
			Gender:    "",
			BirthDate: time.Time{},
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

		req := httptest.NewRequest(http.MethodPost, "/actors", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		res := httptest.NewRecorder()

		app.addActorHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, but got %d", http.StatusCreated, res.Code)
		}
	})

	t.Run("DuplicateName", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			FullName  string    `json:"full_name"`
			Gender    string    `json:"gender"`
			BirthDate time.Time `json:"birth_date"` // RFC3339
		}{
			FullName:  "Mock Actor 1",
			Gender:    "password123",
			BirthDate: time.Time{},
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

		req := httptest.NewRequest(http.MethodPost, "/actors", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(auth.Name, auth.Password)

		res := httptest.NewRecorder()

		app.addActorHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, but got %d", http.StatusUnprocessableEntity, res.Code)
		}
	})
}

func TestGetActorHandler(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		id := int64(1)

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		req := httptest.NewRequest(http.MethodGet, "/actors/", nil)
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

		app.getActorHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, res.Code)
		}

		var respBody struct {
			Actor data.Actor `json:"actor"`
		}
		err := json.NewDecoder(res.Body).Decode(&respBody)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if respBody.Actor.FullName != "Mock Actor 1" {
			t.Errorf("expected user name %q, but got %q", "Mock Actor 1", respBody.Actor.FullName)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		id := int64(-1)

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		req := httptest.NewRequest(http.MethodGet, "/actors/", nil)
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

		app.getActorHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, but got %d", http.StatusNotFound, res.Code)
		}
	})
}

func TestGetActorsHandler(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/actors", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(auth.Name, auth.Password)

	res := httptest.NewRecorder()

	app.getActorsHandler(res, req)

	// Check the response status code
	if res.Code != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, res.Code)
	}

	var respBody struct {
		Actors []data.Actor `json:"actors"`
	}
	err := json.NewDecoder(res.Body).Decode(&respBody)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(respBody.Actors) != 2 {
		t.Errorf("expected actors count %d, but got %d", 2, len(respBody.Actors))
	}
}

func TestDeleteActorHandler(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		id := int64(1)

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		req := httptest.NewRequest(http.MethodDelete, "/actors/", nil)
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

		app.deleteActorHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, res.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		id := int64(-1)

		auth := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "admin",
			Password: "password123",
		}

		req := httptest.NewRequest(http.MethodDelete, "/actors/", nil)
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

		app.deleteActorHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, but got %d", http.StatusNotFound, res.Code)
		}
	})
}

func TestUpdateActorHandler(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		id := int64(1)

		input := struct {
			FullName  string    `json:"full_name"`
			Gender    string    `json:"gender"`
			BirthDate time.Time `json:"birth_date"` // RFC3339
		}{
			FullName: "John Doe",
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

		req := httptest.NewRequest(http.MethodPatch, "/actors/", strings.NewReader(string(reqBody)))
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

		app.updateActorHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, res.Code)
		}

		var respBody struct {
			Actor data.Actor `json:"actor"`
		}
		err = json.NewDecoder(res.Body).Decode(&respBody)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if respBody.Actor.FullName != input.FullName {
			t.Errorf("expected user name %q, but got %q", input.FullName, respBody.Actor.FullName)
		}
	})
}
