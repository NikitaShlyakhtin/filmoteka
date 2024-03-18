package main

import (
	"context"
	"filmoteka/internal/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContextGetUser(t *testing.T) {
	app := &application{}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, &data.User{ID: 1, Name: "John Doe"})
	req = req.WithContext(ctx)

	user := app.contextGetUser(req)

	expectedUser := &data.User{ID: 1, Name: "John Doe"}
	if user.ID != expectedUser.ID || user.Name != expectedUser.Name {
		t.Errorf("unexpected user value, got: %v, want: %v", user, expectedUser)
	}
}

func TestContextSetUser(t *testing.T) {
	app := &application{}
	user := &data.User{ID: 1, Name: "John Doe"}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = app.contextSetUser(req, user)

	ctx := req.Context()
	value := ctx.Value(userContextKey)
	if value == nil {
		t.Error("user context value is nil")
	}

	// Assert the user value
	if u, ok := value.(*data.User); !ok || u.ID != user.ID || u.Name != user.Name {
		t.Errorf("unexpected user value, got: %v, want: %v", u, user)
	}
}
