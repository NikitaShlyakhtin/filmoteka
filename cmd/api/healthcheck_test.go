package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	app := &application{
		config: config{
			env: "test",
		},
	}

	req, err := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()

	app.healthcheckHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rec.Code)
	}

	want := `{"status":"available","system_info":{"environment":"test"}}`
	got := strings.Join(strings.Fields(rec.Body.String()), "")

	if got != want {
		t.Errorf("want body to equal %q, got %q", want, got)
	}
}
