# PatitoMedi Architecture

This repository is organized around the architecture described in the README:

- `clients/web`: web client.
- `clients/mobile`: mobile client.
- `services/user`: Golang user/auth service.
- `services/appointments`: Spring Boot appointments service.
- `services/payments`: Django payments service.
- `services/medical-history`: Express.js GraphQL medical history service.
- `services/video-call`: Golang WebRTC signaling service.
- `infra/nginx`: public reverse proxy.
- `infra/kong`: API gateway routes and plugins.
- `infra/postgres`: database-per-service initialization scripts.
- `infra/coturn`: STUN/TURN server configuration for WebRTC.
- `monitoring/prometheus`: metrics scraping config.
- `monitoring/grafana`: Grafana provisioning.

## Local ports

- `80`: Nginx public entrypoint.
- `8000`: Kong proxy.
- `8001`: Kong admin API.
- `5433`: users PostgreSQL.
- `5434`: appointments PostgreSQL.
- `5435`: payments PostgreSQL.
- `27017`: MongoDB medical records.
- `6379`: Redis call sessions.
- `9092`: Kafka.
- `3478`: coturn STUN/TURN.
- `9090`: Prometheus.
- `3000`: Grafana.

## Gateway routes

- `GET /api/users`: user service stub.
- `GET /api/appointments`: appointments service stub.
- `GET /api/payments`: payments service route behind Kong.
- `GET /graphql/medical-history`: medical history GraphQL service stub.
- `GET /ws/video`: video signaling service stub.

The current compose file uses lightweight HTTP echo containers for some domain services, but `services/payments` is implemented and exposed through Kong at `/api/payments/*`. Replace the remaining stubs with a `build:` section when each real service implementation exists.
