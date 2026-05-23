# Plan De Desarrollo: Appointments Service

## Objetivo

Implementar agenda medica, disponibilidad y ciclo de vida de citas.

## Fases

1. Crear proyecto Spring Boot con config por env.
2. Implementar migraciones para slots y appointments.
3. Crear endpoints de disponibilidad.
4. Crear, confirmar, reagendar y cancelar citas.
5. Consumir eventos de pago para confirmar flujo.
6. Publicar eventos de cambios de cita.

## Entregables Minimos

- Servicio reemplazando el stub en compose mediante `build: ./services/appointments`.
- CRUD minimo de slots.
- Creacion y consulta de citas.
- Validacion de conflictos de horario.
- Eventos Kafka de cita.
- `/health` y `/metrics`.

## Dependencias

- PostgreSQL `appointments-db`.
- Kafka.
- User Service para validar pacientes y medicos.
- Payments Service para estado de pago.

## Endpoints Esperados

- `GET /slots`
- `POST /slots`
- `POST /`
- `GET /{id}`
- `PATCH /{id}/confirm`
- `PATCH /{id}/reschedule`
- `PATCH /{id}/cancel`
- `GET /patients/{patientId}`
- `GET /doctors/{doctorId}`

Ruta gateway: `/api/appointments`.

## Eventos Kafka

Publica:

- `appointment-created`
- `appointment-confirmed`
- `appointment-rescheduled`
- `appointment-cancelled`

Consume:

- `payment-confirmed`
- `payment-failed`
- `call-ended`

## Criterios De Aceptacion

- No se pueden crear dos citas en el mismo slot.
- Una cita creada emite `appointment-created`.
- Una cita pagada puede pasar a confirmada.
- Las cancelaciones liberan disponibilidad cuando aplique.
- `/health` y `/metrics` responden.

## Orden Recomendado

1. Bootstrap Spring Boot.
2. Modelo de slots.
3. Modelo de citas.
4. Integracion Kafka.
5. Observabilidad y tests.
