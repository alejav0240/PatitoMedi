package main

import (
	"testing"
	"time"
)

func TestSignAndParseToken(t *testing.T) {
	a := &app{cfg: config{JWTSecret: "test-secret", JWTIssuer: "test-issuer", JWTTTL: time.Hour}}
	user := userResponse{ID: "user-1", Email: "ana@example.com", Role: rolePatient}

	token, err := a.signToken(user)
	if err != nil {
		t.Fatalf("signToken returned error: %v", err)
	}

	claims, err := a.parseToken(token)
	if err != nil {
		t.Fatalf("parseToken returned error: %v", err)
	}
	if claims.Subject != user.ID || claims.Email != user.Email || claims.Role != user.Role {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseTokenRejectsWrongSecret(t *testing.T) {
	signer := &app{cfg: config{JWTSecret: "one-secret", JWTIssuer: "test-issuer", JWTTTL: time.Hour}}
	verifier := &app{cfg: config{JWTSecret: "another-secret", JWTIssuer: "test-issuer", JWTTTL: time.Hour}}

	token, err := signer.signToken(userResponse{ID: "user-1", Email: "ana@example.com", Role: rolePatient})
	if err != nil {
		t.Fatalf("signToken returned error: %v", err)
	}

	if _, err := verifier.parseToken(token); err == nil {
		t.Fatal("expected token to be rejected with wrong secret")
	}
}

func TestParseTokenRejectsExpiredToken(t *testing.T) {
	a := &app{cfg: config{JWTSecret: "test-secret", JWTIssuer: "test-issuer", JWTTTL: -time.Hour}}

	token, err := a.signToken(userResponse{ID: "user-1", Email: "ana@example.com", Role: rolePatient})
	if err != nil {
		t.Fatalf("signToken returned error: %v", err)
	}

	if _, err := a.parseToken(token); err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}
