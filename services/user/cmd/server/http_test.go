package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequireAuthRejectsMissingBearerToken(t *testing.T) {
	a := &app{cfg: config{JWTSecret: "test-secret", JWTIssuer: "test-issuer", JWTTTL: time.Hour}}
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()

	a.requireAuth(func(w http.ResponseWriter, r *http.Request, c claims) {
		t.Fatal("handler should not be called when token is missing")
	})(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRequireAuthAcceptsValidBearerToken(t *testing.T) {
	a := &app{cfg: config{JWTSecret: "test-secret", JWTIssuer: "test-issuer", JWTTTL: time.Hour}}
	token, err := a.signToken(userResponse{ID: "user-1", Email: "ana@example.com", Role: rolePatient})
	if err != nil {
		t.Fatalf("signToken returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	called := false
	a.requireAuth(func(w http.ResponseWriter, r *http.Request, c claims) {
		called = true
		if c.Subject != "user-1" {
			t.Fatalf("unexpected subject %q", c.Subject)
		}
		w.WriteHeader(http.StatusNoContent)
	})(rec, req)

	if !called {
		t.Fatal("expected handler to be called")
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusNoContent)
	}
}
