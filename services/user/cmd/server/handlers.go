package main

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (a *app) health(w http.ResponseWriter, r *http.Request) {
	if err := a.store.Ping(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "database": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "user"})
}

func (a *app) registerPatient(w http.ResponseWriter, r *http.Request) {
	a.register(w, r, rolePatient)
}

func (a *app) registerDoctor(w http.ResponseWriter, r *http.Request) {
	a.register(w, r, roleDoctor)
}

func (a *app) register(w http.ResponseWriter, r *http.Request, role string) {
	var req registerRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	req = normalizeRegisterRequest(req)
	if message := validateRegisterRequest(req, role); message != "" {
		writeError(w, http.StatusBadRequest, message)
		return
	}
	if _, _, err := a.store.FindUserByEmail(r.Context(), req.Email); err == nil {
		writeError(w, http.StatusConflict, "email already exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user, err := a.store.CreateUser(r.Context(), role, req, string(hash))
	if err != nil {
		if isUniqueViolation(err) {
			writeError(w, http.StatusConflict, "email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	a.metrics.registrations.Add(1)
	a.producer.Publish(r.Context(), a.cfg.KafkaUserTopic, user.ID, map[string]any{"event": "user-registered", "user": user})
	writeJSON(w, http.StatusCreated, user)
}

func (a *app) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	email := normalizeEmail(req.Email)
	if email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, passwordHash, err := a.store.FindUserByEmail(r.Context(), email)
	if err != nil {
		a.metrics.loginsFailed.Add(1)
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
		a.metrics.loginsFailed.Add(1)
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := a.signToken(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not sign token")
		return
	}

	sessionID := newID()
	expiresAt := time.Now().Add(a.cfg.JWTTTL)
	if err := a.store.CreateSession(r.Context(), sessionID, user.ID, user.Role, expiresAt); err != nil {
		writeError(w, http.StatusInternalServerError, "could not create session")
		return
	}

	a.producer.Publish(r.Context(), "session-created", sessionID, map[string]any{"event": "session-created", "session_id": sessionID, "user_id": user.ID, "role": user.Role})
	a.metrics.loginsOK.Add(1)
	writeJSON(w, http.StatusOK, loginResponse{Token: token, User: user})
}

func (a *app) logout(w http.ResponseWriter, r *http.Request, c claims) {
	a.producer.Publish(r.Context(), "session-ended", c.Subject, map[string]any{"event": "session-ended", "user_id": c.Subject, "role": c.Role})
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *app) me(w http.ResponseWriter, r *http.Request, c claims) {
	user, err := a.store.FindUserByID(r.Context(), c.Subject, c.Role)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (a *app) updateMe(w http.ResponseWriter, r *http.Request, c claims) {
	var req struct {
		FullName  string `json:"full_name"`
		Specialty string `json:"specialty,omitempty"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}

	req.FullName = trim(req.FullName)
	req.Specialty = trim(req.Specialty)
	if req.FullName == "" {
		writeError(w, http.StatusBadRequest, "full_name is required")
		return
	}
	if c.Role == roleDoctor && req.Specialty == "" {
		writeError(w, http.StatusBadRequest, "specialty is required for doctors")
		return
	}

	user, err := a.store.UpdateUser(r.Context(), c.Subject, c.Role, req.FullName, req.Specialty)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	a.producer.Publish(r.Context(), "user-updated", user.ID, map[string]any{"event": "user-updated", "user": user})
	writeJSON(w, http.StatusOK, user)
}

func (a *app) listDoctors(w http.ResponseWriter, r *http.Request) {
	doctors, err := a.store.ListDoctors(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list doctors")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"doctors": doctors})
}

func (a *app) getDoctor(w http.ResponseWriter, r *http.Request) {
	user, err := a.store.FindUserByID(r.Context(), r.PathValue("id"), roleDoctor)
	if err != nil {
		writeError(w, http.StatusNotFound, "doctor not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}
