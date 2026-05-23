# Plan De Observabilidad

## Objetivo

Dar visibilidad operativa a PatitoMedi mediante metricas, dashboards, logs estructurados y alertas basicas.

## Fases

1. Levantar Prometheus y Grafana por compose.
2. Exponer `/metrics` en cada servicio real.
3. Crear dashboards por infraestructura y dominio.
4. Estandarizar logs JSON con request id.
5. Agregar alertas para errores, latencia y salud de llamadas.

## Entregables Minimos

- Prometheus operativo en `http://localhost:9090`.
- Grafana operativo en `http://localhost:3000`.
- Datasource Prometheus provisionado.
- Metricas de gateway y servicios.
- Dashboard inicial por servicio cuando exista implementacion real.

## Dependencias

- Kong debe tener plugin Prometheus habilitado.
- Cada microservicio debe publicar `/metrics`.
- Video Call Service debe publicar metricas especificas de WebRTC.

## Interfaces Esperadas

- Prometheus: `/targets`, `/graph`.
- Grafana: dashboard de servicios y dashboard de infraestructura.
- Servicios: `/metrics` y `/health`.

## Eventos Kafka Relevantes

Las metricas deben permitir observar volumen y errores derivados de:

- `appointment-created`
- `payment-confirmed`
- `record-updated`
- `call-started`
- `call-ended`

## Criterios De Aceptacion

- Prometheus ve sus targets.
- Grafana carga datasource automaticamente.
- Cada servicio real expone latencia, errores, throughput y salud.
- Video Call Service reporta peers conectados, salas activas y fallas ICE.

## Orden Recomendado

1. Validar Prometheus y Grafana.
2. Instrumentar gateway.
3. Instrumentar servicios segun se implementen.
4. Crear dashboards y alertas.
