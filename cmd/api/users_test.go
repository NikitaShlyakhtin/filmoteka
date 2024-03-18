package main

import (
	"encoding/json"
	"filmoteka/internal/data"
	"filmoteka/internal/jsonlog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCreateUserHandler(t *testing.T) {
	t.Run("ValidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "John Doe",
			Password: "password123",
		}

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()

		app.createUserHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusCreated {
			t.Errorf("expected status code %d, but got %d", http.StatusCreated, res.Code)
		}

		var respBody struct {
			User data.User `json:"user"`
		}
		err = json.NewDecoder(res.Body).Decode(&respBody)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if respBody.User.Name != input.Name {
			t.Errorf("expected user name %q, but got %q", input.Name, respBody.User.Name)
		}
	})

	t.Run("InvalidInput", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "",
			Password: "",
		}

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()

		app.createUserHandler(res, req)

		// Check the response status code
		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, but got %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("DuplicateName", func(t *testing.T) {
		app := &application{
			models: data.NewMockModels(),
			logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
		}

		input := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}{
			Name:     "user",
			Password: "password123",
		}

		reqBody, err := json.Marshal(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()

		app.createUserHandler(res, req)

		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, but got %d", http.StatusUnprocessableEntity, res.Code)
		}
	})
}
