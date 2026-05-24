package main

import (
	"encoding/json"
	"net/http"
)

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.health)
	mux.HandleFunc("GET /metrics", a.metricsHandler)
	mux.HandleFunc("POST /register/patient", a.registerPatient)
	mux.HandleFunc("POST /register/doctor", a.registerDoctor)
	mux.HandleFunc("POST /login", a.login)
	mux.HandleFunc("POST /logout", a.requireAuth(a.logout))
	mux.HandleFunc("GET /me", a.requireAuth(a.me))
	mux.HandleFunc("PATCH /me", a.requireAuth(a.updateMe))
	mux.HandleFunc("GET /doctors", a.listDoctors)
	mux.HandleFunc("GET /doctors/{id}", a.getDoctor)
	return a.instrument(mux)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "Formato JSON inválido: "+err.Error())
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
