package main

import (
	"errors"
	"filmoteka/internal/jsonlog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestServerErrorResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	err = errors.New("test error")
	app.serverErrorResponse(res, req, err)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, res.Code)
	}

	want := `{"error":"the server encountered a problem and could not process your request"}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}
}

func TestErrorResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	app.errorResponse(res, req, http.StatusOK, "success")

	if res.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, res.Code)
	}

	want := `{"error":"success"}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}
}

func TestNotFoundResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	app.notFoundResponse(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, res.Code)
	}

	want := `{"error":"the requested resource could not be found"}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}
}

func TestMethodNotAllowedResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	app.methodNotAllowedResponse(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status code %d, got %d", http.StatusMethodNotAllowed, res.Code)
	}

	want := `{"error":"the GET method is not supported for this resource"}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}
}

func TestBadRequestResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	err = errors.New("bad request")
	app.badRequestResponse(res, req, err)

	if res.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, res.Code)
	}

	want := `{"error":"bad request"}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}
}

func TestInvalidAuthenticationCredentialsResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	app.invalidAuthenticationCredentialsResponse(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, res.Code)
	}

	want := `{"error":"invalid or missing authentication credentials"}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}

	// Check if the "WWW-Authenticate" header is set to "Basic"
	if res.Header().Get("WWW-Authenticate") != "Basic" {
		t.Errorf("expected WWW-Authenticate header to be set to 'Basic', got %q", res.Header().Get("WWW-Authenticate"))
	}
}

func TestFailedValidationResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	errors := map[string]string{
		"field1": "error1",
		"field2": "error2",
	}

	app.failedValidationResponse(res, req, errors)

	if res.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status code %d, got %d", http.StatusUnprocessableEntity, res.Code)
	}

	want := `{"error":{"field1":"error1","field2":"error2"}}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}
}

func TestRateLimitExceededResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	app.rateLimitExceededResponse(res, req)

	if res.Code != http.StatusTooManyRequests {
		t.Errorf("expected status code %d, got %d", http.StatusTooManyRequests, res.Code)
	}

	want := `{"error":"rate limit exceeded"}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}
}

func TestInvalidCredentialsResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	app.invalidCredentialsResponse(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, res.Code)
	}

	want := `{"error":"invalid authentication credentials"}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}
}

func TestAuthenticationRequiredResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	app.authenticationRequiredResponse(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, res.Code)
	}

	want := `{"error":"you must be authenticated to access this resource"}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}
}

func TestNotPermittedResponse(t *testing.T) {
	app := &application{
		logger: jsonlog.New(os.Stdout, jsonlog.LevelInfo),
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	app.notPermittedResponse(res, req)

	if res.Code != http.StatusForbidden {
		t.Errorf("expected status code %d, got %d", http.StatusForbidden, res.Code)
	}

	want := `{"error":"your user account doesn't have the necessary permissions to access this resource"}`
	want = strings.Join(strings.Fields(want), "")
	got := strings.Join(strings.Fields(res.Body.String()), "")
	if want != got {
		t.Errorf("expected response body %q, got %q", want, got)
	}
}
