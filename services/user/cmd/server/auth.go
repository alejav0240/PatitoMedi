package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

func (a *app) signToken(user userResponse) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	payload := claims{
		Subject: user.ID,
		Email:   user.Email,
		Role:    user.Role,
		Issuer:  a.cfg.JWTIssuer,
		Expires: time.Now().Add(a.cfg.JWTTTL).Unix(),
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	unsigned := base64.RawURLEncoding.EncodeToString(headerJSON) + "." + base64.RawURLEncoding.EncodeToString(payloadJSON)
	sig := hmacSHA256(unsigned, a.cfg.JWTSecret)
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

func (a *app) parseToken(token string) (claims, error) {
	var c claims
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return c, errors.New("invalid token")
	}

	unsigned := parts[0] + "." + parts[1]
	expected := base64.RawURLEncoding.EncodeToString(hmacSHA256(unsigned, a.cfg.JWTSecret))
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return c, errors.New("invalid signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(payload, &c); err != nil {
		return c, err
	}
	if c.Issuer != a.cfg.JWTIssuer || c.Expires < time.Now().Unix() {
		return c, errors.New("expired or invalid issuer")
	}
	return c, nil
}

func (a *app) requireAuth(next func(http.ResponseWriter, *http.Request, claims)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth || token == "" {
			writeError(w, http.StatusUnauthorized, "bearer token is required")
			return
		}
		c, err := a.parseToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		next(w, r, c)
	}
}

func hmacSHA256(message string, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return mac.Sum(nil)
}
