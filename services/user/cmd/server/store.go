package main

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type userStore struct {
	db *sql.DB
}

func (s userStore) Ping(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return s.db.PingContext(pingCtx)
}

func (s userStore) CreateUser(ctx context.Context, role string, req registerRequest, passwordHash string) (userResponse, error) {
	id := newID()
	user := userResponse{ID: id, Email: req.Email, FullName: req.FullName, Role: role, Specialty: req.Specialty}

	if role == rolePatient {
		err := s.db.QueryRowContext(
			ctx,
			`INSERT INTO patients (id, email, password_hash, full_name) VALUES ($1, $2, $3, $4) RETURNING created_at`,
			id, req.Email, passwordHash, req.FullName,
		).Scan(&user.CreatedAt)
		return user, err
	}

	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO doctors (id, email, password_hash, full_name, specialty) VALUES ($1, $2, $3, $4, $5) RETURNING created_at`,
		id, req.Email, passwordHash, req.FullName, req.Specialty,
	).Scan(&user.CreatedAt)
	return user, err
}

func (s userStore) FindUserByEmail(ctx context.Context, email string) (userResponse, string, error) {
	user, passwordHash, err := s.findPatientByEmail(ctx, email)
	if err == nil {
		return user, passwordHash, nil
	}
	return s.findDoctorByEmail(ctx, email)
}

func (s userStore) findPatientByEmail(ctx context.Context, email string) (userResponse, string, error) {
	var user userResponse
	var passwordHash string
	user.Role = rolePatient
	err := s.db.QueryRowContext(ctx, `SELECT id::text, email, password_hash, full_name, created_at FROM patients WHERE email=$1`, email).
		Scan(&user.ID, &user.Email, &passwordHash, &user.FullName, &user.CreatedAt)
	return user, passwordHash, err
}

func (s userStore) findDoctorByEmail(ctx context.Context, email string) (userResponse, string, error) {
	var user userResponse
	var passwordHash string
	user.Role = roleDoctor
	err := s.db.QueryRowContext(ctx, `SELECT id::text, email, password_hash, full_name, specialty, created_at FROM doctors WHERE email=$1`, email).
		Scan(&user.ID, &user.Email, &passwordHash, &user.FullName, &user.Specialty, &user.CreatedAt)
	return user, passwordHash, err
}

func (s userStore) FindUserByID(ctx context.Context, id string, role string) (userResponse, error) {
	var user userResponse
	user.Role = role
	if role == roleDoctor {
		err := s.db.QueryRowContext(ctx, `SELECT id::text, email, full_name, specialty, created_at FROM doctors WHERE id=$1`, id).
			Scan(&user.ID, &user.Email, &user.FullName, &user.Specialty, &user.CreatedAt)
		return user, err
	}
	err := s.db.QueryRowContext(ctx, `SELECT id::text, email, full_name, created_at FROM patients WHERE id=$1`, id).
		Scan(&user.ID, &user.Email, &user.FullName, &user.CreatedAt)
	return user, err
}

func (s userStore) UpdateUser(ctx context.Context, id string, role string, fullName string, specialty string) (userResponse, error) {
	var user userResponse
	user.Role = role
	if role == roleDoctor {
		err := s.db.QueryRowContext(ctx, `UPDATE doctors SET full_name=$1, specialty=$2, updated_at=now() WHERE id=$3 RETURNING id::text, email, full_name, specialty, created_at`, fullName, specialty, id).
			Scan(&user.ID, &user.Email, &user.FullName, &user.Specialty, &user.CreatedAt)
		return user, err
	}
	err := s.db.QueryRowContext(ctx, `UPDATE patients SET full_name=$1, updated_at=now() WHERE id=$2 RETURNING id::text, email, full_name, created_at`, fullName, id).
		Scan(&user.ID, &user.Email, &user.FullName, &user.CreatedAt)
	return user, err
}

func (s userStore) ListDoctors(ctx context.Context) ([]userResponse, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id::text, email, full_name, specialty, created_at FROM doctors ORDER BY full_name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	doctors := []userResponse{}
	for rows.Next() {
		var doctor userResponse
		doctor.Role = roleDoctor
		if err := rows.Scan(&doctor.ID, &doctor.Email, &doctor.FullName, &doctor.Specialty, &doctor.CreatedAt); err != nil {
			return nil, err
		}
		doctors = append(doctors, doctor)
	}
	return doctors, rows.Err()
}

func (s userStore) CreateSession(ctx context.Context, id string, userID string, userRole string, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO sessions (id, user_id, user_role, expires_at) VALUES ($1, $2, $3, $4)`, id, userID, userRole, expiresAt)
	return err
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}
