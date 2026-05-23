# Appointments Service

Servicio responsable de agenda médica, disponibilidad y gestión de citas.

## Stack previsto

- Framework: Spring Boot.
- Protocolo público: REST.
- Base de datos: PostgreSQL `appointments-db`.
- Eventos: Apache Kafka.

## Responsabilidades

- Administrar horarios y slots disponibles de médicos.
- Crear, confirmar, reagendar y cancelar citas.
- Validar conflictos de agenda.
- Coordinar el flujo de cita con pagos y video llamada.
- Publicar eventos relacionados con cambios de estado de citas.

## API esperada

Ruta base por gateway: `/api/appointments`.

- `GET /slots`: buscar disponibilidad.
- `POST /slots`: crear slots para un médico.
- `POST /`: crear cita.
- `GET /{id}`: consultar cita.
- `PATCH /{id}/confirm`: confirmar cita.
- `PATCH /{id}/reschedule`: reagendar cita.
- `PATCH /{id}/cancel`: cancelar cita.
- `GET /patients/{patientId}`: listar citas de un paciente.
- `GET /doctors/{doctorId}`: listar citas de un médico.

## Datos

Tablas iniciales:

- `appointments`
- `slots`

El servicio mantiene su propio modelo de agenda. Los identificadores de pacientes y médicos vienen del User Service, pero no se hacen joins directos contra su base.

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

## Observabilidad

Debe exponer `/metrics` para Prometheus con métricas de citas creadas, cancelaciones, latencia de búsqueda de slots y conflictos de agenda.
