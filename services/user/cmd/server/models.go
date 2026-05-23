package main

import "time"

const (
	rolePatient = "patient"
	roleDoctor  = "doctor"
)

type registerRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FullName  string `json:"full_name"`
	Specialty string `json:"specialty,omitempty"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	Specialty string    `json:"specialty,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

type claims struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Issuer  string `json:"iss"`
	Expires int64  `json:"exp"`
}

func normalizeRegisterRequest(req registerRequest) registerRequest {
	req.Email = normalizeEmail(req.Email)
	req.FullName = trim(req.FullName)
	req.Specialty = trim(req.Specialty)
	return req
}

func validateRegisterRequest(req registerRequest, role string) string {
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return "email, password and full_name are required"
	}
	if len(req.Password) < 8 {
		return "password must be at least 8 characters"
	}
	if role == roleDoctor && req.Specialty == "" {
		return "specialty is required for doctors"
	}
	return ""
}
