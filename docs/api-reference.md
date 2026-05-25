# PatitoMedi — API Reference

Todos los servicios se acceden a través de Nginx (`http://localhost`) → Kong (`http://localhost:8000`).

---

## User Service

**Stack:** Go · PostgreSQL · Kafka  
**Gateway base:** `/api/users`  
**Estado:** ✅ Implementado

### Endpoints

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| `POST` | `/api/users/register/patient` | No | Registrar paciente |
| `POST` | `/api/users/register/doctor` | No | Registrar médico |
| `POST` | `/api/users/login` | No | Login, retorna JWT |
| `POST` | `/api/users/logout` | JWT | Cerrar sesión |
| `GET` | `/api/users/me` | JWT | Perfil del usuario autenticado |
| `PATCH` | `/api/users/me` | JWT | Actualizar perfil |
| `GET` | `/api/users/doctors` | No | Listar médicos |
| `GET` | `/api/users/doctors/{id}` | No | Obtener médico por ID |
| `GET` | `/api/users/health` | No | Health check |
| `GET` | `/api/users/metrics` | No | Métricas Prometheus |

### Request bodies

**POST /register/patient**
```json
{ "email": "ana@example.com", "password": "password123", "full_name": "Ana Pérez" }
```

**POST /register/doctor**
```json
{ "email": "dr@example.com", "password": "password123", "full_name": "Dr. Gómez", "specialty": "Cardiología" }
```

**POST /login**
```json
{ "email": "ana@example.com", "password": "password123" }
```
Respuesta: `{ "token": "<jwt>", "user": { "id", "email", "full_name", "role", "created_at" } }`

**PATCH /me**
```json
{ "full_name": "Ana Pérez López" }
```
Para médicos también acepta `"specialty"`.

### Auth
Header: `Authorization: Bearer <token>`

---

## Appointments Service

**Stack:** Java 21 · Spring Boot 3.3 · PostgreSQL · Flyway · Kafka  
**Gateway base:** `/api/appointments`  
**Estado:** ✅ Implementado

### Slots

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/api/appointments/slots` | Crear slot de disponibilidad |
| `GET` | `/api/appointments/slots` | Listar slots. Params: `doctorId` (UUID), `available` (bool, default `true`) |

**POST /slots — body**
```json
{
  "doctorId": "<uuid>",
  "startsAt": "2026-06-10T10:00:00Z",
  "endsAt":   "2026-06-10T10:30:00Z"
}
```

### Citas

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/api/appointments/` | Crear cita (ocupa el slot) |
| `GET` | `/api/appointments/{id}` | Obtener cita por ID |
| `PATCH` | `/api/appointments/{id}/confirm` | Confirmar cita |
| `PATCH` | `/api/appointments/{id}/reschedule` | Reagendar a otro slot |
| `PATCH` | `/api/appointments/{id}/cancel` | Cancelar (libera el slot) |
| `GET` | `/api/appointments/patients/{patientId}` | Citas de un paciente |
| `GET` | `/api/appointments/doctors/{doctorId}` | Citas de un médico |

**POST / — body**
```json
{
  "patientId": "<uuid>",
  "doctorId":  "<uuid>",
  "slotId":    "<uuid>",
  "notes":     "Primera consulta"
}
```

**PATCH /{id}/reschedule — body**
```json
{ "slotId": "<nuevo-slot-uuid>" }
```

### Observabilidad

| Ruta | Descripción |
|------|-------------|
| `GET /api/appointments/actuator/health` | Health + estado DB |
| `GET /api/appointments/actuator/prometheus` | Métricas Prometheus |

### Códigos de respuesta

| Código | Situación |
|--------|-----------|
| `201` | Recurso creado |
| `400` | Campos requeridos faltantes |
| `404` | Slot o cita no encontrada |
| `409` | Slot ya ocupado |

---

## Payments Service

**Stack:** Python · Django · PostgreSQL · Kafka  
**Gateway base:** `/api/payments`  
**Estado:** ✅ Implementado

El servicio responde a través de Kong en `/api/payments/*` y, dentro del contenedor, expone las mismas rutas sin ese prefijo. Para pruebas locales directas contra el contenedor, usa `/health`, `/metrics`, `/invoices`, `/transactions`, `/refunds` y `/webhooks/provider`.

### Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/api/payments/invoices` | Crear factura para una cita |
| `GET` | `/api/payments/invoices/{id}` | Obtener factura |
| `POST` | `/api/payments/transactions` | Procesar pago (simulado) |
| `GET` | `/api/payments/transactions/{id}` | Obtener transacción |
| `POST` | `/api/payments/refunds` | Crear reembolso |
| `POST` | `/api/payments/webhooks/provider` | Webhook de proveedor externo |
| `GET` | `/health` | Health check |
| `GET` | `/metrics` | Métricas Prometheus |

**POST /invoices — body**
```json
{
  "appointment_id": "<uuid>",
  "patient_id":     "<uuid>",
  "doctor_id":      "<uuid>",
  "amount":         "50.00",
  "currency":       "USD",
  "metadata":       {}
}
```

**POST /transactions — body**
```json
{
  "invoice_id":         "<uuid>",
  "simulate_outcome":    "approved",
  "provider_reference":  "prov_abc123",
  "provider_payload":    {}
}
```
`simulate_outcome`: `"approved"` | `"failed"`

