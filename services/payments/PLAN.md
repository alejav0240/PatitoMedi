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
