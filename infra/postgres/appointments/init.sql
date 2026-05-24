-- Baseline schema. Flyway migrations run after this.
CREATE TABLE IF NOT EXISTS slots (
    id           UUID        PRIMARY KEY,
    doctor_id    UUID        NOT NULL,
    starts_at    TIMESTAMPTZ NOT NULL,
    ends_at      TIMESTAMPTZ NOT NULL,
    is_available BOOLEAN     NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS appointments (
    id         UUID        PRIMARY KEY,
    patient_id UUID        NOT NULL,
    doctor_id  UUID        NOT NULL,
    slot_id    UUID        REFERENCES slots(id),
    starts_at  TIMESTAMPTZ NOT NULL,
    ends_at    TIMESTAMPTZ NOT NULL,
    status     TEXT        NOT NULL DEFAULT 'pending',
    notes      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_slots_doctor_available ON slots(doctor_id, is_available);
CREATE INDEX IF NOT EXISTS idx_appointments_patient   ON appointments(patient_id);
CREATE INDEX IF NOT EXISTS idx_appointments_doctor    ON appointments(doctor_id);
