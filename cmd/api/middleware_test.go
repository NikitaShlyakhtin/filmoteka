package main

import (
	"filmoteka/internal/jsonlog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLogRequest(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.logRequest(handler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRecoverPanic(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.recoverPanic(handler).ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestRateLimit(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	app.config.limiter.enabled = true
	app.config.limiter.rps = 2
	app.config.limiter.burst = 4

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.rateLimit(handler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}
