# Payments Service

Servicio responsable de cobros, facturación y transacciones asociadas a citas médicas.

## Stack previsto

- Framework: Django.
- Protocolo público: REST.
- Base de datos: PostgreSQL `payments-db`.
- Eventos: Apache Kafka.

## Responsabilidades

- Crear facturas para citas.
- Procesar pagos mediante un proveedor externo.
- Registrar transacciones y auditoría financiera.
- Confirmar o rechazar pagos para continuar el flujo de cita.
- Publicar eventos financieros para otros servicios.

## API esperada

Ruta base por gateway: `/api/payments`.

- `POST /invoices`: crear factura.
- `GET /invoices/{id}`: consultar factura.
- `POST /transactions`: iniciar pago.
- `GET /transactions/{id}`: consultar transacción.
- `POST /webhooks/provider`: recibir webhooks del proveedor de pagos.
- `POST /refunds`: solicitar devolución.

## Datos

Tablas iniciales:

- `invoices`
- `transactions`

Este servicio es el dueño de datos financieros. Otros servicios deben reaccionar mediante eventos, no consultar esta base directamente.

## Eventos Kafka

Publica:

- `invoice-created`
- `payment-confirmed`
- `payment-failed`
- `refund-created`

Consume:

- `appointment-created`
- `appointment-cancelled`

## Observabilidad

Debe exponer `/metrics` para Prometheus con métricas de pagos aprobados, pagos fallidos, latencia del proveedor externo y montos procesados.
