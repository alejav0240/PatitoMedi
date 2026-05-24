# Appointments Service

Servicio responsable de agenda médica, disponibilidad de slots y ciclo de vida de citas en PatitoMedi.

## Stack

- Lenguaje: Java 21.
- Framework: Spring Boot 3.3.
- Protocolo público: REST.
- Base de datos: PostgreSQL `appointments-db`.
- Migraciones: Flyway.
- Eventos: Apache Kafka.

## Responsabilidades

- Administrar slots de disponibilidad de médicos.
- Crear, confirmar, reagendar y cancelar citas.
- Garantizar que no se creen dos citas en el mismo slot (transacción con bloqueo).
- Liberar slots al cancelar o reagendar.
- Publicar eventos de cambio de estado de citas.
- Consumir eventos de pago para confirmar o cancelar citas automáticamente.

## API

Ruta base por gateway: `/api/appointments`.

### Slots

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/slots` | Listar slots. Params: `doctorId` (UUID), `available` (bool, default `true`) |
| `POST` | `/slots` | Crear slot para un médico |

### Citas

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/` | Crear cita |
| `GET` | `/{id}` | Consultar cita |
| `PATCH` | `/{id}/confirm` | Confirmar cita |
| `PATCH` | `/{id}/reschedule` | Reagendar cita a otro slot |
| `PATCH` | `/{id}/cancel` | Cancelar cita |
| `GET` | `/patients/{patientId}` | Listar citas de un paciente |
| `GET` | `/doctors/{doctorId}` | Listar citas de un médico |

## Datos

Tablas gestionadas por Flyway:

- `slots` — disponibilidad de médicos.
- `appointments` — citas con referencia al slot.

El servicio es dueño de su base de datos. Los `patient_id` y `doctor_id` son UUIDs provenientes del User Service, pero no se hacen joins directos contra su base.

## Eventos Kafka

Publica:

- `appointment-created`
- `appointment-confirmed`
- `appointment-rescheduled`
- `appointment-cancelled`

Consume:

- `payment-confirmed` → confirma la cita automáticamente.
- `payment-failed` → cancela la cita automáticamente.
- `call-ended` → registrado para uso futuro.

## Observabilidad

- `GET /actuator/health` — estado del servicio y conexión a base de datos.
- `GET /actuator/prometheus` — métricas en formato Prometheus.

## Desarrollo local

```bash
docker compose build appointments-service
docker compose up -d appointments-db kafka appointments-service
```

Variables de entorno:

| Variable | Default | Descripción |
|----------|---------|-------------|
| `PORT` | `8080` | Puerto del servidor |
| `DATABASE_URL` | `jdbc:postgresql://appointments-db:5432/appointments` | URL JDBC |
| `DATABASE_USER` | `appointments_app` | Usuario PostgreSQL |
| `DATABASE_PASSWORD` | `appointments_pass` | Contraseña PostgreSQL |
| `KAFKA_BROKERS` | `kafka:9092` | Brokers Kafka separados por coma |

## Ejemplos

Crear un slot:

```bash
curl -X POST http://localhost/api/appointments/slots \
  -H 'Content-Type: application/json' \
  -d '{
    "doctorId": "00000000-0000-0000-0000-000000000001",
    "startsAt": "2026-06-01T10:00:00Z",
    "endsAt":   "2026-06-01T10:30:00Z"
  }'
```

Crear una cita:

```bash
curl -X POST http://localhost/api/appointments/ \
  -H 'Content-Type: application/json' \
  -d '{
    "patientId": "00000000-0000-0000-0000-000000000002",
    "doctorId":  "00000000-0000-0000-0000-000000000001",
    "slotId":    "<slot-id>",
    "notes":     "Primera consulta"
  }'
```

Confirmar una cita:

```bash
curl -X PATCH http://localhost/api/appointments/<id>/confirm
```

Reagendar:

```bash
curl -X PATCH http://localhost/api/appointments/<id>/reschedule \
  -H 'Content-Type: application/json' \
  -d '{"slotId": "<nuevo-slot-id>"}'
```

Cancelar:

```bash
curl -X PATCH http://localhost/api/appointments/<id>/cancel
```
