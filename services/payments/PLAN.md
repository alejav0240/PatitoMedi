# Payments Service — PLAN

## Objetivo

Implementar el `Payments Service` como un servicio Django para cobros, facturación y transacciones asociadas a citas médicas, siguiendo el stack documentado en `docs/architecture.md` y el plan del propio servicio.

## Análisis del proyecto

- Lenguaje: Django.
- Protocolo público: REST.
- Base de datos: PostgreSQL `payments-db`.
- Mensajería: Apache Kafka.
- Requisitos mínimos: endpoints `/invoices`, `/transactions`, `/webhooks/provider`, `/refunds`, `/health`, `/metrics`.
- Observabilidad: métricas Prometheus con pagos aprobados, fallidos, latencia y montos.
- Despliegue local: `Dockerfile` y servicio `payments-service` en `docker-compose.yml`.

## Dependencias internas

- `Appointments Service`: para asociar pagos a citas.
- `Kafka`: para eventos `appointment-created`, `appointment-cancelled`, `invoice-created`, `payment-confirmed`, `payment-failed`, `refund-created`.

## API — Especificación mínima

1) POST /api/payments/invoices
- Crea una factura asociada a una cita.
- Requiere `appointment_id`, `amount` y `currency`.

2) GET /api/payments/invoices/{id}
- Consulta una factura existente.

3) POST /api/payments/transactions
- Inicia el pago simulado de una factura.
- Requiere `invoice_id` y admite `simulate_outcome` (`approved` o `failed`).

4) GET /api/payments/transactions/{id}
- Consulta una transacción existente.

5) POST /api/payments/webhooks/provider
- Recibe notificaciones del proveedor simulado.

6) POST /api/payments/refunds
- Solicita una devolución para una transacción.

7) GET /health
- Verifica conectividad con PostgreSQL.

8) GET /metrics
- Expone métricas Prometheus.

## Modelo de datos (mínimo)

- Invoice: `id`, `appointment_id`, `patient_id`, `doctor_id`, `amount`, `currency`, `status`, `created_at`, `updated_at`.
- Transaction: `id`, `invoice_id`, `provider_reference`, `amount`, `currency`, `status`, `created_at`, `updated_at`.
- Refund: `id`, `transaction_id`, `amount`, `reason`, `status`, `created_at`.

## Eventos Kafka

- Consumir: `appointment-created`, `appointment-cancelled`.
- Producir: `invoice-created`, `payment-confirmed`, `payment-failed`, `refund-created`.

## Docker y Compose

- `Dockerfile`: imagen basada en `python:3.11-slim`, instalar dependencias de `requirements.txt`, exponer puerto `8080`.
- En `docker-compose.yml` cambiar el stub por `build: ./services/payments` y mantener el nombre `payments-service`.

## Healthchecks y métricas

- `/health` debe comprobar la conexión a la BD y retornar 200.
- `/metrics` expone métricas Prometheus (usando `prometheus_client`).

## Tests y criterios de aceptación

- `docker compose config --quiet` no debe fallar después de reemplazar el stub por el build local.
- Endpoints `/api/payments/invoices`, `/api/payments/transactions`, `/api/payments/refunds` y `/api/payments/webhooks/provider` devuelven respuestas JSON válidas.
- `GET /health` devuelve 200 con la base de datos accesible.

## Roadmap de implementación (short-term)

1. Crear `services/payments/PLAN.md` y la especificación (este documento).
2. Scaffold mínimo Django con endpoints `/invoices`, `/transactions`, `/webhooks/provider`, `/refunds`, `/health`, `/metrics`.
3. Añadir `Dockerfile` y `requirements.txt`.
4. Reemplazar el stub de `payments-service` en `docker-compose.yml` por `build: ./services/payments`.
5. Integrar consumo de eventos de Appointments y emisión de eventos financieros.

## Notas

- Puerto interno del servicio: `8080` para coincidir con Kong y Prometheus.
- Mantener variables sensibles en `.env` y `env_file` en `docker-compose`.

---

Archivo creado automáticamente por la planificación de la Fase 4.
# Plan De Desarrollo: Payments Service

## Objetivo

Implementar facturacion y pagos, empezando con un proveedor simulado para desbloquear el MVP.

## Fases

1. Crear proyecto Django con config por env.
2. Implementar modelos de invoices y transactions.
3. Crear endpoints de factura y pago simulado.
4. Publicar eventos de pago aprobado o fallido.
5. Preparar webhook para proveedor real.
6. Agregar observabilidad y pruebas.

## Entregables Minimos

- Servicio reemplazando el stub en compose mediante `build: ./services/payments`.
- Factura asociada a cita.
- Pago simulado configurable como aprobado o fallido.
- Registro de transaccion.
- Eventos `payment-confirmed` y `payment-failed`.
- `/health` y `/metrics`.

## Dependencias

- PostgreSQL `payments-db`.
- Kafka.
- Appointments Service para citas creadas.

## Endpoints Esperados

- `POST /invoices`
- `GET /invoices/{id}`
- `POST /transactions`
- `GET /transactions/{id}`
- `POST /webhooks/provider`
- `POST /refunds`

Ruta gateway: `/api/payments`.

## Eventos Kafka

Publica:

- `invoice-created`
- `payment-confirmed`
- `payment-failed`
- `refund-created`

Consume:

- `appointment-created`
- `appointment-cancelled`

## Criterios De Aceptacion

- Una cita puede generar una factura.
- Un pago simulado aprobado emite `payment-confirmed`.
- Un pago simulado fallido emite `payment-failed`.
- Las transacciones quedan auditadas.
- `/health` y `/metrics` responden.

## Orden Recomendado

1. Bootstrap Django.
2. Modelos y migraciones.
3. API de facturas.
4. API de transacciones.
5. Eventos Kafka.
6. Observabilidad y tests.
