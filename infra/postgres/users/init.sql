CREATE TABLE IF NOT EXISTS roles (
  name TEXT PRIMARY KEY,
  description TEXT NOT NULL
);

INSERT INTO roles (name, description)
VALUES
  ('patient', 'Paciente de la plataforma'),
  ('doctor', 'Medico de la plataforma')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS patients (
  id UUID PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  full_name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS doctors (
  id UUID PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  full_name TEXT NOT NULL,
  specialty TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  user_role TEXT NOT NULL REFERENCES roles(name),
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_patients_email ON patients (email);
CREATE INDEX IF NOT EXISTS idx_doctors_email ON doctors (email);
CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions (user_id, user_role);
