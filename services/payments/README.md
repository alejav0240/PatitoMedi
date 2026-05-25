# Payments Service

Microservicio de pagos de PatitoMedi. Gestiona facturas, transacciones, reembolsos y webhooks simulados para desbloquear el flujo de citas médicas mientras se integra un proveedor real.

## Resumen técnico

- Framework: Django 4.2.
- Protocolo público: REST.
- Base de datos: PostgreSQL `payments-db`.
- Eventos: Apache Kafka.
- Observabilidad: `/health` y `/metrics` con Prometheus.
- Puerto interno del contenedor: `8080`.

## Rol en la arquitectura

El servicio vive detrás de Kong en la ruta base `/api/payments` y se ejecuta como `payments-service` dentro de Docker Compose. Es el dueño de los datos financieros y no debe ser consultado directamente por otros servicios para leer su base de datos.

Flujo principal:

1. Appointments publica `appointment-created` o `appointment-cancelled`.
2. Payments crea o cancela la factura asociada.
3. El pago simulado produce `payment-confirmed` o `payment-failed`.
4. Otros servicios reaccionan a esos eventos.

## Stack y runtime

- `python:3.11-slim` como imagen base.
- Dependencias principales: `django`, `psycopg2-binary`, `kafka-python`, `prometheus_client`.
- Contenedor de app: `payments-service`.
- Base de datos local: `payments-db` en el puerto `5435` del host.

## Variables de entorno

El servicio usa variables estándar por entorno y se integra con `docker-compose.yml`.

- `DJANGO_DEBUG`: activa modo debug.
- `DJANGO_SECRET_KEY`: clave de Django para desarrollo.
- `PAYMENTS_DB_ENGINE`: por defecto `django.db.backends.postgresql`.
- `PAYMENTS_DB_NAME`: por defecto `payments`.
- `PAYMENTS_DB_USER`: por defecto `payments_app`.
- `PAYMENTS_DB_PASSWORD`: por defecto `payments_pass`.
- `PAYMENTS_DB_HOST`: por defecto `payments-db`.
- `PAYMENTS_DB_PORT`: por defecto `5432`.
- `KAFKA_BROKERS`: por defecto `kafka:9092`.
- `KAFKA_ENABLED`: activa o desactiva publicación/consumo Kafka.
- `PAYMENTS_CONSUMER_ENABLED`: activa el consumidor en background.

## API pública

Ruta base por gateway: `/api/payments`.

### Facturas

- `POST /api/payments/invoices`: crea una factura asociada a una cita.
- `GET /api/payments/invoices/{id}`: consulta una factura.

Request mínimo para crear una factura:

```json
{
	"appointment_id": "apt-1",
	"patient_id": "pat-1",
	"doctor_id": "doc-1",
	"amount": "120.50",
	"currency": "PEN"
}
```

### Transacciones

- `POST /api/payments/transactions`: inicia el pago simulado de una factura.
- `GET /api/payments/transactions/{id}`: consulta una transacción.

Request mínimo para simular pago:

```json
{
	"invoice_id": "uuid-de-la-factura",
	"simulate_outcome": "approved"
}
```

Valores admitidos para `simulate_outcome`:

- `approved`
- `failed`

### Webhook del proveedor

- `POST /api/payments/webhooks/provider`: recibe eventos de pago o reembolso.

Eventos admitidos por el stub actual:

- `payment.succeeded`
- `payment.confirmed`
- `payment.failed`
- `refund.succeeded`

### Reembolsos

- `POST /api/payments/refunds`: solicita una devolución para una transacción.

Request mínimo:

```json
{
	"transaction_id": "uuid-de-la-transaccion",
	"reason": "customer_request"
}
```

### Operación

- `GET /health`: comprueba conectividad con PostgreSQL.
- `GET /metrics`: expone métricas Prometheus.

## Modelo de datos

Tablas iniciales administradas por el servicio:

- `invoices`
- `transactions`
- `refunds`

### Invoice

- `id`
- `appointment_id`
- `patient_id`
- `doctor_id`
- `amount`
- `currency`
- `status`
- `external_reference`
- `metadata`
- `created_at`
- `updated_at`

### Transaction

- `id`
- `invoice_id`
- `provider_reference`
- `amount`
- `currency`
- `status`
- `provider_payload`
- `created_at`
- `updated_at`

### Refund

- `id`
- `transaction_id`
- `amount`
- `reason`
- `status`
- `provider_reference`
- `provider_payload`
- `created_at`

## Eventos Kafka

El servicio publica y consume eventos del dominio de pagos.

### Consume

- `appointment-created`
- `appointment-cancelled`

### Publica

- `invoice-created`
- `payment-confirmed`
- `payment-failed`
- `refund-created`

## Observabilidad

`/metrics` expone métricas como:

- requests HTTP por endpoint, método y status.
- total de facturas creadas.
- total de transacciones por estado.
- total de reembolsos.
- latencia de operaciones simuladas del proveedor.

`/health` realiza una comprobación simple contra la base de datos y retorna `200` si está disponible.

## Ejecución local con Docker

El servicio está preparado para correr como parte del stack completo:

```bash
docker compose up -d payments-db kafka kafka-init payments-service
docker compose logs -f payments-service
```

La ruta pública queda disponible a través de Kong y Nginx según el `docker-compose.yml` del repositorio.

## Ejecución local sin Docker

Si se quiere correr el servicio directamente:

```bash
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
export PAYMENTS_DB_ENGINE=django.db.backends.sqlite3
export PAYMENTS_DB_NAME=:memory:
python manage.py migrate --noinput
python manage.py test payments_app -v 2
python manage.py runserver 0.0.0.0:8080
```

## Pruebas

La suite principal valida:

- `GET /health`
- `GET /metrics`
- flujo de factura -> transacción -> reembolso
- consumo de eventos de `Appointments`

Ejecutar con Docker:

```bash
docker compose run --rm \
	-e PAYMENTS_DB_ENGINE=django.db.backends.sqlite3 \
	-e PAYMENTS_DB_NAME=:memory: \
	payments-service python manage.py test payments_app -v 2
```

## Endpoints de referencia

- `POST /api/payments/invoices`
- `GET /api/payments/invoices/{id}`
- `POST /api/payments/transactions`
- `GET /api/payments/transactions/{id}`
- `POST /api/payments/webhooks/provider`
- `POST /api/payments/refunds`
- `GET /health`
- `GET /metrics`

## Notas de implementación

- El consumidor de Kafka se activa en background cuando `KAFKA_ENABLED=true` y `PAYMENTS_CONSUMER_ENABLED=true`.
- Si Kafka no está disponible, el servicio sigue levantando y registra un aviso en logs.
- La app usa migraciones de Django y un `Dockerfile` propio para integrarse con Compose.