**POST /refunds — body**
```json
{
  "transaction_id":   "<uuid>",
  "amount":           "50.00",
  "reason":           "Cita cancelada",
  "provider_reference": "refund_abc123",
  "provider_payload":  {}
}
```

**POST /webhooks/provider — body**
```json
{
  "event_type": "payment.confirmed",
  "payload": {
    "transaction_id": "<uuid>",
    "provider_reference": "prov_abc123"
  }
}
```
`event_type`: `payment.succeeded` | `payment.confirmed` | `payment.failed` | `refund.succeeded`

### Webhook behavior

- `payment.succeeded` y `payment.confirmed`: marcan la transacción como `approved` y la factura asociada como `paid`.
- `payment.failed`: marca la transacción como `failed`.
- `refund.succeeded`: marca el reembolso como `processed`.

### Statuses

| Entidad | Valores |
|---------|---------|
| Invoice | `issued` → `paid` / `cancelled` |
| Transaction | `pending` → `approved` / `failed` / `refunded` |
| Refund | `requested` → `processed` / `failed` |

### Observaciones de integración

- `POST /api/payments/invoices` crea la factura y publica `invoice-created`.
- `POST /api/payments/transactions` publica `payment-confirmed` o `payment-failed` según `simulate_outcome`.
- `POST /api/payments/refunds` publica `refund-created`.
- `GET /health` comprueba conectividad con PostgreSQL.
- `GET /metrics` expone contadores HTTP, facturas, transacciones, reembolsos y latencia del simulador.
## Medical History Service

**Stack:** Node.js · Express · Apollo Server · MongoDB · Kafka  
**Gateway base:** `/graphql/medical-history`  
**Protocolo:** GraphQL (POST)  
**Estado:** ✅ Implementado

Más detalles y ejemplos de queries/mutations: [Medical History docs](docs/medical-history.md)

### Endpoint

```
POST /graphql/medical-history
Content-Type: application/json
```

### Queries

**Listar registros de un paciente**
```graphql
query {
  patientRecords(patientId: "<id>") {
    id
    patientId
    bloodType
    allergies
    chronicDiseases
    medications
    notes
    createdAt
    updatedAt
  }
}
```

**Obtener registro por ID**
```graphql
query {
  record(id: "<id>") {
    id
    patientId
    bloodType
    allergies
    chronicDiseases
    medications
    notes
  }
}
```

**Health check GraphQL**
```graphql
query {
  healthGraphQL
}
```

### Mutations

**Crear registro clínico**
```graphql
mutation {
  createRecord(input: {
    patientId:       "<id>"
    bloodType:       "O+"
    allergies:       ["penicilina"]
    chronicDiseases: ["hipertensión"]
    medications:     ["enalapril 10mg"]
    notes:           ["Paciente estable"]
  }) {
    id
    createdAt
  }
}
```

**Actualizar registro clínico**
```graphql
mutation {
  updateRecord(id: "<id>", input: {
    patientId:   "<id>"
    medications: ["enalapril 20mg"]
    notes:       ["Dosis ajustada"]
  }) {
    id
    updatedAt
  }
}
```

### Observabilidad

| Ruta | Descripción |
|------|-------------|
| `GET /health` | Health check REST |

---

## Video Call Service

**Stack:** Go · WebSocket · Redis · coturn · Kafka  
**Gateway base:** `/ws/video`  
**Estado:** ✅ Implementado

Más detalles, flujo completo y configuración ICE: [Video Call docs](video-call.md)

### Protocolo WebSocket

Conexión: `ws://localhost:8000/ws/video`

### Mensajes (JSON)

| Tipo | Dirección | Descripción |
|------|-----------|-------------|
| `join-room` | Cliente → Servidor | Unirse a sala de una cita |
| `leave-room` | Cliente → Servidor | Abandonar sala |
| `offer` | Cliente → Servidor | Oferta SDP WebRTC (enrutada al peer `to`) |
| `answer` | Cliente → Servidor | Respuesta SDP WebRTC (enrutada al peer `to`) |
| `ice-candidate` | Cliente ↔ Servidor | Candidato ICE (enrutado al peer `to`) |
| `call-ended` | Cliente → Servidor | Finalizar llamada |
| `error` | Servidor → Cliente | Error de señalización |

**join-room**
```json
{ "type": "join-room", "appointmentId": "<uuid>", "userId": "<uuid>" }
```

**offer / answer**
```json
{ "type": "offer", "sdp": "<sdp-string>", "to": "<peer-userId>" }
```

**ice-candidate**
```json
{ "type": "ice-candidate", "candidate": "<candidate-string>", "sdpMid": "0", "sdpMLineIndex": 0, "to": "<peer-userId>" }
```

### Observabilidad

| Ruta | Descripción |
|------|-------------|
| `GET /health` | `{"status":"ok","service":"video-call"}` |
| `GET /metrics` | `video_active_rooms`, `video_active_participants`, `video_calls_total` |

---

## Resumen de puertos locales

| Servicio | Puerto directo | Ruta gateway |
|----------|---------------|--------------|
| Nginx | `80` | — |
| Kong proxy | `8000` | — |
| Kong admin | `8001` | — |
| User Service | interno | `/api/users` |
| Appointments Service | interno | `/api/appointments` |
| Payments Service | interno | `/api/payments` |
| Medical History | interno | `/graphql/medical-history` |
| Video Call | interno | `/ws/video` |
| Prometheus | `9090` | — |
| Grafana | `3000` | — |
