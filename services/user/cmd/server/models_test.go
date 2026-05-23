package main

import "testing"

func TestNormalizeRegisterRequest(t *testing.T) {
	req := normalizeRegisterRequest(registerRequest{
		Email:     "  ANA@EXAMPLE.COM ",
		FullName:  " Ana Perez ",
		Specialty: " Cardiologia ",
	})

	if req.Email != "ana@example.com" {
		t.Fatalf("unexpected email: %q", req.Email)
	}
	if req.FullName != "Ana Perez" {
		t.Fatalf("unexpected full name: %q", req.FullName)
	}
	if req.Specialty != "Cardiologia" {
		t.Fatalf("unexpected specialty: %q", req.Specialty)
	}
}

func TestValidateRegisterRequest(t *testing.T) {
	tests := []struct {
		name string
		req  registerRequest
		role string
		want string
	}{
		{
			name: "patient ok",
			req:  registerRequest{Email: "ana@example.com", Password: "password123", FullName: "Ana"},
			role: rolePatient,
		},
		{
			name: "doctor requires specialty",
			req:  registerRequest{Email: "dr@example.com", Password: "password123", FullName: "Dr"},
			role: roleDoctor,
			want: "specialty is required for doctors",
		},
		{
			name: "short password",
			req:  registerRequest{Email: "ana@example.com", Password: "short", FullName: "Ana"},
			role: rolePatient,
			want: "password must be at least 8 characters",
		},
		{
			name: "missing required",
			req:  registerRequest{Email: "", Password: "password123", FullName: "Ana"},
			role: rolePatient,
			want: "email, password and full_name are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateRegisterRequest(tt.req, tt.role); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
