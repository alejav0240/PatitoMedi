# Plan De Infraestructura

## Objetivo

Proveer una base local estable para gateway, bases de datos, mensajeria, cache, STUN/TURN y proxy publico.

## Fases

1. Parametrizar `docker-compose.yml` con `.env`.
2. Agregar healthchecks para infraestructura critica.
3. Mantener stubs de microservicios hasta que exista codigo real.
4. Preparar rutas Kong para REST, GraphQL y WebSocket.
5. Agregar topicos Kafka definidos y scripts de inicializacion cuando los consumidores existan.

## Entregables Minimos

- Compose validable con `docker compose config --quiet`.
- `.env.example` con puertos y credenciales locales.
- Nginx con soporte para WebSocket.
- Kong declarativo con CORS, rate limiting y Prometheus.
- PostgreSQL por servicio, MongoDB, Redis, Kafka y coturn.

## Dependencias

- Los servicios reales reemplazaran stubs `hashicorp/http-echo`.
- Kong depende de que los servicios esten disponibles en la red `backend`.
- Video Call Service depende de Redis y coturn.
- Prometheus depende de endpoints `/metrics`.

## Interfaces Esperadas

- Entrada publica: `http://localhost`.
- Kong proxy: `http://localhost:8000`.
- Kong admin: `http://localhost:8001`.
- Kafka interno: `kafka:9092`.
- Redis interno: `redis:6379`.
- coturn local: `localhost:3478`.

## Eventos Kafka Relevantes

- `user-registered`
- `appointment-created`
- `payment-confirmed`
- `record-updated`
- `call-started`
- `call-ended`

## Criterios De Aceptacion

- Los contenedores de infraestructura reportan estado healthy donde aplique.
- Nginx enruta hacia Kong.
- Kong enruta a los stubs actuales.
- Las bases de datos crean sus esquemas iniciales.

## Orden Recomendado

1. Validar compose y variables.
2. Levantar infraestructura.
3. Probar gateway.
4. Reemplazar servicios stub por builds reales servicio por servicio.
