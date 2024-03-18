package main

import (
	"context"
	"filmoteka/internal/validator"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestReadIDParam(t *testing.T) {
	app := new(application)

	req := httptest.NewRequest(http.MethodGet, "/test/", nil)
	params := httprouter.Params{
		httprouter.Param{
			Key:   "id",
			Value: "123",
		},
	}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))

	id, err := app.readIDParam(req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if id != 123 {
		t.Errorf("expected id to be 123, got %d", id)
	}
}

func TestWriteJSON(t *testing.T) {
	app := new(application)

	rec := httptest.NewRecorder()
	data := envelope{"status": "available"}
	headers := http.Header{"X-Test": []string{"true"}}

	err := app.writeJSON(rec, http.StatusOK, data, headers)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status to be %d, got %d", http.StatusOK, rec.Code)
	}

	want := `{"status":"available"}`
	got := strings.Join(strings.Fields(rec.Body.String()), "")
	if got != want {
		t.Errorf("expected body to equal %q, got %q", want, got)
	}

	if rec.Header().Get("X-Test") != "true" {
		t.Errorf("expected header X-Test to equal %q, got %q", "true", rec.Header().Get("X-Test"))
	}
}

func TestReadJSON(t *testing.T) {
	app := new(application)

	t.Run("Valid JSON", func(t *testing.T) {
		body := `{"name": "John", "age": 30}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := app.readJSON(httptest.NewRecorder(), req, &data)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if data.Name != "John" {
			t.Errorf("expected name to be %q, got %q", "John", data.Name)
		}

		if data.Age != 30 {
			t.Errorf("expected age to be %d, got %d", 30, data.Age)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := `{"name": "John", "age": "thirty"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := app.readJSON(httptest.NewRecorder(), req, &data)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestReadInt(t *testing.T) {
	app := new(application)

	t.Run("Empty Query String", func(t *testing.T) {
		qs := url.Values{}
		key := "count"
		defaultValue := 10

		v := validator.New()

		result := app.readInt(qs, key, defaultValue, v)
		if result != defaultValue {
			t.Errorf("expected result to be %d, got %d", defaultValue, result)
		}

		if !v.Valid() {
			t.Error("expected no validation errors")
		}
	})

	t.Run("Non-integer Value", func(t *testing.T) {
		qs := url.Values{"count": []string{"abc"}}
		key := "count"
		defaultValue := 10

		v := validator.New()

		result := app.readInt(qs, key, defaultValue, v)
		if result != defaultValue {
			t.Errorf("expected result to be %d, got %d", defaultValue, result)
		}

		if v.Valid() {
			t.Error("expected validation errors")
		}

		errors := v.Errors
		if len(errors) != 1 {
			t.Errorf("expected 1 validation error, got %d", len(errors))
		}

		expectedError := "must be an integer value"
		if errors[key] != expectedError {
			t.Errorf("expected validation error: %s, got: %s", expectedError, errors[key])
		}
	})

	t.Run("Valid Integer Value", func(t *testing.T) {
		qs := url.Values{"count": []string{"20"}}
		key := "count"
		defaultValue := 10

		v := validator.New()

		result := app.readInt(qs, key, defaultValue, v)
		if result != 20 {
			t.Errorf("expected result to be %d, got %d", 20, result)
		}

		if !v.Valid() {
			t.Error("expected no validation errors")
		}
	})
}

func TestReadString(t *testing.T) {
	app := new(application)

	t.Run("Empty Query String", func(t *testing.T) {
		qs := url.Values{}
		key := "name"
		defaultValue := "John"

		result := app.readString(qs, key, defaultValue)
		if result != defaultValue {
			t.Errorf("expected result to be %q, got %q", defaultValue, result)
		}
	})

	t.Run("Non-empty Query String", func(t *testing.T) {
		qs := url.Values{"name": []string{"Alice"}}
		key := "name"
		defaultValue := "John"

		result := app.readString(qs, key, defaultValue)
		if result != "Alice" {
			t.Errorf("expected result to be %q, got %q", "Alice", result)
		}
	})
}
