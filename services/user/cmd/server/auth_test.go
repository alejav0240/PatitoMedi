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
		t.Fatalf("Error signing token: %v", err)
	}

	claims, err := a.parseToken(token)
	if err != nil {
		t.Fatalf("Error parsing token: %v", err)
	}
	if claims.Subject != user.ID || claims.Email != user.Email || claims.Role != user.Role {
		t.Fatalf("Error Inesperado : %+v", claims)
	}
}

func TestParseTokenRejectsWrongSecret(t *testing.T) {
	signer := &app{cfg: config{JWTSecret: "one-secret", JWTIssuer: "test-issuer", JWTTTL: time.Hour}}
	verifier := &app{cfg: config{JWTSecret: "another-secret", JWTIssuer: "test-issuer", JWTTTL: time.Hour}}

	token, err := signer.signToken(userResponse{ID: "user-1", Email: "ana@example.com", Role: rolePatient})
	if err != nil {
		t.Fatalf("Error signing token: %v", err)
	}

	if _, err := verifier.parseToken(token); err == nil {
		t.Fatal("Se esperaba que el token fuera rechazado con una clave secreta incorrecta.")
	}
}

func TestParseTokenRejectsExpiredToken(t *testing.T) {
	a := &app{cfg: config{JWTSecret: "test-secret", JWTIssuer: "test-issuer", JWTTTL: -time.Hour}}

	token, err := a.signToken(userResponse{ID: "user-1", Email: "ana@example.com", Role: rolePatient})
	if err != nil {
		t.Fatalf("Error signing token: %v", err)
	}

	if _, err := a.parseToken(token); err == nil {
		t.Fatal("Se esperaba que el token expirado fuera rechazado")
	}
}
