# Plan Maestro De Desarrollo

## Objetivo

Construir PatitoMedi como una plataforma de telemedicina basada en microservicios, con gateway central, bases de datos por servicio, eventos Kafka, video llamada WebRTC y observabilidad desde el inicio.

## Fases

1. Preparar base tecnica local: variables, healthchecks, documentacion y compose estable.
2. Implementar User Service para identidad, JWT, pacientes y medicos.
3. Implementar Appointments Service para agenda, slots y citas.
4. Implementar Payments Service con pagos simulados antes de integrar proveedor real.
5. Implementar Medical History Service con GraphQL y MongoDB.
6. Implementar Video Call Service con WebSocket, Redis y coturn.
7. Implementar cliente web MVP y dejar mobile para una fase posterior.
8. Endurecer seguridad, pruebas, metricas, dashboards y despliegue.

## Entregables Minimos

- Infraestructura local reproducible con Docker Compose.
- Contratos HTTP, GraphQL, WebSocket y Kafka documentados.
- Servicios con `/health` y `/metrics`.
- Flujo MVP: registro, login, agenda, cita, pago simulado, historial y sala de video.
- Documentacion para desarrollo local y planes por modulo.

## Dependencias

- User Service desbloquea autenticacion y entidades base.
- Appointments Service depende de usuarios existentes.
- Payments Service depende de citas creadas.
- Medical History Service depende de pacientes y medicos autorizados.
- Video Call Service depende de citas confirmadas.

## Criterios De Aceptacion

- `docker compose config --quiet` no reporta errores.
- El entorno local levanta con stubs mientras no existan servicios reales.
- Cada modulo tiene un `PLAN.md` con fases, interfaces y criterios de aceptacion.
- Las rutas del gateway responden a traves de Nginx/Kong.

## Orden Recomendado

1. Infraestructura local.
2. User Service.
3. Appointments Service.
4. Gateway auth y politicas Kong.
5. Payments Service.
6. Medical History Service.
7. Video Call Service.
8. Cliente web.
9. Observabilidad completa y pruebas end-to-end.
